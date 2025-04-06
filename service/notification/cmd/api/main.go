package main

import (
	"crypto/tls"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/db/postgres"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/platform/pkg/utils"
	"github.com/GCFactory/dbo-system/service/notification/internal/server"
	"github.com/GCFactory/dbo-system/service/notification/pkg/kafka"
	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/opentracing/opentracing-go"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"
)

//	@Title			Users Service
//	@Version		0.1.0
//	@description	Service for users data

//	@contact.name	Rueie
//	@contact.email

//	@license.name	MIT License

//	@BasePath	/api/v1/users/

func main() {
	log.Println("Starting api server")

	configPath := utils.GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	appLogger := logger.NewServerLogger(&config.Config{
		Logger: cfg.Logger,
	})

	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Env: %s, SSL: %v", cfg.Version, cfg.Logger.Level, cfg.Env, cfg.HTTPServer.SSL)

	psqlDB, err := postgres.NewPsqlDB(&config.Config{
		Postgres: cfg.Postgres,
	})
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	} else {
		appLogger.Infof("Postgres connected, Status: %#v", psqlDB.Stats())
	}
	defer psqlDB.Close()

	// driver объект подключения по типу psql
	driver, err := migratePostgres.WithInstance(psqlDB.DB, &migratePostgres.Config{
		MigrationsTable:       "\"schema_migration\"",
		MigrationsTableQuoted: true,
	})
	if err != nil {
		appLogger.Fatalf("Cannot create migration driver: %s", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://migration/postgres",
		cfg.Postgres.PostgresqlDbname, driver)
	if err != nil {
		appLogger.Fatalf("Error on initiate migration: %s", err)
	}

	status := migration.Up()
	if status != nil {
		appLogger.Infof("Migration status: %s", status)
	}
	appLogger.Info("Migration completed")

	var rmqConn *amqp.Connection
	var rmqCh *amqp.Channel

	connect := func() error {
		appLogger.Info("Connecting to RMQ")
		rmqUrl := "amqp://" + cfg.RabbitMQ.User + ":" + cfg.RabbitMQ.Password + "@" +
			cfg.RabbitMQ.Host + ":" + cfg.RabbitMQ.Port

		rmqConn, err = amqp.Dial(rmqUrl)
		if err != nil {
			appLogger.Fatalf("Connecting error to RMQ: %s", err)
			return err
		}
		appLogger.Info("Connecting to RMQ success")

		appLogger.Info("Open RMQ channel")
		rmqCh, err = rmqConn.Channel()
		if err != nil {
			appLogger.Fatalf("Open RMQ channel error: %s", err)
			return err
		}
		appLogger.Info("Open RMQ channel success")

		appLogger.Info("Creating RMQ queue")
		_, err = rmqCh.QueueDeclare(
			cfg.RabbitMQ.Queue, // name
			true,               // durable (сохранять очередь при перезапуске сервера)
			false,              // delete when unused
			false,              // exclusive
			false,              // no-wait
			nil,                // arguments
		)

		if err != nil {
			appLogger.Fatalf("Create RMQ queue error: %s", err)
		} else {
			appLogger.Info("Create RMQ queue success")
		}

		return err
	}

	// Первоначальное подключение
	if err := connect(); err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	// Обработка закрытия соединения
	go func() {
		for {
			reason := <-rmqConn.NotifyClose(make(chan *amqp.Error))
			log.Printf("RabbitMQ connection closed: %v", reason)

			// Пытаемся переподключиться
			for {
				time.Sleep(5 * time.Second)
				if err := connect(); err == nil {
					log.Println("Reconnected to RabbitMQ successfully")
					break
				} else {
					log.Printf("Failed to reconnect to RabbitMQ: %s", err)
				}
			}
		}
	}()

	appLogger.Info("Register RMQ consumer")
	// Потребление сообщений
	msgChan, err := rmqCh.Consume(
		cfg.RabbitMQ.Queue, // queue
		"",                 // consumer
		false,              // auto-ack (false - ручное подтверждение)
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		appLogger.Fatalf("Register RMQ consumer error: %s", err)
		return
	}
	appLogger.Info("Register RMQ consumer success")

	smtpAuth := smtp.PlainAuth(
		cfg.NotificationSmtp.NickName,
		cfg.NotificationSmtp.User,
		cfg.NotificationSmtp.Password,
		cfg.NotificationSmtp.Host)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Проверять сертификат (лучше `false` в продакшене)
		ServerName:         cfg.NotificationSmtp.Host,
	}

	appLogger.Info("Open tls connection with smtp")
	tlsConn, err := tls.Dial("tcp", cfg.NotificationSmtp.Host+":"+cfg.NotificationSmtp.Port, tlsConfig)
	if err != nil {
		appLogger.Fatalf("Open tls connection with smtp error: %s", err)
		return
	}
	defer tlsConn.Close()
	appLogger.Info("Open tls connection with smtp success")

	appLogger.Info("Create smtp client")
	smtpClient, err := smtp.NewClient(tlsConn, cfg.NotificationSmtp.Host)
	if err != nil {
		appLogger.Fatalf("Create smtp client error: %s", err)
		return
	}
	defer smtpClient.Close()
	appLogger.Info("Create smtp client success")

	appLogger.Info("Smtp auth")
	if err = smtpClient.Auth(smtpAuth); err != nil {
		appLogger.Fatalf("Smtp auth error: %s", err)
		return
	}
	appLogger.Info("Smtp auth success")

	kc, err := kafka.NewKafkaConsumer(strings.Split(cfg.KafkaConsumer.Brokers, ";"), cfg.KafkaConsumer.GroupID)
	if err != nil {
		appLogger.Fatal(err)
	}
	appLogger.Infof("Kafka consumer with group '%s' connected", cfg.KafkaConsumer.GroupID)

	kp := kafka.NewKafkaProducer(cfg, appLogger)

	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: cfg.Jaeger.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           cfg.Jaeger.LogSpans,
			LocalAgentHostPort: cfg.Jaeger.Host,
		},
	}

	tracer, closer, err := jaegerCfgInstance.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
	if err != nil {
		log.Fatal("cannot create tracer", err)
	}
	appLogger.Info("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	appLogger.Info("Opentracing connected")

	//Run server
	s := server.NewServer(cfg, kc, kp, psqlDB, msgChan, smtpClient, appLogger)
	if err = s.Run(); err != nil {
		appLogger.Fatal(err)
	}

}

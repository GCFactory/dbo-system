package main

import (
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/db/postgres"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/platform/pkg/utils"
	"github.com/GCFactory/dbo-system/service/notification/internal/server"
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

	appLogger.Info("Connecting to RMQ")
	rmqUrl := "amqp://" +
		cfg.RabbitMQ.User + ":" +
		cfg.RabbitMQ.Password + "@" +
		cfg.RabbitMQ.Host + ":" +
		cfg.RabbitMQ.Port
	rmqConn, err := amqp.Dial(rmqUrl)
	if err != nil {
		appLogger.Fatalf("Connecting error to RMQ: %s", err)
		return
	}
	defer rmqConn.Close()
	appLogger.Info("Connecting to RMQ success")

	appLogger.Info("Open RMQ channel")
	rmqCh, err := rmqConn.Channel()
	if err != nil {
		appLogger.Fatalf("Open RMQ channel error: %s", err)
		return
	}
	defer rmqCh.Close()
	appLogger.Info("Open RMQ channel success")

	appLogger.Info("Creating RMQ queue")
	rmqQueue, err := rmqCh.QueueDeclare(
		cfg.RabbitMQ.Queue, // name
		false,              // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		appLogger.Fatalf("Create RMQ queue error: %s", err)
		return
	}
	appLogger.Info("Create RMQ queue success")

	appLogger.Info("Register RMQ consumer")
	msgChan, err := rmqCh.Consume(
		rmqQueue.Name, // queue
		"",            // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		appLogger.Fatalf("Register RMQ consumer error: %s", err)
		return
	}
	appLogger.Info("Register RMQ consumer success")

	appLogger.Info("Log in to SMTP server")
	smtpAuth := smtp.PlainAuth(
		cfg.NotificationSmtp.NickName,
		cfg.NotificationSmtp.User,
		cfg.NotificationSmtp.Password,
		cfg.NotificationSmtp.Host)
	err = smtp.SendMail(
		cfg.NotificationSmtp.Host+":"+cfg.NotificationSmtp.Port,
		smtpAuth,
		cfg.NotificationSmtp.NickName,
		[]string{
			cfg.NotificationSmtp.NickName,
		},
		[]byte("Auth test msg"))
	if err != nil {
		appLogger.Fatalf("Log in to SMTP server test msg error: %s", err)
		return
	}
	appLogger.Info("Log in to SMTP server success")

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
	s := server.NewServer(cfg, psqlDB, msgChan, smtpAuth, appLogger)
	if err = s.Run(); err != nil {
		appLogger.Fatal(err)
	}

}

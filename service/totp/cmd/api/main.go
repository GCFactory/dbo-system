package main

import (
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/db/postgres"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/platform/pkg/utils"
	"github.com/GCFactory/dbo-system/service/totp/internal/server"
	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"log"
	"os"
)

//	@Title			TOTP Service
//	@Version		0.1.0
//	@description	Service for generating totp codes

//	@contact.name	Rueie
//	@contact.email

//	@license.name	MIT License

//	@BasePath	/api/v1/totp/

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

	appLogger := logger.NewServerLogger(cfg)

	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Env: %s, SSL: %v", cfg.Version, cfg.Logger.Level, cfg.Env, cfg.HTTPServer.SSL)

	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	} else {
		appLogger.Infof("Postgres connected, Status: %#v", psqlDB.Stats())
	}
	defer psqlDB.Close()

	// driver объект подключения по типу psql
	driver, err := migratePostgres.WithInstance(psqlDB.DB, &migratePostgres.Config{
		MigrationsTable:       "\"totp_service\".\"schema_migration\"",
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
	s := server.NewServer(cfg, psqlDB, appLogger)
	if err = s.Run(); err != nil {
		appLogger.Fatal(err)
	}
}

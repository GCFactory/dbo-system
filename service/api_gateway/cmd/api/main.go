package main

import (
	"context"
	"fmt"
	platformConfig "github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/platform/pkg/utils"
	"github.com/GCFactory/dbo-system/service/api_gateway/config"
	"github.com/GCFactory/dbo-system/service/api_gateway/internal/server"
	"github.com/opentracing/opentracing-go"
	redis "github.com/redis/go-redis/v9"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"log"
	"os"
	"time"
)

//	@Title			Registration Service
//	@Version		0.1.0
//	@description	Service for registration user slots

//	@contact.name	Rueie
//	@contact.email

//	@license.name	MIT License

//	@BasePath	/api/v1/registration/

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

	appLogger := logger.NewServerLogger(&platformConfig.Config{
		Logger: cfg.Logger,
	})

	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Env: %s, SSL: %v", cfg.Version, cfg.Logger.Level, cfg.Env, cfg.HTTPServer.SSL)

	redis := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.RedisAddr,
		Password:     cfg.Redis.RedisPassword,
		DB:           cfg.Redis.DB,
		Username:     cfg.Redis.User,
		MaxRetries:   cfg.Redis.MaxRetries,
		DialTimeout:  time.Duration(1000000000 * cfg.Redis.DialTimeout),
		ReadTimeout:  time.Duration(1000000000 * cfg.Redis.Timeout),
		WriteTimeout: time.Duration(1000000000 * cfg.Redis.Timeout),
	})

	if err := redis.Ping(context.Background()).Err(); err != nil {
		fmt.Printf("failed to connect to redis server: %s\n", err.Error())
		return
	}

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
	s := server.NewServer(cfg, redis, appLogger)
	if err = s.Run(); err != nil {
		appLogger.Fatal(err)
	}
}

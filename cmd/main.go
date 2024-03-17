package main

import (
	"fmt"
	logGo "log"
	"notification-service/configs"
	notificationHandler "notification-service/internal/modules/notification/handlers"
	notificationUsecases "notification-service/internal/modules/notification/usecases"
	orderRepoQuery "notification-service/internal/modules/order/repositories/queries"
	"notification-service/internal/pkg/apm"
	"notification-service/internal/pkg/databases/mongodb"
	graceful "notification-service/internal/pkg/gs"
	"notification-service/internal/pkg/helpers"
	kafkaConfluent "notification-service/internal/pkg/kafka/confluent"
	"notification-service/internal/pkg/log"
	"notification-service/internal/pkg/redis"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.elastic.co/apm/module/apmfiber"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// @BasePath	/
func main() {
	// Init Config
	configs.InitConfig()
	// Init Elastic APM Config or DD Apm
	if configs.GetConfig().Datadog.DatadogEnabled == "true" {
		tracer.Start(
			tracer.WithAgentAddr(fmt.Sprintf("%s:%s", configs.GetConfig().Datadog.DatadogHost, configs.GetConfig().Datadog.DatadogPort)),
			tracer.WithEnv(configs.GetConfig().Datadog.DatadogEnv),
			tracer.WithService(configs.GetConfig().Datadog.DatadogService),
		)
		defer tracer.Stop()
	} else {
		// temp handling until we move all to Datadog
		apm.InitConnection(configs.GetConfig().APMElastic.APMUrl, configs.GetConfig().APMElastic.APMSecretToken)
	}
	// Init MongoDB Connection
	mongo := mongodb.MongoImpl{}
	mongo.SetCollections(&mongo)
	mongo.InitConnection(configs.GetConfig().MongoDB.MongoMasterDBUrl, configs.GetConfig().MongoDB.MongoSlaveDBUrl)

	// Init Logger
	logZap := log.SetupLogger()
	log.Init(logZap)

	// Init Kafka Config
	kafkaConfluent.InitKafkaConfig(configs.GetConfig().Kafka.KafkaUrl, configs.GetConfig().Kafka.KafkaUsername, configs.GetConfig().Kafka.KafkaPassword)

	// Init instance fiber
	app := fiber.New(fiber.Config{
		BodyLimit: 30 * 1024 * 1024,
	})
	app.Use(apmfiber.Middleware(apmfiber.WithTracer(apm.GetTracer())))
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(pprof.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
		// Format:       `${time} {"router_activity" : [${status},"${latency}","${method}","${path}"], "query_param":${queryParams}, "body_param":${body}}` + "\n",
		TimeInterval: time.Millisecond,
		TimeFormat:   "02-01-2006 15:04:05",
		TimeZone:     "Indonesia/Jakarta",
	}))
	shutdownDelay, _ := strconv.Atoi(configs.GetConfig().ShutDownDelay)
	// graceful shutdown setup
	gs := &graceful.GracefulShutdown{
		Timeout:        5 * time.Second,
		GracefulPeriod: time.Duration(shutdownDelay) * time.Second,
	}
	app.Get("/healthz", gs.LivenessCheck)
	app.Get("/readyz", gs.ReadinessCheck)
	gs.Enable(app)

	setHttp(app, gs)
	setConfluentEvents(app, gs)
	//=== listen port ===//
	if err := app.Listen(fmt.Sprintf(":%s", configs.GetConfig().ServicePort)); err != nil {
		logGo.Fatal(err)
	}
}

func setHttp(app *fiber.App, gs *graceful.GracefulShutdown) {
	// set module if use http handler
}

func setConfluentEvents(app *fiber.App, gs *graceful.GracefulShutdown) {
	// Init redis
	redisClient := redis.InitConnection(configs.GetConfig().Redis.RedisDB, configs.GetConfig().Redis.RedisHost, configs.GetConfig().Redis.RedisPort,
		configs.GetConfig().Redis.RedisPassword, configs.GetConfig().Redis.RedisAppConfig)
	logger := log.GetLogger()
	mongoSlaveClient := mongodb.NewMongoDBLogger(mongodb.GetSlaveConn(), mongodb.GetMasterDBName(), logger)
	kafkaProducer, err := kafkaConfluent.NewProducer(kafkaConfluent.GetConfig().GetKafkaConfig(configs.GetConfig().ServiceName, true), logger)
	if err != nil {
		panic(err)
	}
	gs.Register(
		graceful.FnWithError(redisClient.Close),
		kafkaProducer,
		mongoSlaveClient,
	)

	orderQueryMongodbRepo := orderRepoQuery.NewQueryMongodbRepository(mongoSlaveClient, logger)

	// mailHelper := helpers.New(configs.GetConfig().Email.SmtpHost, configs.GetConfig().Email.SmtpPort, configs.GetConfig().Email.EmailUsername, configs.GetConfig().Email.EmailPassword)
	mailHelper := &helpers.Mail{
		Email:    configs.GetConfig().Email.EmailUsername,
		Password: configs.GetConfig().Email.EmailPassword,
		SmtpHost: configs.GetConfig().Email.SmtpHost,
		SmtpPort: configs.GetConfig().Email.SmtpPort,
	}

	notificationCommandUsecase := notificationUsecases.NewCommandUsecase(logger, mailHelper, orderQueryMongodbRepo)

	notificationHandler.InitNotificationEventConflHandler(notificationCommandUsecase, logger)
}

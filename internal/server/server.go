package server

import (
	"fmt"
	config "go-email/config"
	repository "go-email/internal/database"
	delivery "go-email/internal/delivery/grpc"
	"go-email/internal/mailer"
	db "go-email/pkg/database"
	logger "go-email/pkg/logger"
	mail "go-email/pkg/mailer"
	pb "go-email/pkg/proto"
	rb "go-email/pkg/rabbitmq"
	"net"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	consumer "go-email/internal/delivery/rabbitmq"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

// :TODO
// jaeger
// email validator +
// rabbitmq
// metrics +
// devops
// db +
// logger
// config +
// add mode (debug, prod)
// tests
// scripts

func Run() {
	// loading config from env variables
	conf, err := config.LoadConfigFromEnv()

	if err != nil {
		log.Printf("Unable to load config from env %s", err.Error())
	}

	// Init logger
	logger.InitLogger()

	// Settings for SMTP server
	d := mail.NewMailDialer(conf)
	mailer := mailer.NewMailer(d)

	// Init rabbitmq
	rabbitConnection := rb.NewRabbitMQ(conf)
	cons := consumer.NewConsumer(rabbitConnection, mailer, conf)

	go func() {
		err := cons.Consume(conf.Rabbit.ConsumePool)

		if err != nil {
			log.Printf("Failed to consume from rabbitmq %s", err.Error())
		}
	}()

	// Init metrics
	go func() {
		e := echo.New()
		e.GET("/", echo.WrapHandler(promhttp.Handler()))
		log.Fatal(e.Start(":8242"))
	}()

	// Init database
	dbConn := db.NewDatabase(conf)
	repo := repository.NewRepository(dbConn)

	// Jaeger tracing

	// Implementing grpc server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50001))
	if err != nil {
		log.Errorf("failed to start server %v", err)
	}

	log.WithFields(log.Fields{
		"port": lis.Addr(),
	}).Debug("Server started successfully")

	s := grpc.NewServer()
	pb.RegisterMailerServiceServer(s, delivery.NewServer(conf, mailer, repo))

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

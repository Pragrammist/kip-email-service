package server

import (
	"fmt"
	config "go-email/config"
	delivery "go-email/internal/delivery/grpc"
	"go-email/internal/mailer"
	logger "go-email/pkg/logger"
	mail "go-email/pkg/mailer"
	pb "go-email/pkg/proto"
	"net"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

func Run() {
	port, _ := strconv.Atoi(os.Getenv("PORT"))


	// Paste credentials here
	conf := &config.Config{Smtp: config.Smtp{
		Host:     "",
		Port:     0,
		User:     "",
		Password: "",
	}}

	// Init logger
	logger.InitLogger()

	// Settings for SMTP server
	d := mail.NewMailDialer(conf)
	mailer := mailer.NewMailer(d)

	// Init metrics
	go func() {
		e := echo.New()
		e.GET("/", echo.WrapHandler(promhttp.Handler()))
		log.Fatal(e.Start(":8242"))
	}()

	// Init database
	//dbConn := db.NewDatabase(conf)
	//repo := repository.NewRepository(dbConn)

	// Jaeger tracing

	// Implementing grpc server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Errorf("failed to start server %v", err)
	}

	log.WithFields(log.Fields{
		"port": lis.Addr(),
	}).Debug("Server started successfully")

	s := grpc.NewServer()
	pb.RegisterMailerServiceServer(s, delivery.NewServer(conf, mailer))

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

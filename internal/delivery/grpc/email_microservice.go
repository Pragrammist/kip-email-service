package delivery

import (
	"context"
	"errors"
	"fmt"
	"go-email/config"
	"go-email/internal/mailer"
	"go-email/internal/models"
	"go-email/internal/validator"
	pb "go-email/pkg/proto"
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	emailsSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_emails_sended_successfully",
		Help: "Successfully sended emails count",
	})

	emailsFailure = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_emails_sended_failure",
		Help: "Failed emails count",
	})

	emailsSavedSuccessfully = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_emails_saved_to_database_successfully",
		Help: "Count of successfully saved to database emails",
	})

	emailsSavedFailure = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_emails_saved_to_database_failure",
		Help: "Count of failure saved to database emails",
	})
)

type Server struct {
	pb.UnimplementedMailerServiceServer
	mailer *mailer.Mailer
	cfg    *config.Config
}

func NewServer(cfg *config.Config, mailer *mailer.Mailer) *Server {
	return &Server{cfg: cfg, mailer: mailer}
}

func (s *Server) SendEmails(ctx context.Context, r *pb.EmailRequest) (*pb.EmailResponse, error) {
	email := &models.Email{
		From:        s.cfg.Smtp.User,
		To:          r.GetTo(),
		Body:        string(r.GetBody()),
		Subject:     r.GetSubject(),
		ContentType: r.GetContentType(),
	}

	for _, reciever := range r.GetTo() {
		if !validator.ValidateEmail(reciever) {
			return nil, errors.New("Unable to validate email")
		}
	}

	log.Println("Inside send emails method")

	if err := s.mailer.SendEmails(email); err != nil {
		emailsFailure.Inc()
		fmt.Println(err)
		return nil, err
	}

	//if err := s.repo.CreateEmail(email); err != nil {
	//  emailsSavedFailure.Inc()
	//	return nil, err
	//}

	emailsSuccess.Inc()
	emailsSavedSuccessfully.Inc()

	return &pb.EmailResponse{}, nil
}

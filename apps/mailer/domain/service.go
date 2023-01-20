package domain

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/sultanaliev-s/kiteps/pkg/logging"
)

type (
	Service interface {
		SendMail(ctx context.Context, mail Mail) error
		Ping(ctx context.Context) error
	}

	validator interface {
		ValidateStruct(s any) error
	}

	service struct {
		addr     string
		password string
		from     string
		port     string

		validator validator
		logger    *logging.Logger
	}
)

func NewService(
	addr string,
	password string,
	from string,
	port string,
	validator validator,
	logger *logging.Logger,
) Service {
	// todo: validate input

	return &service{
		addr:      addr,
		password:  password,
		from:      from,
		port:      port,
		validator: validator,
		logger:    logger,
	}
}

func (s service) SendMail(ctx context.Context, mail Mail) error {
	if err := s.validator.ValidateStruct(mail); err != nil {
		s.logger.Debug("Mailer.SendMail() - validating", logging.Error("error", err))
		return fmt.Errorf("invalid mail: %w", err)
	}

	message := fmt.Sprintf(
		"To: %s\r\nSubject: %s \r\nMIME-version: 1.0;\nContent-Type: %s; charset=\"UTF-8\";\n\n\r\n%s",
		mail.Recipient, mail.Subject, mail.MIMEType, mail.Body,
	)

	auth := smtp.PlainAuth("", s.from, s.password, s.addr)

	if err := smtp.SendMail(
		s.addr+":"+s.port,
		auth,
		s.from,
		[]string{mail.Recipient},
		[]byte(message),
	); err != nil {
		s.logger.Debug("Mailer.SendMail() - sending mail", logging.Error("error", err))
		return fmt.Errorf("failed to send mail: %w", err)
	}

	return nil
}

func (s service) Ping(ctx context.Context) error {
	// TODO: add checks
	return nil
}

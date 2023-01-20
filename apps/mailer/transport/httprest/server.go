package httprest

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sultanaliev-s/kiteps/apps/mailer/domain"
	"github.com/sultanaliev-s/kiteps/pkg/health"
	"github.com/sultanaliev-s/kiteps/pkg/logging"
	"github.com/sultanaliev-s/kiteps/pkg/validation"
)

type Server struct {
	server *http.Server

	validator *validation.Validator
	logger    *logging.Logger
	service   domain.Service
}

func NewServer(service domain.Service, logger *logging.Logger, validator *validation.Validator, address string) *Server {
	return &Server{
		server: &http.Server{
			Addr: address,
		},
		validator: validator,
		logger:    logger,
		service:   service,
	}
}

func (s Server) Start() error {
	router := echo.New()
	router.Use(s.logger.NewEchoMiddleware)
	router.POST("/mail", s.handleMailerSend)
	router.Any("/health*", echo.WrapHandler(health.NewHTTPHandler("mailer", []health.Checker{
		// TODO: move checks out of Server.Start function
		func(ctx context.Context) health.Check {
			check := health.Check{
				Name:     "mailerSNMP",
				Status:   "UP",
				Critical: true,
			}
			if err := s.service.Ping(ctx); err != nil {
				check.Status = "DOWN"
				check.Message = err.Error()
			}
			return check
		},
	})))
	router.GET("/routes", func(c echo.Context) error {
		return c.JSON(http.StatusOK, router.Routes())
	})

	s.server.Handler = router

	return s.server.ListenAndServe()
}

func (s Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

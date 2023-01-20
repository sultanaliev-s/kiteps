package httprest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sultanaliev-s/kiteps/apps/mailer/domain"
)

func (s *Server) handleMailerSend(ctx echo.Context) error {
	request := domain.Mail{}

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	if err := s.service.SendMail(ctx.Request().Context(), request); err != nil {
		if errs := s.validator.UnpackErrors(err); errs != nil {
			return ctx.JSON(http.StatusBadRequest, errs)
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "mail sent",
	})
}

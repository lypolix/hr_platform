package ginhandler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hr-platform-mosprom/internal/core/application/port"
	"github.com/hr-platform-mosprom/internal/core/domain"
)

type (
	universityHandlers struct {
		universityService port.UniversityService
		logger            *slog.Logger
		validator         *validator.Validate
	}

	universityWithTokenResponse struct {
		Login string `json:"login"`
		Title string `json:"title"`
		INN   string `json:"inn"`
		Token string `json:"token"`
	}

	signInRequest struct {
		Login    string `json:"login" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	signUpRequest struct {
		Login    string `json:"login" validate:"required"`
		Password string `json:"password" validate:"required"`
		Title    string `json:"title" validate:"required"`
		INN      string `json:"inn" validate:"required"`
	}
)

func RegisterUniversityHandlers(
	engine *gin.Engine,
	userService port.UniversityService,
	logger *slog.Logger,
	validator *validator.Validate,
) {
	handlers := universityHandlers{userService, logger, validator}

	group := engine.Group("/universities")
	group.POST("/sign-in", handlers.SingIn)
	group.POST("/sign-up", handlers.SignUp)
}

func (h *universityHandlers) SingIn(c *gin.Context) {
	ctx := c.Request.Context()

	var request signInRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		h.logger.ErrorContext(ctx, "error parsing request body", "err", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "error parsing request body"})

		return
	}

	err = h.validator.StructCtx(ctx, request)
	if err != nil {
		h.logger.ErrorContext(ctx, "error validating body", "err", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "error validating body"})

		return
	}

	universityWithTokenResult, err := h.universityService.SignIn(ctx, port.SignInUniversityData{
		Login:    request.Login,
		Password: request.Password,
	})
	if err != nil {
		h.logger.ErrorContext(ctx, "error during sign in", "err", err)

		if errors.Is(err, domain.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}

		return
	}

	response := universityWithTokenResponse{
		Login: universityWithTokenResult.Login,
		Title: universityWithTokenResult.Title,
		INN:   universityWithTokenResult.INN,
		Token: universityWithTokenResult.Token,
	}

	c.JSON(http.StatusOK, response)
}

func (h *universityHandlers) SignUp(c *gin.Context) {
	ctx := c.Request.Context()

	var request signUpRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		h.logger.ErrorContext(ctx, "error parsing request body", "err", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "error parsing request body"})

		return
	}

	err = h.validator.StructCtx(ctx, request)
	if err != nil {
		h.logger.ErrorContext(ctx, "error validating body", "err", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "error validating body"})

		return
	}

	universityWithTokenResult, err := h.universityService.SignUp(ctx, port.SignUpUniversityData{
		Login:    request.Login,
		Password: request.Password,
		Title:    request.Title,
		INN:      request.INN,
	})
	if err != nil {
		h.logger.ErrorContext(ctx, "error during sign up", "err", err)

		if errors.Is(err, domain.ErrConflict) {
			c.JSON(http.StatusConflict, gin.H{"message": "login already taken"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}

		return
	}

	response := universityWithTokenResponse{
		Login: universityWithTokenResult.Login,
		Title: universityWithTokenResult.Title,
		INN:   universityWithTokenResult.INN,
		Token: universityWithTokenResult.Token,
	}

	c.JSON(http.StatusOK, response)
}

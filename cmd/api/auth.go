package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mgiks/gopher-social/internal/mailer"
	"github.com/mgiks/gopher-social/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72" format:"password"`
}

type UserWithToken struct {
	User  store.User `json:"user"`
	Token string     `json:"token"`
}

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	apiResponse{data=UserWithToken}
//	@Failure		400		{object}	apiError
//	@Failure		500		{object}	apiError
//	@Security		ApiKeyAuth
//	@Router			/auth/user [post]
func (app application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.validator.ValidateJSON(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	plainToken := uuid.New().String()

	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	err := app.store.Users.CreateAndInvite(r.Context(), user, hashToken, app.config.mail.exp)
	if err != nil {
		switch err {
		case store.ErrDuplicateEmail, store.ErrDuplicateUsername:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)

	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	isProdEnv := app.config.env == "production"

	emailSender, err := app.mailer.NewSender(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err = app.retry(3, emailSender); err != nil {
		app.logger.Errorw("error sending welcome email", "error", err)

		if err := app.store.Users.Delete(r.Context(), user.ID); err != nil {
			app.logger.Errorw("error deleting user", "error", err)
		}

		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, UserWithToken{User: *user, Token: plainToken}); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app application) retry(retryCount int, sender mailer.Sender) error {
	for i := range retryCount {
		receiverEmail, err := sender.Send()
		if err != nil {
			app.logger.Warnw("Failed to send email", "receiverEmail", receiverEmail, "attempt", fmt.Sprintf("%v of %v", i+1, retryCount))
			app.logger.Warn("Error:", err.Error())

			// exponential backoff
			secsToWait := math.Pow(float64(2), float64(i))
			time.Sleep(time.Second * time.Duration(secsToWait))
			continue
		}
		app.logger.Info("Email sent succesfully")
		return nil
	}
	return fmt.Errorf("failed to send email after %d attempts", retryCount)
}

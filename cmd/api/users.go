package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mgiks/gopher-social/internal/store"
)

// getUserHandler godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		200		{object}	apiResponse{data=store.User}
//	@Failure		404		{object}	apiError	"User not found"
//	@Failure		500		{object}	apiError
//	@Security		ApiKeyAuth
//	@Router			/users/{userID} [get]
func (app application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

type FollowUserPayload struct {
	UserID int64 `json:"user_id" validate:"required"`
}

// followUserHandler godoc
//
//	@Summary		Follows a user
//	@Description	Follows a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	int	true	"User ID"
//	@Success		200		"User followed"
//	@Failure		400		{object}	apiError	"User payload missing"
//	@Failure		404		{object}	apiError	"User not found"
//	@Failure		409		{object}	apiError	"Tried following a user more than once"
//	@Failure		500		{object}	apiError
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [put]
func (app application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followee := getUserFromContext(r.Context())

	// TODO: Revert back to auth userID from ctx after auth
	var payload FollowUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.Followers.Follow(r.Context(), followee.ID, payload.UserID); err != nil {
		switch {
		case errors.Is(err, store.ErrConflict):
			app.conflictResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// unfollowUserHandler godoc
//
//	@Summary		Unfollows a user
//	@Description	Unfollows a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	int	true	"User ID"
//	@Success		200		"User unfollowed"
//	@Failure		400		{object}	apiError	"User payload missing"
//	@Failure		404		{object}	apiError	"User not found"
//	@Failure		500		{object}	apiError
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [put]
func (app application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followee := getUserFromContext(r.Context())

	// TODO: Revert back to auth userID from ctx after auth
	var payload FollowUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.Followers.Unfollow(r.Context(), followee.ID, payload.UserID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type userKey string

var userCtx userKey = "user"

func (app application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "userID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		user, err := app.store.Users.GetByID(r.Context(), id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(ctx context.Context) store.User {
	user := ctx.Value(userCtx).(store.User)
	return user
}

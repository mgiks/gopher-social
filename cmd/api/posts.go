package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mgiks/gopher-social/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

// createPostHandler godoc
//
//	@Summary		Creates a post
//	@Description	Creates a post using title, content and tags
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreatePostPayload	true	"Payload for post creation"
//	@Success		201		{object}	apiResponse{data=store.Post}
//	@Failure		400		{object}	apiError
//	@Failure		404		{object}	apiError
//	@Failure		500		{object}	apiError
//	@Security		ApiKeyAuth
//	@Router			/posts/ [post]
func (app application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())

	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.validator.ValidateJSON(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  user.ID,
	}

	if err := app.store.Posts.Create(r.Context(), post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

// getPostHandler godoc
//
//	@Summary		Gets a post
//	@Description	Gets a post by id
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int	true	"Post ID"
//	@Success		200		{object}	apiResponse{data=store.Post}
//	@Failure		400		{object}	apiError
//	@Failure		404		{object}	apiError
//	@Failure		500		{object}	apiError
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID} [get]
func (app application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r.Context())

	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

// deletePostHandler godoc
//
//	@Summary		Deletes a post
//	@Description	Deletes a post by id
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int			true	"Post ID"
//	@Success		204		string		{}			"Post deleted"
//	@Failure		400		{object}	apiError	"Post not found"
//	@Failure		500		{object}	apiError
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID} [delete]
func (app application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.store.Posts.DeleteByID(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdatePostPayload struct {
	Title   *string  `json:"title" validate:"omitempty,max=100"`
	Content *string  `json:"content" validate:"omitempty,max=1000"`
	Tags    []string `json:"tags" validate:"omitempty"`
}

// updatePostHandler godoc
//
//	@Summary		Updates a post
//	@Description	Updates a post by id using title or content or tags
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int					true	"Post ID"
//	@Param			payload	body		UpdatePostPayload	true	"Payload for post updating"
//	@Success		200		{object}	apiResponse{data=store.Post}
//	@Failure		400		{object}	apiError
//	@Failure		404		{object}	apiError
//	@Failure		409		{object}	apiError
//	@Failure		500		{object}	apiError
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID}/ [patch]
func (app application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r.Context())

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.validator.ValidateJSON(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if payload.Tags != nil {
		post.Tags = payload.Tags
	}

	if err := app.store.Posts.Update(r.Context(), &post); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.conflictResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

type postKey string

const postCtx postKey = "post"

func (app application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		post, err := app.store.Posts.GetByID(r.Context(), id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromContext(ctx context.Context) store.Post {
	post := ctx.Value(postCtx).(store.Post)
	return post
}

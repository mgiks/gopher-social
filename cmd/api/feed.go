package main

import (
	"net/http"

	"github.com/mgiks/gopher-social/internal/store"
)

// getUserFeedHandler godoc
//
//	@Summary		Gets an authenticated user's feed
//	@Description	Gets a user's feed
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int			false	"Post limit"				minimum(1)				maximum(20)	default(20)
//	@Param			offset	query		int			false	"Post offset"				minimum(0)				default(0)
//	@Param			sort	query		string		false	"Sorting order"				Enums(asc, desc)		default(desc)
//	@Param			tags	query		[]string	false	"Tags to search by"			collectionFormat(csv)	maxlength(5)	example("new", "post")
//	@Param			search	query		string		false	"Seach query"				maxlength(100)
//	@Param			since	query		string		false	"Post creation date cutoff"	format(date-time)	example(2006-01-02T15:04:05Z)
//	@Param			until	query		string		false	"Post creation date cutoff"	format(date-time)	example(2006-01-02T15:04:05Z)
//	@Success		200		{object}	apiResponse{Data=[]store.PostWithMetadata}
//	@Failure		400		{object}	apiError
//	@Failure		500		{object}	apiError
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (app application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
		Tags:   []string{},
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.validator.ValidateJSON(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// TODO: Change 1 to userID after auth
	feed, err := app.store.Posts.GetUserFeed(r.Context(), 1, fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}

package main

import (
	"net/http"

	"github.com/mgiks/gopher-social/internal/store"
)

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

	if err := validateJSON(fq); err != nil {
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

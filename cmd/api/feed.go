package main

import (
	"net/http"
)

func (app application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Change 1 to userID after auth
	feed, err := app.store.Posts.GetUserFeed(r.Context(), 1)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}

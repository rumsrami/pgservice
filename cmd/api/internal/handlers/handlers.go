package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/rumsrami/pgservice/internal/data"
	"gopkg.in/matryer/respond.v1"
)

// getTitle
func getTitle(db *sqlx.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		titleFromURL := chi.URLParam(r, "title")

		title, err := data.Read.Info(context.Background(), db, titleFromURL)
		if err != nil {
			if err == data.ErrNotFound {
				respond.With(w, r, http.StatusNotFound, "title not found")
				return
			}
		}

		respond.With(w, r, http.StatusOK, title)
	}
	return fn
}

func postTitle(db *sqlx.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {

		var reqInfo data.Info
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&reqInfo); err != nil {
			respond.With(w, r, http.StatusBadRequest, "unexpected request body")
			return
		}

		addedInfo, err := data.Create.Info(context.Background(), db, reqInfo.Title)
		if err != nil {
			if err == data.ErrNotFound {
				respond.With(w, r, http.StatusNotFound, nil)
				return
			}
		}

		respond.With(w, r, http.StatusOK, addedInfo)
	}
	return fn
}

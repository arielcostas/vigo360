/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package admin

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"vigo360.es/new/internal/database"
)

var db *sqlx.DB

func InitRouter() *mux.Router {
	db = database.GetDB()

	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/admin/", http.RedirectHandler("/admin/login", 302)).Methods("GET")
	router.Handle("/admin/login", appHandler(viewLogin)).Methods("GET")
	router.Handle("/admin/login", appHandler(doLogin)).Methods("POST")
	router.Handle("/admin/logout", appHandler(logoutPage)).Methods("GET")

	router.Handle("/admin/dashboard", appHandler(viewDashboard)).Methods("GET")
	router.Handle("/admin/post", appHandler(listPosts)).Methods("GET")
	router.Handle("/admin/post", appHandler(createPost)).Methods("POST")

	router.Handle("/admin/post/{id}", appHandler(postEditor)).Methods("GET")
	router.Handle("/admin/post/{id}", appHandler(editPost)).Methods("POST")
	router.Handle("/admin/post/{postid}/delete", appHandler(deletePost)).Methods("GET")

	router.Handle("/admin/series", appHandler(listSeries)).Methods("GET")
	router.Handle("/admin/series", appHandler(createSeries)).Methods("POST")

	router.Handle("/admin/sesiones", appHandler(listSessions)).Methods("GET")

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		t.ExecuteTemplate(w, "_404.html", struct{}{})
	})

	return router
}
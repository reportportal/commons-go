package server

import (
	"github.com/go-chi/chi"
	"github.com/gorilla/handlers"
	"github.com/reportportal/commons-go/commons"
	"github.com/reportportal/commons-go/conf"
	"net/http"
	"os"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func ExampleRpServer() {
	rpConf, _ := conf.LoadConfig(nil, nil)
	rp := New(rpConf, commons.GetBuildInfo())

	rp.AddRoute(func(router *chi.Mux) {
		router.Get("/ping", func(w http.ResponseWriter, rq *http.Request) {
			commons.WriteJSON(http.StatusOK, Person{"av", 20}, w)
		})
	})

	rp.StartServer()

}

func ExampleRpServer_StartServer() {

	rpConf, _ := conf.LoadConfig(nil,
		map[string]string{"AuthServerURL": "http://localhost:9998/sso/me"})

	srv := New(rpConf, commons.GetBuildInfo())

	srv.AddRoute(func(mux *chi.Mux) {
		mux.Use(func(next http.Handler) http.Handler {
			return handlers.LoggingHandler(os.Stdout, next)
		})

		secured := chi.NewMux()
		secured.Use(RequireRole("USER", rpConf.Get("AuthServerURL")))

		me := func(w http.ResponseWriter, rq *http.Request) {
			commons.WriteJSON(http.StatusOK, rq.Context().Value("user"), w)

		}
		secured.HandleFunc("/me", me)

		mux.Handle("/", secured)

	})

	srv.StartServer()
}

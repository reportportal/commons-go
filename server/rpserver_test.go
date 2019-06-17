package server

import (
	"github.com/go-chi/chi"
	"github.com/gorilla/handlers"
	"github.com/reportportal/commons-go/commons"
	"github.com/reportportal/commons-go/conf"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func ExampleRpServer() {
	rpConf := conf.EmptyConfig()
	_ = conf.LoadConfig(rpConf)
	rp := New(rpConf, commons.GetBuildInfo())

	rp.WithRouter(func(router *chi.Mux) {
		router.Get("/ping", func(w http.ResponseWriter, rq *http.Request) {
			if err := WriteJSON(http.StatusOK, Person{"av", 20}, w); err != nil {
				logrus.Error(err)
			}
		})
	})

	rp.StartServer()

}

func ExampleRpServer_StartServer() {
	rpConf := conf.EmptyConfig()
	authServerURL := "http://localhost:9998/sso/me"
	_ = conf.LoadConfig(rpConf)

	srv := New(rpConf, commons.GetBuildInfo())

	srv.WithRouter(func(mux *chi.Mux) {
		mux.Use(func(next http.Handler) http.Handler {
			return handlers.LoggingHandler(os.Stdout, next)
		})

		secured := chi.NewMux()
		secured.Use(RequireRole("USER", authServerURL))

		me := func(w http.ResponseWriter, rq *http.Request) {
			if err := WriteJSON(http.StatusOK, rq.Context().Value("user"), w); err != nil {
				logrus.Error(err)
			}

		}
		secured.HandleFunc("/me", me)

		mux.Handle("/", secured)

	})

	srv.StartServer()
}

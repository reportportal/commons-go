package server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/reportportal/commons-go/commons"
	"github.com/reportportal/commons-go/conf"
	"github.com/reportportal/commons-go/registry"
	"log"
	"net/http"
	"strconv"
)

//RpServer represents ReportPortal micro-service instance
type RpServer struct {
	mux     *chi.Mux
	cfg     *conf.RpConfig
	Sd      registry.ServiceDiscovery
	hChecks []HealthCheck
}

//New creates new instance of RpServer struct
func New(cfg *conf.RpConfig, buildInfo *commons.BuildInfo) *RpServer {

	var sd registry.ServiceDiscovery
	switch cfg.Registry {
	case conf.Eureka:
		sd = registry.NewEureka(cfg)
	case conf.Consul:
		cfg.Consul.Tags = append(cfg.Consul.Tags, "statusPageUrlPath=/info", "healthCheckUrlPath=/health")
		sd = registry.NewConsul(cfg)
	}

	srv := &RpServer{
		mux:     chi.NewMux(),
		cfg:     cfg,
		Sd:      sd,
		hChecks: []HealthCheck{},
	}

	srv.mux.Use(middleware.Recoverer)

	srv.mux.Get("/health", func(w http.ResponseWriter, rq *http.Request) {

		//collect status results
		errs := []string{}
		for _, hc := range srv.hChecks {
			if err := hc.Check(); nil != err {
				errs = append(errs, err.Error())
			}
		}

		rs := map[string]interface{}{}
		status := http.StatusOK
		if len(errs) > 0 {
			rs["status"] = "DOWN"
			rs["errors"] = errs
			status = http.StatusBadRequest
		} else {
			rs["status"] = "UP"
		}

		commons.WriteJSON(status, rs, w)
	})

	bi := map[string]interface{}{"build": buildInfo}
	srv.mux.Get("/info", func(w http.ResponseWriter, rq *http.Request) {
		commons.WriteJSON(200, bi, w)

	})
	return srv
}

//AddRoute gives access to GIN router to add route and perform other modifications
func (srv *RpServer) AddRoute(f func(router *chi.Mux)) {
	f(srv.mux)
}

//AddHealthCheck adds health check
func (srv *RpServer) AddHealthCheck(h HealthCheck) {
	srv.hChecks = append(srv.hChecks, h)
}

//StartServer starts HTTP server
func (srv *RpServer) StartServer() {

	if nil != srv.Sd {
		registry.Register(srv.Sd)
	}
	// listen and server on mentioned port
	log.Printf("Starting on port %d", srv.cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(srv.cfg.Server.Port), srv.mux))
}

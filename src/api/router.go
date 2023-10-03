package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/apoloa/bjournal/src/service"
)

type Router struct {
	port       int
	logService *service.LogService
	router     *http.ServeMux
}

func NewRouter(port int, logService *service.LogService) *Router {
	return &Router{port: port, logService: logService, router: http.NewServeMux()}
}

func (r *Router) Init() {
	r.router.HandleFunc("/api/log/today", func(writer http.ResponseWriter, request *http.Request) {
		day, err := r.logService.ReadDay(time.Now())
		if err != nil {
			return
		}
		err = json.NewEncoder(writer).Encode(day)
		if err != nil {
			return
		}

	})
	r.router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
}

func (r *Router) Start() {
	srv := &http.Server{
		Handler: r.router,
		Addr:    fmt.Sprintf("127.0.0.1:%v", r.port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		return
	}
}

package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"sync"
)

type handler struct {
	mu       *sync.RWMutex
	handlers map[string]http.HandlerFunc
}

func main() {
	host := "0.0.0.0"
	port := "9999"
	if err := execute(host, port); err != nil {
		os.Exit(1)
	}
}
func execute(host string, port string) (err error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/banner.getAll", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte("demo data"))
		if err != nil {
			log.Println(err)
		}
	})
	srv := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: mux,
	}
	return srv.ListenAndServe()
}

func (mux *http.ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h, _ := mux.Handler(r)
	h.ServeHTTP(w, r)
}

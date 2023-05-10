package main

import (
	"http/cmd/app"
	"http/pkg/banners"
	"net"
	"net/http"
	"os"
)

//type handler struct {
//	mu       *sync.RWMutex
//	handlers map[string]http.HandlerFunc
//}

func main() {
	host := "0.0.0.0"
	port := "9999"
	if err := execute(host, port); err != nil {
		os.Exit(1)
	}
}
func execute(host string, port string) (err error) {
	mux := http.NewServeMux()
	bannersSvc := banners.NewService()
	server := app.NewServer(mux, bannersSvc)
	srv := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: server,
	}
	server.Init()
	return srv.ListenAndServe()
}

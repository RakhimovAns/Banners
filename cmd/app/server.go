package app

import (
	"encoding/json"
	"http/pkg/banners"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	mux        *http.ServeMux
	bannersSvc *banners.Service
}

func NewServer(mux *http.ServeMux, bannersSvc *banners.Service) *Server {
	return &Server{mux: mux, bannersSvc: bannersSvc}
}
func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) handleGetBannerByID(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam, 10, 54)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item, err := s.bannersSvc.ByID(request.Context(), id)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Println(err)
	}
}
func (s *Server) handleGetAllBanner(writer http.ResponseWriter, request *http.Request) {
	items := s.bannersSvc.All()
	for _, banner := range items {
		data, err := json.Marshal(banner)
		if err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			continue
		}
		writer.Header().Set("Content-Type", "application/json")
		_, err = writer.Write(data)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
func (s *Server) handleGetBannerToSave(writer http.ResponseWriter, request *http.Request) {
	idParam := request.FormValue("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	title := request.FormValue("title")
	content := request.FormValue("content")
	button := request.FormValue("button")
	link := request.FormValue("link")
	Ban := banners.Banner{
		Title:   title,
		Content: content,
		Button:  button,
		Link:    link,
	}
	err = s.bannersSvc.Save(request.Context(), &Ban, id, request)
}
func (s *Server) handleRemoveByID(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	item, err := s.bannersSvc.ByID(request.Context(), id)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Println(err)
	}
	err = s.bannersSvc.DeleteByID(request.Context(), id)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
func (s *Server) Init() {
	s.mux.HandleFunc("/banners.save", s.handleGetBannerToSave)
	s.mux.HandleFunc("/banners.getByID", s.handleGetBannerByID)
	s.mux.HandleFunc("/banners.getAll", s.handleGetAllBanner)
	s.mux.HandleFunc("/banners.removeById", s.handleRemoveByID)
}

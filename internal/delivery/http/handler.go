package delivery

import (
	"github.com/mehrdod/purehttpserver/internal/service"
	"net/http"
	"strconv"
)

type Handler struct {
	services   *service.Services
	counterTTL int
}

func NewHandler(services *service.Services, counterTTl int) *Handler {
	return &Handler{
		services:   services,
		counterTTL: counterTTl,
	}
}

func (h *Handler) Init() *http.ServeMux {
	routerMux := http.NewServeMux()

	routerMux.Handle("/counter", http.HandlerFunc(h.RequestInfo))

	return routerMux
}

func (h *Handler) RequestInfo(w http.ResponseWriter, r *http.Request) {

	reqInfo, err := h.services.RequestInfo.Get()
	if err != nil {
		w.Write([]byte("error: " + err.Error()))
		return
	}

	w.Write([]byte("Number of requests in a last " + strconv.Itoa(h.counterTTL) + " seconds are: " + strconv.Itoa(reqInfo.LastRequestsNum)))
}

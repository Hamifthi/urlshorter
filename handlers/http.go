package handlers

import (
	"net/http"
	"urlshortner"
	"urlshortner/config"
)

type HTTPHandler struct {
	Service urlshortner.UrlService
}

func (h *HTTPHandler) ShortenLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("only support post request\n"))
		return
	}
	url := r.URL.Query().Get("url")
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("need to send the url that must be shorten\n"))
		return
	}
	key, err := h.Service.CreateShortLink(url)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("can't create short link with provided key\n"))
		return
	}
	toClickUrl := config.REDIRECT_URL + key
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(toClickUrl + "\n"))
}

func (h *HTTPHandler) GetOriginalUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("only support get request\n"))
		return
	}
	key := r.URL.Query().Get("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("you must send the key\n"))
		return
	}
	url, err := h.Service.GetUrl(key)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error in getting original url with provided key\n"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(url.URL + "\n"))
}

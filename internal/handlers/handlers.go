package handler

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

var randSrc = rand.New(rand.NewSource(time.Now().UnixNano()))

type Handler struct {
	Log          *zap.Logger
	urlStore     sync.Map
	reverseStore sync.Map
	domainCount  sync.Map
}

func NewHandler(log *zap.Logger) *Handler {
	handler := Handler{
		Log: log,
	}

	return &handler
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortCode(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[randSrc.Intn(len(letterBytes))]
	}
	return string(b)
}

func (h *Handler) ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	parsedURL, err := url.ParseRequestURI(payload.URL)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	if tmp, ok := h.reverseStore.Load(payload.URL); ok {
		code := tmp.(string)
		respondWithJSON(w, map[string]string{"short_url": "/" + code})
		return
	}

	code := generateShortCode(6)
	h.urlStore.Store(code, payload.URL)
	h.reverseStore.Store(payload.URL, code)

	domain := parsedURL.Hostname()
	val, loaded := h.domainCount.Load(domain)
	if !loaded {
		var count int64 = 1
		h.domainCount.Store(domain, &count)
	} else {
		countPtr := val.(*int64)
		atomic.AddInt64(countPtr, 1)
	}

	respondWithJSON(w, map[string]string{"short_url": "/" + code})
}

func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/")
	if tmp, ok := h.urlStore.Load(code); ok {
		http.Redirect(w, r, tmp.(string), http.StatusFound)
		return
	}
	http.NotFound(w, r)
}

func (h *Handler) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	type domainEntry struct {
		Domain string `json:"domain"`
		Count  int64  `json:"count"`
	}

	entries := make([]domainEntry, 0)
	h.domainCount.Range(func(key, value any) bool {
		entries = append(entries, domainEntry{key.(string), *value.(*int64)})
		return true
	})

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Count > entries[j].Count
	})

	if len(entries) > 3 {
		entries = entries[:3]
	}
	respondWithJSON(w, entries)
}

func respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

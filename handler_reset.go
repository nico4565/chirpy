package main

import (
	"net/http"
	"sync/atomic"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits = atomic.Int32{}
	w.Write([]byte("OK\n"))
}

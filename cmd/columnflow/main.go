package main

import (
	"github.com/wyw14/solo-039/internal/api"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := os.Getenv("COLUMNFLOW_ADDR")
	if addr == "" {
		addr = "127.0.0.1:21246"
	}
	server := api.NewServer()
	log.Fatal(http.ListenAndServe(addr, server.Handler()))
}

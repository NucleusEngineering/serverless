package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/helloworlddan/tortune/tortune"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, tortune.HitMe())
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

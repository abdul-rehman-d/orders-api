package main

import (
	"fmt"
	"net/http"
)

func main() {

	server := &http.Server{
		Addr: ":3000",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, world!"))
		}),
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to start server %v\n", err)
	}
}

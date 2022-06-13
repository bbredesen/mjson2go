package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Development Server listening on localhost:9090")

	err := http.ListenAndServe(":9090", http.FileServer(http.Dir("../../")))
	if err != nil {
		fmt.Println("Failed to start server", err)
		return
	}
}

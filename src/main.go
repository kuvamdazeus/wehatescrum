package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
  generate_summary()
  
  PORT := os.Getenv("PORT")
  if PORT == "" {
    PORT = "8080"
  }

	http.HandleFunc("/summary", fetchSummaryJson)
  http.HandleFunc("/trigger", triggerSummaryGenerate)

	fmt.Println("Starting server at port", PORT)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	generateSummary(SummaryOpts{
		date: time.Now(),
		duration: 24 * time.Hour,
	})

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	http.HandleFunc("/summary", fetchSummaryJson)

	fmt.Println("Starting server at port", PORT)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
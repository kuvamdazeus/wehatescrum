package main

import (
	"fmt"
	"net/http"
)

func fetchSummaryJson(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Add("Content-Type", "application/json")

	summary, err := getSummary()
	if err != nil {
		fmt.Println(err)
			response := Response {
			Message: "some error occured while fetching summary!",
			Data: nil,
		}
		response.WriteResponse(writer, http.StatusInternalServerError)
		return
	}

	response := Response {
		Message: "summary retrieved successfully",
		Data: summary,
	}
	response.WriteResponse(writer, http.StatusOK)
}

func triggerSummaryGenerate(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Add("Content-Type", "application/json")

	generateSummary(true)

	writer.WriteHeader(http.StatusOK)
}
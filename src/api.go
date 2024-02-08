package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func fetchSummaryJson(writer http.ResponseWriter, _ *http.Request) {
  writer.Header().Add("Content-Type", "application/json")

  summaryFile, err := os.ReadFile("summary.json")
  if err != nil {
    fmt.Println(err)
    errResponse, _ := json.Marshal(ErrResponse{
      Message: "Oops, something went wrong while reading summary.json!",
    })
    writer.Write(errResponse)

    return
  }
  
  var summary Summary
  json.Unmarshal(summaryFile, &summary)

  successResponse, _ := json.Marshal(SuccessResponse{
    Message: "",
    Data: summary,
  })
  writer.Write(successResponse)
}

func triggerSummaryGenerate(writer http.ResponseWriter, _ *http.Request) {
  writer.Header().Add("Content-Type", "application/json")
  
  generate_summary()

  writer.WriteHeader(http.StatusOK)
}
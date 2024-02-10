package main

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string
	Data any
}

func (res Response) WriteResponse(writer http.ResponseWriter, statusCode int) {
	bytesRes, _ := json.Marshal(res)
	writer.WriteHeader(statusCode)
	writer.Write(bytesRes)
}
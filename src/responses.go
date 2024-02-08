package main

type ErrResponse struct {
  Message string
}

type SuccessResponse struct {
  Message string
  Data any
}
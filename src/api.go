package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func fetchSummaryJson(writer http.ResponseWriter, r *http.Request) {
	writer.Header().Add("Content-Type", "application/json")

	query := r.URL.Query()
	force, _ := strconv.ParseBool(query.Get("force"))
	date := query.Get("date")
	duration := query.Get("duration")

	atDateTime, dateErr := time.Parse(getConstants().timeLayout, date)
	if date != "" && dateErr != nil {
		res := Response{
			Message: fmt.Sprint("Date parsing error occured:", dateErr),
			Data: nil,
		}
		res.WriteResponse(writer, http.StatusBadRequest)
		return
	}

	opts := SummaryOpts{
		force: force,
	}
	
	// get cleaned default from value in opts
	if dateErr != nil {
		opts.date = time.Now()
	} else {
		opts.date = atDateTime
	}

	// get cleaned duration value in opts
	if duration == "all" {
		var zeroDate time.Time
		opts.duration = opts.date.Sub(zeroDate)
	} else if duration == "" {
		opts.duration = time.Hour * 24
	} else {
		n, err := strconv.Atoi(duration[:len(duration)-1])
		strconv.ParseInt()
		if err != nil {
			res := Response{
				Message: fmt.Sprint(err),
				Data: nil,
			}
			res.WriteResponse(writer, http.StatusBadRequest)
			return
		}

		suffix := string(duration[len(duration)-1])
		if suffix == "h" {
			opts.duration = time.Duration(n) * time.Hour
		} else if suffix == "d" {
			opts.duration = time.Duration(n) * 24 * time.Hour
		} else if suffix == "w" {
			opts.duration = time.Duration(n) * 7 * 24 * time.Hour
		} else if suffix == "m" {
			opts.duration = time.Duration(n) * 30 * 24 * time.Hour
		} else {
			res := Response{
				Message: "invalid duration! only formats like '2d' (2 days), '4w' (4 weeks), '6m' (6 months), '12h' (12 hours) are allowed",
				Data: nil,
			}
			res.WriteResponse(writer, http.StatusBadRequest)
			return
		}
	}

	summary, err := generateSummary(opts)
	if err != nil {
		fmt.Println(err)
		response := Response {
			Message: fmt.Sprintf("some error occured while fetching summary! (%s)", err),
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
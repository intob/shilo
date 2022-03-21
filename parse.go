package main

import (
	"bytes"
	"net/http"
)

type data struct {
	content      *bytes.Buffer
	contentExt   string
	outExt       string
	outRes       string
	progressAddr string
}

func parseRequest(r *http.Request) (*data, error) {
	var buf bytes.Buffer
	buf.ReadFrom(r.Body)
	r.Body.Close()

	d := &data{
		content:      &buf,
		contentExt:   r.Header.Get("x-in-ext"),
		outExt:       r.Header.Get("x-out-ext"),
		outRes:       r.Header.Get("x-out-res"),
		progressAddr: r.Header.Get("x-progress-addr"),
	}

	return d, nil
}

package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

func parseContentTypeHeader(r *http.Request) *contentType {
	h := r.Header.Get("Content-Type")
	fields := strings.Split(h, ";")
	ct := &contentType{
		mimeType: fields[0],
	}
	if len(fields) > 1 {
		b := strings.Split(fields[1], "=")
		if len(b) > 1 {
			ct.boundary = b[1]
		}
	}
	return ct
}

func parseData(r *http.Request) (*data, error) {
	ct := parseContentTypeHeader(r)
	mr := multipart.NewReader(r.Body, ct.boundary)
	defer r.Body.Close()

	data := &data{
		outType: "webm",
	}

	var buf bytes.Buffer

	for {
		p, eof := mr.NextPart()
		if p == nil {
			break
		}
		k := p.FormName()
		if k != "" { // option
			buf.ReadFrom(p)
			v := string(buf.Bytes())
			buf.Reset()
			switch k {
			case "outType":
				data.outType = v
			case "width":
				width, err := strconv.Atoi(v)
				if err != nil {
					return nil, err
				}
				data.width = width
			case "height":
				height, err := strconv.Atoi(v)
				if err != nil {
					return nil, err
				}
				data.height = height
			}

		} else { // file
			var buf bytes.Buffer
			buf.ReadFrom(p)
			data.content = &buf
			data.fileName = p.FileName()
		}

		if eof != nil {
			break
		}
	}
	return data, nil
}

package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/intob/shilo/ffmpeg"
)

type contentType struct {
	mimeType string
	boundary string
}

type data struct {
	fileName string
	content  *bytes.Buffer
	outType  string
	width    int
	height   int
}

func handle(w http.ResponseWriter, r *http.Request) {
	data, err := parseData(r)
	if err != nil {
		log.Println("error parsing data:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ext := path.Ext(data.fileName)
	tmp, err := writeTempFile(data.content, "shilo*"+ext)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	in := tmp.Name()
	defer os.Remove(in)
	defer log.Println("removed", tmp.Name())
	log.Println("created", in)

	out := strings.Replace(in, ext, "_out."+data.outType, 1)
	defer os.Remove(out)
	defer log.Println("removed", out)

	cmd := ffmpeg.Scale(in, out, data.width, data.height)
	cmd.Stdout = os.Stdout //io.MultiWriter(logW, os.Stdout)
	cmd.Stderr = os.Stdout //io.MultiWriter(logW, os.Stdout)

	err = cmd.Start()
	if err != nil {
		log.Println("ffmpeg failed to start:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = cmd.Wait()
	if err != nil {
		log.Println("ffmpeg failed:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, out)
}

func writeTempFile(content *bytes.Buffer, pattern string) (*os.File, error) {
	tmp, err := os.CreateTemp(os.Getenv("TEMP_DIR"), pattern)
	if err != nil {
		log.Println("error creating temp file:", err)
		return tmp, err
	}
	defer tmp.Close()

	_, err = content.WriteTo(tmp)
	if err != nil {
		log.Println("error writing content:", err)
		return tmp, err
	}

	return tmp, nil
}

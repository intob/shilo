package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/intob/shilo/ffmpeg"
)

func handle(w http.ResponseWriter, r *http.Request) {
	data, err := parseRequest(r)
	if err != nil {
		log.Println("error parsing data:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tmp, err := writeTempFile(data.content, "shilo*"+data.contentExt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	in := tmp.Name()
	defer os.Remove(in)
	defer log.Println("removed", in)
	log.Println("created", in)

	out := strings.Replace(in, data.contentExt, "_out"+data.outExt, 1)
	defer os.Remove(out)
	defer log.Println("removed", out)

	cmd := ffmpeg.Scale(r.Context(), in, out, data.outRes)
	progW, conn, err := getProgressWriters(cmd, data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if conn != nil {
		defer conn.Close()
	}
	cmd.Stdout = progW
	cmd.Stderr = progW

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

func getProgressWriters(cmd *exec.Cmd, data *data) (io.Writer, net.Conn, error) {
	if data.progressAddr == "" {
		return os.Stdout, nil, nil
	}
	// try to connect to progress socket
	c, err := net.Dial("tcp", data.progressAddr)
	if err != nil {
		err = errors.New("failed to connect to progress socket:" + err.Error())
		return os.Stdout, nil, err
	}
	return io.MultiWriter(c, os.Stdout), c, nil
}

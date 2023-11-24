package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/stream/play", play)
	mux.HandleFunc("/stream/status", status)

	log.Print("Listening...")
	log.Fatal(http.ListenAndServe(":8090", mux))
}

func ping(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("pong"))
}

func play(w http.ResponseWriter, req *http.Request) {
	urlStream := req.URL.Query().Get("url")
	if urlStream == "" {
		log.Print("Url is not valid")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	cmd := fmt.Sprintf("ffprobe -v quiet -show_entries format_tags=StreamTitle -of default=nw=1:nk=1 %s", urlStream)
	out, err := executeCommand(cmd)
	if err != nil {
		log.Printf("Command execution failed: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write(out)
}

func status(w http.ResponseWriter, req *http.Request) {
	urlStream := req.URL.Query().Get("url")
	if urlStream == "" {
		log.Print("Url is not valid")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	duration := 10
	durationStr := req.URL.Query().Get("duration")
	if durationStr != "" {
		var err error
		duration, err = strconv.Atoi(durationStr)
		if err != nil {
			log.Printf("Command execution failed: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	cmd := fmt.Sprintf("ffmpeg -t %d -i %s -af volumedetect -f null - 2>&1 | grep mean_volume | cut -d ' ' -f 5", duration, urlStream)
	out, err := executeCommand(cmd)
	if err != nil {
		log.Printf("Command execution failed: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	feetStr := strings.TrimSpace(string(out))
	meanVolume, err := strconv.ParseFloat(feetStr, 64)
	if err != nil {
		log.Printf("Error parsing float: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if meanVolume < -7 {
		w.Write([]byte("online"))
		return
	}
	w.Write([]byte("offline"))
}

func executeCommand(cmd string) ([]byte, error) {
	command := exec.Command("sh", "-c", cmd)
	var out bytes.Buffer
	command.Stdout = &out

	err := command.Run()
	return out.Bytes(), err
}

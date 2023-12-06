package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var logger *slog.Logger

const urlPattern string = `^((ftp|http|https):\/\/)?(\S+(:\S*)?@)?((([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(((([a-z\x{00a1}-\x{ffff}0-9]+-?-?_?)*[a-z\x{00a1}-\x{ffff}0-9]+)\.)?)?(([a-z\x{00a1}-\x{ffff}0-9]+-?-?_?)*[a-z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-z\x{00a1}-\x{ffff}]{2,}))?)|localhost)(:(\d{1,5}))?((\/|\?|#)[^\s]*)?$`

func main() {
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	// logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger = slog.New(slog.NewJSONHandler(logFile, nil))

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/stream/play", play)
	mux.HandleFunc("/stream/status", status)

	logger.Info("Listening...")
	if err := http.ListenAndServe(":8090", mux); err != nil {
		logger.Error("Server failed to start", "error", err)
	}
}

func ping(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("pong"))
}

func play(w http.ResponseWriter, req *http.Request) {
	urlStream := req.URL.Query().Get("url")
	if !matches(urlStream, urlPattern) {
		logger.Info("Url is not valid", "url", urlStream)
		encodeResponse(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	cmd := fmt.Sprintf("ffprobe -v quiet -show_entries format_tags=StreamTitle -of default=nw=1:nk=1 %s", urlStream)
	out, err := executeCommand(cmd)
	if err != nil {
		logger.Error("Command ffprobe execution failed", "error", err, "url", urlStream)
		encodeResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	play := strings.Trim(string(out), "\n")
	encodeResponse(w, http.StatusOK, play)
}

func status(w http.ResponseWriter, req *http.Request) {
	urlStream := req.URL.Query().Get("url")
	if !matches(urlStream, urlPattern) {
		logger.Info("Url is not valid", "url", urlStream)
		encodeResponse(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	duration := 10
	volume := -70.0

	handleParseError := func(param string, err error) {
		logger.Error("Command execution failed", "error", err, param, req.URL.Query().Get(param))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	if durationStr := req.URL.Query().Get("duration"); durationStr != "" {
		var err error
		if duration, err = strconv.Atoi(durationStr); err != nil {
			handleParseError("duration", err)
			return
		}
	}

	if volumeStr := req.URL.Query().Get("volume"); volumeStr != "" {
		var err error
		if volume, err = strconv.ParseFloat(volumeStr, 64); err != nil {
			handleParseError("volume", err)
			return
		}
	}

	cmd := fmt.Sprintf("ffmpeg -t %d -i %s -af volumedetect -f null - 2>&1 | grep mean_volume | cut -d ' ' -f 5", duration, urlStream)
	out, err := executeCommand(cmd)
	if err != nil {
		logger.Error("Command ffmpeg execution failed", "error", err, "url", urlStream)
		encodeResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	feetStr := strings.TrimSpace(string(out))
	meanVolume, err := strconv.ParseFloat(feetStr, 64)
	if err != nil {
		logger.Error("Error parsing float", "error", err)
		encodeResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	if meanVolume > volume {
		encodeResponse(w, http.StatusOK, "online")
		return
	}
	encodeResponse(w, http.StatusOK, "offline")
}

func executeCommand(cmd string) ([]byte, error) {
	command := exec.Command("sh", "-c", cmd)
	var out bytes.Buffer
	command.Stdout = &out

	err := command.Run()
	return out.Bytes(), err
}

func encodeResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{
		Code:    statusCode,
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Error encode json data", "error", err)
	}
}

func matches(str, pattern string) bool {
	match, _ := regexp.MatchString(pattern, str)
	return match
}

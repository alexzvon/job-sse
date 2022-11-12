package handler

import (
	"encoding/json"
	"fmt"
	"github.com/alexzvon/job-sse/internal/splitter"
	"net/http"
)

type SseHandler struct {
	s *splitter.SseSplit
}

func NewSseHendler(sf *splitter.SseSplit) *SseHandler {
	return &SseHandler{s: sf}
}

func (h *SseHandler) ListenSseHandler(res http.ResponseWriter, req *http.Request) {
	cm := make(chan string, 10)

	h.setSseHeaders(res)
	h.s.Subscribe(cm)
	defer h.s.Unsubsctibe(cm)

	for {
		select {
		case word, ok := <-cm:
			if !ok {
				return
			}

			if word == "" {
				return
			}

			sw := SseWord{word}

			if err := h.writeSseHandler(res, sw); err != nil {
				return
			}
		}
	}
}

func (h *SseHandler) setSseHeaders(res http.ResponseWriter) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "text/handler-stream; charset=utf-8")
	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Connection", "keep-alive")
}

func (h *SseHandler) writeSseHandler(res http.ResponseWriter, word SseWord) error {
	b, err := json.Marshal(word)
	if err != nil {
		return err
	}

	_, err = res.Write(b)
	if err != nil {
		return err
	}

	_, err = res.Write([]byte("\n\n"))
	if err != nil {
		return err
	}

	res.(http.Flusher).Flush()

	return nil
}

func (h *SseHandler) SaySseHandler(res http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := req.Body.Close(); err != nil {
			fmt.Println("close req.Body: ", err)
			return
		}
	}()

	if req.Method != "POST" {
		return
	}

	var sw SseWord

	err := json.NewDecoder(req.Body).Decode(&sw)
	if err != nil {
		fmt.Println("Decode: ", err)
		return
	}

	h.s.Publish(sw.Word)

	return
}

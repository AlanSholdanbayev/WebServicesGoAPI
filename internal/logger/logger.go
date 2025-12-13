package logger

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Config struct {
	SeqURL    string
	SeqAPIKey string
}

// SeqWriter отправляет логи в Seq
type SeqWriter struct {
	Endpoint string
	ApiKey   string
	Client   *http.Client
}

func NewSeqWriter(endpoint, apiKey string) *SeqWriter {
	if endpoint != "" && endpoint[len(endpoint)-1] == '/' {
		endpoint = endpoint[:len(endpoint)-1]
	}
	return &SeqWriter{
		Endpoint: endpoint + "/api/events/raw",
		ApiKey:   apiKey,
		Client:   &http.Client{Timeout: 5 * time.Second},
	}
}

func (s *SeqWriter) Write(p []byte) (n int, err error) {
	if s.Endpoint == "" {
		return len(p), nil
	}

	payload := map[string]interface{}{
		"Events": []json.RawMessage{p},
	}
	data, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", s.Endpoint, bytes.NewReader(data))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	if s.ApiKey != "" {
		req.Header.Set("X-Seq-ApiKey", s.ApiKey)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		os.Stderr.WriteString("Seq error: " + err.Error() + "\n")
		return len(p), nil
	}
	defer resp.Body.Close()

	return len(p), nil
}

// LoggerWrapper оборачивает zerolog.Logger
type LoggerWrapper struct {
	*zerolog.Logger
}

func New(cfg *Config) *LoggerWrapper {
	seq := NewSeqWriter(cfg.SeqURL, cfg.SeqAPIKey)

	console := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	mw := zerolog.MultiLevelWriter(console, seq)
	l := zerolog.New(mw).With().Timestamp().Logger()

	return &LoggerWrapper{Logger: &l}
}

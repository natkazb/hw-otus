package internalhttp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/storage" //nolint
)

var (
	ErrNoDay      = `{"error": "day parameter is required"}`
	ErrInvalidDay = `{"error": "invalid date format, use YYYY-MM-DD"}`
	DayListFormat = "2006-01-02"
)

type EventHandler struct {
	app Application
}

func (h *EventHandler) List(startData, endData time.Time, w http.ResponseWriter, r *http.Request) {
	items, err := h.app.ListEvents(startData, endData)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(items); err != nil {
		jsonError(w, fmt.Sprintf("encoding error %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (h *EventHandler) ListDay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	day := r.URL.Query().Get("day")
	if day == "" {
		jsonError(w, ErrNoDay, http.StatusBadRequest)
		return
	}

	dayParsed, err := time.Parse(DayListFormat, day)
	if err != nil {
		jsonError(w, ErrInvalidDay, http.StatusBadRequest)
		return
	}

	h.List(dayParsed, dayParsed.Add(24*time.Hour), w, r)
}

func (h *EventHandler) ListWeek(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	day := r.URL.Query().Get("day")
	if day == "" {
		jsonError(w, ErrNoDay, http.StatusBadRequest)
		return
	}

	dayParsed, err := time.Parse(DayListFormat, day)
	if err != nil {
		jsonError(w, ErrInvalidDay, http.StatusBadRequest)
		return
	}

	h.List(dayParsed, dayParsed.Add(7*24*time.Hour), w, r)
}

func (h *EventHandler) ListMonth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	day := r.URL.Query().Get("day")
	if day == "" {
		jsonError(w, ErrNoDay, http.StatusBadRequest)
		return
	}

	dayParsed, err := time.Parse(DayListFormat, day)
	if err != nil {
		jsonError(w, ErrInvalidDay, http.StatusBadRequest)
		return
	}
	endDate := dayParsed.AddDate(0, 1, 0)

	h.List(dayParsed, endDate, w, r)
}

func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var event storage.Event
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		jsonError(w, fmt.Sprintf("read body ( %s )", err.Error()), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(data, &event); err != nil {
		jsonError(w, fmt.Sprintf("parse json ( %s )", err.Error()), http.StatusBadRequest)
		return
	}

	if err := h.app.CreateEvent(event); err != nil {
		jsonError(w, fmt.Sprintf("can not create ( %s )", err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"message": "Event created successfully",
		"event":   event,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		jsonError(w, fmt.Sprintf("encode response ( %s )", err.Error()), http.StatusInternalServerError)
	}
}

func jsonError(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, message)))
}

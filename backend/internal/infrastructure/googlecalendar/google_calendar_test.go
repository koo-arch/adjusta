package googlecalendar

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func TestUpdateOrCreateGoogleEventUpdatesExistingEventRepeatedly(t *testing.T) {
	t.Parallel()

	updateCalls := 0
	insertCalls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPut:
			updateCalls++
			if r.URL.Path != "/calendars/candidate-calendar/events/google-event" {
				t.Errorf("unexpected update path: %s", r.URL.Path)
			}
			fmt.Fprint(w, `{"id":"google-event"}`)
		case http.MethodPost:
			insertCalls++
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"error":{"code":500,"message":"insert should not be called"}}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestCalendarClient(t, server.URL)
	manager := NewGoogleCalendarManager()
	event := manager.ConvertToCalendarEvent(
		stringPointer("google-event"),
		"候補日程",
		"Tokyo",
		"description",
		time.Date(2026, 7, 20, 10, 0, 0, 0, time.UTC),
		time.Date(2026, 7, 20, 11, 0, 0, 0, time.UTC),
	)

	for i := 0; i < 2; i++ {
		updated, err := manager.UpdateOrCreateGoogleEvent(client, "candidate-calendar", event)
		if err != nil {
			t.Fatalf("UpdateOrCreateGoogleEvent returned error: %v", err)
		}
		if updated.Id != "google-event" {
			t.Fatalf("unexpected google event id: %s", updated.Id)
		}
	}

	if updateCalls != 2 {
		t.Fatalf("expected 2 update calls, got %d", updateCalls)
	}
	if insertCalls != 0 {
		t.Fatalf("expected no insert calls, got %d", insertCalls)
	}
}

func TestUpdateOrCreateGoogleEventRecreatesMissingEvent(t *testing.T) {
	t.Parallel()

	updateCalls := 0
	insertCalls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPut:
			updateCalls++
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"error":{"code":404,"message":"not found"}}`)
		case http.MethodPost:
			insertCalls++
			if r.URL.Path != "/calendars/candidate-calendar/events" {
				t.Errorf("unexpected insert path: %s", r.URL.Path)
			}
			fmt.Fprint(w, `{"id":"recreated-google-event"}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestCalendarClient(t, server.URL)
	manager := NewGoogleCalendarManager()
	event := &calendar.Event{Id: "missing-google-event", Summary: "候補日程"}

	recreated, err := manager.UpdateOrCreateGoogleEvent(client, "candidate-calendar", event)
	if err != nil {
		t.Fatalf("UpdateOrCreateGoogleEvent returned error: %v", err)
	}
	if recreated.Id != "recreated-google-event" {
		t.Fatalf("unexpected recreated google event id: %s", recreated.Id)
	}
	if updateCalls != 1 || insertCalls != 1 {
		t.Fatalf("expected one update and one insert, got update=%d insert=%d", updateCalls, insertCalls)
	}
}

func newTestCalendarClient(t *testing.T, endpoint string) *Client {
	t.Helper()

	service, err := calendar.NewService(
		context.Background(),
		option.WithEndpoint(endpoint+"/"),
		option.WithoutAuthentication(),
	)
	if err != nil {
		t.Fatalf("failed to create calendar service: %v", err)
	}
	return &Client{Service: service}
}

func stringPointer(value string) *string {
	return &value
}

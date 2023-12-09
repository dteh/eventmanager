package eventmanager

import (
	"os"
	"testing"
	"time"
)

var e *eventManager

func TestMain(m *testing.M) {
	e = NewEventManager("http://localhost:8080/ingest")
	os.Exit(m.Run())
}

func TestAddEvent(t *testing.T) {
	s := Event{
		Group: "test",
		Site:  "test",
		Key:   "test_key",
		Type:  STATUS_EVENT,
		Date:  time.Now().Format(time.RFC3339),
	}
	s2 := Event{
		Group: "test",
		Site:  "test",
		Key:   "test_key",
		Type:  STATUS_EVENT,
		Date:  time.Now().Format(time.RFC3339),
	}
	s3 := Event{
		Group: "test",
		Site:  "test2",
		Key:   "test_key",
		Type:  STATUS_EVENT,
		Date:  time.Now().Format(time.RFC3339),
	}

	s4 := Event{
		Group:   "test",
		Site:    "test",
		Key:     "test_key",
		Type:    NOTIFICATION_EVENT,
		Date:    time.Now().Format(time.RFC3339),
		Message: "TEST_MESSAGE",
	}
	s5 := Event{
		Group:   "test",
		Site:    "test",
		Key:     "test_key",
		Type:    NOTIFICATION_EVENT,
		Date:    time.Now().Format(time.RFC3339),
		Message: "TEST_MESSAGE",
	}
	e.AddEvent(s)
	e.AddEvent(s2)
	e.AddEvent(s3)
	e.AddEvent(s4)
	e.AddEvent(s5)

	if len(e.DumpAllEvents()) != 3 {
		t.Error("Expected 3 events")
	}
}

func TestSubmitEvent(t *testing.T) {
	s := Event{
		Group:   "test",
		Site:    "test",
		Key:     "ASDF",
		Type:    NOTIFICATION_EVENT,
		Date:    time.Now().Format(time.RFC3339),
		Message: "TEST_NOTIFICATION_MESSAGE",
	}
	e.AddEvent(s)

	err := e.SubmitEvents()
	if err != nil {
		t.Error(err)
	}
}

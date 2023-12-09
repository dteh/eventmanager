package eventmanager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type eventManager struct {
	events            *sync.Map
	ingestionEndpoint string
}

// Create a new event manager
func NewEventManager(ingestionEndpoint string) *eventManager {
	return &eventManager{
		events:            &sync.Map{},
		ingestionEndpoint: ingestionEndpoint,
	}
}

// Get the global event manager
func GetEventManager() (*eventManager, error) {
	if globalManager == nil {
		return nil, fmt.Errorf("event manager not initialized")
	}

	return globalManager, nil
}

// Initialize the global event manager
func InitializeEventManager(ingestionEndpoint string) {
	globalManager = NewEventManager(ingestionEndpoint)
	go globalManager.Start()
}

var globalManager *eventManager

type EventType string

const (
	STATUS_EVENT       EventType = "STATUS"
	NOTIFICATION_EVENT EventType = "NOTIFICATION"
)

type Event struct {
	Group   string    `json:"group"`
	Site    string    `json:"site"`
	Key     string    `json:"key"`
	Type    EventType `json:"type"`
	Date    string    `json:"date"`
	Message string    `json:"message,omitempty"`
}

func (e Event) EventKey() string {
	return e.Group + e.Site + e.Key + string(e.Type) + e.Message
}

type Keyable interface {
	EventKey() string
}

func (e *eventManager) AddEvent(event Keyable) {
	e.events.Store(event.EventKey(), event)
}

func (e *eventManager) DumpAllEvents() []any {
	events := []any{}
	e.events.Range(func(key, value any) bool {
		events = append(events, value)
		e.events.Delete(key)
		return true
	})
	return events
}

func (e *eventManager) SubmitEvents() error {
	events := e.DumpAllEvents()
	if len(events) == 0 {
		return nil
	}
	b, err := json.Marshal(events)
	if err != nil {
		return err
	}
	resp, err := http.Post(e.ingestionEndpoint, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("submit events got status code %d", resp.StatusCode)
	}
	return nil
}

func (e *eventManager) Start() {
	for {
		time.Sleep(5 * time.Minute)
		err := e.SubmitEvents()
		if err != nil {
			fmt.Println("unable to submit events:", err)
		}
	}
}

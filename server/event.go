package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	CopyAndPaste = "copyAndPaste"
	TimeTaken    = "timeTaken"
	Resize       = "resize"
)

type ScreenResizeEvent struct {
	WebsiteUrl string    `json:"websiteUrl"`
	SessionId  string    `json:"sessionId"`
	ResizeFrom Dimension `json:"resizeFrom"`
	ResizeTo   Dimension `json:"resizeTo"`
}

type CopyPasteEvent struct {
	WebsiteUrl string `json:"websiteUrl"`
	SessionId  string `json:"sessionId"`
	Pasted     bool   `json:"pasted"`
	FieldName  string `json:"fieldName"`
}

type TimeTakenEvent struct {
	WebsiteUrl string `json:"websiteUrl"`
	SessionId  string `json:"sessionId"`
	TimeTaken  int    `json:"timeTaken"`
}

type Event struct {
	WebsiteUrl         string
	SessionId          string
	ResizeFrom         Dimension
	ResizeTo           Dimension
	CopyAndPaste       map[string]bool
	FormCompletionTime int
}

type Dimension struct {
	Width  string `json:"width"`
	Height string `json:"height"`
}

func (e Event) String() string {
	return fmt.Sprintf("Website URL: %s\nSession ID: %s\nResized from: (%s, %s)\nResized to: (%s, %s)\nCopy and pasted fields: %v\nForm completion time: %d\n",
		e.WebsiteUrl, e.SessionId, e.ResizeFrom.Width, e.ResizeFrom.Height, e.ResizeTo.Width, e.ResizeTo.Height, e.CopyAndPaste, e.FormCompletionTime)
}

func PopulateCopyPasteEvent(e Event, b []byte) (Event, error) {
	copyPasteEvent := CopyPasteEvent{}
	err := json.Unmarshal(b, &copyPasteEvent)
	if err != nil {
		return Event{}, errors.New("Unable to unmarshal copy paste event data")
	}

	if copyPasteEvent.WebsiteUrl == "" || copyPasteEvent.FieldName == "" {
		return Event{}, errors.New("Invalid data passed for copy and paste event")
	}

	if copyPasteEvent.FieldName != "email" && copyPasteEvent.FieldName != "cardNumber" && copyPasteEvent.FieldName != "CVV" {
		return Event{}, errors.New("Invalid field name passed in event")
	}

	if e.CopyAndPaste == nil {
		e.CopyAndPaste = make(map[string]bool)
		e.CopyAndPaste["email"] = false
		e.CopyAndPaste["cardNumber"] = false
		e.CopyAndPaste["CVV"] = false
	}

	if e.SessionId == "" {
		e.SessionId = copyPasteEvent.SessionId
	}

	e.WebsiteUrl = copyPasteEvent.WebsiteUrl
	e.CopyAndPaste[copyPasteEvent.FieldName] = copyPasteEvent.Pasted
	return e, nil
}

func PopulateResizeEvent(e Event, b []byte) (Event, error) {
	if e.ResizeFrom.Width != "" {
		return e, nil
	}

	screenResizeEvent := ScreenResizeEvent{}
	err := json.Unmarshal(b, &screenResizeEvent)
	if err != nil {
		return Event{}, errors.New("Unable to unmarshal resize event data")
	}

	if screenResizeEvent.WebsiteUrl == "" || screenResizeEvent.ResizeFrom.Height == "" || screenResizeEvent.ResizeFrom.Width == "" || screenResizeEvent.ResizeTo.Height == "" || screenResizeEvent.ResizeTo.Width == "" {
		return Event{}, errors.New("Invalid data passed for resize event")
	}

	if e.SessionId == "" {
		e.SessionId = screenResizeEvent.SessionId
	}

	e.WebsiteUrl = screenResizeEvent.WebsiteUrl
	e.ResizeFrom = screenResizeEvent.ResizeFrom
	e.ResizeTo = screenResizeEvent.ResizeTo
	return e, nil
}

func PopulateTimeTakenEvent(e Event, b []byte) (Event, error) {
	if e.FormCompletionTime != 0 {
		return e, nil
	}	

	timeTakenEvent := TimeTakenEvent{}
	err := json.Unmarshal(b, &timeTakenEvent)
	if err != nil {
		return Event{}, errors.New("Unable to unmarshal time taken event data")
	}

	if timeTakenEvent.TimeTaken == 0 {
		return Event{}, errors.New("Invalid data passed for time taken event")
	}

	if e.SessionId == "" {
		e.SessionId = timeTakenEvent.SessionId
	}

	if e.WebsiteUrl == "" {
		e.WebsiteUrl = timeTakenEvent.WebsiteUrl
	}

	e.FormCompletionTime = timeTakenEvent.TimeTaken
	return e, nil
}

func ParseEvent(eventType string, e Event, b []byte) (Event, error) {
	switch eventType {
	case CopyAndPaste:
		return PopulateCopyPasteEvent(e, b)
	case TimeTaken:
		return PopulateTimeTakenEvent(e, b)
	case Resize:
		return PopulateResizeEvent(e, b)
	default:
		return Event{}, errors.New(fmt.Sprintf("Invalid event type: %s", eventType))
	}
}

func IsComplete(e Event) bool {
	return e.FormCompletionTime != 0
}

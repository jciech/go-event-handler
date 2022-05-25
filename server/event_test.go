package main

import (
	"testing"
)

func TestNewEventCopyPaste(t *testing.T) {
	testJson := []byte(`{"eventType": "copyAndPaste","sessionId": "123","websiteUrl": "websiteUrl","fieldName": "CVV","pasted": true}`)
	e := Event{}
	e, err := ParseEvent(CopyAndPaste, e, testJson)
	if e.SessionId != "123" || e.WebsiteUrl != "websiteUrl" || !e.CopyAndPaste["CVV"] {
		t.Fatalf("Failed to parse copy and paste event: %s", err)
	}
}

func TestNewEventResize(t *testing.T) {
	testJson := []byte(`{
	"eventType": "resize",
	"sessionId": "123",
	"websiteUrl": "websiteUrl",
	"resizeFrom": {"width": "20", "height": "25"},
	"resizeTo": {"width": "30", "height": "40"}
	}`)
	e := Event{}
	e, err := ParseEvent(Resize, e, testJson)
	if e.SessionId != "123" || e.WebsiteUrl != "websiteUrl" || e.ResizeFrom.Width != "20" || e.ResizeFrom.Height != "25" || e.ResizeTo.Width != "30" || e.ResizeTo.Height != "40" {
		t.Fatalf("Failed to parse resize event: %s", err)
	}
}

func TestNewEventTimeTaken(t *testing.T) {
	testJson := []byte(`{
		"eventType": "timeTaken",
		"sessionId": "123",
		"websiteUrl": "websiteUrl",
		"timeTaken": 100
	}`)
	e := Event{}
	e, err := ParseEvent(TimeTaken, e, testJson)
	if e.SessionId != "123" || e.WebsiteUrl != "websiteUrl" || e.FormCompletionTime != 100 {
		t.Fatalf("Failed to parse time taken event: %s", err)
	}
}


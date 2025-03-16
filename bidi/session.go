package bidi

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// Session represents a BiDi session with the browser.
type Session struct {
	cdp      CDPSession
	script   *Script
	handlers map[string][]func(json.RawMessage)
	mu       sync.RWMutex
	closed   bool
}

// NewSession creates a new BiDi session.
func NewSession(ws WebSocket, sessionID, targetID string) *Session {
	cdp := NewCDPSession(ws, sessionID, targetID)
	session := &Session{
		cdp:      cdp,
		script:   nil,
		handlers: make(map[string][]func(json.RawMessage)),
		mu:       sync.RWMutex{},
		closed:   false,
	}
	session.script = NewScript(cdp)

	return session
}

// Script returns the script executor for this session.
func (s *Session) Script() ScriptEvaluator {
	return s.script
}

// Subscribe adds an event handler for the specified event type.
func (s *Session) Subscribe(eventType string, handler func(json.RawMessage)) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.handlers[eventType] = append(s.handlers[eventType], handler)
}

// Unsubscribe removes all handlers for the specified event type.
func (s *Session) Unsubscribe(eventType string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.handlers, eventType)
}

// HandleEvent processes an incoming event.
func (s *Session) HandleEvent(eventType string, data json.RawMessage) {
	s.mu.RLock()
	handlers := s.handlers[eventType]
	s.mu.RUnlock()

	for _, handler := range handlers {
		handler(data)
	}
}

// Close closes the session and its associated resources.
func (s *Session) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true
	s.handlers = nil

	if err := s.cdp.Close(); err != nil {
		return fmt.Errorf("failed to close CDP session: %w", err)
	}

	return nil
}

// Execute sends a command to the browser and waits for the response.
func (s *Session) Execute(ctx context.Context, method string, params interface{}) (json.RawMessage, error) {
	result, err := s.cdp.Execute(ctx, method, params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

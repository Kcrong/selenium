package bidi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

// BrowserError represents an error returned by the browser.
type BrowserError struct {
	Detail  interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
}

func (e *BrowserError) Error() string {
	return fmt.Sprintf("BrowserError<code=%d message=%s> %v", e.Code, e.Message, e.Detail)
}

// CDPSessionImpl represents a CDP session implementation.
type CDPSessionImpl struct {
	ws        WebSocket
	inflight  map[int64]chan json.RawMessage
	sessionID string
	targetID  string
	nextID    atomic.Int64
	mu        sync.RWMutex
	closed    bool
}

// NewCDPSession creates a new CDP session.
func NewCDPSession(ws WebSocket, sessionID, targetID string) CDPSession {
	//nolint:exhaustruct // Initialize memory-sensitive fields only for better readability.
	return &CDPSessionImpl{
		ws:        ws,
		sessionID: sessionID,
		targetID:  targetID,
		inflight:  make(map[int64]chan json.RawMessage),
	}
}

var ErrSessionClosed = errors.New("session is closed")

// Execute sends a command to the browser and waits for the response.
func (s *CDPSessionImpl) Execute(ctx context.Context, method string, params interface{}) ([]byte, error) {
	if s.closed {
		return nil, ErrSessionClosed
	}

	id := s.nextID.Add(1)
	cmd := struct {
		Params    interface{} `json:"params,omitempty"`
		Method    string      `json:"method"`
		SessionID string      `json:"sessionId,omitempty"`
		ID        int64       `json:"id"`
	}{
		ID:        id,
		Method:    method,
		Params:    params,
		SessionID: s.sessionID,
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal command: %w", err)
	}

	resultCh := make(chan json.RawMessage, 1)

	s.mu.Lock()
	s.inflight[id] = resultCh
	s.mu.Unlock()

	if err := s.ws.Send(data); err != nil {
		s.mu.Lock()
		delete(s.inflight, id)
		s.mu.Unlock()

		return nil, fmt.Errorf("failed to send command: %w", err)
	}

	select {
	case result := <-resultCh:
		return result, nil
	case <-ctx.Done():
		s.mu.Lock()
		delete(s.inflight, id)
		s.mu.Unlock()

		return nil, ctx.Err()
	}
}

// Close closes the CDP session.
func (s *CDPSessionImpl) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true
	for _, ch := range s.inflight {
		close(ch)
	}

	s.inflight = nil

	return s.ws.Close()
}

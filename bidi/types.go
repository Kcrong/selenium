package bidi

import "context"

// WebSocket interface defines the methods required for WebSocket communication.
type WebSocket interface {
	Send(data []byte) error
	Receive() ([]byte, error)
	Close() error
}

// CDPSession represents a CDP session with the browser.
type CDPSession interface {
	Execute(ctx context.Context, method string, params interface{}) ([]byte, error)
	Close() error
}

// ScriptEvaluator defines the interface for script evaluation.
type ScriptEvaluator interface {
	// EvaluateScript evaluates a JavaScript script and returns the result
	EvaluateScript(ctx context.Context, script string, args []interface{}) (*ScriptResult, error)
	// CallFunction calls a JavaScript function and returns the result
	CallFunction(ctx context.Context, functionDeclaration string, args []interface{}) (*ScriptResult, error)
}

package bidi

import (
	"context"
	"encoding/json"
	"fmt"
)

// ScriptResult represents the result of a script execution.
type ScriptResult struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

// Script provides methods for script execution.
type Script struct {
	session CDPSession
}

// NewScript creates a new Script instance.
func NewScript(session CDPSession) *Script {
	return &Script{
		session: session,
	}
}

// EvaluateScript evaluates a JavaScript script in the current browsing context.
func (s *Script) EvaluateScript(ctx context.Context, script string, args []interface{}) (*ScriptResult, error) {
	params := map[string]interface{}{
		"expression": script,
		"arguments":  args,
	}

	result, err := s.session.Execute(ctx, "script.evaluate", params)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate script: %w", err)
	}

	var scriptResult ScriptResult
	if err := json.Unmarshal(result, &scriptResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal script result: %w", err)
	}

	return &scriptResult, nil
}

// CallFunction calls a JavaScript function in the current browsing context.
func (s *Script) CallFunction(
	ctx context.Context, functionDeclaration string, args []interface{},
) (*ScriptResult, error) {
	params := map[string]interface{}{
		"functionDeclaration": functionDeclaration,
		"arguments":           args,
	}

	result, err := s.session.Execute(ctx, "script.callFunction", params)
	if err != nil {
		return nil, fmt.Errorf("failed to call function: %w", err)
	}

	var scriptResult ScriptResult
	if err := json.Unmarshal(result, &scriptResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal function result: %w", err)
	}

	return &scriptResult, nil
}

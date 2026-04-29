package ctxerror

import (
	"encoding/json"
	"strings"
	"testing"
)

type nilReceiverError struct {
	message string
}

func (err *nilReceiverError) Error() string {
	return err.message
}

type functionFieldError struct {
	message string
	format  func(string) string
}

func (err functionFieldError) Error() string {
	return err.format(err.message)
}

func TestCtxErrorMarshalJSONHandlesTypedNilError(t *testing.T) {
	var err *nilReceiverError

	ctxErr := CtxError{
		Message: "outer",
		ErrorS:  "fallback error",
		ErrorI:  err,
	}

	bytes, marshalErr := json.Marshal(ctxErr)
	if marshalErr != nil {
		t.Fatalf("json.Marshal returned error: %v", marshalErr)
	}

	if !strings.Contains(string(bytes), `"error":"fallback error"`) {
		t.Fatalf("expected fallback error in JSON, got %s", string(bytes))
	}
}

func TestCtxErrorTraceErrorHandlesTypedNilError(t *testing.T) {
	var err *nilReceiverError

	trace := CtxErrorTrace{Trace: []CtxError{{
		Message: "outer",
		ErrorS:  "fallback error",
		ErrorI:  err,
	}}}

	got := trace.Error()
	if !strings.Contains(got, `"error":"fallback error"`) {
		t.Fatalf("expected fallback error in trace, got %s", got)
	}
}

func TestWrapHandlesTypedNilError(t *testing.T) {
	var err *nilReceiverError

	trace := Wrap(err, "outer")
	if trace == nil {
		t.Fatal("expected trace for non-nil error interface")
	}

	if got := trace.Error(); !strings.Contains(got, `"message":"outer"`) {
		t.Fatalf("expected wrapped trace, got %s", got)
	}
}

func TestWrapHandlesErrorWithFunctionField(t *testing.T) {
	err := functionFieldError{
		message: "inner",
		format: func(message string) string {
			return "formatted " + message
		},
	}

	trace := Wrap(err, "outer")
	if trace == nil {
		t.Fatal("expected trace")
	}

	got := trace.Error()
	if !strings.Contains(got, `"error":"formatted inner"`) {
		t.Fatalf("expected formatted error in trace, got %s", got)
	}
}

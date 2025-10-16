package customerrors

import (
	"errors"
	"fmt"
)

var (
	ErrContentTypeAJ   = errors.New("content type must be application/json")
	ErrEmptyUserID     = errors.New("user_id must not be empty")
	ErrEmptyID         = errors.New("id must not be empty")
	ErrEmptyTitle      = errors.New("title must not be empty")
	ErrEmptyDate       = errors.New("date must not be empty")
	ErrUserIDNotFound  = errors.New("user id not found")
	ErrEventIDNotFound = errors.New("event id not found")
)

// AppError — custom error with context, HTTP status and trace support
type AppError struct {
	Op         string // operation / context
	Err        error  // underlying error
	HTTPStatus int    // HTTP status code
	TraceID    string // trace id for distributed systems
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Op, e.Err)
	}
	return e.Op
}

// Unwrap supports errors.Is / errors.As
func (e *AppError) Unwrap() error {
	return e.Err
}

// Wrap wraps an existing error with operation context
func Wrap(err error, op string) error {
	if err == nil {
		return nil
	}
	return &AppError{
		Op:  op,
		Err: err,
	}
}

// Wrapf wraps an existing error with formatted operation context
func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return &AppError{
		Op:  fmt.Sprintf(format, args...),
		Err: err,
	}
}

// New creates a new AppError without underlying error
func New(op string, httpStatus int) *AppError {
	return &AppError{
		Op:         op,
		HTTPStatus: httpStatus,
	}
}

// WithTrace sets the trace ID for distributed tracing
func (e *AppError) WithTrace(traceID string) *AppError {
	e.TraceID = traceID
	return e
}

// Helpers for type checking

// IsAppError checks if an error is of type AppError
func IsAppError(err error) bool {
	var e *AppError
	return errors.As(err, &e)
}

// ExtractOp extracts the Op field from AppError if possible
func ExtractOp(err error) string {
	var e *AppError
	if errors.As(err, &e) {
		return e.Op
	}
	return ""
}

// ExtractHTTPStatus extracts HTTPStatus from AppError, returns 500 if not AppError
func ExtractHTTPStatus(err error) int {
	var e *AppError
	if errors.As(err, &e) && e.HTTPStatus != 0 {
		return e.HTTPStatus
	}
	return 500
}

// ExtractTraceID extracts TraceID from AppError
func ExtractTraceID(err error) string {
	var e *AppError
	if errors.As(err, &e) {
		return e.TraceID
	}
	return ""
}

/*
import (
	"github.com/rs/zerolog/log"
	"myapp/internal/customerrors"
	"errors"
)

func HandleRequest() {
	err := doSomething()
	if err != nil {
		var appErr *customerrors.AppError
		if errors.As(err, &appErr) {
			log.Error().
				Str("op", appErr.Op).
				Int("http_status", appErr.HTTPStatus).
				Str("trace_id", appErr.TraceID).
				Err(appErr.Err).
				Msg("operation failed")
		} else {
			log.Error().Err(err).Msg("unknown error")
		}
	}
}

func doSomething() error {
	// simulate an error
	return customerrors.Wrapf(customerrors.New("db connection failed", 503), "load user %d", 123)
}
*/

/*
{
  "level":"error",
  "op":"load user 123",
  "http_status":503,
  "trace_id":"abcd-1234",
  "error":"db connection failed",
  "message":"operation failed"
}
*/

package apidiags

import (
	"strings"

	"github.com/go-openapi/jsonpointer"
)

// Severity indicates whether the diagnostic is advisory or fatal.
type Severity string

const (
	// DiagnosticError is a Severity used when the diagnostic should be
	// treated as fatal.
	DiagnosticError Severity = "error"
	// DiagnosticWarning is a Severity used when the diagnostic should be
	// treated as advisory.
	DiagnosticWarning Severity = "warning"
)

// Code indicates the type of failure or situation that a diagnostic is
// communicating.
type Code string

const (
	// CodeAccessDenied indicates that the user does not have sufficient
	// permissions to perform that action.
	CodeAccessDenied Code = "access_denied"
	// CodeInsufficient indicates that the value provided was somehow
	// insufficient; too low, not long enough, not enough items, etc.
	CodeInsufficient Code = "insufficient"
	// CodeOverflow indicates that the value provided exceed some sort of
	// bounds; too high, too long, too many items, etc.
	CodeOverflow Code = "overflow"
	// CodeInvalidValue indicates that an unacceptable value was supplied.
	CodeInvalidValue Code = "invalid_value"
	// CodeInvalidFormat indicates that the value came in a format that we
	// don't know how to read; the wrong type, the wrong encoding, etc.
	CodeInvalidFormat Code = "invalid_format"
	// CodeMissing indicates that the value was expected but wasn't
	// present.
	CodeMissing Code = "missing"
	// CodeNotFound indicates that the resource specified by the value
	// wasn't found.
	CodeNotFound Code = "not_found"
	// CodeConflict means two or more values cannot be used together, or
	// cannot be set to those values.
	CodeConflict Code = "conflict"
	// CodeActOfGod indicates that something outside the user's control
	// happened that disrupted the request, and they should try again.
	CodeActOfGod Code = "act_of_god"
	// CodeDeprecated indicates that the field or value is deprecated and
	// may be removed or have its behavior changed in a future API update.
	CodeDeprecated Code = "deprecated"
)

type pointerType string

const (
	bodyPointer   pointerType = "body"
	urlPointer    pointerType = "url"
	headerPointer pointerType = "header"
)

// Diagnostic supplies information about the API and its status to the caller.
// It can be used to inform callers about errors and to warn them of future
// deprecations.
type Diagnostic struct {
	Severity Severity  `json:"severity"`
	Code     Code      `json:"code"`
	Pointers []Pointer `json:"pointers,omitempty"`
}

// Pointer indicates what part of a request prompted a diagnostic. The Field
// indicates whether the value was in the body, the URL, or a header. The Path
// is a JSON pointer to the value in that field that prompted the diagnostic.
//
// Pointers are used to give consumers more information about what specifically
// went wrong, and to help consumers craft UIs that can prompt users to fix the
// problem.
type Pointer struct {
	Field pointerType `json:"field"`
	Path  string      `json:"path"`
}

func buildPointerPath(segments []string) string {
	if len(segments) < 1 {
		return ""
	}
	sanitized := make([]string, 0, len(segments))
	for _, segment := range segments {
		sanitized = append(sanitized, jsonpointer.Escape(segment))
	}
	return "/" + strings.Join(sanitized, "/")
}

// NewBodyPointer returns a Pointer to a value in the request body consisting
// of the passed segments. It will safely escape the segments as necessary.
func NewBodyPointer(segments ...string) Pointer {
	return Pointer{
		Field: bodyPointer,
		Path:  buildPointerPath(segments),
	}
}

// NewURLPointer returns a Pointer to a value in the request URL consisting of
// the passed segments. It will safely escape the segments as necessary.
func NewURLPointer(segments ...string) Pointer {
	return Pointer{
		Field: urlPointer,
		Path:  buildPointerPath(segments),
	}
}

// NewHeaderPointer returns a Pointer to a value in the request headers
// consisting of the passed segments. It will safely escape the segments as
// necessary.
func NewHeaderPointer(segments ...string) Pointer {
	return Pointer{
		Field: headerPointer,
		Path:  buildPointerPath(segments),
	}
}

package apidiags

import (
	"bytes"
	"encoding/json"
	"fmt"
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

// Diagnostic supplies information about the API and its status to the caller.
// It can be used to inform callers about errors and to warn them of future
// deprecations.
type Diagnostic struct {
	Severity Severity `json:"severity"`
	Code     Code     `json:"code"`
	Paths    []Steps  `json:"path,omitempty"`
}

// Steps are a collection of transforms or accesses that point to
// ever-more-specific parts of a request.
type Steps []Step

// AddStep appends a Step to the Steps, pointing to another level of
// specificity for the request.
func (steps Steps) AddStep(step Step) Steps {
	steps = append(steps, step)
	return steps
}

// UnmarshalJSON turns a JSON-encoded set of bytes into Steps.
func (steps *Steps) UnmarshalJSON(in []byte) error {
	var genSteps []genericStep
	dec := json.NewDecoder(bytes.NewBuffer(in))
	dec.UseNumber()
	err := dec.Decode(&genSteps)
	if err != nil {
		return err
	}
	results := make(Steps, 0, len(genSteps))
	for pos, step := range genSteps {
		switch step.Kind {
		case "body":
			results = results.AddStep(BodyStep{})
		case "header":
			if step.Value == nil {
				return fmt.Errorf("error parsing step %d: no value", pos)
			}
			header, ok := (*step.Value).(string)
			if !ok {
				return fmt.Errorf("error parsing step %d: wanted string, got %T", pos, *step.Value)
			}
			results = results.AddStep(HeaderStep(header))
		case "url_param":
			if step.Value == nil {
				return fmt.Errorf("error parsing step %d: no value", pos)
			}
			param, ok := (*step.Value).(string)
			if !ok {
				return fmt.Errorf("error parsing step %d: wanted string, got %T", pos, *step.Value)
			}
			results = results.AddStep(URLParamStep(param))
		case "array_index":
			if step.Value == nil {
				return fmt.Errorf("error parsing step %d: no value", pos)
			}
			index, ok := (*step.Value).(json.Number)
			if !ok {
				return fmt.Errorf("error parsing step %d: wanted json.Number, got %T", pos, *step.Value)
			}
			idx, err := index.Int64()
			if err != nil {
				return fmt.Errorf("error parsing step %d: %w", pos, err)
			}
			results = results.AddStep(ArrayIndexStep(idx))
		case "object_property":
			if step.Value == nil {
				return fmt.Errorf("error parsing step %d: no value", pos)
			}
			property, ok := (*step.Value).(string)
			if !ok {
				return fmt.Errorf("error parsing step %d: wanted string, got %T", pos, *step.Value)
			}
			results = results.AddStep(ObjectPropertyStep(property))
		case "string_index":
			if step.Value == nil {
				return fmt.Errorf("error parsing step %d: no value", pos)
			}
			index, ok := (*step.Value).(json.Number)
			if !ok {
				return fmt.Errorf("error parsing step %d: wanted json.Number, got %T", pos, *step.Value)
			}
			idx, err := index.Int64()
			if err != nil {
				return fmt.Errorf("error parsing step %d: %w", pos, err)
			}
			results = results.AddStep(StringIndexStep(idx))
		default:
			return fmt.Errorf("error parsing step %d: unexpected step kind %q with value type %T", pos, step.Kind, step.Value)
		}
	}
	*steps = results
	return nil
}

// MarshalJSON turns Steps into a JSON-encoded set of bytes.
func (steps Steps) MarshalJSON() ([]byte, error) {
	genSteps := make([]genericStep, 0, len(steps))
	for pos, step := range steps {
		switch value := step.(type) {
		case BodyStep:
			genSteps = append(genSteps, genericStep{Kind: "body"})
		case HeaderStep:
			val := any(string(value))
			genSteps = append(genSteps, genericStep{Kind: "header", Value: &val})
		case URLParamStep:
			val := any(string(value))
			genSteps = append(genSteps, genericStep{Kind: "url_param", Value: &val})
		case ArrayIndexStep:
			val := any(int64(value))
			genSteps = append(genSteps, genericStep{Kind: "array_index", Value: &val})
		case ObjectPropertyStep:
			val := any(string(value))
			genSteps = append(genSteps, genericStep{Kind: "object_property", Value: &val})
		case StringIndexStep:
			val := any(int64(value))
			genSteps = append(genSteps, genericStep{Kind: "string_index", Value: &val})
		default:
			return nil, fmt.Errorf("unknown step type %T for step %d", step, pos)
		}
	}
	return json.Marshal(genSteps)
}

type genericStep struct {
	Kind  string `json:"kind,omitempty"`
	Value *any   `json:"value,omitempty"`
}

// Step is a single transform or access that points to a more specific part of
// the request. It should only ever be created by using the *Step functions.
type Step interface {
	step()
}

// BodyStep is a Step that selects the body of a request.
type BodyStep struct{}

func (BodyStep) step() {}

// HeaderStep is a Step that specifies a single header on a request.
type HeaderStep string

func (HeaderStep) step() {}

// URLParamStep is a Step that specifies a single URL parameter on a request.
type URLParamStep string

func (URLParamStep) step() {}

// ArrayIndexStep is a Step that specifies a single element within an array.
type ArrayIndexStep int64

func (ArrayIndexStep) step() {}

// ObjectPropertyStep is a Step that specifies a single property within an
// object.
type ObjectPropertyStep string

func (ObjectPropertyStep) step() {}

// StringIndexStep is a Step that specifies a single character within a string.
type StringIndexStep int64

func (StringIndexStep) step() {}

// BodyPath returns Steps that point to the body of the request.
func BodyPath() Steps {
	return Steps{BodyStep{}}
}

// HeaderPath returns Steps that point to the specified header of the request.
func HeaderPath(header string) Steps {
	return Steps{HeaderStep(header)}
}

// URLParamPath returns steps that point to the specified URL parameter of the
// request.
func URLParamPath(param string) Steps {
	return Steps{URLParamStep(param)}
}

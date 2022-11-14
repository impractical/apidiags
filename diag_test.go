package apidiags

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nsf/jsondiff"
)

func TestMarshalJSON(t *testing.T) {
	t.Parallel()

	type testCase struct {
		steps    Steps
		expected string
	}

	cases := map[string]testCase{
		"no-steps": {
			expected: "[]",
		},
		"body-step": {
			steps:    BodyPath(),
			expected: `[{"kind": "body"}]`,
		},
		"body-prop": {
			steps:    BodyPath().AddStep(ObjectPropertyStep("foo")),
			expected: `[{"kind": "body"}, {"kind": "object_property", "value": "foo"}]`,
		},
		"body-prop-arrayIndex": {
			steps: BodyPath().
				AddStep(ObjectPropertyStep("foo")).
				AddStep(ArrayIndexStep(0)),
			expected: `[{"kind": "body"}, {"kind": "object_property", "value": "foo"}, {"kind": "array_index", "value": 0}]`,
		},
		"body-prop-arrayIndex-stringIndex": {
			steps: BodyPath().
				AddStep(ObjectPropertyStep("foo")).
				AddStep(ArrayIndexStep(1)).
				AddStep(StringIndexStep(0)),
			expected: `[{"kind": "body"}, {"kind": "object_property", "value": "foo"}, {"kind": "array_index", "value": 1}, {"kind": "string_index", "value": 0}]`,
		},
		"header-step": {
			steps:    HeaderPath("foo"),
			expected: `[{"kind": "header", "value": "foo"}]`,
		},
		"urlParam-step": {
			steps:    URLParamPath("foo"),
			expected: `[{"kind": "url_param", "value": "foo"}]`,
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := json.Marshal(tc.steps)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			opts := jsondiff.DefaultConsoleOptions()
			match, diff := jsondiff.Compare([]byte(tc.expected), result, &opts)
			if match != jsondiff.FullMatch {
				t.Errorf("Unexpected result: %s", diff)
			}
			if match > jsondiff.NoMatch {
				t.Logf("first argument: %s", tc.expected)
				t.Logf("second argument: %s", result)
			}
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input    string
		expected Steps
	}

	cases := map[string]testCase{
		"no-steps": {
			input:    "[]",
			expected: Steps{},
		},
		"body-step": {
			input:    `[{"kind": "body"}]`,
			expected: BodyPath(),
		},
		"body-prop": {
			input:    `[{"kind": "body"}, {"kind": "object_property", "value": "foo"}]`,
			expected: BodyPath().AddStep(ObjectPropertyStep("foo")),
		},
		"body-prop-arrayIndex": {
			input: `[{"kind": "body"}, {"kind": "object_property", "value": "foo"}, {"kind": "array_index", "value": 0}]`,
			expected: BodyPath().
				AddStep(ObjectPropertyStep("foo")).
				AddStep(ArrayIndexStep(0)),
		},
		"body-prop-arrayIndex-stringIndex": {
			input: `[{"kind": "body"}, {"kind": "object_property", "value": "foo"}, {"kind": "array_index", "value": 1}, {"kind": "string_index", "value": 0}]`,
			expected: BodyPath().
				AddStep(ObjectPropertyStep("foo")).
				AddStep(ArrayIndexStep(1)).
				AddStep(StringIndexStep(0)),
		},
		"header-step": {
			input:    `[{"kind": "header", "value": "foo"}]`,
			expected: HeaderPath("foo"),
		},
		"urlParam-step": {
			input:    `[{"kind": "url_param", "value": "foo"}]`,
			expected: URLParamPath("foo"),
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var result Steps
			err := json.Unmarshal([]byte(tc.input), &result)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if diff := cmp.Diff(tc.expected, result); diff != "" {
				t.Fatalf("unexpected results (-wanted, +got): %s", diff)
			}
		})
	}
}

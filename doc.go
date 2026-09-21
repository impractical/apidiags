// Package apidiags provides types and helpers for formatting HTTP API
// diagnostic responses.
//
// Diagnostic responses consist of three things: a severity, a code, and one or
// more paths.
//
// The severity indicates whether a diagnostic is an error (must be remedied
// before the request can succeed) or a warning (the request succeeded, but
// there's relevant information to know about that success).
//
// The code indicates the type of diagnostic being returned. See the comments
// on each of the code constants for information on the various codes.
//
// The paths indicate the parts of the request that contributed to that
// specific code being returned.
package apidiags

# apidiags

Package `apidiags` provides types to construct API-level Diagnostics,
information that can be returned with a response to specify errors or warnings.
They are shamelessly "inspired" by Martin Atkins' work on
[zcl](https://github.com/zclconf/go-zcl) and the Diagnostics introduced there.

A Diagnostic indicates some information about an API request that should be
returned to the caller. It has a Severity (indicating whether it should be
treated as an error or warning), a Code (indicating the information being
communicated), and Paths (indicating the part(s) of the request that triggered
the Diagnostic). Each time a Diagnostic would be triggered, even if it uses the
same Code, should result in a new Diagnostic being added to the request rather
than adding another Path to the Diagnostic. The ability to specify multiple
Paths in a single Diagnostic is intended to allow Diagnostics to describe a
conflict between two parts of a request, or other similar situations where a
single part of a request is insufficient for indicating what went wrong.

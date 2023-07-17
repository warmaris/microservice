// Package errors provides typed errors which all we need in our projects.
// There are several types implementing Error interface, with constructors. So we can use type assertion or errors.As
// in server implementation to define response status. Also, it can be used in tests.
// See [server/error](../server/error.go)
package errors

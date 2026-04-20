package goscript

import "fmt"

// Option models an optional value without forcing nil checks everywhere.
type Option struct {
	Value interface{}
	Ok    bool
}

// Some creates a populated option.
func Some(value interface{}) Option {
	return Option{Value: value, Ok: true}
}

// None creates an empty option.
func None() Option {
	return Option{Ok: false}
}

// UnwrapOr returns a fallback when the option is empty.
func (o Option) UnwrapOr(fallback interface{}) interface{} {
	if o.Ok {
		return o.Value
	}
	return fallback
}

// Expect returns the contained value or panics with the supplied message.
func (o Option) Expect(message string) interface{} {
	if !o.Ok {
		panic(message)
	}
	return o.Value
}

// Map transforms the contained value if present.
func (o Option) Map(transform func(interface{}) interface{}) Option {
	if !o.Ok || transform == nil {
		return o
	}
	return Some(transform(o.Value))
}

// Result models a success or failure outcome.
type Result struct {
	Value interface{}
	Err   error
}

// Ok creates a successful result.
func Ok(value interface{}) Result {
	return Result{Value: value}
}

// ErrResult creates a failed result.
func ErrResult(err error) Result {
	return Result{Err: err}
}

// IsOk reports whether the result succeeded.
func (r Result) IsOk() bool {
	return r.Err == nil
}

// IsErr reports whether the result failed.
func (r Result) IsErr() bool {
	return r.Err != nil
}

// UnwrapOr returns the value or a fallback when failed.
func (r Result) UnwrapOr(fallback interface{}) interface{} {
	if r.Err != nil {
		return fallback
	}
	return r.Value
}

// Expect returns the value or panics with context.
func (r Result) Expect(message string) interface{} {
	if r.Err != nil {
		panic(fmt.Sprintf("%s: %v", message, r.Err))
	}
	return r.Value
}

// Map transforms the value when the result succeeded.
func (r Result) Map(transform func(interface{}) interface{}) Result {
	if r.Err != nil || transform == nil {
		return r
	}
	return Ok(transform(r.Value))
}


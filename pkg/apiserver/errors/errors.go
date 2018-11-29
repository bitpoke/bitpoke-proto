/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package errors

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StatusError is an error intended for consumption by apiserver
type StatusError struct {
	Code int
	Msg  string
}

var (
	// NotFound is 'not found' error for API Server
	NotFound = StatusError{
		Code: 404,
		Msg:  "not found",
	}
	// AlreadyExists is 'already exists' error for API Server
	AlreadyExists = StatusError{
		Code: 302,
		Msg:  "already exists",
	}
	// InternalError is an internal error for API Server
	InternalError = StatusError{
		Code: 500,
		Msg:  "internal error",
	}
)

// APIStatus is exposed by errors that can be converted to an api.Status object
// for finer grained details.
type APIStatus interface {
	Status() metav1.Status
}

var _ error = &StatusError{}

// Error implements the Error interface
func (e *StatusError) Error() string {
	return e.Msg
}

// NewApiserverError returns to user a meaningfull error
func NewApiserverError(err error) *StatusError {
	reason := reasonForError(err)
	switch reason {
	case metav1.StatusReasonNotFound:
		return &NotFound
	case metav1.StatusReasonAlreadyExists:
		return &AlreadyExists
	default:
		return &InternalError
	}
}

// reasonForError returns the HTTP status for a particular error.
func reasonForError(err error) metav1.StatusReason {
	switch t := err.(type) {
	case APIStatus:
		return t.Status().Reason
	}
	return metav1.StatusReasonUnknown
}

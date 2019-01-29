/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package status

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var k8sErrorMap = map[metav1.StatusReason]codes.Code{
	metav1.StatusReasonNotFound:      codes.NotFound,
	metav1.StatusReasonAlreadyExists: codes.AlreadyExists,
	metav1.StatusReasonForbidden:     codes.PermissionDenied,
}

// StatusError provides mapping and context between various statuses and grpc status
// nolint: golint
type StatusError struct {
	*status.Status
	errors []error
}

func (s *StatusError) Error() string {
	return s.Err().Error()
}

// GRPCStatus implements the gRPC status interface
func (s *StatusError) GRPCStatus() *status.Status {
	return s.Status
}

// Errors returns a list of underlying errors
func (s *StatusError) Errors() []error {
	return s.errors
}

// Because adds an errors to the cause for the status
func (s *StatusError) Because(err ...error) *StatusError {
	s.errors = append(s.errors, err...)
	return s
}

// FromError converts an error to a gRPC status error
func FromError(err error) *StatusError {
	st := &StatusError{
		Status: status.New(codes.Internal, "internal error"),
		errors: []error{err},
	}

	switch e := err.(type) {
	case *k8serrors.StatusError:
		if code, codeExists := k8sErrorMap[e.Status().Reason]; codeExists {
			st.Status = status.New(code, code.String())
		}
	}
	return st
}

// InvalidArgumentf returns a gRPC InvalidArgument status with a custom formatted message
func InvalidArgumentf(format string, a ...interface{}) *StatusError {
	return &StatusError{
		Status: status.New(codes.InvalidArgument, fmt.Sprintf(format, a...)),
	}
}

// InvalidArgument returns a gRPC InvalidArgument status
func InvalidArgument() *StatusError {
	return InvalidArgumentf(codes.InvalidArgument.String())
}

// Unauthenticatedf returns a gRPC Unauthenticated status with a custom formatted message
func Unauthenticatedf(format string, a ...interface{}) *StatusError {
	return &StatusError{
		Status: status.New(codes.Unauthenticated, fmt.Sprintf(format, a...)),
	}
}

// Unauthenticated returns a gRPC Unauthenticated status
func Unauthenticated() *StatusError {
	return Unauthenticatedf(codes.Unauthenticated.String())
}

// NotFoundf returns a gRPC NotFound status with a custom formatted message
func NotFoundf(format string, a ...interface{}) *StatusError {
	return &StatusError{
		Status: status.New(codes.NotFound, fmt.Sprintf(format, a...)),
	}
}

// NotFound returns a gRPC NotFound status
func NotFound() *StatusError {
	return NotFoundf(codes.NotFound.String())
}

// InternalErrorf returns a gRPC InternalError status with a custom formatted message
func InternalErrorf(format string, a ...interface{}) *StatusError {
	return &StatusError{
		Status: status.New(codes.Internal, fmt.Sprintf(format, a...)),
	}
}

// InternalError returns a gRPC InternalError status
func InternalError() *StatusError {
	return InternalErrorf(codes.Internal.String())
}

/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package status

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Error converts an k8s api error to a gRPC status error
func Error(err error) error {
	return FromError(err).Err()
}

// FromError converts an k8s api error to it's corresponding gRPC status
func FromError(err error) *status.Status {
	reason := reasonForError(err)
	switch reason {
	case metav1.StatusReasonNotFound:
		return status.New(codes.NotFound, "not found")
	case metav1.StatusReasonAlreadyExists:
		return status.New(codes.AlreadyExists, "already exists")
	default:
		return status.New(codes.Internal, "internal error")
	}
}

// reasonForError returns the HTTP status for a particular error.
func reasonForError(err error) metav1.StatusReason {
	switch t := err.(type) {
	case k8serrors.APIStatus:
		return t.Status().Reason
	}
	return metav1.StatusReasonUnknown
}

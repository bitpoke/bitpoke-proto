/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package metadata

import (
	"context"
)

var (
	// FakeSubject is the subject field from Claims for Fake Auth
	FakeSubject = "subject"
	// FakeEmail is the email field from Claims for Fake Auth
	FakeEmail = "email"
	// FakeVerified is the verified field from Claims for Fake Auth
	FakeVerified = true
)

// FakeAuth puts fake claims in context
func FakeAuth(ctx context.Context) (context.Context, error) {
	var claims Claims
	claims.Subject = FakeSubject
	claims.Email = FakeEmail
	claims.Verified = FakeVerified

	newCtx := context.WithValue(ctx, AuthTokenContextKey, claims)

	return newCtx, nil
}

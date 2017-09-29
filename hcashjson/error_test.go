// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2015-2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package hcashjson_test

import (
	"testing"

	"github.com/HcashOrg/hcashd/hcashjson"
)

// TestErrorCodeStringer tests the stringized output for the ErrorCode type.
func TestErrorCodeStringer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   hcashjson.ErrorCode
		want string
	}{
		{hcashjson.ErrDuplicateMethod, "ErrDuplicateMethod"},
		{hcashjson.ErrInvalidUsageFlags, "ErrInvalidUsageFlags"},
		{hcashjson.ErrInvalidType, "ErrInvalidType"},
		{hcashjson.ErrEmbeddedType, "ErrEmbeddedType"},
		{hcashjson.ErrUnexportedField, "ErrUnexportedField"},
		{hcashjson.ErrUnsupportedFieldType, "ErrUnsupportedFieldType"},
		{hcashjson.ErrNonOptionalField, "ErrNonOptionalField"},
		{hcashjson.ErrNonOptionalDefault, "ErrNonOptionalDefault"},
		{hcashjson.ErrMismatchedDefault, "ErrMismatchedDefault"},
		{hcashjson.ErrUnregisteredMethod, "ErrUnregisteredMethod"},
		{hcashjson.ErrNumParams, "ErrNumParams"},
		{hcashjson.ErrMissingDescription, "ErrMissingDescription"},
		{0xffff, "Unknown ErrorCode (65535)"},
	}

	// Detect additional error codes that don't have the stringer added.
	if len(tests)-1 != int(hcashjson.TstNumErrorCodes) {
		t.Errorf("It appears an error code was added without adding an " +
			"associated stringer test")
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		result := test.in.String()
		if result != test.want {
			t.Errorf("String #%d\n got: %s want: %s", i, result,
				test.want)
			continue
		}
	}
}

// TestError tests the error output for the Error type.
func TestError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   hcashjson.Error
		want string
	}{
		{
			hcashjson.Error{Message: "some error"},
			"some error",
		},
		{
			hcashjson.Error{Message: "human-readable error"},
			"human-readable error",
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		result := test.in.Error()
		if result != test.want {
			t.Errorf("Error #%d\n got: %s want: %s", i, result,
				test.want)
			continue
		}
	}
}

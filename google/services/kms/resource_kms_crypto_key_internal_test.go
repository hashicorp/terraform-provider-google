// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package kms

import (
	"testing"
	"time"
)

func TestCryptoKeyNextRotationCalculation(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	period, _ := time.ParseDuration("1000000s")

	expected := now.Add(period).Format(time.RFC3339Nano)

	timestamp, err := kmsCryptoKeyNextRotation(now, "1000000s")

	if err != nil {
		t.Fatalf("unexpected failure parsing time %s and duration 1000s: %s", now, err.Error())
	}

	if expected != timestamp {
		t.Fatalf("expected %s to equal %s", timestamp, expected)
	}
}

func TestCryptoKeyNextRotationCalculation_validation(t *testing.T) {
	t.Parallel()

	_, errs := validateKmsCryptoKeyRotationPeriod("86399s", "rotation_period")

	if len(errs) == 0 {
		t.Fatalf("Periods of less than a day should be invalid")
	}

	_, errs = validateKmsCryptoKeyRotationPeriod("100000.0000000001s", "rotation_period")

	if len(errs) == 0 {
		t.Fatalf("Numbers with more than 9 fractional digits are invalid")
	}
}

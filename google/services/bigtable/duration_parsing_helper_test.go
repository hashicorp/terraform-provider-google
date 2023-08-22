// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigtable_test

import (
	"math"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/services/bigtable"
)

var parseDurationTests = []struct {
	in   string
	want time.Duration
}{
	// simple
	{"0", 0},
	{"5s", 5 * time.Second},
	{"30s", 30 * time.Second},
	{"1478s", 1478 * time.Second},
	// sign
	{"-5s", -5 * time.Second},
	{"+5s", 5 * time.Second},
	{"-0", 0},
	{"+0", 0},
	// decimal
	{"5.0s", 5 * time.Second},
	{"5.6s", 5*time.Second + 600*time.Millisecond},
	{"5.s", 5 * time.Second},
	{".5s", 500 * time.Millisecond},
	{"1.0s", 1 * time.Second},
	{"1.00s", 1 * time.Second},
	{"1.004s", 1*time.Second + 4*time.Millisecond},
	{"1.0040s", 1*time.Second + 4*time.Millisecond},
	{"100.00100s", 100*time.Second + 1*time.Millisecond},
	// different units
	{"10ns", 10 * time.Nanosecond},
	{"11us", 11 * time.Microsecond},
	{"12µs", 12 * time.Microsecond}, // U+00B5
	{"12μs", 12 * time.Microsecond}, // U+03BC
	{"13ms", 13 * time.Millisecond},
	{"14s", 14 * time.Second},
	{"15m", 15 * time.Minute},
	{"16h", 16 * time.Hour},
	{"5d", 5 * 24 * time.Hour},
	// composite durations
	{"3h30m", 3*time.Hour + 30*time.Minute},
	{"10.5s4m", 4*time.Minute + 10*time.Second + 500*time.Millisecond},
	{"-2m3.4s", -(2*time.Minute + 3*time.Second + 400*time.Millisecond)},
	{"1h2m3s4ms5us6ns", 1*time.Hour + 2*time.Minute + 3*time.Second + 4*time.Millisecond + 5*time.Microsecond + 6*time.Nanosecond},
	{"39h9m14.425s", 39*time.Hour + 9*time.Minute + 14*time.Second + 425*time.Millisecond},
	// large value
	{"52763797000ns", 52763797000 * time.Nanosecond},
	// more than 9 digits after decimal point, see https://golang.org/issue/6617
	{"0.3333333333333333333h", 20 * time.Minute},
	// 9007199254740993 = 1<<53+1 cannot be stored precisely in a float64
	{"9007199254740993ns", (1<<53 + 1) * time.Nanosecond},
	// largest duration that can be represented by int64 in nanoseconds
	{"9223372036854775807ns", (1<<63 - 1) * time.Nanosecond},
	{"9223372036854775.807us", (1<<63 - 1) * time.Nanosecond},
	{"9223372036s854ms775us807ns", (1<<63 - 1) * time.Nanosecond},
	{"-9223372036854775808ns", -1 << 63 * time.Nanosecond},
	{"-9223372036854775.808us", -1 << 63 * time.Nanosecond},
	{"-9223372036s854ms775us808ns", -1 << 63 * time.Nanosecond},
	// largest negative value
	{"-9223372036854775808ns", -1 << 63 * time.Nanosecond},
	// largest negative round trip value, see https://golang.org/issue/48629
	{"-2562047h47m16.854775808s", -1 << 63 * time.Nanosecond},
	// huge string; issue 15011.
	{"0.100000000000000000000h", 6 * time.Minute},
	// This value tests the first overflow check in leadingFraction.
	{"0.830103483285477580700h", 49*time.Minute + 48*time.Second + 372539827*time.Nanosecond},
}

func TestParseDuration(t *testing.T) {
	for _, tc := range parseDurationTests {
		d, err := bigtable.ParseDuration(tc.in)
		if err != nil || d != tc.want {
			t.Errorf("bigtable.ParseDuration(%q) = %v, %v, want %v, nil", tc.in, d, err, tc.want)
		}
	}
}

var parseDurationErrorTests = []struct {
	in     string
	expect string
}{
	// invalid
	{"", `""`},
	{"3", `"3"`},
	{"-", `"-"`},
	{"s", `"s"`},
	{".", `"."`},
	{"-.", `"-."`},
	{".s", `".s"`},
	{"+.s", `"+.s"`},
	{"\x85\x85", `"\x85\x85"`},
	{"\xffff", `"\xffff"`},
	{"hello \xffff world", `"hello \xffff world"`},
	{"\uFFFD", `"\xef\xbf\xbd"`},                                             // utf8.RuneError
	{"\uFFFD hello \uFFFD world", `"\xef\xbf\xbd hello \xef\xbf\xbd world"`}, // utf8.RuneError
	// overflow
	{"9223372036854775810ns", `"9223372036854775810ns"`},
	{"9223372036854775808ns", `"9223372036854775808ns"`},
	{"-9223372036854775809ns", `"-9223372036854775809ns"`},
	{"9223372036854776us", `"9223372036854776us"`},
	{"3000000h", `"3000000h"`},
	{"9223372036854775.808us", `"9223372036854775.808us"`},
	{"9223372036854ms775us808ns", `"9223372036854ms775us808ns"`},
}

func TestParseDurationErrors(t *testing.T) {
	for _, tc := range parseDurationErrorTests {
		_, err := bigtable.ParseDuration(tc.in)
		if err == nil {
			t.Errorf("bigtable.ParseDuration(%q) = _, nil, want _, non-nil", tc.in)
		} else if !strings.Contains(err.Error(), tc.expect) {
			t.Errorf("bigtable.ParseDuration(%q) = _, %q, error does not contain %q", tc.in, err, tc.expect)
		}
	}
}

func TestParseDurationRoundTrip(t *testing.T) {
	// https://golang.org/issue/48629
	max0 := time.Duration(math.MaxInt64)
	max1, err := bigtable.ParseDuration(max0.String())
	if err != nil || max0 != max1 {
		t.Errorf("round-trip failed: %d => %q => %d, %v", max0, max0.String(), max1, err)
	}

	min0 := time.Duration(math.MinInt64)
	min1, err := bigtable.ParseDuration(min0.String())
	if err != nil || min0 != min1 {
		t.Errorf("round-trip failed: %d => %q => %d, %v", min0, min0.String(), min1, err)
	}

	for i := 0; i < 100; i++ {
		// Resolutions finer than milliseconds will result in
		// imprecise round-trips.
		d0 := time.Duration(rand.Int31()) * time.Millisecond
		s := d0.String()
		d1, err := bigtable.ParseDuration(s)
		if err != nil || d0 != d1 {
			t.Errorf("round-trip failed: %d => %q => %d, %v", d0, s, d1, err)
		}
	}
}

func BenchmarkParseDuration(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bigtable.ParseDuration("9007199254.740993ms")
		bigtable.ParseDuration("9007199254740993ns")
	}
}

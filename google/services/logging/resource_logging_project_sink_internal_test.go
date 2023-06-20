// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestLoggingProjectSink_bigqueryOptionCustomizedDiff(t *testing.T) {
	t.Parallel()

	type LoggingProjectSink struct {
		BigqueryOptions      int
		UniqueWriterIdentity bool
	}
	cases := map[string]struct {
		ExpectedError bool
		After         LoggingProjectSink
	}{
		"no biquery options with false unique writer identity": {
			ExpectedError: false,
			After: LoggingProjectSink{
				BigqueryOptions:      0,
				UniqueWriterIdentity: false,
			},
		},
		"no biquery options with true unique writer identity": {
			ExpectedError: false,
			After: LoggingProjectSink{
				BigqueryOptions:      0,
				UniqueWriterIdentity: true,
			},
		},
		"biquery options with false unique writer identity": {
			ExpectedError: true,
			After: LoggingProjectSink{
				BigqueryOptions:      1,
				UniqueWriterIdentity: false,
			},
		},
		"biquery options with true unique writer identity": {
			ExpectedError: false,
			After: LoggingProjectSink{
				BigqueryOptions:      1,
				UniqueWriterIdentity: true,
			},
		},
	}

	for tn, tc := range cases {
		d := &tpgresource.ResourceDiffMock{
			After: map[string]interface{}{
				"bigquery_options.#":     tc.After.BigqueryOptions,
				"unique_writer_identity": tc.After.UniqueWriterIdentity,
			},
		}
		err := resourceLoggingProjectSinkCustomizeDiffFunc(d)
		hasError := err != nil
		if tc.ExpectedError != hasError {
			t.Errorf("%v: expected has error %v, but was %v", tn, tc.ExpectedError, hasError)
		}
	}
}

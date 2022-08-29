package google

import (
	"fmt"
	"strings"
	"testing"

	"google.golang.org/api/compute/v1"
)

const (
	topLevelMsg             = "Top-level message."
	localizedMsgTmpl        = "LocalizedMessage%d message"
	helpLinkDescriptionTmpl = "Help%dLink%d Description"
	helpLinkUrlTmpl         = "https://help%d.com/link%d"
)

var locales = []string{"en-US", "es-US", "es-ES", "es-MX", "de-DE"}

func buildOperationError(numLocalizedMsg int, numHelpWithLinks []int) compute.OperationError {
	opError := &compute.OperationErrorErrors{Message: topLevelMsg}
	opErrorErrors := []*compute.OperationErrorErrors{opError}

	for n := 1; n <= numLocalizedMsg; n++ {
		opError.ErrorDetails = append(opError.ErrorDetails,
			&compute.OperationErrorErrorsErrorDetails{
				LocalizedMessage: &compute.LocalizedMessage{
					Locale:  locales[n-1%len(locales)],
					Message: formatLocalizedMsg(n),
				},
			})
	}

	for i := 0; i < len(numHelpWithLinks); i++ {
		errorDetail := &compute.OperationErrorErrorsErrorDetails{
			Help: &compute.Help{},
		}

		for nLinks := 1; nLinks <= numHelpWithLinks[i]; nLinks++ {
			desc, url := formatLink(i+1, nLinks)
			errorDetail.Help.Links = append(errorDetail.Help.Links, &compute.HelpLink{
				Description: desc,
				Url:         url,
			})
		}

		opError.ErrorDetails = append(opError.ErrorDetails, errorDetail)
	}

	return compute.OperationError{Errors: opErrorErrors}

}

func omitAlways(numLocalizedMsg int, numHelpWithLinks []int) []string {
	var omits []string

	for n := 2; n <= numLocalizedMsg; n++ {
		omits = append(omits, fmt.Sprintf("LocalizedMessage%d", n))
	}

	for i := 0; i < len(numHelpWithLinks); i++ {
		for j := maxLinks(i); j < numHelpWithLinks[i]; j++ {
			desc, url := formatLink(i+1, j+1)
			omits = append(omits, desc, url)
		}
	}

	return omits

}

func maxLinks(helpIndex int) int {
	if helpIndex == 0 {
		return 1
	}

	return 0
}

func formatLocalizedMsg(localizedMsgNum int) string {
	return fmt.Sprintf(localizedMsgTmpl, localizedMsgNum)
}

func formatLink(helpNum, linkNum int) (string, string) {
	return fmt.Sprintf(helpLinkDescriptionTmpl, helpNum, linkNum), fmt.Sprintf(helpLinkUrlTmpl, helpNum, linkNum)
}

func TestComputeOperationError_Error(t *testing.T) {
	testCases := []struct {
		name           string
		input          compute.OperationError
		expectContains []string
		expectOmits    []string
	}{
		{
			name:  "MessageOnly",
			input: buildOperationError(0, []int{}),
			expectContains: []string{
				"Top-level",
			},
			expectOmits: append(omitAlways(0, []int{}), []string{
				"LocalizedMessage1",
				"Help1Link1 Description",
				"https://help1.com/link1",
			}...),
		},
		{
			name:  "WithLocalizedMessageAndNoHelp",
			input: buildOperationError(1, []int{}),
			expectContains: []string{
				"Top-level",
				"LocalizedMessage1",
			},
			expectOmits: append(omitAlways(1, []int{}), []string{
				"Help1Link1 Description",
				"https://help1.com/link1",
			}...),
		},
		{
			name:  "WithLocalizedMessageAndHelp",
			input: buildOperationError(1, []int{1}),
			expectContains: []string{
				"Top-level",
				"LocalizedMessage1",
				"Help1Link1 Description",
				"https://help1.com/link1",
			},
			expectOmits: append(omitAlways(1, []int{1}), []string{}...),
		},
		{
			name:  "WithNoLocalizedMessageAndHelp",
			input: buildOperationError(0, []int{1}),
			expectContains: []string{
				"Top-level",
				"Help1Link1 Description",
				"https://help1.com/link1",
			},
			expectOmits: append(omitAlways(0, []int{1}), []string{
				"LocalizedMessage1",
			}...),
		},
		{
			name:  "WithLocalizedMessageAndHelpWithTwoLinks",
			input: buildOperationError(1, []int{2}),
			expectContains: []string{
				"Top-level",
				"LocalizedMessage1",
				"Help1Link1 Description",
				"https://help1.com/link1",
			},
			expectOmits: append(omitAlways(1, []int{2}), []string{}...),
		},
		// The case below should never happen because the server should just send multiple links
		// but the protobuf defition would allow it, so testing anyway.
		{
			name:  "WithLocalizedMessageAndTwoHelpsWithTwoLinks",
			input: buildOperationError(1, []int{2, 2}),
			expectContains: []string{
				"Top-level",
				"LocalizedMessage1",
				"Help1Link1 Description",
				"https://help1.com/link1",
			},
			expectOmits: append(omitAlways(1, []int{2, 2}), []string{}...),
		},
		// This should never happen because the server should never respond with the messages for
		// two locales at once, but should rather take the locale as input to the API and serve
		// the appropriate message for that locale. However, the protobuf defition would allow it,
		// so we'll test for it. The second message in the list would be ignored.
		{
			name:  "WithTwoLocalizedMessageAndHelp",
			input: buildOperationError(2, []int{1}),
			expectContains: []string{
				"Top-level",
				"LocalizedMessage1",
				"Help1Link1 Description",
				"https://help1.com/link1",
			},
			expectOmits: append(omitAlways(2, []int{1}), []string{}...),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ComputeOperationError(tc.input)
			str := err.Error()

			for _, contains := range tc.expectContains {
				if !strings.Contains(str, contains) {
					t.Errorf("expected\n%s\nto contain, %q, and did not", str, contains)
				}
			}

			for _, omits := range tc.expectOmits {
				if strings.Contains(str, omits) {
					t.Errorf("expected\n%s\nnot to contain, %q, and did not", str, omits)
				}
			}
		})
	}
}

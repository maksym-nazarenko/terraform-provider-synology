package api

import (
	"fmt"
	"strings"
)

type ErrorDescriber interface {
	// ErrorSummaries returns information about all known errors.
	ErrorSummaries() []ErrorSummary
}

// SynologyError defines a structure for error object returned by Synology API.
// It is a high-level error for a particular API family.
type SynologyError struct {
	Code    int
	Summary string
	// Errors is a collection of detailed errors for a concrete API request.
	Errors []ErrorItem
}

// ErrorItem defines detailed request error.
type ErrorItem struct {
	Code    int
	Summary string
	Details ErrorFields
}

type ErrorSummary map[int]string

// ErrorFields defines extra fields for particular detailed error.
type ErrorFields map[string]interface{}

var GlobalErrors ErrorSummary = ErrorSummary{
	100: "Unknown error",
	101: "No parameter of API, method or version",
	102: "The requested API does not exist",
	103: "The requested method does not exist",
	104: "The requested version does not support the functionality",
	105: "The logged in session does not have permission",
	106: "Session timeout",
	107: "Session interrupted by duplicate login",
	119: "SID not found",
}

func (se SynologyError) Error() string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("[%d] %s", se.Code, se.Summary))
	if len(se.Errors) > 0 {
		buf.WriteString("\n\tDetails:")
	}

	for _, e := range se.Errors {
		detailedFields := []string{}
		buf.WriteString(fmt.Sprintf("\n\t\t[%d] %s", e.Code, e.Summary))
		if len(e.Details) > 0 {
			for k, v := range e.Details {
				detailedFields = append(detailedFields, k+": "+fmt.Sprintf("%v", v))
			}
			buf.WriteString(": [" + strings.Join(detailedFields, ",") + "]")
		}
	}

	return buf.String()
}

// DescribeError translates error code to human-readable summary text.
// It accepts error code and number of summary maps to look in.
// First summary with this code wins.
func DescribeError(code int, summaries ...ErrorSummary) string {
	for _, summaryMap := range summaries {
		if summary, ok := summaryMap[code]; ok {
			return summary
		}
	}

	return "Unknown error code"
}

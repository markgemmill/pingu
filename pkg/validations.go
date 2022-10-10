package pkg

import (
	"fmt"
	"regexp"
)

type Assertion interface {
	Assert(test *UrlResult) (bool, string)
	Name() string
}

type StatusCodeAssertion struct {
	Expected int
}

func NewStatusCodeAssertion(expected int) *StatusCodeAssertion {
	return &StatusCodeAssertion{
		Expected: expected,
	}
}

func (a *StatusCodeAssertion) Name() string {
	return "Status Code Assertion"
}

func (a *StatusCodeAssertion) Assert(test *UrlResult) (bool, string) {
	pass := a.Expected == test.StatusCode
	if !pass {
		return pass, fmt.Sprintf("expecting status of %d, but received %d", a.Expected, test.StatusCode)
	}
	return pass, ""
}

type ContentAssertion struct {
	Regex string
}

func NewContentAssertion(regex string) *ContentAssertion {
	return &ContentAssertion{
		Regex: regex,
	}
}

func (a *ContentAssertion) Name() string {
	return "Content Assertion"
}

func (a *ContentAssertion) Assert(test *UrlResult) (bool, string) {
	re := regexp.MustCompile(a.Regex)
	result := re.FindStringIndex(test.Content)
	if result == nil {
		return false, "does not contain the expected text"
	}
	return true, ""
}

func BuildAssertions(expectedStatus int, expectedContent string) *[]*Assertion {
	assertions := []*Assertion{}

	var si Assertion
	si = NewStatusCodeAssertion(expectedStatus)
	assertions = append(assertions, &si)

	if expectedContent != "" {
		var i Assertion
		i = NewContentAssertion(expectedContent)
		assertions = append(assertions, &i)
	}

	return &assertions
}

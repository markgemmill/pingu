package pkg

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type UrlResult struct {
	StatusCode int
	Content    string
	Fail       bool
}

func UrlFetch(url string) UrlResult {
	result := UrlResult{Fail: false}

	resp, err := http.Get(url)
	if err != nil {
		result.Fail = true
		return result
	}

	result.StatusCode = resp.StatusCode

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	result.Content = string(body)
	resp.Body.Close()
	return result
}

type UrlCheck struct {
	Url        string
	Assertions []*Assertion
	Pass       bool
	Errors     []string
}

func NewUrlCheck(url string, assertions []*Assertion) *UrlCheck {
	return &UrlCheck{
		Url:        url,
		Assertions: assertions,
		Pass:       false,
	}
}

func PassFail(pass bool) string {
	if pass {
		return Green(PASS)
	}
	return Red(FAIL)
}

func ErrMsg(errMsg string) string {
	if errMsg != "" {
		return fmt.Sprintf(": %s", errMsg)
	}
	return ""
}

func (u *UrlCheck) Test() {
	console.Trace("Fetching url: %s\n", u.Url)
	result := UrlFetch(u.Url)

	console.Indent()

	if result.Fail {
		console.Trace("Failed to fech url!")
		u.Pass = false
		u.Errors = append(u.Errors, "Could not fetch url.")
		return
	}

	for _, assertion := range u.Assertions {
		assert := *assertion
		//console.Trace("Running assertion...\n")
		passed, errMsg := assert.Assert(&result)
		console.Trace("%s %s%s\n", assert.Name(), PassFail(passed), ErrMsg(errMsg))

		u.Pass = passed
		if passed == false {
			console.Print("GET %s %s.\n", u.Url, errMsg)
			u.Errors = append(u.Errors, errMsg)
			break
		}
	}

	console.Dedent()
}

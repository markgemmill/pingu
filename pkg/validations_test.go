package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStatusCodeAssertionPass(t *testing.T) {

	urlResult := UrlResult{
		StatusCode: 200,
		Content:    "active",
		Fail:       false,
	}

	test := NewStatusCodeAssertion(200)
	pass, errMsg := test.Assert(&urlResult)

	assert.True(t, pass)
	assert.Equal(t, errMsg, "")
}

func TestStatusCodeAssertionFail(t *testing.T) {

	urlResult := UrlResult{
		StatusCode: 200,
		Content:    "active",
		Fail:       false,
	}

	test := NewStatusCodeAssertion(401)
	pass, errMsg := test.Assert(&urlResult)

	assert.False(t, pass)
	assert.Equal(t, errMsg, "expecting status of 401, but received 200")
}

func TestContentAssertionPass(t *testing.T) {

	urlResult := UrlResult{
		StatusCode: 200,
		Content:    "active",
		Fail:       false,
	}

	test := NewContentAssertion("active")
	pass, errMsg := test.Assert(&urlResult)

	assert.True(t, pass)
	assert.Equal(t, errMsg, "")
}

func TestContentAssertionFail(t *testing.T) {

	urlResult := UrlResult{
		StatusCode: 200,
		Content:    "inactive",
		Fail:       false,
	}

	test := NewContentAssertion("^active$")
	pass, errMsg := test.Assert(&urlResult)

	assert.False(t, pass)
	assert.Equal(t, errMsg, "does not contain the expected text")
}

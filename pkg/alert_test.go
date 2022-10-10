package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseEmailAddresses(t *testing.T) {
	assert.Equal(t, []string{"mgemmill@mail.com"}, ParseEmailAddresses("mgemmill@mail.com"))
	assert.Equal(t, []string{"mgemmill@mail.com", "schen@mailing.com"}, ParseEmailAddresses("mgemmill@mail.com;schen@mailing.com"))
	assert.Equal(t, []string{}, ParseEmailAddresses(""))
}

package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSluggifyUrl(t *testing.T) {
	assert.Equal(t, "https_markgemmill_com", sluggifyUrl("https://markgemmill.com"))
}

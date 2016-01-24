package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryParserError(t *testing.T) {
	q := queryParseError{
		StatusCode: 400,
		Message:    "Test Message",
	}

	assert.Equal(t, 400, q.StatusCode, "StatusCode Should be 400")
	assert.Equal(t, "Test Message", q.Error(), "Should be 'Test Message'")

}

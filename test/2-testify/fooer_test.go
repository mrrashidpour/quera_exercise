package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFooer(t *testing.T) {
	result := Fooer(3)
	assert.Equal(t, "Foo", result, "3 باید Foo برگرداند")

	result2 := Fooer(4)
	assert.Equal(t, "4", result2, "4 باید '4' برگرداند")
}

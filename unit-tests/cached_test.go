package unit_tests

import (
	"testing"

	"github.com/atrico-go/templates/unit-tests/generated"
	. "github.com/atrico-go/testing/assert"
	"github.com/atrico-go/testing/is"
)

func Test_Cached_GetValue(t *testing.T) {
	// Arrange
	val := randGen.Int()
	count := 0
	cache := generated.CachedIntValue{Creator: mockCreator(val, &count)}

	// Act
	result1 := cache.GetValue()
	result2 := cache.GetValue()
	result3 := cache.GetValue()

	// Assert
	Assert(t).That(result1, is.EqualTo(val), "result1")
	Assert(t).That(result2, is.EqualTo(val), "result2")
	Assert(t).That(result3, is.EqualTo(val), "result3")
	Assert(t).That(count, is.EqualTo(1), "called once")
}

func Test_Cached_Reset(t *testing.T) {
	// Arrange
	val := randGen.Int()
	count := 0
	cache := generated.CachedIntValue{Creator: mockCreator(val, &count)}

	// Act
	result1 := cache.GetValue()
	cache.Reset()
	result2 := cache.GetValue()
	cache.Reset()
	result3 := cache.GetValue()

	// Assert
	Assert(t).That(result1, is.EqualTo(val), "result1")
	Assert(t).That(result2, is.EqualTo(val), "result2")
	Assert(t).That(result3, is.EqualTo(val), "result3")
	Assert(t).That(count, is.EqualTo(3), "called 3x")
}

func mockCreator(value int, called *int) func() int {
	return func() int {
		*called++
		return value
	}
}

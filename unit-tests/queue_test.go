package unit_tests

import (
	"fmt"
	"testing"

	"github.com/atrico-go/templates/unit-tests/generated"
	. "github.com/atrico-go/testing/assert"
	"github.com/atrico-go/testing/is"
	"github.com/atrico-go/testing/random"
)

var randGen = random.NewValueGenerator()
var emptyPanic = "queue is empty"

func Test_Queue_Empty(t *testing.T) {
	// Arrange

	// Act
	queue := generated.NewIntQueue()
	isEmpty := queue.IsEmpty()
	count := queue.Count()
	peek := PanicCatcher(func() { queue.Peek() })
	pop := PanicCatcher(func() { queue.Pop() })

	// Assert
	Assert(t).That(isEmpty, is.True, "IsEmpty")
	Assert(t).That(count, is.EqualTo(0), "Count == 0")
	Assert(t).That(peek, is.EqualTo(emptyPanic), "peek panic")
	Assert(t).That(pop, is.EqualTo(emptyPanic), "pop panic")
}

func Test_Queue_InitialValues(t *testing.T) {
	// Arrange
	val1, val2 := randGen.Int(), randGen.Int()
	// Act
	queue := generated.NewIntQueue(val1, val2)
	fmt.Printf("Create with: %v,%v\n", val1, val2)
	count0 := queue.Count()
	peek1 := queue.Peek()
	pop1 := queue.Pop()
	count1 := queue.Count()
	peek2 := queue.Peek()
	pop2 := queue.Pop()
	count2 := queue.Count()

	// Assert
	Assert(t).That(count0, is.EqualTo(2), "Initial Count")
	Assert(t).That(peek1, is.EqualTo(val1), "Peek 1")
	Assert(t).That(pop1, is.EqualTo(val1), "Pop 1")
	Assert(t).That(count1, is.EqualTo(1), "Count after pop")
	Assert(t).That(peek2, is.EqualTo(val2), "Peek 2")
	Assert(t).That(pop2, is.EqualTo(val2), "Pop 2")
	Assert(t).That(count2, is.EqualTo(0), "Count after 2x pop")
}

func Test_Queue_PushAndPop(t *testing.T) {
	// Arrange
	val1, val2 := randGen.Int(), randGen.Int()

	// Act
	queue := generated.NewIntQueue()
	count0 := queue.Count()
	queue.Push(val1)
	fmt.Printf("Push: %v\n", val1)
	count1 := queue.Count()
	queue.Push(val2)
	fmt.Printf("Push: %v\n", val2)
	count2 := queue.Count()
	pop1 := queue.Pop()
	pop2 := queue.Pop()
	count3 := queue.Count()

	// Assert
	Assert(t).That(count0, is.EqualTo(0), "Initial Count")
	Assert(t).That(count1, is.EqualTo(1), "Count after push")
	Assert(t).That(count2, is.EqualTo(2), "Count after 2x push")
	Assert(t).That(pop1, is.EqualTo(val1), "Pop 1")
	Assert(t).That(pop2, is.EqualTo(val2), "Pop 2")
	Assert(t).That(count3, is.EqualTo(0), "Count after 2x pop")
}

func Test_Queue_Sort(t *testing.T) {
	// Arrange
	queue := generated.NewIntQueue(3, 1, 2)

	// Act
	queue.Sort(func(i, j int) bool { return i < j })
	pop1 := queue.Pop()
	pop2 := queue.Pop()
	pop3 := queue.Pop()

	// Assert
	Assert(t).That(pop1, is.EqualTo(1), "Pop 1")
	Assert(t).That(pop2, is.EqualTo(2), "Pop 2")
	Assert(t).That(pop3, is.EqualTo(3), "Pop 3")
}

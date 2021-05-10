package unit_tests

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/atrico-go/templates/unit-tests/generated"
	. "github.com/atrico-go/testing/assert"
	"github.com/atrico-go/testing/is"
)

var timeout = 250 * time.Millisecond

func Test_MtQueue_Empty(t *testing.T) {
	// Arrange

	// Act
	queue := generated.NewIntMtQueue()
	isEmpty := queue.IsEmpty()
	count := queue.Count()
	_,peek := queue.Peek(timeout)
	_,pop := queue.Pop(timeout)

	// Assert
	Assert(t).That(isEmpty, is.True, "IsEmpty")
	Assert(t).That(count, is.EqualTo(0), "Count == 0")
	Assert(t).That(peek, is.False, "peek timeout")
	Assert(t).That(pop, is.False, "pop timeout")
}

func Test_MtQueue_InitialValues(t *testing.T) {
	// Arrange
	val1, val2 := randGen.Int(), randGen.Int()
	// Act
	queue := generated.NewIntMtQueue(val1, val2)
	fmt.Printf("Create with: %v,%v\n", val1, val2)
	count0 := queue.Count()
	peek1,peekOk1 := queue.Peek(timeout)
	pop1,popOk1 := queue.Pop(timeout)
	count1 := queue.Count()
	peek2,peekOk2 := queue.Peek(timeout)
	pop2,popOk2 := queue.Pop(timeout)
	count2 := queue.Count()

	// Assert
	Assert(t).That(count0, is.EqualTo(2), "Initial Count")
	Assert(t).That(peek1, is.EqualTo(val1), "Peek 1")
	Assert(t).That(peekOk1, is.True, "Peek 1 OK")
	Assert(t).That(pop1, is.EqualTo(val1), "Pop 1")
	Assert(t).That(popOk1, is.True, "Pop 1 OK")
	Assert(t).That(count1, is.EqualTo(1), "Count after pop")
	Assert(t).That(peek2, is.EqualTo(val2), "Peek 2")
	Assert(t).That(peekOk2, is.True, "Peek 2 OK")
	Assert(t).That(pop2, is.EqualTo(val2), "Pop 2")
	Assert(t).That(popOk2, is.True, "Pop 2 OK")
	Assert(t).That(count2, is.EqualTo(0), "Count after 2x pop")
}

func Test_MtQueue_PushAndPop(t *testing.T) {
	// Arrange
	val1, val2 := randGen.Int(), randGen.Int()

	// Act
	queue := generated.NewIntMtQueue()
	count0 := queue.Count()
	queue.Push(val1)
	fmt.Printf("Push: %v\n", val1)
	count1 := queue.Count()
	queue.Push(val2)
	fmt.Printf("Push: %v\n", val2)
	count2 := queue.Count()
	pop1,_ := queue.Pop(timeout)
	pop2,_ := queue.Pop(timeout)
	count3 := queue.Count()

	// Assert
	Assert(t).That(count0, is.EqualTo(0), "Initial Count")
	Assert(t).That(count1, is.EqualTo(1), "Count after push")
	Assert(t).That(count2, is.EqualTo(2), "Count after 2x push")
	Assert(t).That(pop1, is.EqualTo(val1), "Pop 1")
	Assert(t).That(pop2, is.EqualTo(val2), "Pop 2")
	Assert(t).That(count3, is.EqualTo(0), "Count after 2x pop")
}

func Test_MtQueue_Sort(t *testing.T) {
	// Arrange
	queue := generated.NewIntMtQueue(3, 1, 2)

	// Act
	queue.Sort(func(i, j int) bool { return i < j })
	pop1,_ := queue.Pop(timeout)
	pop2,_ := queue.Pop(timeout)
	pop3,_ := queue.Pop(timeout)

	// Assert
	Assert(t).That(pop1, is.EqualTo(1), "Pop 1")
	Assert(t).That(pop2, is.EqualTo(2), "Pop 2")
	Assert(t).That(pop3, is.EqualTo(3), "Pop 3")
}

func Test_MtQueue_EmptyEvent(t *testing.T) {
	// Arrange
	queue := generated.NewIntMtQueue(randGen.Int())

	// Act
	empty1 := queue.WaitUntilEmpty(timeout)
	queue.Pop(timeout)
	empty2:= queue.WaitUntilEmpty(timeout)

	// Assert
	Assert(t).That(empty1, is.False, "Empty 1")
	Assert(t).That(empty2, is.True, "Empty 2")
}

func Test_MtQueue_PeekWithDelayedPush(t *testing.T) {
	// Arrange
	val := randGen.Int()
	queue := generated.NewIntMtQueue()
	var peek int
	var peekOk bool
	wg := sync.WaitGroup{}

	// Act
	wg.Add(1)
	go func() {
		peek, peekOk = queue.Peek(time.Minute)
		wg.Done()
	}()
	queue.Push(val)
	wg.Wait()

	// Assert
	Assert(t).That(peek, is.EqualTo(val), "Pop")
	Assert(t).That(peekOk, is.True, "Pop OK")
}

func Test_MtQueue_PopWithDelayedPush(t *testing.T) {
	// Arrange
	val := randGen.Int()
	queue := generated.NewIntMtQueue()
	var pop int
	var popOk bool
	wg := sync.WaitGroup{}

	// Act
	wg.Add(1)
	go func() {
		pop,popOk = queue.Pop(time.Minute)
		wg.Done()
	}()
	queue.Push(val)
	wg.Wait()

	// Assert
	Assert(t).That(pop, is.EqualTo(val), "Pop")
	Assert(t).That(popOk, is.True, "Pop OK")
}
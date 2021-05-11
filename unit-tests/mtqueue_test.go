package unit_tests

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/atrico-go/templates/unit-tests/generated"
	. "github.com/atrico-go/testing/assert"
	"github.com/atrico-go/testing/is"
)

func Test_MtQueue_Empty(t *testing.T) {
	// Arrange

	// Act
	queue := generated.IntMtQueue{}
	isEmpty := queue.IsEmpty()
	count := queue.Count()
	_,peek := queue.Peek(createContext())
	_,pop := queue.Pop(createContext())

	// Assert
	Assert(t).That(isEmpty, is.True, "IsEmpty")
	Assert(t).That(count, is.EqualTo(0), "Count == 0")
	Assert(t).That(peek, is.EqualTo(context.DeadlineExceeded), "peek timeout")
	Assert(t).That(pop, is.EqualTo(context.DeadlineExceeded), "pop timeout")
}

func Test_MtQueue_InitialValues(t *testing.T) {
	// Arrange
	val1, val2 := randGen.Int(), randGen.Int()
	// Act
	queue := generated.MakeIntMtQueue(val1, val2)
	fmt.Printf("Create with: %v,%v\n", val1, val2)
	count0 := queue.Count()
	peek1,peekErr1 := queue.Peek(createContext())
	pop1,popErr1 := queue.Pop(createContext())
	count1 := queue.Count()
	peek2,peekErr2 := queue.Peek(createContext())
	pop2,popErr2 := queue.Pop(createContext())
	count2 := queue.Count()

	// Assert
	Assert(t).That(count0, is.EqualTo(2), "Initial Count")
	Assert(t).That(peek1, is.EqualTo(val1), "Peek 1")
	Assert(t).That(peekErr1, is.Nil, "Peek 1 OK")
	Assert(t).That(pop1, is.EqualTo(val1), "Pop 1")
	Assert(t).That(popErr1, is.Nil, "Pop 1 OK")
	Assert(t).That(count1, is.EqualTo(1), "Count after pop")
	Assert(t).That(peek2, is.EqualTo(val2), "Peek 2")
	Assert(t).That(peekErr2, is.Nil, "Peek 2 OK")
	Assert(t).That(pop2, is.EqualTo(val2), "Pop 2")
	Assert(t).That(popErr2, is.Nil, "Pop 2 OK")
	Assert(t).That(count2, is.EqualTo(0), "Count after 2x pop")
}

func Test_MtQueue_PushAndPop(t *testing.T) {
	// Arrange
	val1, val2 := randGen.Int(), randGen.Int()

	// Act
	queue := generated.IntMtQueue{}
	count0 := queue.Count()
	queue.Push(val1)
	fmt.Printf("Push: %v\n", val1)
	count1 := queue.Count()
	queue.Push(val2)
	fmt.Printf("Push: %v\n", val2)
	count2 := queue.Count()
	pop1,_ := queue.Pop(createContext())
	pop2,_ := queue.Pop(createContext())
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
	queue := generated.MakeIntMtQueue(3, 1, 2)

	// Act
	queue.Sort(func(i, j int) bool { return i < j })
	pop1,_ := queue.Pop(createContext())
	pop2,_ := queue.Pop(createContext())
	pop3,_ := queue.Pop(createContext())

	// Assert
	Assert(t).That(pop1, is.EqualTo(1), "Pop 1")
	Assert(t).That(pop2, is.EqualTo(2), "Pop 2")
	Assert(t).That(pop3, is.EqualTo(3), "Pop 3")
}

func Test_MtQueue_SortEmptyQueue(t *testing.T) {
	// Arrange
	queue := generated.IntMtQueue{}

	// Act
	queue.Sort(func(i, j int) bool { return i < j })
	isEmpty := queue.IsEmpty()

	// Assert
	Assert(t).That(isEmpty, is.True, "Empty (no sort panic)")
}

func Test_MtQueue_EmptyEvent(t *testing.T) {
	// Arrange
	queue := generated.MakeIntMtQueue(randGen.Int())

	// Act
	empty1 := queue.WaitUntilEmpty(createContext())
	queue.Pop(createContext())
	empty2:= queue.WaitUntilEmpty(createContext())

	// Assert
	Assert(t).That(empty1, is.EqualTo(context.DeadlineExceeded), "Empty 1")
	Assert(t).That(empty2, is.Nil, "Empty 2")
}

func Test_MtQueue_PeekWithDelayedPush(t *testing.T) {
	// Arrange
	val := randGen.Int()
	queue := generated.IntMtQueue{}
	var peek int
	var peekErr error
	wg := sync.WaitGroup{}

	// Act
	wg.Add(1)
	go func() {
		peek, peekErr = queue.Peek(createNoReturnContext())
		wg.Done()
	}()
	queue.Push(val)
	wg.Wait()

	// Assert
	Assert(t).That(peek, is.EqualTo(val), "Pop")
	Assert(t).That(peekErr, is.Nil, "Pop OK")
}

func Test_MtQueue_PopWithDelayedPush(t *testing.T) {
	// Arrange
	val := randGen.Int()
	queue := generated.IntMtQueue{}
	var pop int
	var popErr error
	wg := sync.WaitGroup{}

	// Act
	wg.Add(1)
	go func() {
		pop, popErr = queue.Pop(createNoReturnContext())
		wg.Done()
	}()
	queue.Push(val)
	wg.Wait()

	// Assert
	Assert(t).That(pop, is.EqualTo(val), "Pop")
	Assert(t).That(popErr, is.Nil, "Pop OK")
}
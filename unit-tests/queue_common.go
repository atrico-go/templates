package unit_tests

import (
	"context"
	"time"
)

func createContext() context.Context {
	timeout := 500 * time.Millisecond
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	return ctx
}

func createNoReturnContext() context.Context {
	timeout := time.Minute
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	return ctx
}

// +build ignore (Used to create examples only, call this before tests)

package main

import (
	"fmt"
	"os"

	"github.com/atrico-go/templates"
	"github.com/atrico-go/templates/collections"
)

func makeFile(name string) string {
	return fmt.Sprintf("unit-tests/generated/%s.go", name)
}

func main() {
	// Queue
	{
		fileDetails := templates.FileDetails{
			TargetFile: makeFile("queue"),
		}
		details := collections.QueueTemplateDetails{
			ElementType: "int",
			TypeName:    "IntQueue",
		}
		fmt.Printf("Creating %s (%s)\n", details.TypeName, fileDetails.TargetFile)
		err := collections.CreateQueue(fileDetails, details)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
	// Multi thread Queue
	{
		fileDetails := templates.FileDetails{
			TargetFile: makeFile("mt_queue"),
		}
		details := collections.QueueTemplateDetails{
			ElementType: "int",
			TypeName:    "IntMtQueue",
		}
		fmt.Printf("Creating %s (%s)\n", details.TypeName, fileDetails.TargetFile)
		err := collections.CreateMultiThreadQueue(fileDetails, details)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

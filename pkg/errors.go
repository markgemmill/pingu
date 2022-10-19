package pkg

import (
	"fmt"
	"os"
)

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func ExitOnError(err error, message string) {
	if err != nil {
		fmt.Printf("Error: %s", err)
		fmt.Printf("%s", message)
		os.Exit(1)
	}
}

func IgnoreOnError(err error) {
	if err != nil {
		fmt.Printf("Ignoring error: %s\n", err)
	}
}

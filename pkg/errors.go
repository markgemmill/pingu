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
		os.Exit(1)
	}
}

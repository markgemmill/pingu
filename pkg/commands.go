package pkg

import (
	"errors"
	"fmt"
	"strings"
)

func CheckCommand(url string, expectedStatus int, expectedContent string, storeName string, console *Console) (*StoreRecord, error) {

	//console = NewConsole(verbose)

	assertions := BuildAssertions(expectedStatus, expectedContent)

	//console.Indent()
	//console.Trace("Assertion count: %v\n", len(*assertions))
	//console.Dedent()

	urlCheck := NewUrlCheck(url, *assertions)

	urlCheck.Test()

	store := NewStore(url, storeName)
	store.Read()

	if urlCheck.Pass == true {
		console.Print("%s GET %s\n", Green(PASS), url)
		store.Save(PASS, "")
		store.Write()
		return nil, nil
	}

	b := strings.Builder{}

	console.Indent()
	for _, msg := range urlCheck.Errors {
		console.Print("%s: %s\n", Red(FAIL), msg)
		fmt.Fprintf(&b, "%s; ", msg)
	}
	console.Dedent()

	store.Save(FAIL, b.String())
	store.Write()

	return &store.Data.Current, errors.New("Url check failed.")

}

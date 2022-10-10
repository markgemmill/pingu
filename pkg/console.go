package pkg

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
	"sync"
)

/*
		0 - Print, Error
	    1 - Info
		2 - Debug
		3 - Trace
*/

var Red func(opt ...interface{}) string
var Green func(opt ...interface{}) string
var Yellow func(opt ...interface{}) string

func init() {
	Red = color.New(color.FgRed).SprintFunc()
	Green = color.New(color.FgGreen).SprintFunc()
	Yellow = color.New(color.FgYellow).SprintFunc()
}

type Console struct {
	Verbosity    int
	IndentAmount int
}

var console Console
var once sync.Once

func NewConsole(verbosity int) *Console {
	once.Do(func() {
		console = Console{
			Verbosity: verbosity,
		}
	})
	return &console
}

func (c *Console) Indent() {
	c.IndentAmount += 1
}

func (c *Console) Dedent() {
	c.IndentAmount -= 1
	if c.IndentAmount < 0 {
		c.IndentAmount = 0
	}
}

func (c *Console) indentMessage(message string) string {
	if c.IndentAmount == 0 {
		return message
	}

	indention := strings.Repeat("  ", c.IndentAmount)

	return indention + message
}

func (c *Console) Print(message string, opt ...interface{}) {
	if c.Verbosity >= 0 {
		fmt.Printf(c.indentMessage(message), opt...)
	}
}

func (c *Console) Info(message string, opt ...interface{}) {
	if c.Verbosity >= 1 {
		fmt.Printf(c.indentMessage(message), opt...)
	}
}

func (c *Console) Debug(message string, opt ...interface{}) {
	if c.Verbosity >= 2 {
		fmt.Printf(c.indentMessage(message), opt...)
	}
}

func (c *Console) Trace(message string, opt ...interface{}) {
	if c.Verbosity >= 3 {
		fmt.Printf(c.indentMessage(message), opt...)
	}
}

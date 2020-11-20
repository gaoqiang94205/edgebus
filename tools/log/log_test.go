package log

import (
	"fmt"
	"os"
	"testing"
)

func TestLogPath(t *testing.T) {
	fmt.Print(getCurrentDirectory())
	fmt.Println(os.Args[1])
}

package utils

import (
    "testing"
    "fmt"
)

func TestDir(t *testing.T) {
    list, err := Dir("*.go")
    fmt.Println(list)
    if err != nil {
       t.Errorf("Dir errored")
    }
}

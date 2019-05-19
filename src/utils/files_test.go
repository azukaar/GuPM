package utils

import (
    "testing"
)

func TestDir(t *testing.T) {
    list, err := Dir("*.go")
    if err != nil {
       t.Errorf("Dir errored")
    }
}

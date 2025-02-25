//go:build dev

package main

import (
	"os"
)

var static = os.DirFS(".")

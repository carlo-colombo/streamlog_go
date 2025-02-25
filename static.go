//go:build !dev

package main

import "embed"

//go:embed app/dist/app/browser
var static embed.FS

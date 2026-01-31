package main

import (
	"github.com/dickus/dreadnotes/internal/config"
)

func main() {
	config.LoadConfig()
	config.ReadFile()
}


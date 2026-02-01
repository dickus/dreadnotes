package main

import (
	"github.com/dickus/dreadnotes/internal/config"
	"github.com/dickus/dreadnotes/internal/utils"
)

func main() {
	config.LoadConfig()
	utils.ArgsParser()
}


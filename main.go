/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/cjairm/devgita/cmd"
	"github.com/cjairm/devgita/internal/embedded"
)

func main() {
	// Set the default extractor function for devgita app
	embedded.DefaultExtractor = ExtractEmbeddedConfigs

	cmd.Execute()
}

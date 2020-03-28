package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/petarov/nenu/config"
)

const heart = "\u2764"

var redColor = color.New(color.FgRed)
var templatesPath = flag.String("t", "web/blazer", "Path to HTML template to use")
var configPath = flag.String("c", "config.yml", "Path to YAML configuration file")
var contentPath = flag.String("b", "", "Path to blog contents. A directory that contains markdown (.md) files")

func verifyPath(path string, what string) {
	if len(path) < 1 {
		redColor.Println(what + " path not specified")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		redColor.Println(what + " path not found")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	cyan := color.New(color.FgCyan).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("%s %s (%s) %s v%s - Tiny Static Blog Builder\n", red(heart), config.AppNameEN, cyan(config.AppNameBG),
		red(heart), config.VERSION)

	flag.Parse()

	verifyPath(*templatesPath, "Templates")
	verifyPath(*configPath, "Config")
	verifyPath(*contentPath, "Blog contents")

	yml, err := config.ParseYMLConfig(*configPath)
	if err != nil {
		fmt.Printf("Failed loading configuration! %v\n", red(err))
		os.Exit(1)
	}

	fmt.Printf("CFG: %v\n", yml)
}

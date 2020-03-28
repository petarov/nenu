package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/petarov/nenu/config"
)

const heart = "\u2764"

func main() {
	cyan := color.New(color.FgCyan).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("%s %s (%s) %s v%s - Tiny Static Blog Builder\n", red(heart), config.AppNameEN, cyan(config.AppNameBG), red(heart), config.VERSION)

	args := os.Args[1:]
	if len(args) < 1 {
		color.New(color.FgRed).Printf("Insufficient arguments!")
		os.Exit(1)
	}

	yml, err := config.ParseYMLConfig(args[0])
	if err != nil {
		fmt.Printf("Failed loading configuration! %v\n", red(err))
		os.Exit(1)
	}

	fmt.Printf("CFG: %v\n", yml)
}

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/petarov/nenu/config"
	"github.com/petarov/nenu/engine"
)

var (
	heart    = "\u2764"
	redColor = color.New(color.FgRed)
	cyan     = color.New(color.FgCyan).SprintFunc()
	red      = color.New(color.FgRed).SprintFunc()
)

func init() {
	flag.StringVar(&config.ConfigPath, "c", "config.yml", "Path to YAML configuration file")
	flag.StringVar(&config.TemplatePath, "t", "web/blazer", "Path to HTML template to use")
	flag.StringVar(&config.PostsPath, "b", "", "Path to a directory with markdown (.md) posts")
	flag.StringVar(&config.OutputPath, "o", "", "Path to where to write the generated HTML website files")
	flag.StringVar(&config.TempPath, "tmp", os.TempDir(), "Temporary path used during content generation")
}

func verifyPath(path string, what string, mustExist bool) {
	if len(path) < 1 {
		redColor.Println(what + " path not specified")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if mustExist {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			redColor.Println(what + " path not found")
			flag.PrintDefaults()
			os.Exit(1)
		}
	}
}

func main() {
	fmt.Printf("%s %s (%s) %s v%s - Tiny Static Blog Builder\n", red(heart), config.AppNameEN, cyan(config.AppNameBG),
		red(heart), config.VERSION)

	flag.Parse()

	verifyPath(config.ConfigPath, "Templates", true)
	verifyPath(config.TemplatePath, "Config", true)
	verifyPath(config.PostsPath, "Blog posts", true)
	verifyPath(config.OutputPath, "Output path", false)

	yml, err := config.ParseYMLConfig(config.ConfigPath)
	if err != nil {
		fmt.Printf("Failed loading configuration! %v\n", red(err))
		os.Exit(1)
	}
	config.YMLConfig = yml

	if err = engine.Spew(); err != nil {
		fmt.Printf("Failed parsing contents! %v\n", red(err))
		os.Exit(1)
	}

	fmt.Printf("Web site contents written to %s\n", cyan(config.OutputPath))
}

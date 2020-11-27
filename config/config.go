package config

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/url"

	"github.com/goccy/go-yaml"
)

type YML struct {
	Site struct {
		Title       string `yaml:"title"`
		Description string `yaml:"description"`
		URL         string `yaml:"url"`
	}
	Content struct {
		FontLink string `yaml:"font_link"`
		FontName string `yaml:"font_name"`
		Timezone string `yaml:"timezone"`
		Locale   string `yaml:"locale"`
		Rss      bool   `yaml:"rss"`
	}
	Footer struct {
		Copyright     string `yaml:"copyright"`
		CopyrightHTML template.HTML
		Twitter       string `yaml:"twitter"`
		ShowBuilder   bool   `yaml:"show_builder"`
	}
	Locales map[string]string
}

func validate(yml *YML) (err error) {
	if len(yml.Site.Title) == 0 {
		return errors.New("site.title cannot be empty")
	}
	if _, err = url.ParseRequestURI(yml.Site.URL); err != nil {
		return err
	}
	if len(yml.Content.Timezone) == 0 {
		return errors.New("content.timezone cannot be empty")
	}
	if _, err = url.ParseRequestURI(yml.Content.FontLink); err != nil {
		return err
	}
	return nil
}

// ParseYMLConfig parses config.yml
func ParseYMLConfig(filepath string) (*YML, error) {
	fmt.Printf("Parsing config file %s\n", filepath)

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var yml YML

	err = yaml.Unmarshal(data, &yml)
	if err != nil {
		return nil, err
	}

	if err = validate(&yml); err != nil {
		return nil, err
	}

	yml.Locales = Locales["en_US"]

	if Locales[yml.Content.Locale] != nil {
		yml.Locales = Locales[yml.Content.Locale]
	} else {
		// default
		yml.Locales = Locales["en_US"]
	}

	yml.Footer.CopyrightHTML = template.HTML(yml.Footer.Copyright)

	return &yml, nil
}

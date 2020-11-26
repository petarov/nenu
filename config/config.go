package config

import (
	"errors"
	"fmt"
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
		Timezone string `yaml:"timezone"`
		Locale   string `yaml:"locale"`
		Rss      bool   `yaml:"rss"`
	}
	Footer struct {
		Copyright   string `yaml:"copyright"`
		Twitter     string `yaml:"twitter"`
		ShowBuilder bool   `yaml:"show_builder"`
	}
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

	return &yml, nil
}

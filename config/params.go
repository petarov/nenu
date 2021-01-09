package config

import (
	"time"
)

var (
	// ConfigPath config.yml file path
	ConfigPath string
	// ThemePath path to web template to use for the blog
	ThemePath string
	// PostsPath path to markdown (.md) posts
	PostsPath string
	// OutputPath path to where the site contents will be written
	OutputPath string
	// TempPath temporary path used during content generation
	TempPath string
	// YMLConfig configuration YAML file
	YMLConfig *YML
	// TimeZoneLocation inited from the YMLConfig
	TimeZoneLocation *time.Location
)

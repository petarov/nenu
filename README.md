nenu
===================

[![build](https://github.com/petarov/nenu/workflows/CI%20Build/badge.svg)](https://github.com/petarov/nenu/actions?query=workflow%3A%22CI+Build%22)
[![goreport](https://goreportcard.com/badge/github.com/petarov/nenu)](https://goreportcard.com/report/github.com/petarov/nenu)

`пепи` - A tiny static site generator for your journal.

  * Binary executable without additional dependencies
  * Drop-in replacement for Jekyll markdown posts
  * One YAML config file
  * Archive, atom feed and custom themes support
  
See [Demo](http://petarov.github.io/nenu/)

# Installation

[Download binaries](https://github.com/petarov/nenu/releases) for Linux, macOS or Windows.

Create a new configuration file:

    cp ./config.yml.tpl config.yml

Make sure the `url`, `title` and `description` are specified.

# Usage

Generate web site contents into an output folder called `_site` using the markdown files from `_posts`. By default, the website uses the  `blazer` theme:

    nenu_linux_amd64 -p _posts -o _site

Specify another theme to use:

    nenu_linux_amd64 -p _posts -o _site -t themes/my-custom-theme

## Markdown post meta params

  * `title: <value>` - Post title
  * `subtitle: <value>` - Post subtitle. Optional.
  * `summary: <value>` - Post meta description. Optional.
  * `publish: <true|false>` - Default `true`. Set to `false` to skip generating this post. Optional.

See examples in [test-data](test-data).

# Development

пепи is a tool that should remain as small as possible. Adding new features is nice, but not really the goal of the project.

Refer to the [TODO](TODO.md) list to check what's ~~coming~~ missing.

If you'd like to add your own language locale, check `config/locales.go`.

# License

[MIT](LICENSE)

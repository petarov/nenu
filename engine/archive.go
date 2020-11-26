package engine

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/petarov/nenu/config"
)

type archivePageData struct {
	*config.YML
	Posts []PostMeta
	// Post an always null hack that helps with the header template if condition
	Post *PostMeta
}

// SpewArchive generate site posts archive
func SpewArchive(meta []*PostMeta, templates *Templates) error {
	fmt.Print("| Generating archive ...\n")

	var f *os.File
	f, err := os.Create(filepath.Join(config.TempPath, "archive.html"))
	if err != nil {
		return err
	}
	defer f.Close()

	metas := make([]PostMeta, 0)
	for _, p := range meta {
		if p.Publish {
			metas = append(metas, *p)
		}
	}

	pd := &archivePageData{config.YMLConfig, metas, nil}

	// header
	if err = templates.Header.Execute(f, pd); err != nil {
		return err
	}
	// body
	if err = templates.Archive.Execute(f, pd); err != nil {
		return err
	}
	// footer
	if err = templates.Footer.Execute(f, pd); err != nil {
		return err
	}

	return nil
}

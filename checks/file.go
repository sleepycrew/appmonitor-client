package checks

import (
	"errors"
	"fmt"
	"github.com/sleepycrew/appmonitor-client/pkg/data"
	. "github.com/sleepycrew/appmonitor-client/pkg/data/result"
	"io/fs"
	"os"
)

// TODO improve output

type FileCheckSettings struct {
	Exists bool
	Dir    bool
	Link   bool
}

type FileCheck struct {
	name        string
	description string
	filePath    string
	settings    FileCheckSettings
}

func (c FileCheck) GetName() string {
	return c.name
}

func (c FileCheck) GetDescription() *string {
	return &c.description
}

type fileCheckMapping = func(settings FileCheckSettings, info os.FileInfo) (Result, string)

func wrapFileCheckMapping(name string, fun func(settings FileCheckSettings, info os.FileInfo) bool) fileCheckMapping {
	return func(settings FileCheckSettings, info os.FileInfo) (Result, string) {
		if fun(settings, info) {
			return OK, name
		} else {
			return Error, name
		}
	}
}

var fileCheckMappings = []fileCheckMapping{
	wrapFileCheckMapping("dir", func(settings FileCheckSettings, info os.FileInfo) bool {
		return settings.Dir == info.IsDir()
	}),
	wrapFileCheckMapping("link", func(settings FileCheckSettings, info os.FileInfo) bool {
		isLink := info.Mode() == os.ModeSymlink
		return settings.Link == isLink
	}),
}

func (c FileCheck) checkFile() (Result, string) {
	if fs.ValidPath(c.filePath) {
		return Unknown, "path is invalid"
	}

	stat, err := os.Stat(c.filePath)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if !c.settings.Exists {
				return OK, "File does not exists"
			}
		}

		return Unknown, fmt.Sprint("could not get file statistics: %i", err)
	}

	for _, mapping := range fileCheckMappings {
		result, name := mapping(c.settings, stat)
		if result != OK {
			return Error, fmt.Sprintf("file does not meet requirement: %s", name)
		}
	}

	return OK, "File exists and meets requirements"
}

func (c FileCheck) RunCheck(results chan<- data.ClientCheck) {
	result, value := c.checkFile()

	results <- data.ClientCheck{
		Name:        c.GetName(),
		Value:       value,
		Description: c.description,
		Result:      int(result),
	}
}

func NewFileCheck(name string, description string, filePath string, settings FileCheckSettings) FileCheck {
	return FileCheck{
		name,
		description,
		filePath,
		settings,
	}
}

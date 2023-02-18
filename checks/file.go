package checks

import (
	"errors"
	"fmt"
	"github.com/sleepycrew/appmonitor-client/pkg/check"
	"github.com/sleepycrew/appmonitor-client/pkg/data/result"
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
	filePath string
	settings FileCheckSettings
}

type fileCheckMapping = func(settings FileCheckSettings, info os.FileInfo) (result.Code, string)

func wrapFileCheckMapping(name string, fun func(settings FileCheckSettings, info os.FileInfo) bool) fileCheckMapping {
	return func(settings FileCheckSettings, info os.FileInfo) (result.Code, string) {
		if fun(settings, info) {
			return result.OK, name
		} else {
			return result.Error, name
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

func (c FileCheck) checkFile() (result.Code, string) {
	if fs.ValidPath(c.filePath) {
		return result.Unknown, "path is invalid"
	}

	stat, err := os.Stat(c.filePath)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if !c.settings.Exists {
				return result.OK, "File does not exists"
			}
		}

		return result.Unknown, fmt.Sprint("could not get file statistics: %i", err)
	}

	for _, mapping := range fileCheckMappings {
		res, name := mapping(c.settings, stat)
		if res != result.OK {
			return result.Error, fmt.Sprintf("file does not meet requirement: %s", name)
		}
	}

	return result.OK, "File exists and meets requirements"
}

func (c FileCheck) RunCheck(output chan<- check.Result) {
	result, value := c.checkFile()

	output <- check.Result{
		Value:  value,
		Result: result,
	}
}

func NewFileCheck(filePath string, settings FileCheckSettings) FileCheck {
	return FileCheck{
		filePath,
		settings,
	}
}

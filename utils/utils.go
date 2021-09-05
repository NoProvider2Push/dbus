package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// storagePaths provides a basic cache for StoragePath
var storagePaths = map[string]string{}

// StoragePath appName only recommends using something that can be a filename for now
func StoragePath(appName string) string {
	if a, ok := storagePaths[appName]; ok {
		return a
	}

	basedir := os.Getenv("NP2P_STORAGE_PATH")
	if len(basedir) == 0 {
		basedir = os.Getenv("XDG_CONFIG_HOME")
		if len(basedir) == 0 {
			basedir = os.Getenv("HOME")
			if len(basedir) == 0 {
				basedir = "./" // FIXME: set to cwd if dunno wth is going on
			}
			basedir = filepath.Join(basedir, ".config")
		}
		basedir = filepath.Join(basedir, "unifiedpush", "distributors")
		err := os.MkdirAll(basedir, 0o700)
		if err != nil {
			basedir = "./"
			// FIXME idk wth to do when there's an error here
		}
	}
	finalFilename := filepath.Join(basedir, appName)
	storagePaths[appName] = finalFilename
	return finalFilename
}

var Log Logger

type Logger struct {
}

func (Logger) Debugln(inps ...interface{}) {
	if os.Getenv("UP_NP2P_DEBUG") == "true" || strings.HasPrefix(os.Args[0], "/tmp/go-build") {
		log.Println(inps...)
	}
}

func (Logger) Infoln(inps ...interface{}) {
	log.Println(inps...)
}

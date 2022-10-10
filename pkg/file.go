package pkg

import (
	"github.com/markgemmill/appdirs"
	"github.com/spf13/afero"
)

var dirs appdirs.AppDirs
var fs afero.Fs

func init() {
	dirs = appdirs.NewAppDirs("pingu", "")
	fs = afero.NewOsFs()

	err := fs.MkdirAll(dirs.UserDataDir(), 0777)
	if err != nil {
		panic(err)
	}
}

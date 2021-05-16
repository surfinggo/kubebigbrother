package fileorcreate

import (
	"github.com/pkg/errors"
	"github.com/spongeprojects/magicconch"
	"io"
	"os"
)

// Ensure ensures the file exist, create it with a template if not.
func Ensure(filePath, tmplPath string) error {
	_, err := os.Stat(filePath)
	if err == nil {
		return nil // already exist
	}
	if !os.IsNotExist(err) {
		return errors.Wrapf(err, "os.Stat: %s error", filePath)
	}

	_, err = os.Stat(tmplPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.Errorf("template not exist: %s", tmplPath)
		}
		return errors.Wrapf(err, "os.Stat: %s error", tmplPath)
	}

	tmpl, err := os.Open(tmplPath)
	if err != nil {
		return errors.Wrapf(err, "os.Open: %s error", tmplPath)
	}
	defer func() {
		magicconch.Must(tmpl.Close())
	}()

	file, err := os.Create(filePath)
	if err != nil {
		return errors.Wrapf(err, "os.Open: %s error", filePath)
	}
	defer func() {
		magicconch.Must(file.Close())
	}()
	_, err = io.Copy(file, tmpl)
	if err != nil {
		return errors.Wrap(err, "io.Copy error")
	}
	return nil
}

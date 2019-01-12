package storage

import (
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/pkg/errors"
)

type FormattedFS struct {
	fs     boshsys.FileSystem
	parser Storage
}

var _ Storage = FormattedFS{}

func NewFormattedFS(fs boshsys.FileSystem, parser Storage) FormattedFS {
	return FormattedFS{
		fs:     fs,
		parser: parser,
	}
}

func (s FormattedFS) Get(path string, get interface{}) error {
	absPath, err := s.fs.ExpandPath(path)
	if err != nil {
		return err
	}

	if s.fs.FileExists(absPath) {
		bytes, err := s.fs.ReadFile(absPath)
		if err != nil {
			return errors.Wrapf(err, "Reading config file '%s'", absPath)
		}

		err = s.parser.Get(string(bytes), get)
		if err != nil {
			return errors.Wrap(err, "Parsing config")
		}
	}

	return nil
}

func (s FormattedFS) Put(path string, put interface{}) (string, error) {
	absPath, err := s.fs.ExpandPath(path)
	if err != nil {
		return "", err
	}

	bytes, err := s.parser.Put("", put)
	if err != nil {
		return "", errors.Wrap(err, "Serializing config")
	}

	err = s.fs.WriteFileString(absPath, bytes)
	if err != nil {
		return "", errors.Wrapf(err, "Writing config file '%s'", absPath)
	}

	return path, nil
}

package storage

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
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
			return bosherr.WrapErrorf(err, "Reading config file '%s'", absPath)
		}

		err = s.parser.Get(string(bytes), get)
		if err != nil {
			return bosherr.WrapError(err, "Parsing config")
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
		return "", bosherr.WrapError(err, "Serializing config")
	}

	err = s.fs.WriteFileString(absPath, bytes)
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Writing config file '%s'", absPath)
	}

	return path, nil
}

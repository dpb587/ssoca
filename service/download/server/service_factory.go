package server

import (
	"os"
	"path/filepath"

	boshcrypto "github.com/cloudfoundry/bosh-utils/crypto"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	"github.com/dpb587/ssoca/config"
	"github.com/dpb587/ssoca/server/service"
	svc "github.com/dpb587/ssoca/service/download"
	svcconfig "github.com/dpb587/ssoca/service/download/config"
)

type ServiceFactory struct {
	fs boshsys.FileSystem
}

var _ service.ServiceFactory = ServiceFactory{}

func NewServiceFactory(fs boshsys.FileSystem) ServiceFactory {
	return ServiceFactory{
		fs: fs,
	}
}

func (f ServiceFactory) Type() string {
	return svc.Service{}.Type()
}

func (f ServiceFactory) Create(name string, options map[string]interface{}) (service.Service, error) {
	var cfg svcconfig.Config

	config.RemarshalYAML(options, &cfg)

	paths, err := f.fs.Glob(cfg.Glob)
	if err != nil {
		return nil, bosherr.WrapError(err, "Globbing")
	}

	cfg.Paths = []svcconfig.PathConfig{}

	for _, path := range paths {
		stat, err := f.fs.Stat(path)
		if err != nil {
			return nil, bosherr.WrapError(err, "Stat file")
		}

		file, err := f.fs.OpenFile(path, os.O_RDONLY, 0)
		if err != nil {
			return nil, bosherr.WrapError(err, "Opening file for digest")
		}

		defer file.Close()

		digest, err := boshcrypto.DigestAlgorithmSHA1.CreateDigest(file)
		if err != nil {
			return nil, bosherr.WrapError(err, "Creating file digest")
		}

		cfg.Paths = append(
			cfg.Paths,
			svcconfig.PathConfig{
				Name:   filepath.Base(path),
				Path:   path,
				Size:   stat.Size(),
				Digest: digest.String(),
			},
		)
	}

	return NewService(name, cfg, f.fs), nil
}

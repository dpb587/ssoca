package server

import (
	"os"
	"path/filepath"
	"strings"

	boshcrypto "github.com/cloudfoundry/bosh-utils/crypto"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/config"
	"github.com/dpb587/ssoca/server/service"
	svc "github.com/dpb587/ssoca/service/download"
	svcconfig "github.com/dpb587/ssoca/service/download/server/config"
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
		return nil, errors.Wrap(err, "Globbing")
	}

	cfg.Paths = []svcconfig.PathConfig{}

	for _, path := range paths {
		stat, err := f.fs.Stat(path)
		if err != nil {
			return nil, errors.Wrap(err, "Stat file")
		}

		file, err := f.fs.OpenFile(path, os.O_RDONLY, 0)
		if err != nil {
			return nil, errors.Wrap(err, "Opening file for digest")
		}

		defer file.Close()

		digestSHA1, err := boshcrypto.DigestAlgorithmSHA1.CreateDigest(file)
		if err != nil {
			return nil, errors.Wrap(err, "Creating sha1 digest")
		}

		digestSHA256, err := boshcrypto.DigestAlgorithmSHA256.CreateDigest(file)
		if err != nil {
			return nil, errors.Wrap(err, "Creating sha256 digest")
		}

		digestSHA512, err := boshcrypto.DigestAlgorithmSHA512.CreateDigest(file)
		if err != nil {
			return nil, errors.Wrap(err, "Creating sha512 digest")
		}

		cfg.Paths = append(
			cfg.Paths,
			svcconfig.PathConfig{
				Name: filepath.Base(path),
				Path: path,
				Size: stat.Size(),
				Digest: svcconfig.PathDigestConfig{
					SHA1:   digestSHA1.String(),
					SHA256: strings.SplitN(digestSHA256.String(), ":", 2)[1],
					SHA512: strings.SplitN(digestSHA512.String(), ":", 2)[1],
				},
			},
		)
	}

	return NewService(name, cfg, f.fs), nil
}

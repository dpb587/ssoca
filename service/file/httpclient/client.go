package httpclient

import (
	"fmt"
	"io"
	"net/url"

	"github.com/cheggaaa/pb"
	boshcrypto "github.com/cloudfoundry/bosh-utils/crypto"
	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/httpclient"
	"github.com/dpb587/ssoca/service/file/api"
)

func New(baseclient httpclient.Client, service string) (Client, error) {
	if baseclient == nil {
		return nil, errors.New("client is nil")
	}

	return &client{
		client:  baseclient,
		service: service,
	}, nil
}

type client struct {
	client  httpclient.Client
	service string
}

func (c client) GetList() (api.ListResponse, error) {
	out := api.ListResponse{}
	path := fmt.Sprintf("/%s/list", c.service)

	err := c.client.APIGet(path, &out)
	if err != nil {
		return out, errors.Wrapf(err, "getting %s", path)
	}

	return out, nil
}

func (c client) GetMetadata() (api.MetadataResponse, error) {
	out := api.MetadataResponse{}
	path := fmt.Sprintf("/%s/metadata", c.service)

	err := c.client.APIGet(path, &out)
	if err != nil {
		return out, errors.Wrapf(err, "getting %s", path)
	}

	return out, nil
}

func (c client) Download(name string, target io.ReadWriteSeeker, downloadStatus *pb.ProgressBar) error {
	list, err := c.GetList()
	if err != nil {
		return errors.Wrap(err, "listing artifacts")
	}

	for _, file := range list.Files {
		if file.Name != name {
			continue
		}

		return c.download(file, target, downloadStatus)
	}

	return fmt.Errorf("file is not known: %s", name)
}

func (c *client) download(file api.ListFileResponse, target io.ReadWriteSeeker, downloadStatus *pb.ProgressBar) error {
	path := fmt.Sprintf("/%s/get?name=%s", c.service, url.QueryEscape(file.Name))

	res, err := c.client.Get(path)
	if err != nil {
		return errors.Wrap(err, "getting file")
	}

	downloadStatus.Total = file.Size
	downloadStatus.SetUnits(pb.U_BYTES)
	downloadStatus.Start()
	defer downloadStatus.Finish()

	_, err = io.Copy(target, downloadStatus.NewProxyReader(res.Body))
	if err != nil {
		return errors.Wrap(err, "streaming to file")
	}

	var algo boshcrypto.Algorithm
	var hash string

	if file.Digest.SHA512 != "" {
		algo = boshcrypto.DigestAlgorithmSHA512
		hash = fmt.Sprintf("sha512:%s", file.Digest.SHA512)
	} else if file.Digest.SHA256 != "" {
		algo = boshcrypto.DigestAlgorithmSHA256
		hash = fmt.Sprintf("sha256:%s", file.Digest.SHA256)
	} else if file.Digest.SHA1 != "" {
		algo = boshcrypto.DigestAlgorithmSHA1
		hash = file.Digest.SHA1
	} else {
		return errors.New("no digest available to verify download")
	}

	digest, err := algo.CreateDigest(target)
	if err != nil {
		return errors.Wrap(err, "creating digest")
	}

	if digest.String() != hash {
		return fmt.Errorf("expected digest '%s' but got '%s'", hash, digest.String())
	}

	return nil
}

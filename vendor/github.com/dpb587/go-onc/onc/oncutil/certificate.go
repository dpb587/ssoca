package oncutil

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// seems Encode is not natively supported https://godoc.org/golang.org/x/crypto/pkcs12
func ConvertKeyPairToPKCS12(certificate []byte, privateKey []byte) ([]byte, error) {
	ch, err := ioutil.TempFile("", "pkcs12-crt-*")
	if err != nil {
		return nil, errors.Wrap(err, "creating certificate file")
	}

	defer os.Remove(ch.Name())
	defer ch.Close()

	_, err = ch.Write(certificate)
	if err != nil {
		return nil, errors.Wrap(err, "writing certificate")
	}

	if err = ch.Close(); err != nil {
		return nil, errors.Wrap(err, "closing certificate")
	}

	kh, err := ioutil.TempFile("", "pkcs12-key-*")
	if err != nil {
		return nil, errors.Wrap(err, "creating private key file")
	}

	defer os.Remove(kh.Name())
	defer kh.Close()

	_, err = kh.Write(privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "writing private key")
	}

	if err = kh.Close(); err != nil {
		return nil, errors.Wrap(err, "closing private key")
	}

	oh, err := ioutil.TempFile("", "pkcs12-out-*")
	if err != nil {
		return nil, errors.Wrap(err, "creating output file")
	}

	defer os.Remove(oh.Name())
	defer oh.Close()

	cmd := exec.Command(
		"openssl",
		"pkcs12",
		"-export",
		"-inkey", kh.Name(),
		"-in", ch.Name(),
		"-out", oh.Name(),
		"-passout", "pass:",
	)

	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("openssl exited with error [%s]: stderr: %s", err, stderr.Bytes())
	}

	output, err := ioutil.ReadAll(oh)
	if err != nil {
		return nil, errors.Wrap(err, "reading pkcs12 file")
	}

	return output, nil
}

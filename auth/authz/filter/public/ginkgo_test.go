package public_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPublic(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "github.com/dpb587/ssoca/auth/authz/filter/public")
}

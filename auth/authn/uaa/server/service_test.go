package server_test

import (
	. "github.com/dpb587/ssoca/auth/authn/uaa/server"

	"github.com/dpb587/ssoca/server/service"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("Service", func() {
	Describe("interface", func() {
		It("github.com/dpb587/ssoca/server/service.AuthService", func() {
			var _ service.AuthService = (*Service)(nil)
		})
	})
})

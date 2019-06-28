package service_test

import (
	. "github.com/dpb587/ssoca/server/service"
	"github.com/dpb587/ssoca/server/service/servicefakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {
	Describe("GetAuthServices", func() {
		It("works", func() {
			subject := NewDefaultManager()

			svc1 := &servicefakes.FakeAuthService{}
			svc1.NameReturns("svc1")
			subject.Add(svc1)

			svc2 := &servicefakes.FakeService{}
			svc2.NameReturns("svc2")
			subject.Add(svc2)

			svc3 := &servicefakes.FakeAuthService{}
			svc3.NameReturns("svc3")
			subject.Add(svc3)

			svc4 := &servicefakes.FakeService{}
			svc4.NameReturns("svc4")
			subject.Add(svc4)

			svc5 := &servicefakes.FakeAuthService{}
			svc5.NameReturns("svc5")
			subject.Add(svc5)

			res := GetAuthServices(subject)
			Expect(res).To(HaveLen(3))
		})
	})
})

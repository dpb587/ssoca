package dynamicvalue_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDynamicvalue(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "github.com/dpb587/ssoca/server/service/dynamicvalue")
}

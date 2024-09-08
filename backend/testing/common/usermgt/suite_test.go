package usermgt_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gleak"
)

var _ = BeforeSuite(func() {
	IgnoreGinkgoParallelClient()
})

func TestUserManagement(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "common/user package")
}

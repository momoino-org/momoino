package showmgt_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gleak"
)

var _ = BeforeSuite(func() {
	IgnoreGinkgoParallelClient()
})

func TestShowManagement(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Show Management")
}

package httpsrv_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gleak"
)

var _ = BeforeSuite(func() {
	IgnoreGinkgoParallelClient()
})

func TestHTTPSrvModule(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "httpsrv suite")
}

package versions_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gleak"
)

var _ = BeforeSuite(func() {
	IgnoreGinkgoParallelClient()
})

func TestDBMigration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DB Migration")
}

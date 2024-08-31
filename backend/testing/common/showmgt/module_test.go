package showmgt_test

import (
	"wano-island/common/showmgt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
)

var _ = Describe("[module.go]", func() {
	Context("when initializing the show management module", func() {
		It("should return an fx.Option", func() {
			Expect(showmgt.NewShowMgtModule()).To(BeAssignableToTypeOf(fx.Module("")))
		})
	})
})

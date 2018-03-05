package backend

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cleanup API", func() {
	It("should succeed", func() {
		code, err := Cleanup()
		Expect(err).ShouldNot(HaveOccurred(), "Cleanup returned an error")
		Expect(code).Should(Equal(http.StatusNoContent))
	})
})

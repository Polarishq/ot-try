package backend

import (
	"github.com/Polarishq/ot-try/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Health API", func() {
	It("should return healthy", func() {
		health, err := Health()

		Expect(err).ShouldNot(HaveOccurred(), "Health check returned an error")

		healthy := true
		si := models.ServiceInfo{
			Name:    &Constants.ComponentName,
			Version: &Constants.Version,
		}
		h := models.Health{
			Healthy:     &healthy,
			ServiceInfo: &si,
			Revision:    &GitShaCommit,
		}

		Expect(h).To(Equal(health))
	})
})

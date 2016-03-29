package requirements_test

import (
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/models"
	. "github.com/theophoric/cf-cli/cf/requirements"
	testconfig "github.com/theophoric/cf-cli/testhelpers/configuration"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TargetedSpaceRequirement", func() {
	var (
		config core_config.ReadWriter
	)

	BeforeEach(func() {
		config = testconfig.NewRepositoryWithDefaults()
	})

	Context("when the user has targeted a space", func() {
		It("succeeds", func() {
			req := NewTargetedSpaceRequirement(config)
			err := req.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("when the user does not have a space targeted", func() {
		It("errors", func() {
			config.SetSpaceFields(models.SpaceFields{})

			err := NewTargetedSpaceRequirement(config).Execute()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("No space targeted"))
		})
	})
})

package requirements_test

import (
	testapi "github.com/theophoric/cf-cli/cf/api/fakes"
	"github.com/theophoric/cf-cli/cf/errors"
	"github.com/theophoric/cf-cli/cf/models"
	. "github.com/theophoric/cf-cli/cf/requirements"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ServiceInstanceRequirement", func() {
	Context("when a service instance with the given name can be found", func() {
		It("succeeds", func() {
			instance := models.ServiceInstance{}
			instance.Name = "my-service"
			instance.Guid = "my-service-guid"
			repo := &testapi.FakeServiceRepository{}
			repo.FindInstanceByNameReturns(instance, nil)

			req := NewServiceInstanceRequirement("my-service", repo)

			err := req.Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(repo.FindInstanceByNameArgsForCall(0)).To(Equal("my-service"))
			Expect(req.GetServiceInstance()).To(Equal(instance))
		})
	})

	Context("when a service instance with the given name can't be found", func() {
		It("errors", func() {
			repo := &testapi.FakeServiceRepository{}
			repo.FindInstanceByNameReturns(models.ServiceInstance{}, errors.NewModelNotFoundError("Service instance", "my-service"))
			err := NewServiceInstanceRequirement("foo", repo).Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Service instance my-service not found"))
		})
	})
})

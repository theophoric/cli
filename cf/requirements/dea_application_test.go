package requirements_test

import (
	"errors"

	"github.com/theophoric/cf-cli/cf/models"
	"github.com/theophoric/cf-cli/cf/requirements"

	testApplication "github.com/theophoric/cf-cli/cf/api/applications/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeaApplication", func() {
	var (
		req     requirements.DEAApplicationRequirement
		appRepo *testApplication.FakeApplicationRepository
		appName string
	)

	BeforeEach(func() {
		appName = "fake-app-name"
		appRepo = &testApplication.FakeApplicationRepository{}
		req = requirements.NewDEAApplicationRequirement(appName, appRepo)
	})

	Describe("GetApplication", func() {
		It("returns an empty application", func() {
			Expect(req.GetApplication()).To(Equal(models.Application{}))
		})

		Context("when the requirement has been executed", func() {
			BeforeEach(func() {
				app := models.Application{}
				app.Guid = "fake-app-guid"
				appRepo.ReadReturns(app, nil)

				req.Execute()
			})

			It("returns the application", func() {
				Expect(req.GetApplication().Guid).To(Equal("fake-app-guid"))
			})
		})
	})

	Describe("Execute", func() {
		Context("when the returned application is a Diego application", func() {
			BeforeEach(func() {
				app := models.Application{}
				app.Diego = true
				appRepo.ReadReturns(app, nil)
			})

			It("fails with error", func() {
				err := req.Execute()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("The app is running on the Diego backend, which does not support this command."))
			})
		})

		Context("when the returned application is not a Diego application", func() {
			BeforeEach(func() {
				app := models.Application{}
				app.Diego = false
				appRepo.ReadReturns(app, nil)
			})

			It("succeeds", func() {
				err := req.Execute()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when finding the application results in an error", func() {
			BeforeEach(func() {
				appRepo.ReadReturns(models.Application{}, errors.New("find-err"))
			})

			It("fails with error", func() {
				err := req.Execute()
				Expect(err.Error()).To(ContainSubstring("find-err"))
			})
		})
	})
})

package quotas_test

import (
	"net/http"

	"github.com/theophoric/cf-cli/cf/api/quotas"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/models"

	"github.com/theophoric/cf-cli/testhelpers/cloud_controller_gateway"
	testconfig "github.com/theophoric/cf-cli/testhelpers/configuration"

	"github.com/onsi/gomega/ghttp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CloudControllerQuotaRepository", func() {
	var (
		ccServer   *ghttp.Server
		configRepo core_config.ReadWriter
		repo       quotas.CloudControllerQuotaRepository
	)

	BeforeEach(func() {
		ccServer = ghttp.NewServer()
		configRepo = testconfig.NewRepositoryWithDefaults()
		configRepo.SetApiEndpoint(ccServer.URL())
		gateway := cloud_controller_gateway.NewTestCloudControllerGateway(configRepo)
		repo = quotas.NewCloudControllerQuotaRepository(configRepo, gateway)
	})

	AfterEach(func() {
		ccServer.Close()
	})

	Describe("FindByName", func() {
		BeforeEach(func() {
			ccServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v2/quota_definitions"),
					ghttp.RespondWith(http.StatusOK, `{
						"next_url": "/v2/quota_definitions?page=2",
						"resources": [
							{
								"metadata": { "guid": "my-quota-guid" },
								"entity": {
									"name": "my-remote-quota",
									"memory_limit": 1024,
									"instance_memory_limit": -1,
									"total_routes": 123,
									"total_services": 321,
									"non_basic_services_allowed": true,
									"app_instance_limit": 7
								}
							}
						]
					}`),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v2/quota_definitions", "page=2"),
					ghttp.RespondWith(http.StatusOK, `{
						"resources": [
							{
								"metadata": { "guid": "my-quota-guid2" },
								"entity": { "name": "my-remote-quota2", "memory_limit": 1024 }
							},
							{
								"metadata": { "guid": "my-quota-guid3" },
								"entity": { "name": "my-remote-quota3", "memory_limit": 1024 }
							}
						]
					}`),
				),
			)
		})

		It("Finds Quota definitions by name", func() {
			quota, err := repo.FindByName("my-remote-quota")
			Expect(err).NotTo(HaveOccurred())
			Expect(ccServer.ReceivedRequests()).To(HaveLen(2))
			Expect(quota).To(Equal(models.QuotaFields{
				Guid:                    "my-quota-guid",
				Name:                    "my-remote-quota",
				MemoryLimit:             1024,
				InstanceMemoryLimit:     -1,
				RoutesLimit:             123,
				ServicesLimit:           321,
				NonBasicServicesAllowed: true,
				AppInstanceLimit:        7,
			}))
		})
	})

	Describe("FindAll", func() {
		var (
			quotas []models.QuotaFields
			err    error
		)

		BeforeEach(func() {
			ccServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v2/quota_definitions"),
					ghttp.RespondWith(http.StatusOK, `{
						"next_url": "/v2/quota_definitions?page=2",
						"resources": [
							{
								"metadata": { "guid": "my-quota-guid" },
								"entity": {
									"name": "my-remote-quota",
									"memory_limit": 1024,
									"instance_memory_limit": -1,
									"total_routes": 123,
									"total_services": 321,
									"non_basic_services_allowed": true,
									"app_instance_limit": 7
								}
							}
						]
					}`),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v2/quota_definitions", "page=2"),
					ghttp.RespondWith(http.StatusOK, `{
						"resources": [
							{
								"metadata": { "guid": "my-quota-guid2" },
								"entity": { "name": "my-remote-quota2", "memory_limit": 1024 }
							},
							{
								"metadata": { "guid": "my-quota-guid3" },
								"entity": { "name": "my-remote-quota3", "memory_limit": 1024 }
							}
						]
					}`),
				),
			)
		})

		JustBeforeEach(func() {
			quotas, err = repo.FindAll()
			Expect(err).NotTo(HaveOccurred())
		})

		It("finds all Quota definitions", func() {
			Expect(ccServer.ReceivedRequests()).To(HaveLen(2))
			Expect(quotas).To(HaveLen(3))
			Expect(quotas[0].Guid).To(Equal("my-quota-guid"))
			Expect(quotas[0].Name).To(Equal("my-remote-quota"))
			Expect(quotas[0].MemoryLimit).To(Equal(int64(1024)))
			Expect(quotas[0].RoutesLimit).To(Equal(123))
			Expect(quotas[0].ServicesLimit).To(Equal(321))
			Expect(quotas[0].AppInstanceLimit).To(Equal(7))

			Expect(quotas[1].Guid).To(Equal("my-quota-guid2"))
			Expect(quotas[2].Guid).To(Equal("my-quota-guid3"))
		})

		It("defaults missing app instance limit to -1 (unlimited)", func() {
			Expect(quotas[1].AppInstanceLimit).To(Equal(-1))
		})
	})

	Describe("AssignQuotaToOrg", func() {
		BeforeEach(func() {
			ccServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/v2/organizations/my-org-guid"),
					ghttp.VerifyJSON(`{"quota_definition_guid":"my-quota-guid"}`),
					ghttp.RespondWith(http.StatusCreated, nil),
				),
			)
		})

		It("sets the quota for an organization", func() {
			err := repo.AssignQuotaToOrg("my-org-guid", "my-quota-guid")
			Expect(ccServer.ReceivedRequests()).To(HaveLen(1))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Create", func() {
		BeforeEach(func() {
			ccServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/v2/quota_definitions"),
					ghttp.VerifyJSON(`{
					"name": "not-so-strict",
					"non_basic_services_allowed": false,
					"total_services": 1,
					"total_routes": 12,
					"memory_limit": 123,
					"app_instance_limit": 42,
					"instance_memory_limit": 0
				}`),
					ghttp.RespondWith(http.StatusCreated, nil),
				),
			)
		})

		It("creates a new quota with the given name", func() {
			quota := models.QuotaFields{
				Name:             "not-so-strict",
				ServicesLimit:    1,
				RoutesLimit:      12,
				MemoryLimit:      123,
				AppInstanceLimit: 42,
			}

			err := repo.Create(quota)
			Expect(err).NotTo(HaveOccurred())
			Expect(ccServer.ReceivedRequests()).To(HaveLen(1))
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			ccServer.AppendHandlers(
				ghttp.CombineHandlers(

					ghttp.VerifyRequest("PUT", "/v2/quota_definitions/my-quota-guid"),
					ghttp.VerifyJSON(`{
					"guid": "my-quota-guid",
					"non_basic_services_allowed": false,
					"name": "amazing-quota",
					"total_services": 1,
					"total_routes": 12,
					"memory_limit": 123,
					"app_instance_limit": 42,
					"instance_memory_limit": 0
				}`),
					ghttp.RespondWith(http.StatusOK, nil),
				),
			)
		})

		It("updates an existing quota", func() {
			quota := models.QuotaFields{
				Guid:             "my-quota-guid",
				Name:             "amazing-quota",
				ServicesLimit:    1,
				RoutesLimit:      12,
				MemoryLimit:      123,
				AppInstanceLimit: 42,
			}

			err := repo.Update(quota)
			Expect(err).NotTo(HaveOccurred())
			Expect(ccServer.ReceivedRequests()).To(HaveLen(1))
		})
	})

	Describe("Delete", func() {
		BeforeEach(func() {
			ccServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/v2/quota_definitions/my-quota-guid"),
					ghttp.RespondWith(http.StatusNoContent, nil),
				),
			)
		})

		It("deletes the quota with the given name", func() {
			err := repo.Delete("my-quota-guid")
			Expect(err).NotTo(HaveOccurred())
			Expect(ccServer.ReceivedRequests()).To(HaveLen(1))
		})
	})
})

package models_test

import (
	"github.com/theophoric/cf-cli/cf/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RouteSummary", func() {
	Describe("URL", func() {
		var (
			r    models.RouteSummary
			host string
			path string
		)

		AfterEach(func() {
			host = ""
			path = ""
		})

		JustBeforeEach(func() {
			r = models.RouteSummary{
				Host: host,
				Domain: models.DomainFields{
					Name: "the-domain",
				},
				Path: path,
			}
		})

		Context("when the host is blank", func() {
			BeforeEach(func() {
				host = ""
			})

			It("returns the domain", func() {
				Expect(r.URL()).To(Equal("the-domain"))
			})

			Context("when the path is present", func() {
				BeforeEach(func() {
					path = "/the-path"
				})

				It("returns the domain and path", func() {
					Expect(r.URL()).To(Equal("the-domain/the-path"))
				})
			})
		})

		Context("when the host is not blank", func() {
			BeforeEach(func() {
				host = "the-host"
			})

			It("returns the host and domain", func() {
				Expect(r.URL()).To(Equal("the-host.the-domain"))
			})

			Context("when the path is present", func() {
				BeforeEach(func() {
					path = "/the-path"
				})

				It("returns the host and domain and path", func() {
					Expect(r.URL()).To(Equal("the-host.the-domain/the-path"))
				})
			})
		})
	})
})

package domain_test

import (
	"github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/testhelpers/configuration"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDomain(t *testing.T) {
	config := configuration.NewRepositoryWithDefaults()
	i18n.T = i18n.Init(config)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Domain Suite")
}

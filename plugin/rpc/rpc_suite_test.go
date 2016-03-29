package rpc_test

import (
	"github.com/theophoric/cf-cli/plugin/rpc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var rpcService *rpc.CliRpcService

func TestRpc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rpc Suite")
}

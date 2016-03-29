package net

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/errors"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/cf/trace"
)

type errorResponse struct {
	Name    string
	Message string
}

func errorHandler(statusCode int, body []byte) error {
	response := errorResponse{}
	err := json.Unmarshal(body, &response)
	if err != nil {
		return errors.NewHttpError(http.StatusInternalServerError, "", "")
	}

	return errors.NewHttpError(statusCode, response.Name, response.Message)
}

func NewRoutingApiGateway(config core_config.Reader, clock func() time.Time, ui terminal.UI, logger trace.Printer) Gateway {
	gateway := newGateway(errorHandler, config, ui, logger)
	gateway.Clock = clock
	gateway.PollingEnabled = true
	return gateway
}

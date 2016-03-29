package net

import (
	"encoding/json"

	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/errors"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/cf/trace"
)

type uaaErrorResponse struct {
	Code        string `json:"error"`
	Description string `json:"error_description"`
}

var uaaErrorHandler = func(statusCode int, body []byte) error {
	response := uaaErrorResponse{}
	json.Unmarshal(body, &response)

	if response.Code == "invalid_token" {
		return errors.NewInvalidTokenError(response.Description)
	}

	return errors.NewHttpError(statusCode, response.Code, response.Description)
}

func NewUAAGateway(config core_config.Reader, ui terminal.UI, logger trace.Printer) Gateway {
	return newGateway(uaaErrorHandler, config, ui, logger)
}

package service

import (
	"github.com/theophoric/cf-cli/cf/api"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/errors"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/flags"
)

type DeleteService struct {
	ui                 terminal.UI
	config             core_config.Reader
	serviceRepo        api.ServiceRepository
	serviceInstanceReq requirements.ServiceInstanceRequirement
}

func init() {
	command_registry.Register(&DeleteService{})
}

func (cmd *DeleteService) MetaData() command_registry.CommandMetadata {
	fs := make(map[string]flags.FlagSet)
	fs["f"] = &flags.BoolFlag{ShortName: "f", Usage: T("Force deletion without confirmation")}

	return command_registry.CommandMetadata{
		Name:        "delete-service",
		ShortName:   "ds",
		Description: T("Delete a service instance"),
		Usage: []string{
			T("CF_NAME delete-service SERVICE_INSTANCE [-f]"),
		},
		Flags: fs,
	}
}

func (cmd *DeleteService) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	if len(fc.Args()) != 1 {
		cmd.ui.Failed(T("Incorrect Usage. Requires an argument\n\n") + command_registry.Commands.CommandUsage("delete-service"))
	}

	reqs := []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
	}

	return reqs
}

func (cmd *DeleteService) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.serviceRepo = deps.RepoLocator.GetServiceRepository()
	return cmd
}

func (cmd *DeleteService) Execute(c flags.FlagContext) {
	serviceName := c.Args()[0]

	if !c.Bool("f") {
		if !cmd.ui.ConfirmDelete(T("service"), serviceName) {
			return
		}
	}

	cmd.ui.Say(T("Deleting service {{.ServiceName}} in org {{.OrgName}} / space {{.SpaceName}} as {{.CurrentUser}}...",
		map[string]interface{}{
			"ServiceName": terminal.EntityNameColor(serviceName),
			"OrgName":     terminal.EntityNameColor(cmd.config.OrganizationFields().Name),
			"SpaceName":   terminal.EntityNameColor(cmd.config.SpaceFields().Name),
			"CurrentUser": terminal.EntityNameColor(cmd.config.Username()),
		}))

	instance, apiErr := cmd.serviceRepo.FindInstanceByName(serviceName)

	switch apiErr.(type) {
	case nil:
	case *errors.ModelNotFoundError:
		cmd.ui.Ok()
		cmd.ui.Warn(T("Service {{.ServiceName}} does not exist.", map[string]interface{}{"ServiceName": serviceName}))
		return
	default:
		cmd.ui.Failed(apiErr.Error())
		return
	}

	apiErr = cmd.serviceRepo.DeleteService(instance)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	apiErr = printSuccessMessageForServiceInstance(serviceName, cmd.serviceRepo, cmd.ui)
	if apiErr != nil {
		cmd.ui.Ok()
	}
}

package user

import (
	"github.com/theophoric/cf-cli/cf"
	"github.com/theophoric/cf-cli/cf/api"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/errors"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/flags"
)

type CreateUser struct {
	ui       terminal.UI
	config   core_config.Reader
	userRepo api.UserRepository
}

func init() {
	command_registry.Register(&CreateUser{})
}

func (cmd *CreateUser) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "create-user",
		Description: T("Create a new user"),
		Usage: []string{
			T("CF_NAME create-user USERNAME PASSWORD"),
		},
	}
}

func (cmd *CreateUser) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	if len(fc.Args()) != 2 {
		usage := command_registry.Commands.CommandUsage("create-user")
		cmd.ui.Failed(T("Incorrect Usage. Requires arguments\n\n") + usage)
	}

	reqs := []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
	}

	return reqs
}

func (cmd *CreateUser) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.userRepo = deps.RepoLocator.GetUserRepository()
	return cmd
}

func (cmd *CreateUser) Execute(c flags.FlagContext) {
	username := c.Args()[0]
	password := c.Args()[1]

	cmd.ui.Say(T("Creating user {{.TargetUser}}...",
		map[string]interface{}{
			"TargetUser":  terminal.EntityNameColor(username),
			"CurrentUser": terminal.EntityNameColor(cmd.config.Username()),
		}))

	err := cmd.userRepo.Create(username, password)
	switch err.(type) {
	case nil:
	case *errors.ModelAlreadyExistsError:
		cmd.ui.Warn("%s", err.Error())
	default:
		cmd.ui.Failed(T("Error creating user {{.TargetUser}}.\n{{.Error}}",
			map[string]interface{}{
				"TargetUser": terminal.EntityNameColor(username),
				"Error":      err.Error(),
			}))
	}

	cmd.ui.Ok()
	cmd.ui.Say(T("\nTIP: Assign roles with '{{.CurrentUser}} set-org-role' and '{{.CurrentUser}} set-space-role'", map[string]interface{}{"CurrentUser": cf.Name}))
}

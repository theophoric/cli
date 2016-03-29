package securitygroup

import (
	"github.com/theophoric/cf-cli/cf/api/security_groups/defaults/running"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/flags"
)

type listRunningSecurityGroups struct {
	ui                       terminal.UI
	runningSecurityGroupRepo running.RunningSecurityGroupsRepo
	configRepo               core_config.Reader
}

func init() {
	command_registry.Register(&listRunningSecurityGroups{})
}

func (cmd *listRunningSecurityGroups) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "running-security-groups",
		Description: T("List security groups in the set of security groups for running applications"),
		Usage: []string{
			"CF_NAME running-security-groups",
		},
	}
}

func (cmd *listRunningSecurityGroups) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	usageReq := requirements.NewUsageRequirement(command_registry.CliCommandUsagePresenter(cmd),
		T("No argument required"),
		func() bool {
			return len(fc.Args()) != 0
		},
	)

	reqs := []requirements.Requirement{
		usageReq,
		requirementsFactory.NewLoginRequirement(),
	}
	return reqs
}

func (cmd *listRunningSecurityGroups) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.configRepo = deps.Config
	cmd.runningSecurityGroupRepo = deps.RepoLocator.GetRunningSecurityGroupsRepository()
	return cmd
}

func (cmd *listRunningSecurityGroups) Execute(context flags.FlagContext) {
	cmd.ui.Say(T("Acquiring running security groups as '{{.username}}'", map[string]interface{}{
		"username": terminal.EntityNameColor(cmd.configRepo.Username()),
	}))

	defaultSecurityGroupsFields, err := cmd.runningSecurityGroupRepo.List()
	if err != nil {
		cmd.ui.Failed(err.Error())
	}

	cmd.ui.Ok()
	cmd.ui.Say("")

	if len(defaultSecurityGroupsFields) > 0 {
		for _, value := range defaultSecurityGroupsFields {
			cmd.ui.Say(value.Name)
		}
	} else {
		cmd.ui.Say(T("No running security groups set"))
	}
}

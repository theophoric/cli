package securitygroup

import (
	"fmt"

	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/flags"

	"github.com/theophoric/cf-cli/cf/api/security_groups"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/models"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
)

type SecurityGroups struct {
	ui                terminal.UI
	securityGroupRepo security_groups.SecurityGroupRepo
	configRepo        core_config.Reader
}

func init() {
	command_registry.Register(&SecurityGroups{})
}

func (cmd *SecurityGroups) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "security-groups",
		Description: T("List all security groups"),
		Usage: []string{
			"CF_NAME security-groups",
		},
	}
}

func (cmd *SecurityGroups) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
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

func (cmd *SecurityGroups) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.configRepo = deps.Config
	cmd.securityGroupRepo = deps.RepoLocator.GetSecurityGroupRepository()
	return cmd
}

func (cmd *SecurityGroups) Execute(c flags.FlagContext) {
	cmd.ui.Say(T("Getting security groups as {{.username}}",
		map[string]interface{}{
			"username": terminal.EntityNameColor(cmd.configRepo.Username()),
		}))

	securityGroups, err := cmd.securityGroupRepo.FindAll()
	if err != nil {
		cmd.ui.Failed(err.Error())
	}

	cmd.ui.Ok()
	cmd.ui.Say("")

	if len(securityGroups) == 0 {
		cmd.ui.Say(T("No security groups"))
		return
	}

	table := terminal.NewTable(cmd.ui, []string{"", T("Name"), T("Organization"), T("Space")})

	for index, securityGroup := range securityGroups {
		if len(securityGroup.Spaces) > 0 {
			cmd.printSpaces(table, securityGroup, index)
		} else {
			table.Add(fmt.Sprintf("#%d", index), securityGroup.Name, "", "")
		}
	}
	table.Print()
}

func (cmd SecurityGroups) printSpaces(table terminal.Table, securityGroup models.SecurityGroup, index int) {
	outputted_index := false

	for _, space := range securityGroup.Spaces {
		if !outputted_index {
			table.Add(fmt.Sprintf("#%d", index), securityGroup.Name, space.Organization.Name, space.Name)
			outputted_index = true
		} else {
			table.Add("", securityGroup.Name, space.Organization.Name, space.Name)
		}
	}
}

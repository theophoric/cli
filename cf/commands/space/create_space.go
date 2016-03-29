package space

import (
	"github.com/theophoric/cf-cli/cf"
	"github.com/theophoric/cf-cli/cf/api"
	"github.com/theophoric/cf-cli/cf/api/organizations"
	"github.com/theophoric/cf-cli/cf/api/space_quotas"
	"github.com/theophoric/cf-cli/cf/api/spaces"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/commands/user"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/errors"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/cf/models"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/flags"
)

type CreateSpace struct {
	ui              terminal.UI
	config          core_config.Reader
	spaceRepo       spaces.SpaceRepository
	orgRepo         organizations.OrganizationRepository
	userRepo        api.UserRepository
	spaceRoleSetter user.SpaceRoleSetter
	spaceQuotaRepo  space_quotas.SpaceQuotaRepository
}

func init() {
	command_registry.Register(&CreateSpace{})
}

func (cmd *CreateSpace) MetaData() command_registry.CommandMetadata {
	fs := make(map[string]flags.FlagSet)
	fs["o"] = &flags.StringFlag{ShortName: "o", Usage: T("Organization")}
	fs["q"] = &flags.StringFlag{ShortName: "q", Usage: T("Quota to assign to the newly created space")}

	return command_registry.CommandMetadata{
		Name:        "create-space",
		Description: T("Create a space"),
		Usage: []string{
			T("CF_NAME create-space SPACE [-o ORG] [-q SPACE-QUOTA]"),
		},
		Flags: fs,
	}
}

func (cmd *CreateSpace) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	if len(fc.Args()) != 1 {
		cmd.ui.Failed(T("Incorrect Usage. Requires an argument\n\n") + command_registry.Commands.CommandUsage("create-space"))
	}

	reqs := []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
	}

	if fc.String("o") == "" {
		reqs = append(reqs, requirementsFactory.NewTargetedOrgRequirement())
	}

	return reqs
}

func (cmd *CreateSpace) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.spaceRepo = deps.RepoLocator.GetSpaceRepository()
	cmd.orgRepo = deps.RepoLocator.GetOrganizationRepository()
	cmd.userRepo = deps.RepoLocator.GetUserRepository()
	cmd.spaceQuotaRepo = deps.RepoLocator.GetSpaceQuotaRepository()

	//get command from registry for dependency
	commandDep := command_registry.Commands.FindCommand("set-space-role")
	commandDep = commandDep.SetDependency(deps, false)
	cmd.spaceRoleSetter = commandDep.(user.SpaceRoleSetter)

	return cmd
}

func (cmd *CreateSpace) Execute(c flags.FlagContext) {
	spaceName := c.Args()[0]
	orgName := c.String("o")
	spaceQuotaName := c.String("q")
	orgGuid := ""
	if orgName == "" {
		orgName = cmd.config.OrganizationFields().Name
		orgGuid = cmd.config.OrganizationFields().Guid
	}

	cmd.ui.Say(T("Creating space {{.SpaceName}} in org {{.OrgName}} as {{.CurrentUser}}...",
		map[string]interface{}{
			"SpaceName":   terminal.EntityNameColor(spaceName),
			"OrgName":     terminal.EntityNameColor(orgName),
			"CurrentUser": terminal.EntityNameColor(cmd.config.Username()),
		}))

	if orgGuid == "" {
		org, apiErr := cmd.orgRepo.FindByName(orgName)
		switch apiErr.(type) {
		case nil:
		case *errors.ModelNotFoundError:
			cmd.ui.Failed(T("Org {{.OrgName}} does not exist or is not accessible", map[string]interface{}{"OrgName": orgName}))
			return
		default:
			cmd.ui.Failed(T("Error finding org {{.OrgName}}\n{{.ErrorDescription}}",
				map[string]interface{}{
					"OrgName":          orgName,
					"ErrorDescription": apiErr.Error(),
				}))
			return
		}

		orgGuid = org.Guid
	}

	var spaceQuotaGuid string
	if spaceQuotaName != "" {
		spaceQuota, err := cmd.spaceQuotaRepo.FindByNameAndOrgGuid(spaceQuotaName, orgGuid)
		if err != nil {
			cmd.ui.Failed(err.Error())
		}
		spaceQuotaGuid = spaceQuota.Guid
	}

	space, err := cmd.spaceRepo.Create(spaceName, orgGuid, spaceQuotaGuid)
	if err != nil {
		if httpErr, ok := err.(errors.HttpError); ok && httpErr.ErrorCode() == errors.SpaceNameTaken {
			cmd.ui.Ok()
			cmd.ui.Warn(T("Space {{.SpaceName}} already exists", map[string]interface{}{"SpaceName": spaceName}))
			return
		}
		cmd.ui.Failed(err.Error())
		return
	}
	cmd.ui.Ok()

	err = cmd.spaceRoleSetter.SetSpaceRole(space, models.SPACE_MANAGER, cmd.config.UserGuid(), cmd.config.Username())
	if err != nil {
		cmd.ui.Failed(err.Error())
		return
	}

	err = cmd.spaceRoleSetter.SetSpaceRole(space, models.SPACE_DEVELOPER, cmd.config.UserGuid(), cmd.config.Username())
	if err != nil {
		cmd.ui.Failed(err.Error())
		return
	}

	cmd.ui.Say(T("\nTIP: Use '{{.CFTargetCommand}}' to target new space",
		map[string]interface{}{
			"CFTargetCommand": terminal.CommandColor(cf.Name + " target -o \"" + orgName + "\" -s \"" + space.Name + "\""),
		}))
}

package organization

import (
	"github.com/theophoric/cf-cli/cf/api/organizations"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/errors"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/cf/models"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/flags"
)

type DeleteOrg struct {
	ui      terminal.UI
	config  core_config.ReadWriter
	orgRepo organizations.OrganizationRepository
	orgReq  requirements.OrganizationRequirement
}

func init() {
	command_registry.Register(&DeleteOrg{})
}

func (cmd *DeleteOrg) MetaData() command_registry.CommandMetadata {
	fs := make(map[string]flags.FlagSet)
	fs["f"] = &flags.BoolFlag{ShortName: "f", Usage: T("Force deletion without confirmation")}

	return command_registry.CommandMetadata{
		Name:        "delete-org",
		Description: T("Delete an org"),
		Usage: []string{
			T("CF_NAME delete-org ORG [-f]"),
		},
		Flags: fs,
	}
}

func (cmd *DeleteOrg) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	if len(fc.Args()) != 1 {
		cmd.ui.Failed(T("Incorrect Usage. Requires an argument\n\n") + command_registry.Commands.CommandUsage("delete-org"))
	}

	reqs := []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
	}

	return reqs
}

func (cmd *DeleteOrg) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.orgRepo = deps.RepoLocator.GetOrganizationRepository()
	return cmd
}

func (cmd *DeleteOrg) Execute(c flags.FlagContext) {
	orgName := c.Args()[0]

	if !c.Bool("f") {
		if !cmd.ui.ConfirmDeleteWithAssociations(T("org"), orgName) {
			return
		}
	}

	cmd.ui.Say(T("Deleting org {{.OrgName}} as {{.Username}}...",
		map[string]interface{}{
			"OrgName":  terminal.EntityNameColor(orgName),
			"Username": terminal.EntityNameColor(cmd.config.Username())}))

	org, apiErr := cmd.orgRepo.FindByName(orgName)

	switch apiErr.(type) {
	case nil:
	case *errors.ModelNotFoundError:
		cmd.ui.Ok()
		cmd.ui.Warn(T("Org {{.OrgName}} does not exist.",
			map[string]interface{}{"OrgName": orgName}))
		return
	default:
		cmd.ui.Failed(apiErr.Error())
		return
	}

	apiErr = cmd.orgRepo.Delete(org.Guid)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	if org.Guid == cmd.config.OrganizationFields().Guid {
		cmd.config.SetOrganizationFields(models.OrganizationFields{})
		cmd.config.SetSpaceFields(models.SpaceFields{})
	}

	cmd.ui.Ok()
	return
}

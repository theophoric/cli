package organization

import (
	"github.com/theophoric/cf-cli/cf/api/organizations"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/cf/models"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/flags"
	"github.com/theophoric/cf-cli/plugin/models"
)

const orgLimit = 0

type ListOrgs struct {
	ui              terminal.UI
	config          core_config.Reader
	orgRepo         organizations.OrganizationRepository
	pluginOrgsModel *[]plugin_models.GetOrgs_Model
	pluginCall      bool
}

func init() {
	command_registry.Register(&ListOrgs{})
}

func (cmd *ListOrgs) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "orgs",
		ShortName:   "o",
		Description: T("List all orgs"),
		Usage: []string{
			"CF_NAME orgs",
		},
	}
}

func (cmd *ListOrgs) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
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

func (cmd *ListOrgs) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.orgRepo = deps.RepoLocator.GetOrganizationRepository()
	cmd.pluginOrgsModel = deps.PluginModels.Organizations
	cmd.pluginCall = pluginCall
	return cmd
}

func (cmd ListOrgs) Execute(fc flags.FlagContext) {
	cmd.ui.Say(T("Getting orgs as {{.Username}}...\n",
		map[string]interface{}{"Username": terminal.EntityNameColor(cmd.config.Username())}))

	noOrgs := true
	table := cmd.ui.Table([]string{T("name")})

	orgs, apiErr := cmd.orgRepo.ListOrgs(orgLimit)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
	}
	for _, org := range orgs {
		table.Add(org.Name)
		noOrgs = false
	}

	table.Print()

	if apiErr != nil {
		cmd.ui.Failed(T("Failed fetching orgs.\n{{.ApiErr}}",
			map[string]interface{}{"ApiErr": apiErr}))
		return
	}

	if noOrgs {
		cmd.ui.Say(T("No orgs found"))
	}

	if cmd.pluginCall {
		cmd.populatePluginModel(orgs)
	}

}

func (cmd *ListOrgs) populatePluginModel(orgs []models.Organization) {
	for _, org := range orgs {
		orgModel := plugin_models.GetOrgs_Model{}
		orgModel.Name = org.Name
		orgModel.Guid = org.Guid
		*(cmd.pluginOrgsModel) = append(*(cmd.pluginOrgsModel), orgModel)
	}
}

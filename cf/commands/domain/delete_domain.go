package domain

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

type DeleteDomain struct {
	ui         terminal.UI
	config     core_config.Reader
	orgReq     requirements.TargetedOrgRequirement
	domainRepo api.DomainRepository
}

func init() {
	command_registry.Register(&DeleteDomain{})
}

func (cmd *DeleteDomain) MetaData() command_registry.CommandMetadata {
	fs := make(map[string]flags.FlagSet)
	fs["f"] = &flags.BoolFlag{ShortName: "f", Usage: T("Force deletion without confirmation")}

	return command_registry.CommandMetadata{
		Name:        "delete-domain",
		Description: T("Delete a domain"),
		Usage: []string{
			T("CF_NAME delete-domain DOMAIN [-f]"),
		},
		Flags: fs,
	}
}

func (cmd *DeleteDomain) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	if len(fc.Args()) != 1 {
		cmd.ui.Failed(T("Incorrect Usage. Requires an argument\n\n") + command_registry.Commands.CommandUsage("delete-domain"))
	}

	loginReq := requirementsFactory.NewLoginRequirement()
	cmd.orgReq = requirementsFactory.NewTargetedOrgRequirement()

	reqs := []requirements.Requirement{
		loginReq,
		cmd.orgReq,
	}

	return reqs
}

func (cmd *DeleteDomain) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.domainRepo = deps.RepoLocator.GetDomainRepository()
	return cmd
}

func (cmd *DeleteDomain) Execute(c flags.FlagContext) {
	domainName := c.Args()[0]
	domain, apiErr := cmd.domainRepo.FindByNameInOrg(domainName, cmd.orgReq.GetOrganizationFields().Guid)

	switch apiErr.(type) {
	case nil:
		if domain.Shared {
			cmd.ui.Failed(T("domain {{.DomainName}} is a shared domain, not an owned domain.",
				map[string]interface{}{
					"DomainName": domainName}))
			return
		}
	case *errors.ModelNotFoundError:
		cmd.ui.Ok()
		cmd.ui.Warn(apiErr.Error())
		return
	default:
		cmd.ui.Failed(T("Error finding domain {{.DomainName}}\n{{.ApiErr}}",
			map[string]interface{}{"DomainName": domainName, "ApiErr": apiErr.Error()}))
		return
	}

	if !c.Bool("f") {
		if !cmd.ui.ConfirmDelete(T("domain"), domainName) {
			return
		}
	}

	cmd.ui.Say(T("Deleting domain {{.DomainName}} as {{.Username}}...",
		map[string]interface{}{
			"DomainName": terminal.EntityNameColor(domainName),
			"Username":   terminal.EntityNameColor(cmd.config.Username())}))

	apiErr = cmd.domainRepo.Delete(domain.Guid)
	if apiErr != nil {
		cmd.ui.Failed(T("Error deleting domain {{.DomainName}}\n{{.ApiErr}}",
			map[string]interface{}{"DomainName": domainName, "ApiErr": apiErr.Error()}))
		return
	}

	cmd.ui.Ok()
}

package route

import (
	"github.com/blang/semver"
	"github.com/theophoric/cf-cli/cf/api"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/errors"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/flags"
)

type DeleteRoute struct {
	ui        terminal.UI
	config    core_config.Reader
	routeRepo api.RouteRepository
	domainReq requirements.DomainRequirement
}

func init() {
	command_registry.Register(&DeleteRoute{})
}

func (cmd *DeleteRoute) MetaData() command_registry.CommandMetadata {
	fs := make(map[string]flags.FlagSet)
	fs["f"] = &flags.BoolFlag{ShortName: "f", Usage: T("Force deletion without confirmation")}
	fs["hostname"] = &flags.StringFlag{Name: "hostname", ShortName: "n", Usage: T("Hostname used to identify the route")}
	fs["path"] = &flags.StringFlag{Name: "path", Usage: T("Path used to identify the route")}

	return command_registry.CommandMetadata{
		Name:        "delete-route",
		Description: T("Delete a route"),
		Usage: []string{
			T("CF_NAME delete-route DOMAIN [--hostname HOSTNAME] [--path PATH] [-f]"),
		},
		Examples: []string{
			"CF_NAME delete-route example.com                              # example.com",
			"CF_NAME delete-route example.com --hostname myhost            # myhost.example.com",
			"CF_NAME delete-route example.com --hostname myhost --path foo # myhost.example.com/foo",
		},
		Flags: fs,
	}
}

func (cmd *DeleteRoute) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	if len(fc.Args()) != 1 {
		cmd.ui.Failed(T("Incorrect Usage. Requires an argument\n\n") + command_registry.Commands.CommandUsage("delete-route"))
	}

	cmd.domainReq = requirementsFactory.NewDomainRequirement(fc.Args()[0])

	requiredVersion, err := semver.Make("2.36.0")
	if err != nil {
		panic(err.Error())
	}

	var reqs []requirements.Requirement

	if fc.String("path") != "" {
		reqs = append(reqs, requirementsFactory.NewMinAPIVersionRequirement("Option '--path'", requiredVersion))
	}

	reqs = append(reqs, []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
		cmd.domainReq,
	}...)

	return reqs
}

func (cmd *DeleteRoute) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.routeRepo = deps.RepoLocator.GetRouteRepository()
	return cmd
}

func (cmd *DeleteRoute) Execute(c flags.FlagContext) {
	var port int

	host := c.String("n")
	path := c.String("path")
	domainName := c.Args()[0]

	url := domainName
	if host != "" {
		url = host + "." + domainName
	}
	if !c.Bool("f") {
		if !cmd.ui.ConfirmDelete("route", url) {
			return
		}
	}

	cmd.ui.Say(T("Deleting route {{.URL}}...", map[string]interface{}{"URL": terminal.EntityNameColor(url)}))

	domain := cmd.domainReq.GetDomain()
	route, err := cmd.routeRepo.Find(host, domain, path, port)
	if err != nil {
		if _, ok := err.(*errors.ModelNotFoundError); ok {
			cmd.ui.Warn(T("Unable to delete, route '{{.URL}}' does not exist.",
				map[string]interface{}{"URL": url}))
			return
		}
		cmd.ui.Failed(err.Error())
		return
	}

	err = cmd.routeRepo.Delete(route.Guid)
	if err != nil {
		cmd.ui.Failed(err.Error())
		return
	}

	cmd.ui.Ok()
}

package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/theophoric/cf-cli/cf"
	"github.com/theophoric/cf-cli/cf/actors/plan_builder"
	"github.com/theophoric/cf-cli/cf/api"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/cf/models"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/cf/ui_helpers"
	"github.com/theophoric/cf-cli/flags"
	"github.com/theophoric/cf-cli/json"
)

type UpdateService struct {
	ui          terminal.UI
	config      core_config.Reader
	serviceRepo api.ServiceRepository
	planBuilder plan_builder.PlanBuilder
}

func init() {
	command_registry.Register(&UpdateService{})
}

func (cmd *UpdateService) MetaData() command_registry.CommandMetadata {
	baseUsage := T("CF_NAME update-service SERVICE_INSTANCE [-p NEW_PLAN] [-c PARAMETERS_AS_JSON] [-t TAGS]")
	paramsUsage := T(`   Optionally provide service-specific configuration parameters in a valid JSON object in-line.
   CF_NAME update-service -c '{"name":"value","name":"value"}'

   Optionally provide a file containing service-specific configuration parameters in a valid JSON object. 
   The path to the parameters file can be an absolute or relative path to a file.
   CF_NAME update-service -c PATH_TO_FILE

   Example of valid JSON object:
   {
      "cluster_nodes": {
         "count": 5,
         "memory_mb": 1024
      }
   }`)
	tagsUsage := T(`   Optionally provide a list of comma-delimited tags that will be written to the VCAP_SERVICES environment variable for any bound applications.`)

	fs := make(map[string]flags.FlagSet)
	fs["p"] = &flags.StringFlag{ShortName: "p", Usage: T("Change service plan for a service instance")}
	fs["c"] = &flags.StringFlag{ShortName: "c", Usage: T("Valid JSON object containing service-specific configuration parameters, provided either in-line or in a file. For a list of supported configuration parameters, see documentation for the particular service offering.")}
	fs["t"] = &flags.StringFlag{ShortName: "t", Usage: T("User provided tags")}

	return command_registry.CommandMetadata{
		Name:        "update-service",
		Description: T("Update a service instance"),
		Usage: []string{
			baseUsage,
			"\n\n",
			paramsUsage,
			"\n\n",
			tagsUsage,
		},
		Examples: []string{
			`CF_NAME update-service mydb -p gold`,
			`CF_NAME update-service mydb -c '{"ram_gb":4}'`,
			`CF_NAME update-service mydb -c ~/workspace/tmp/instance_config.json`,
			`CF_NAME update-service mydb -t "list,of, tags"`,
		},
		Flags: fs,
	}
}

func (cmd *UpdateService) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	if len(fc.Args()) != 1 {
		cmd.ui.Failed(T("Incorrect Usage. Requires an argument\n\n") + command_registry.Commands.CommandUsage("update-service"))
	}

	reqs := []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
		requirementsFactory.NewTargetedSpaceRequirement(),
	}

	if fc.String("p") != "" {
		reqs = append(reqs, requirementsFactory.NewMinAPIVersionRequirement("Updating a plan", cf.UpdateServicePlanMinimumApiVersion))
	}

	return reqs
}

func (cmd *UpdateService) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.serviceRepo = deps.RepoLocator.GetServiceRepository()
	cmd.planBuilder = deps.PlanBuilder
	return cmd
}

func (cmd *UpdateService) Execute(c flags.FlagContext) {
	planName := c.String("p")
	params := c.String("c")

	tagsSet := c.IsSet("t")
	tagsList := c.String("t")

	if planName == "" && params == "" && tagsSet == false {
		cmd.ui.Ok()
		cmd.ui.Say(T("No changes were made"))
		return
	}

	serviceInstanceName := c.Args()[0]
	serviceInstance, err := cmd.serviceRepo.FindInstanceByName(serviceInstanceName)
	if err != nil {
		cmd.ui.Failed(err.Error())
	}

	paramsMap, err := json.ParseJsonFromFileOrString(params)
	if err != nil {
		cmd.ui.Failed(T("Invalid configuration provided for -c flag. Please provide a valid JSON object or path to a file containing a valid JSON object."))
	}

	tags := ui_helpers.ParseTags(tagsList)

	var plan models.ServicePlanFields
	if planName != "" {
		plan, err = cmd.findPlan(serviceInstance, planName)
		if err != nil {
			cmd.ui.Failed(err.Error())
		}
	}

	cmd.printUpdatingServiceInstanceMessage(serviceInstanceName)

	err = cmd.serviceRepo.UpdateServiceInstance(serviceInstance.Guid, plan.Guid, paramsMap, tags)
	if err != nil {
		cmd.ui.Failed(err.Error())
	}
	err = printSuccessMessageForServiceInstance(serviceInstanceName, cmd.serviceRepo, cmd.ui)
	if err != nil {
		cmd.ui.Failed(err.Error())
	}
}

func (cmd *UpdateService) findPlan(serviceInstance models.ServiceInstance, planName string) (plan models.ServicePlanFields, err error) {
	plans, err := cmd.planBuilder.GetPlansForServiceForOrg(serviceInstance.ServiceOffering.Guid, cmd.config.OrganizationFields().Name)
	if err != nil {
		return
	}

	for _, p := range plans {
		if p.Name == planName {
			plan = p
			return
		}
	}
	err = errors.New(T("Plan does not exist for the {{.ServiceName}} service",
		map[string]interface{}{"ServiceName": serviceInstance.ServiceOffering.Label}))
	return
}

func (cmd *UpdateService) printUpdatingServiceInstanceMessage(serviceInstanceName string) {
	cmd.ui.Say(T("Updating service instance {{.ServiceName}} as {{.UserName}}...",
		map[string]interface{}{
			"ServiceName": terminal.EntityNameColor(serviceInstanceName),
			"UserName":    terminal.EntityNameColor(cmd.config.Username()),
		}))
}

func printSuccessMessageForServiceInstance(serviceInstanceName string, serviceRepo api.ServiceRepository, ui terminal.UI) error {
	instance, apiErr := serviceRepo.FindInstanceByName(serviceInstanceName)
	if apiErr != nil {
		return apiErr
	}

	if instance.ServiceInstanceFields.LastOperation.State == "in progress" {
		ui.Ok()
		ui.Say("")
		ui.Say(T("{{.State}} in progress. Use '{{.ServicesCommand}}' or '{{.ServiceCommand}}' to check operation status.",
			map[string]interface{}{
				"State":           strings.Title(instance.ServiceInstanceFields.LastOperation.Type),
				"ServicesCommand": terminal.CommandColor("cf services"),
				"ServiceCommand":  terminal.CommandColor(fmt.Sprintf("cf service %s", serviceInstanceName)),
			}))
	} else {
		ui.Ok()
	}

	return nil
}

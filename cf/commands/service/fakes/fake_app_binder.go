package commands

import (
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/models"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/flags"
)

type FakeAppBinder struct {
	AppsToBind        []models.Application
	InstancesToBindTo []models.ServiceInstance
	Params            map[string]interface{}

	BindApplicationReturns struct {
		Error error
	}
}

func (binder *FakeAppBinder) BindApplication(app models.Application, service models.ServiceInstance, paramsMap map[string]interface{}) error {
	binder.AppsToBind = append(binder.AppsToBind, app)
	binder.InstancesToBindTo = append(binder.InstancesToBindTo, service)
	binder.Params = paramsMap

	return binder.BindApplicationReturns.Error
}

func (binder *FakeAppBinder) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{Name: "bind-service"}
}

func (binder *FakeAppBinder) SetDependency(_ command_registry.Dependency, _ bool) command_registry.Command {
	return binder
}

func (binder *FakeAppBinder) Requirements(_ requirements.Factory, _ flags.FlagContext) []requirements.Requirement {
	return []requirements.Requirement{}
}

func (binder *FakeAppBinder) Execute(_ flags.FlagContext) {}

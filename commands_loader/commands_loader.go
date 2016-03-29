package commands_loader

import (
	"github.com/theophoric/cf-cli/cf/commands"
	"github.com/theophoric/cf-cli/cf/commands/application"
	"github.com/theophoric/cf-cli/cf/commands/buildpack"
	"github.com/theophoric/cf-cli/cf/commands/domain"
	"github.com/theophoric/cf-cli/cf/commands/environmentvariablegroup"
	"github.com/theophoric/cf-cli/cf/commands/featureflag"
	"github.com/theophoric/cf-cli/cf/commands/organization"
	"github.com/theophoric/cf-cli/cf/commands/plugin"
	"github.com/theophoric/cf-cli/cf/commands/plugin_repo"
	"github.com/theophoric/cf-cli/cf/commands/quota"
	"github.com/theophoric/cf-cli/cf/commands/route"
	"github.com/theophoric/cf-cli/cf/commands/routergroups"
	"github.com/theophoric/cf-cli/cf/commands/securitygroup"
	"github.com/theophoric/cf-cli/cf/commands/service"
	"github.com/theophoric/cf-cli/cf/commands/serviceaccess"
	"github.com/theophoric/cf-cli/cf/commands/serviceauthtoken"
	"github.com/theophoric/cf-cli/cf/commands/servicebroker"
	"github.com/theophoric/cf-cli/cf/commands/servicekey"
	"github.com/theophoric/cf-cli/cf/commands/space"
	"github.com/theophoric/cf-cli/cf/commands/spacequota"
	"github.com/theophoric/cf-cli/cf/commands/user"
)

/*******************
This package make a reference to all the command packages
in cf/commands/..., so all init() in the directories will
get initialized

* Any new command packages must be included here for init() to get called
********************/

func Load() {
	_ = commands.Api{}
	_ = application.ListApps{}
	_ = buildpack.ListBuildpacks{}
	_ = domain.CreateDomain{}
	_ = environmentvariablegroup.RunningEnvironmentVariableGroup{}
	_ = featureflag.ShowFeatureFlag{}
	_ = organization.ListOrgs{}
	_ = plugin.Plugins{}
	_ = plugin_repo.RepoPlugins{}
	_ = quota.CreateQuota{}
	_ = route.CreateRoute{}
	_ = routergroups.RouterGroups{}
	_ = securitygroup.ShowSecurityGroup{}
	_ = service.ShowService{}
	_ = serviceauthtoken.ListServiceAuthTokens{}
	_ = serviceaccess.ServiceAccess{}
	_ = servicebroker.ListServiceBrokers{}
	_ = servicekey.ServiceKey{}
	_ = space.CreateSpace{}
	_ = spacequota.SpaceQuota{}
	_ = user.CreateUser{}
}

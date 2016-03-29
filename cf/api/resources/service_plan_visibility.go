package resources

import "github.com/theophoric/cf-cli/cf/models"

type ServicePlanVisibilityResource struct {
	Resource
	Entity models.ServicePlanVisibilityFields
}

func (resource ServicePlanVisibilityResource) ToFields() (fields models.ServicePlanVisibilityFields) {
	fields.Guid = resource.Metadata.Guid
	fields.ServicePlanGuid = resource.Entity.ServicePlanGuid
	fields.OrganizationGuid = resource.Entity.OrganizationGuid
	return
}

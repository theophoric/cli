package resources

import "github.com/theophoric/cf-cli/cf/models"

type DomainResource struct {
	Resource
	Entity DomainEntity
}

type DomainEntity struct {
	Name                   string   `json:"name"`
	OwningOrganizationGuid string   `json:"owning_organization_guid,omitempty"`
	SharedOrganizationsUrl string   `json:"shared_organizations_url,omitempty"`
	RouterGroupGuid        string   `json:"router_group_guid,omitempty"`
	RouterGroupTypes       []string `json:"router_group_types,omitempty"`
	Wildcard               bool     `json:"wildcard"`
}

func (resource DomainResource) ToFields() models.DomainFields {
	privateDomain := resource.Entity.SharedOrganizationsUrl != "" || resource.Entity.OwningOrganizationGuid != ""
	return models.DomainFields{
		Name: resource.Entity.Name,
		Guid: resource.Metadata.Guid,
		OwningOrganizationGuid: resource.Entity.OwningOrganizationGuid,
		Shared:                 !privateDomain,
		RouterGroupGuid:        resource.Entity.RouterGroupGuid,
		RouterGroupTypes:       resource.Entity.RouterGroupTypes,
	}
}

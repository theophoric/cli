package errors

import (
	. "github.com/theophoric/cf-cli/cf/i18n"
)

type ServiceAssociationError struct {
}

func NewServiceAssociationError() error {
	return &ServiceAssociationError{}
}

func (err *ServiceAssociationError) Error() string {
	return T("Cannot delete service instance, service keys and bindings must first be deleted")
}

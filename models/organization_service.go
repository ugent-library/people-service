package models

import "context"

type OrganizationService interface {
	SaveOrganization(context.Context, *Organization) (*Organization, error)
	CreateOrganization(context.Context, *Organization) (*Organization, error)
	UpdateOrganization(context.Context, *Organization) (*Organization, error)
	GetOrganization(context.Context, string) (*Organization, error)
	GetOrganizationsByIdentifier(context.Context, ...URN) ([]*Organization, error)
	DeleteOrganization(context.Context, string) error
	EachOrganization(context.Context, func(*Organization) bool) error
	GetOrganizations(context.Context) ([]*Organization, string, error)
	GetMoreOrganizations(context.Context, string) ([]*Organization, string, error)
}

package models

import (
	"time"
)

type OrganizationMember struct {
	ID          string     `json:"id,omitempty"`
	DateCreated *time.Time `json:"date_created,omitempty"`
	DateUpdated *time.Time `json:"date_updated,omitempty"`
}

func (om OrganizationMember) Dup() *OrganizationMember {
	return &OrganizationMember{
		ID:          om.ID,
		DateCreated: copyTime(om.DateCreated),
		DateUpdated: copyTime(om.DateUpdated),
	}
}

type ByOrganizationMember []*OrganizationMember

func (orgMembers ByOrganizationMember) Len() int {
	return len(orgMembers)
}

func (orgMembers ByOrganizationMember) Swap(i, j int) {
	orgMembers[i], orgMembers[j] = orgMembers[j], orgMembers[i]
}

func (orgMembers ByOrganizationMember) Less(i, j int) bool {
	return orgMembers[i].ID < orgMembers[j].ID
}

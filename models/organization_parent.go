package models

import "time"

type OrganizationParent struct {
	ID          string     `json:"id,omitempty"`
	DateCreated *time.Time `json:"date_created,omitempty"`
	DateUpdated *time.Time `json:"date_updated,omitempty"`
	From        *time.Time `json:"from,omitempty"`
	Until       *time.Time `json:"until,omitempty"`
}

func (op *OrganizationParent) Dup() *OrganizationParent {
	return &OrganizationParent{
		ID:          op.ID,
		DateCreated: copyTime(op.DateCreated),
		DateUpdated: copyTime(op.DateUpdated),
		From:        copyTime(op.From),
		Until:       copyTime(op.Until),
	}
}

type ByOrganizationParent []*OrganizationParent

func (parents ByOrganizationParent) Len() int {
	return len(parents)
}

func (parents ByOrganizationParent) Swap(i, j int) {
	parents[i], parents[j] = parents[j], parents[i]
}

func (parents ByOrganizationParent) Less(i, j int) bool {
	if !parents[i].From.Equal(*parents[j].From) {
		return parents[i].From.Before(*parents[j].From)
	}
	return parents[i].ID < parents[j].ID
}

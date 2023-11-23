package models

import (
	"sort"
	"time"
)

type Organization struct {
	ID          string                `json:"id,omitempty"`
	DateCreated *time.Time            `json:"date_created,omitempty"`
	DateUpdated *time.Time            `json:"date_updated,omitempty"`
	Type        string                `json:"type,omitempty"`
	NameDut     string                `json:"name_dut,omitempty"`
	NameEng     string                `json:"name_eng,omitempty"`
	Parent      []*OrganizationParent `json:"parent,omitempty"`
	Identifier  []*URN                `json:"identifier,omitempty"`
	Acronym     string                `json:"acronym,omitempty"`
}

func (org *Organization) IsStored() bool {
	return org.DateCreated != nil
}

func NewOrganization() *Organization {
	org := &Organization{}
	return org
}

func (org *Organization) AddIdentifier(urn *URN) {
	org.Identifier = append(org.Identifier, urn)
	sort.Sort(ByURN(org.Identifier))
}

func (org *Organization) SetIdentifier(urns ...*URN) {
	sort.Sort(ByURN(urns))
	org.Identifier = urns
}

func (org *Organization) ClearIdentifier() {
	org.Identifier = nil
}

func (org *Organization) GetIdentifierQualifiedValues() []string {
	ids := make([]string, 0, len(org.Identifier))
	for _, id := range org.Identifier {
		ids = append(ids, id.String())
	}
	return ids
}

func (org *Organization) GetIdentifierValues() []string {
	ids := make([]string, 0, len(org.Identifier))
	for _, id := range org.Identifier {
		ids = append(ids, id.Value)
	}
	return ids
}

func (org *Organization) GetIdentifierValueByNS(ns string) string {
	for _, id := range org.Identifier {
		if id.Namespace == ns {
			return id.Value
		}
	}
	return ""
}

func (org *Organization) GetIdentifierByNS(ns string) []*URN {
	urns := []*URN{}
	for _, id := range org.Identifier {
		if id.Namespace == ns {
			urns = append(urns, id)
		}
	}
	return urns
}

func (org *Organization) SetParent(parents ...*OrganizationParent) {
	sort.Sort(ByOrganizationParent(parents))
	org.Parent = parents
}

func (org *Organization) AddParent(parents ...*OrganizationParent) {
	org.Parent = append(org.Parent, parents...)
	sort.Sort(ByOrganizationParent(org.Parent))
}

func (org *Organization) Dup() *Organization {
	newOrg := &Organization{
		ID:          org.ID,
		Type:        org.Type,
		NameDut:     org.NameDut,
		NameEng:     org.NameEng,
		Acronym:     org.Acronym,
		DateCreated: copyTime(org.DateCreated),
		DateUpdated: copyTime(org.DateUpdated),
	}

	for _, id := range org.Identifier {
		newOrg.Identifier = append(newOrg.Identifier, id.Dup())
	}
	for _, op := range org.Parent {
		newOrg.Parent = append(newOrg.Parent, op.Dup())
	}

	return newOrg
}

package models

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Person struct {
	ID                  string                `json:"id,omitempty"`
	Active              bool                  `json:"active,omitempty"`
	DateCreated         *time.Time            `json:"date_created,omitempty"`
	DateUpdated         *time.Time            `json:"date_updated,omitempty"`
	Name                string                `json:"name,omitempty"`
	GivenName           string                `json:"given_name,omitempty"`
	FamilyName          string                `json:"family_name,omitempty"`
	Email               string                `json:"email,omitempty"`
	Token               map[string]string     `json:"token"`
	PreferredGivenName  string                `json:"preferred_given_name,omitempty"`
	PreferredFamilyName string                `json:"preferred_family_name,omitempty"`
	BirthDate           string                `json:"birth_date,omitempty"`
	HonorificPrefix     string                `json:"honorific_prefix,omitempty"`
	Identifier          []*URN                `json:"identifier,omitempty"`
	Organization        []*OrganizationMember `json:"organization,omitempty"`
	JobCategory         []string              `json:"job_category,omitempty"`
	Role                []string              `json:"role,omitempty"`
	Settings            map[string]string     `json:"settings,omitempty"`
	ObjectClass         []string              `json:"object_class,omitempty"`
}

func (person *Person) IsStored() bool {
	return person.DateCreated != nil
}

func NewPerson() *Person {
	p := &Person{}
	return p
}

func NewOrganizationMember(id string) *OrganizationMember {
	return &OrganizationMember{
		ID: id,
	}
}

func (p *Person) SetEmail(email string) {
	p.Email = strings.ToLower(email)
}

func (p *Person) AddIdentifier(urn *URN) {
	p.Identifier = append(p.Identifier, urn)
	sort.Sort(ByURN(p.Identifier))
}

func (p *Person) SetIdentifier(ids ...*URN) {
	sort.Sort(ByURN(ids))
	p.Identifier = ids
}

func (p *Person) ClearIdentifier() {
	p.Identifier = nil
}

func (p *Person) GetIdentifierQualifiedValues() []string {
	ids := make([]string, 0, len(p.Identifier))
	for _, id := range p.Identifier {
		ids = append(ids, id.String())
	}
	return ids
}

func (p *Person) GetIdentifierValues() []string {
	ids := make([]string, 0, len(p.Identifier))
	for _, id := range p.Identifier {
		ids = append(ids, id.Value)
	}
	return ids
}

func (p *Person) GetIdentifierByNS(ns string) []*URN {
	urns := []*URN{}
	for _, id := range p.Identifier {
		if id.Namespace == ns {
			urns = append(urns, id)
		}
	}
	return urns
}

func (p *Person) GetIdentifierValuesByNS(ns string) []string {
	vals := make([]string, 0, len(p.Identifier))
	for _, id := range p.Identifier {
		if id.Namespace == ns {
			vals = append(vals, id.Value)
		}
	}
	return vals
}

func (p *Person) EnsureBiblioID() {
	hasBiblioId := false
	for _, urn := range p.Identifier {
		if urn.Namespace == "biblio_id" {
			hasBiblioId = true
			break
		}
	}
	if !hasBiblioId {
		biblioId := uuid.NewString()
		fmt.Fprintf(os.Stderr, "adding new biblio_id: %s\n", biblioId)
		p.AddIdentifier(NewURN("biblio_id", biblioId))
	}
}

func (p *Person) SetToken(typ string, val string) {
	p.Token[typ] = val
}

func (p *Person) ClearToken() {
	p.Token = map[string]string{}
}

func (p *Person) GetTokenValue(typ string) string {
	return p.Token[typ]
}

func (p *Person) SetRole(role ...string) {
	sort.Strings(role)
	p.Role = role
}

func (p *Person) AddRole(role ...string) {
	p.Role = append(p.Role, role...)
	sort.Strings(p.Role)
}

func (p *Person) SetObjectClass(objectClass ...string) {
	sort.Strings(objectClass)
	p.ObjectClass = objectClass
}

func (p *Person) AddObjectClass(objectClass ...string) {
	p.ObjectClass = append(p.ObjectClass, objectClass...)
	sort.Strings(p.ObjectClass)
}

func (p *Person) SetJobCategory(jobCategory ...string) {
	sort.Strings(jobCategory)
	p.JobCategory = jobCategory
}

func (p *Person) AddJobCategory(jobCategory ...string) {
	p.JobCategory = append(p.JobCategory, jobCategory...)
	sort.Strings(p.JobCategory)
}

func (p *Person) AddOrganizationMember(orgMembers ...*OrganizationMember) {
	p.Organization = append(p.Organization, orgMembers...)
	sort.Sort(ByOrganizationMember(p.Organization))
}

func (p *Person) SetOrganizationMember(orgMembers ...*OrganizationMember) {
	sort.Sort(ByOrganizationMember(orgMembers))
	p.Organization = orgMembers
}

func (p *Person) Dup() *Person {
	newP := &Person{
		ID:                  p.ID,
		DateCreated:         copyTime(p.DateCreated),
		DateUpdated:         copyTime(p.DateUpdated),
		Active:              p.Active,
		Name:                p.Name,
		GivenName:           p.GivenName,
		FamilyName:          p.FamilyName,
		Email:               p.Email,
		PreferredGivenName:  p.PreferredGivenName,
		PreferredFamilyName: p.PreferredFamilyName,
		BirthDate:           p.BirthDate,
		HonorificPrefix:     p.HonorificPrefix,
	}
	newP.Token = map[string]string{}
	for typ, val := range p.Token {
		newP.Token[typ] = val
	}
	for _, id := range p.Identifier {
		newP.AddIdentifier(NewURN(id.Namespace, id.Value))
	}
	for _, orgMember := range p.Organization {
		newP.Organization = append(newP.Organization, orgMember.Dup())
	}
	if p.Settings != nil {
		newP.Settings = make(map[string]string)
		for key, val := range p.Settings {
			newP.Settings[key] = val
		}
	}
	if p.JobCategory != nil {
		newP.JobCategory = append(newP.JobCategory, p.JobCategory...)
	}
	if p.ObjectClass != nil {
		newP.ObjectClass = append(newP.ObjectClass, p.ObjectClass...)
	}
	if p.Role != nil {
		newP.Role = append(newP.Role, p.Role...)
	}

	return newP
}

type ByPerson []*Person

func (people ByPerson) Len() int {
	return len(people)
}

func (people ByPerson) Swap(i, j int) {
	people[i], people[j] = people[j], people[i]
}

func (people ByPerson) Less(i, j int) bool {
	return people[i].DateUpdated.Before(*people[j].DateUpdated)
}

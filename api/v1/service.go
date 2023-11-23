package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/ugent-library/people-service/models"
)

type Service struct {
	repository models.Repository
}

func NewService(repository models.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetPerson(ctx context.Context, req *GetPersonRequest) (*Person, error) {
	person, err := s.repository.GetPerson(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return mapToExternalPerson(person), nil
}

func (s *Service) GetPeopleById(ctx context.Context, req *GetPeopleByIdRequest) (*PersonListResponse, error) {
	urns := make([]models.URN, 0, len(req.ID))
	for _, id := range req.ID {
		urn, err := models.ParseURN(id)
		if err != nil {
			return nil, err
		}
		urns = append(urns, *urn)
	}
	people, err := s.repository.GetPeopleByIdentifier(ctx, urns...)
	if err != nil {
		return nil, err
	}
	res := &PersonListResponse{
		Data: make([]Person, 0, len(people)),
	}
	for _, person := range people {
		res.Data = append(res.Data, *mapToExternalPerson(person))
	}
	return res, nil
}

func (s *Service) GetPeople(ctx context.Context, req *GetPeopleRequest) (*PersonListResponse, error) {
	var people []*models.Person
	var err error
	var cursor string

	if req.Cursor != "" {
		people, cursor, err = s.repository.GetMorePeople(ctx, req.Cursor)
	} else {
		people, cursor, err = s.repository.GetPeople(ctx)
	}
	if err != nil {
		return nil, err
	}

	res := &PersonListResponse{
		Data: make([]Person, 0, len(people)),
	}
	if cursor != "" {
		res.Cursor = NewOptString(cursor)
	}
	for _, person := range people {
		res.Data = append(res.Data, *mapToExternalPerson(person))
	}

	return res, nil
}

func (s *Service) SuggestPeople(ctx context.Context, req *SuggestPeopleRequest) (*PersonListResponse, error) {
	people, err := s.repository.SuggestPeople(ctx, req.Query)
	if err != nil {
		return nil, err
	}

	res := &PersonListResponse{
		Data: make([]Person, 0, len(people)),
	}
	for _, person := range people {
		res.Data = append(res.Data, *mapToExternalPerson(person))
	}

	return res, nil
}

func (s *Service) SetPersonOrcid(ctx context.Context, req *SetPersonOrcidRequest) (*Person, error) {
	if err := s.repository.SetPersonOrcid(ctx, req.ID, req.Orcid); err != nil {
		return nil, err
	}
	person, err := s.repository.GetPerson(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return mapToExternalPerson(person), nil
}

func (s *Service) SetPersonOrcidToken(ctx context.Context, req *SetPersonOrcidTokenRequest) (*Person, error) {
	if err := s.repository.SetPersonOrcidToken(ctx, req.ID, req.OrcidToken); err != nil {
		return nil, err
	}
	person, err := s.repository.GetPerson(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return mapToExternalPerson(person), nil
}

func (s *Service) SetPersonRole(ctx context.Context, req *SetPersonRoleRequest) (*Person, error) {
	if err := s.repository.SetPersonRole(ctx, req.ID, req.Role); err != nil {
		return nil, err
	}
	person, err := s.repository.GetPerson(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return mapToExternalPerson(person), nil
}

func (s *Service) SetPersonSettings(ctx context.Context, req *SetPersonSettingsRequest) (*Person, error) {
	if req.Settings == nil {
		return nil, fmt.Errorf("%w: attribute settings is missing in request body", models.ErrMissingArgument)
	}
	if err := s.repository.SetPersonSettings(ctx, req.ID, req.Settings); err != nil {
		return nil, err
	}
	person, err := s.repository.GetPerson(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return mapToExternalPerson(person), nil
}

func (s *Service) GetOrganization(ctx context.Context, req *GetOrganizationRequest) (*Organization, error) {
	org, err := s.repository.GetOrganization(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return mapToExternalOrganization(org), nil
}

func (s *Service) GetOrganizationsById(ctx context.Context, req *GetOrganizationsByIdRequest) (*OrganizationListResponse, error) {
	urns := make([]models.URN, 0, len(req.ID))
	for _, id := range req.ID {
		urn, err := models.ParseURN(id)
		if err != nil {
			return nil, err
		}
		urns = append(urns, *urn)
	}
	orgs, err := s.repository.GetOrganizationsByIdentifier(ctx, urns...)
	if err != nil {
		return nil, err
	}
	res := &OrganizationListResponse{
		Data: make([]Organization, 0, len(orgs)),
	}
	for _, org := range orgs {
		res.Data = append(res.Data, *mapToExternalOrganization(org))
	}
	return res, nil
}

func (s *Service) GetOrganizations(ctx context.Context, req *GetOrganizationsRequest) (*OrganizationListResponse, error) {
	var organizations []*models.Organization
	var err error
	var cursor string

	if req.Cursor != "" {
		organizations, cursor, err = s.repository.GetMoreOrganizations(ctx, req.Cursor)
	} else {
		organizations, cursor, err = s.repository.GetOrganizations(ctx)
	}
	if err != nil {
		return nil, err
	}

	res := &OrganizationListResponse{
		Data: make([]Organization, 0, len(organizations)),
	}
	if cursor != "" {
		res.Cursor = NewOptString(cursor)
	}
	for _, org := range organizations {
		res.Data = append(res.Data, *mapToExternalOrganization(org))
	}

	return res, nil
}

func (s *Service) SuggestOrganizations(ctx context.Context, req *SuggestOrganizationsRequest) (*OrganizationListResponse, error) {
	orgs, err := s.repository.SuggestOrganizations(ctx, req.Query)
	if err != nil {
		return nil, err
	}

	res := &OrganizationListResponse{
		Data: make([]Organization, 0, len(orgs)),
	}
	for _, org := range orgs {
		res.Data = append(res.Data, *mapToExternalOrganization(org))
	}

	return res, nil
}

func (s *Service) AddPerson(ctx context.Context, p *Person) (*Person, error) {
	var person *models.Person

	if p.ID.Set {
		oldPerson, err := s.repository.GetPerson(ctx, p.ID.Value)
		if errors.Is(err, models.ErrNotFound) {
			return nil, fmt.Errorf("cannot find person record %s to update", p.ID.Value)
		} else if err != nil {
			return nil, err
		}
		person = oldPerson
	} else {
		person = models.NewPerson()
	}

	person.Active = p.Active.Value
	person.BirthDate = p.BirthDate.Value
	person.SetEmail(p.Email.Value)
	person.ExpirationDate = p.ExpirationDate.Value
	person.GivenName = p.GivenName.Value
	person.FamilyName = p.FamilyName.Value
	person.Name = p.Name.Value
	person.SetJobCategory(p.JobCategory...)
	person.SetObjectClass(p.ObjectClass...)
	person.ClearToken()
	for _, token := range p.Token {
		person.AddToken(token.PropertyID, token.Value)
	}
	person.PreferredGivenName = p.PreferredGivenName.Value
	person.PreferredFamilyName = p.PreferredFamilyName.Value
	person.SetRole(p.Role...)
	person.Settings = p.Settings.Value
	person.HonorificPrefix = p.HonorificPrefix.Value

	ids := make([]*models.URN, 0, len(p.Identifier))
	for _, identifier := range p.Identifier {
		ids = append(ids, models.NewURN(identifier.PropertyID, identifier.Value))
	}
	person.SetIdentifier(ids...)

	orgMembers := []*models.OrganizationMember{}
	for _, orgMember := range p.Organization {
		newOrgMember := models.NewOrganizationMember(orgMember.ID)
		orgMembers = append(orgMembers, newOrgMember)
	}
	person.SetOrganizationMember(orgMembers...)

	if newPerson, err := s.repository.SavePerson(ctx, person); err != nil {
		return nil, err
	} else {
		person = newPerson
	}

	return mapToExternalPerson(person), nil
}

func (s *Service) AddOrganization(ctx context.Context, o *Organization) (*Organization, error) {
	var org *models.Organization

	if o.ID.Set {
		oldOrg, err := s.repository.GetOrganization(ctx, o.ID.Value)
		if errors.Is(err, models.ErrNotFound) {
			return nil, fmt.Errorf("cannot find organization record \"%s\" to update", o.ID.Value)
		} else if err != nil {
			return nil, err
		}
		org = oldOrg
	} else {
		org = models.NewOrganization()
	}

	org.Acronym = o.Acronym.Value
	org.NameDut = o.NameDut.Value
	org.NameEng = o.NameEng.Value
	parents := []*models.OrganizationParent{}
	for _, parent := range o.Parent {
		op := models.OrganizationParent{
			ID: parent.ID,
		}
		parents = append(parents, &op)
	}
	org.SetParent(parents...)
	org.Type = o.Type.Value

	ids := make([]*models.URN, 0, len(o.Identifier))
	for _, identifier := range o.Identifier {
		ids = append(ids, models.NewURN(identifier.PropertyID, identifier.Value))
	}
	org.SetIdentifier(ids...)

	if newOrg, err := s.repository.SaveOrganization(ctx, org); err != nil {
		return nil, err
	} else {
		org = newOrg
	}

	return mapToExternalOrganization(org), nil
}

func (s *Service) NewError(ctx context.Context, err error) *ErrorStatusCode {
	if errors.Is(err, models.ErrNotFound) {
		return &ErrorStatusCode{
			StatusCode: 404,
			Response: Error{
				Code:    404,
				Message: err.Error(),
			},
		}
	}
	if errors.Is(err, models.ErrMissingArgument) {
		return &ErrorStatusCode{
			StatusCode: 400,
			Response: Error{
				Code:    400,
				Message: err.Error(),
			},
		}
	}
	if errors.Is(err, models.ErrInvalidReference) {
		return &ErrorStatusCode{
			StatusCode: 400,
			Response: Error{
				Code:    400,
				Message: err.Error(),
			},
		}
	}

	return &ErrorStatusCode{
		StatusCode: 500,
		Response: Error{
			Code:    500,
			Message: err.Error(),
		},
	}
}

func mapToExternalPerson(person *models.Person) *Person {
	p := &Person{}
	p.ID = NewOptString(person.ID)
	p.Active = NewOptBool(person.Active)
	if person.BirthDate != "" {
		p.BirthDate = NewOptString(person.BirthDate)
	}
	p.DateCreated = NewOptDateTime(*person.DateCreated)
	p.DateUpdated = NewOptDateTime(*person.DateUpdated)
	if person.Email != "" {
		p.Email = NewOptString(person.Email)
	}
	if person.ExpirationDate != "" {
		p.ExpirationDate = NewOptString(person.ExpirationDate)
	}
	if person.GivenName != "" {
		p.GivenName = NewOptString(person.GivenName)
	}
	if person.FamilyName != "" {
		p.FamilyName = NewOptString(person.FamilyName)
	}
	if person.Name != "" {
		p.Name = NewOptString(person.Name)
	}
	if person.PreferredGivenName != "" {
		p.PreferredGivenName = NewOptString(person.PreferredGivenName)
	}
	if person.PreferredFamilyName != "" {
		p.PreferredFamilyName = NewOptString(person.PreferredFamilyName)
	}
	p.JobCategory = append(p.JobCategory, person.JobCategory...)
	p.ObjectClass = append(p.ObjectClass, person.ObjectClass...)
	for _, token := range person.Token {
		p.Token = append(p.Token, newPropertyValue(token.Namespace, token.Value))
	}

	for _, orgMember := range person.Organization {
		externalOrgMember := OrganizationMember{
			ID:          orgMember.ID,
			DateCreated: NewOptDateTime(*orgMember.DateCreated),
			DateUpdated: NewOptDateTime(*orgMember.DateUpdated),
		}
		p.Organization = append(p.Organization, externalOrgMember)
	}
	for _, id := range person.Identifier {
		p.Identifier = append(p.Identifier, newPropertyValue(id.Namespace, id.Value))
	}

	p.Role = append(p.Role, person.Role...)
	if person.Settings != nil {
		pSettings := PersonSettings{}
		for k, v := range person.Settings {
			pSettings[k] = v
		}
		p.Settings = NewOptPersonSettings(pSettings)
	}
	if person.HonorificPrefix != "" {
		p.HonorificPrefix = NewOptString(person.HonorificPrefix)
	}

	return p
}

func mapToExternalOrganization(org *models.Organization) *Organization {
	o := &Organization{}
	o.ID = NewOptString(org.ID)
	o.DateCreated = NewOptDateTime(*org.DateCreated)
	o.DateUpdated = NewOptDateTime(*org.DateUpdated)
	if org.Acronym != "" {
		o.Acronym = NewOptString(org.Acronym)
	}
	if org.NameDut != "" {
		o.NameDut = NewOptString(org.NameDut)
	}
	if org.NameEng != "" {
		o.NameEng = NewOptString(org.NameEng)
	}
	for _, id := range org.Identifier {
		o.Identifier = append(o.Identifier, newPropertyValue(id.Namespace, id.Value))
	}
	for _, organizationParent := range org.Parent {
		op := OrganizationParent{
			ID:          organizationParent.ID,
			DateCreated: NewOptDateTime(*organizationParent.DateCreated),
			DateUpdated: NewOptDateTime(*organizationParent.DateUpdated),
		}
		o.Parent = append(o.Parent, op)
	}
	o.Type = NewOptString(org.Type)

	return o
}

func newPropertyValue(propertyID string, value string) PropertyValue {
	return PropertyValue{
		Type:       "PropertyValue",
		PropertyID: propertyID,
		Value:      value,
	}
}

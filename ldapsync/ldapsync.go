package ldapsync

import (
	"context"
	"reflect"
	"slices"
	"sort"

	"github.com/go-ldap/ldap/v3"
	"github.com/ugent-library/people-service/models"
	"github.com/ugent-library/people-service/ugentldap"
	"go.uber.org/zap"
)

type Synchronizer struct {
	repository      models.Repository
	ugentLdapClient *ugentldap.Client
	logger          *zap.SugaredLogger
}

func NewSynchronizer(repo models.Repository, ugentLdapClient *ugentldap.Client, l *zap.SugaredLogger) *Synchronizer {
	return &Synchronizer{
		repository:      repo,
		ugentLdapClient: ugentLdapClient,
		logger:          l,
	}
}

func (si *Synchronizer) Sync(ctx context.Context) error {
	newActiveIDs := []string{}

	err := si.ugentLdapClient.SearchPeople(ctx, PersonQuery, func(ldapEntry *ldap.Entry) error {
		newPerson, err := si.ldapEntryToPerson(ctx, ldapEntry)

		if err != nil {
			return err
		}

		var oldPeople []*models.Person

		historicUgentIDs := newPerson.GetIdentifierByNS("historic_ugent_id")
		if len(historicUgentIDs) > 0 {
			oldPeople, err = si.repository.GetPeopleByIdentifier(ctx, historicUgentIDs...)
			if err != nil {
				return err
			}
		}

		// IMPORTANT: sort inverse by date_updated
		sort.Sort(sort.Reverse(models.ByPerson(oldPeople)))

		if len(oldPeople) == 0 {
			newPerson, err := si.repository.CreatePerson(ctx, newPerson)
			if err != nil {
				return err
			}
			si.logger.Infof("person record %s: created", newPerson.ID)
			newActiveIDs = append(newActiveIDs, newPerson.ID)
		} else {
			// delete older versions with same historic_ugent_id
			if len(oldPeople) > 1 {
				for _, person := range oldPeople[1:] {
					err := si.repository.DeletePerson(ctx, person.ID)
					if err != nil {
						return err
					}
					si.logger.Infof("person record %s: deleted", person.ID)
				}
			}

			// insert updated version
			oldPerson := oldPeople[0]

			newActiveIDs = append(newActiveIDs, oldPerson.ID)

			oldStoredPerson := oldPerson.Dup()
			keepIds := make([]*models.URN, 0, 3)
			for _, id := range oldPerson.Identifier {
				switch id.Namespace {
				case "orcid", "gismo_id", "biblio_id":
					keepIds = append(keepIds, id.Dup())
				}
			}
			oldPerson.ClearIdentifier()
			oldPerson.SetIdentifier(newPerson.Identifier...)
			for _, id := range keepIds {
				oldPerson.AddIdentifier(id)
			}
			oldPerson.EnsureBiblioID() // P.S. also done in repository for other reasons

			oldPerson.Active = true
			oldPerson.BirthDate = newPerson.BirthDate
			oldPerson.Email = newPerson.Email
			oldPerson.GivenName = newPerson.GivenName
			oldPerson.FamilyName = newPerson.FamilyName
			oldPerson.Name = newPerson.Name
			oldPerson.JobCategory = newPerson.JobCategory
			oldPerson.HonorificPrefix = newPerson.HonorificPrefix
			oldPerson.ObjectClass = newPerson.ObjectClass

			// only add organizations not known yet (gismo possibly knows more)
			for _, newOrgMember := range newPerson.Organization {
				found := false
				for _, oldOrgMember := range oldPerson.Organization {
					if oldOrgMember.ID == newOrgMember.ID {
						found = true
						break
					}
				}
				if !found {
					oldPerson.AddOrganizationMember(newOrgMember)
				}
			}

			// prepare for comparison
			if len(oldPerson.Organization) == 0 {
				oldPerson.Organization = nil
			}
			if len(oldPerson.JobCategory) == 0 {
				oldPerson.JobCategory = nil
			}
			if len(oldPerson.ObjectClass) == 0 {
				oldPerson.ObjectClass = nil
			}
			if len(oldPerson.Identifier) == 0 {
				oldPerson.Identifier = nil
			}
			if len(oldPerson.Token) == 0 {
				oldPerson.Token = map[string]string{}
			}

			if reflect.DeepEqual(oldPerson, oldStoredPerson) {
				si.logger.Infof("person record %s: no update", oldPerson.ID)
				return nil
			}

			oldPerson, err := si.repository.SavePerson(ctx, oldPerson)
			if err != nil {
				return err
			}
			si.logger.Infof("person record %s: updated", oldPerson.ID)
		}

		return nil
	})

	if err != nil {
		return err
	}

	si.logger.Infof("processed %d ldap records", len(newActiveIDs))

	// deactivate people
	activeIDs, err := si.repository.GetPersonIDActive(ctx, true)
	if err != nil {
		return err
	}

	for _, activeID := range activeIDs {
		if !slices.Contains(newActiveIDs, activeID) {
			err := si.repository.SetPersonActive(ctx, activeID, false)
			if err != nil {
				si.logger.Errorf("failed to set person record %s to active=false: %s", activeID, err)
			}
			si.logger.Infof("set person record %s to active=false", activeID)
		}
	}

	return err
}

func (si *Synchronizer) ldapEntryToPerson(ctx context.Context, ldapEntry *ldap.Entry) (*models.Person, error) {
	newPerson := models.NewPerson()
	newPerson.Active = true

	depIds := []string{}
	facultyIds := []string{}

	for _, attr := range ldapEntry.Attributes {
		for _, val := range attr.Values {
			switch attr.Name {
			case "uid":
				newPerson.AddIdentifier(models.NewURN("ugent_username", val))
			case "ugentHistoricIDs":
				newPerson.AddIdentifier(models.NewURN("historic_ugent_id", val))
			case "ugentBarcode":
				newPerson.AddIdentifier(models.NewURN("ugent_barcode", val))
			case "ugentPreferredGivenName":
				newPerson.GivenName = val
			case "ugentPreferredSn":
				newPerson.FamilyName = val
			case "displayName":
				newPerson.Name = val
			case "ugentBirthDate":
				newPerson.BirthDate = val
			case "mail":
				newPerson.SetEmail(val)
			case "ugentJobCategory":
				newPerson.AddJobCategory(val)
			case "ugentAddressingTitle":
				newPerson.HonorificPrefix = val
			case "objectClass":
				newPerson.AddObjectClass(val)
			case "ugentFaculty":
				facultyIds = append(facultyIds, val)
			case "departmentNumber":
				depIds = append(depIds, val)
			}
		}
	}

	orgIds := []string{}

	if len(depIds) > 0 {
		orgIds = append(orgIds, depIds...)
	}
	if len(facultyIds) > 0 {
		orgIds = append(orgIds, facultyIds...)
	}

	for _, orgId := range orgIds {
		orgs, err := si.repository.GetOrganizationsByIdentifier(ctx, models.NewURN("biblio_id", orgId))
		if err != nil {
			return nil, err
		}

		var org *models.Organization
		if len(orgs) == 0 {
			si.logger.Infof("adding dummy organization %s for person with name '%s'", orgId, newPerson.Name)
			o, err := si.addDummyOrg(ctx, orgId)
			if err != nil {
				return nil, err
			}
			org = o
		} else {
			org = orgs[0]
		}
		newOrgMember := models.NewOrganizationMember(org.ID)
		newPerson.AddOrganizationMember(newOrgMember)
	}

	if slices.Contains(newPerson.ObjectClass, "ugentFormerEmployee") {
		orgs, err := si.repository.GetOrganizationsByIdentifier(ctx, models.NewURN("biblio_id", "UGent"))
		if err != nil {
			return nil, err
		}
		var org *models.Organization
		if len(orgs) == 0 {
			o, err := si.addDummyOrg(ctx, "UGent")
			if err != nil {
				return nil, err
			}
			org = o
		} else {
			org = orgs[0]
		}
		hasOrg := false
		for _, orgMember := range newPerson.Organization {
			if orgMember.ID == org.ID {
				hasOrg = true
				break
			}
		}
		if !hasOrg {
			newPerson.AddOrganizationMember(models.NewOrganizationMember(org.ID))
		}
	}
	if slices.Contains(newPerson.ObjectClass, "uzEmployee") {
		orgs, err := si.repository.GetOrganizationsByIdentifier(ctx, models.NewURN("biblio_id", "UZGent"))
		if err != nil {
			return nil, err
		}
		var org *models.Organization
		if len(orgs) == 0 {
			o, err := si.addDummyOrg(ctx, "UZGent")
			if err != nil {
				return nil, err
			}
			org = o
		} else {
			org = orgs[0]
		}
		hasOrg := false
		for _, orgMember := range newPerson.Organization {
			if orgMember.ID == org.ID {
				hasOrg = true
				break
			}
		}
		if !hasOrg {
			newPerson.AddOrganizationMember(models.NewOrganizationMember(org.ID))
		}
	}

	return newPerson, nil
}

func (si *Synchronizer) addDummyOrg(ctx context.Context, orgId string) (*models.Organization, error) {
	org := models.NewOrganization()
	org.NameEng = orgId
	org.AddIdentifier(models.NewURN("biblio_id", orgId))
	return si.repository.CreateOrganization(ctx, org)
}

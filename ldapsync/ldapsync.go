package ldapsync

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"slices"
	"sort"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/samber/lo"
	"github.com/ugent-library/people-service/models"
	"github.com/ugent-library/people-service/ugentldap"
)

type Synchronizer struct {
	repository      models.Repository
	ugentLdapClient *ugentldap.Client
}

func NewSynchronizer(repo models.Repository, ugentLdapClient *ugentldap.Client) *Synchronizer {
	return &Synchronizer{
		repository:      repo,
		ugentLdapClient: ugentLdapClient,
	}
}

// TODO: set all to active=false, and then set found records to active=true. Needed: repository transaction
func (si *Synchronizer) Sync(cb func(*models.Person)) error {
	ctx := context.TODO()
	newActiveIDs := []string{}
	err := si.ugentLdapClient.SearchPeople(ldapPersonQuery, func(ldapEntry *ldap.Entry) error {
		newPerson, err := si.ldapEntryToPerson(ldapEntry)
		if err != nil {
			return err
		}

		var oldPeople []*models.Person

		historicUgentIDs := newPerson.GetIdentifierByNS("historic_ugent_id")
		fmt.Fprintf(os.Stderr, "ldap person: name: %s, historic_ugent_id: %+v\n", newPerson.Name, historicUgentIDs)
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
			cb(newPerson)
			newActiveIDs = append(newActiveIDs, newPerson.ID)
		} else {
			// delete older versions with same historic_ugent_id
			if len(oldPeople) > 1 {
				for _, person := range oldPeople[1:] {
					fmt.Fprintf(os.Stderr, "deleting old person %s\n", person.ID)
					err := si.repository.DeletePerson(ctx, person.ID)
					if err != nil {
						return err
					}
				}
			}

			// insert updated version
			oldPerson := oldPeople[0]

			fmt.Fprintf(os.Stderr,
				"found old person %s for %s\n", oldPerson.ID, oldPerson.Name)

			newActiveIDs = append(newActiveIDs, oldPerson.ID)

			oldStoredPerson := oldPerson.Dup()
			var gismoId string
			var orcid string
			for _, id := range oldPerson.Identifier {
				switch id.Namespace {
				case "orcid":
					orcid = id.Value
				case "gismo_id":
					gismoId = id.Value
				}
			}
			oldPerson.ClearIdentifier()
			oldPerson.SetIdentifier(newPerson.Identifier...)
			if gismoId != "" {
				oldPerson.AddIdentifier(models.NewURN("gismo_id", gismoId))
			}
			if orcid != "" {
				oldPerson.AddIdentifier(models.NewURN("orcid", orcid))
			}
			oldPerson.Active = true
			oldPerson.BirthDate = newPerson.BirthDate
			oldPerson.Email = newPerson.Email
			oldPerson.GivenName = newPerson.GivenName
			oldPerson.FamilyName = newPerson.FamilyName
			oldPerson.PreferredGivenName = newPerson.PreferredGivenName
			oldPerson.PreferredFamilyName = newPerson.PreferredFamilyName
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

			if len(oldPerson.Organization) == 0 {
				oldPerson.Organization = nil
			}
			if len(oldPerson.JobCategory) == 0 {
				oldPerson.JobCategory = nil
			}
			if len(oldPerson.ObjectClass) == 0 {
				oldPerson.ObjectClass = nil
			}
			if reflect.DeepEqual(oldPerson, oldStoredPerson) {
				return nil
			}

			oldPerson, err := si.repository.SavePerson(ctx, oldPerson)
			if err != nil {
				return err
			}
			cb(oldPerson)
		}

		return nil
	})

	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "processed %d ldap people\n", len(newActiveIDs))

	// deactivate people
	activeIDs, err := si.repository.GetPersonIDActive(ctx, true)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "found %d active people in db\n", len(activeIDs))

	inactiveIDs := []string{}
	for _, activeID := range activeIDs {
		if !slices.Contains(newActiveIDs, activeID) {
			inactiveIDs = append(inactiveIDs, activeID)
		}
	}
	activeIDs = nil

	fmt.Fprintf(os.Stderr, "found %d active people in db that should be inactive\n", len(inactiveIDs))

	chunkedList := []string{}
	chunkSize := 100
	for len(inactiveIDs) > 0 {
		var id string
		id, inactiveIDs = inactiveIDs[0], inactiveIDs[1:]
		chunkedList = append(chunkedList, id)
		if len(chunkedList) >= chunkSize {
			si.repository.SetPeopleActive(ctx, false, chunkedList...)
			chunkedList = nil
		}
	}
	if len(chunkedList) > 0 {
		si.repository.SetPeopleActive(ctx, false, chunkedList...)
	}

	return err
}

// ldapEntryToPerson maps ldap entry to new Person
func (si *Synchronizer) ldapEntryToPerson(ldapEntry *ldap.Entry) (*models.Person, error) {
	newPerson := models.NewPerson()
	newPerson.Active = true
	ctx := context.TODO()

	for _, attr := range ldapEntry.Attributes {
		for _, val := range attr.Values {
			switch attr.Name {
			case "uid":
				newPerson.AddIdentifier(models.NewURN("ugent_username", val))
			case "ugentHistoricIDs":
				newPerson.AddIdentifier(models.NewURN("historic_ugent_id", val))
			case "ugentBarcode":
				newPerson.AddIdentifier(models.NewURN("ugent_barcode", val))
			case "givenName":
				newPerson.GivenName = val
			case "ugentPreferredGivenName":
				newPerson.PreferredGivenName = val
			case "sn":
				newPerson.FamilyName = val
			case "ugentPreferredSn":
				newPerson.PreferredFamilyName = val
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
			case "departmentNumber":
				realOrgs, err := si.repository.GetOrganizationsByIdentifier(ctx, models.NewURN("ugent_id", val))
				// ignore for now. Maybe tomorrow on the next run
				if err != nil {
					return nil, err
				}
				if len(realOrgs) == 0 {
					continue
				}
				// ugent_id not unique, and some of them are not in use anymore
				// e.g. LW06 used to be "Latijn en Grieks", now "Taalkunde"
				now := time.Now()
				realOrgs = lo.Filter(realOrgs, func(org *models.Organization, index int) bool {
					if org.Type != "department" {
						return false
					}
					var validOrganizationParent *models.OrganizationParent
					for _, oParent := range org.Parent {
						if oParent.From.Before(now) && (oParent.Until == nil || oParent.Until.After(now)) {
							validOrganizationParent = oParent
							break
						}
					}
					return validOrganizationParent != nil
				})
				if len(realOrgs) == 0 {
					continue
				}
				newOrgMember := models.NewOrganizationMember(realOrgs[0].ID)
				newPerson.AddOrganizationMember(newOrgMember)
			}
		}
	}

	return newPerson, nil
}

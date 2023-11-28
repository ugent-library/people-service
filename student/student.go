package student

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/samber/lo"
	"github.com/ugent-library/people-service/models"
	"github.com/ugent-library/people-service/ugentldap"
)

type Importer struct {
	repository      models.Repository
	ugentLdapClient *ugentldap.Client
}

func NewImporter(repo models.Repository, ugentLdapClient *ugentldap.Client) *Importer {
	return &Importer{
		repository:      repo,
		ugentLdapClient: ugentLdapClient,
	}
}

// Each calls callback function with valid models.Person to save
func (si *Importer) Each(cb func(*models.Person) error) error {
	ctx := context.TODO()
	err := si.ugentLdapClient.SearchPeople("(objectClass=ugentStudent)", func(ldapEntry *ldap.Entry) error {
		newPerson, err := si.ldapEntryToPerson(ldapEntry)
		if err != nil {
			return err
		}

		if newPerson.Email == "" {
			fmt.Fprintf(os.Stderr, "ignoring student record without email\n")
			return nil
		}

		oldPeople, err := si.repository.GetPeopleByIdentifier(ctx, newPerson.GetIdentifierByNS("historic_ugent_id")...)
		if err != nil {
			return err
		}

		if len(oldPeople) == 0 {
			if err := cb(newPerson); err != nil {
				return err
			}
		} else {
			oldPerson := oldPeople[0]
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
			oldPerson.Name = newPerson.Name
			oldPerson.JobCategory = newPerson.JobCategory
			oldPerson.HonorificPrefix = newPerson.HonorificPrefix
			oldPerson.ObjectClass = newPerson.ObjectClass
			oldPerson.ExpirationDate = newPerson.ExpirationDate

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
				fmt.Fprintf(os.Stderr, "no changes detected for person %s\n", oldPerson.Email)
				return nil
			}

			if err := cb(oldPerson); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

// ldapEntryToPerson maps ldap entry to new Person
func (si *Importer) ldapEntryToPerson(ldapEntry *ldap.Entry) (*models.Person, error) {
	newPerson := models.NewPerson()
	newPerson.Active = true
	ctx := context.TODO()

	for _, attr := range ldapEntry.Attributes {
		for _, val := range attr.Values {
			switch attr.Name {
			case "uid":
				newPerson.AddIdentifier(models.NewURN("ugent_username", val))
			case "ugentID":
				newPerson.AddIdentifier(models.NewURN("ugent_id", val))
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
			case "ugentExpirationDate":
				newPerson.ExpirationDate = val
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

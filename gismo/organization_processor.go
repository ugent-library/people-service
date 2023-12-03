package gismo

import (
	"context"
	"fmt"
	"time"

	"github.com/ugent-library/people-service/models"
)

type OrganizationProcessor struct {
	repository models.Repository
}

func NewOrganizationProcessor(repo models.Repository) *OrganizationProcessor {
	return &OrganizationProcessor{
		repository: repo,
	}
}

func (op *OrganizationProcessor) Process(buf []byte) (*Message, error) {
	msg, err := ParseOrganizationMessage(buf)
	if err != nil {
		return nil, err
	}

	ctx := context.TODO()

	var org *models.Organization
	orgs, err := op.repository.GetOrganizationsByIdentifier(ctx, models.NewURN("gismo_id", msg.ID))
	if err != nil {
		return nil, fmt.Errorf("%w: unable to fetch organization record: %s", models.ErrFatal, err)
	}
	if len(orgs) == 0 {
		org = models.NewOrganization()
	} else {
		org = orgs[0]
	}

	if msg.Source == "gismo.organization.update" {
		now := time.Now()
		org.NameDut = ""
		org.NameEng = ""
		org.Type = "organization"
		org.Parent = nil
		org.ClearIdentifier()
		org.AddIdentifier(models.NewURN("gismo_id", msg.ID))

		for _, attr := range msg.Attributes {
			withinDateRange := attr.ValidAt(now)
			switch attr.Name {
			case "parent_id":
				orgParentsByGismo, err := op.repository.GetOrganizationsByIdentifier(ctx, models.NewURN("gismo_id", attr.Value))
				if err != nil {
					return nil, fmt.Errorf("%w: unable to query database: %s", models.ErrFatal, err)
				}

				if len(orgParentsByGismo) > 0 {
					org.AddParent(&models.OrganizationParent{
						ID:    orgParentsByGismo[0].ID,
						From:  attr.StartDate,
						Until: attr.EndDate,
					})
				} else {
					orgParentByGismo := models.NewOrganization()
					orgParentByGismo.AddIdentifier(models.NewURN("gismo_id", attr.Value))
					orgParentByGismo, err = op.repository.CreateOrganization(ctx, orgParentByGismo)
					if err != nil {
						return nil, fmt.Errorf("%w: unable to create parent organization: %s", models.ErrFatal, err)
					}
					org.AddParent(&models.OrganizationParent{
						ID:    orgParentByGismo.ID,
						From:  attr.StartDate,
						Until: attr.EndDate,
					})
				}
			case "acronym":
				org.Acronym = attr.Value
			case "name_dut":
				if withinDateRange {
					org.NameDut = attr.Value
				}
			case "name_eng":
				if withinDateRange {
					org.NameEng = attr.Value
				}
			case "type":
				org.Type = attr.Value
			case "ugent_memorialis_id":
				org.AddIdentifier(models.NewURN("ugent_memorialis_id", attr.Value))
			case "code":
				org.AddIdentifier(models.NewURN("ugent_id", attr.Value))
			case "biblio_code":
				org.AddIdentifier(models.NewURN("biblio_id", attr.Value))
			}
		}

		if _, err := op.repository.SaveOrganization(ctx, org); err != nil {
			return nil, fmt.Errorf("%w: unable to save organization record: %s", models.ErrFatal, err)
		}
	} else if msg.Source == "gismo.organization.delete" {
		if org.IsStored() {
			if err := op.repository.DeleteOrganization(ctx, org.ID); err != nil {
				return nil, fmt.Errorf("%w: unable to delete organization record: %s", models.ErrFatal, err)
			}
		}
	}
	return msg, nil
}

package gismo

import (
	"bytes"
	"errors"
	"strings"
	"time"

	"github.com/antchfx/xmlquery"
	"github.com/ugent-library/people-service/models"
)

func ParseOrganizationMessage(buf []byte) (*Message, error) {
	doc, err := Parse(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	node := xmlquery.FindOne(doc, "//cfOrgUnit")

	if node == nil {
		return nil, errors.New("unable to find xml node //cfOrgUnit")
	}

	idNode := xmlquery.FindOne(node, "cfOrgUnitId")

	if idNode == nil {
		return nil, errors.New("unable to find xml node //cfOrgUnit/cfOrgUnitId")
	}

	orgNode := xmlquery.FindOne(doc, "//Body/organizations")

	if orgNode == nil {
		return nil, errors.New("unable to find node //Body/organizations")
	}

	msg := &Message{
		ID:   strings.TrimSpace(idNode.InnerText()),
		Date: orgNode.SelectAttr("date"),
	}

	if node.SelectAttr("action") == "DELETE" {
		msg.Source = "gismo.organization.delete"
	} else {
		msg.Source = "gismo.organization.update"

		acronymNode := xmlquery.FindOne(node, "cfAcro")
		if acronymNode != nil {
			msg.Attributes = append(msg.Attributes, Attribute{
				Name:      "acronym",
				Value:     strings.TrimSpace(acronymNode.InnerText()),
				StartDate: &models.BeginningOfTime,
				EndDate:   nil,
			})
		}

		for _, n := range xmlquery.Find(node, "./cfName") {
			t1, err := time.Parse(time.RFC3339, n.SelectAttr("cfStartDate"))
			if err != nil {
				return nil, err
			}
			t2, err := time.Parse(time.RFC3339, n.SelectAttr("cfEndDate"))
			if err != nil {
				return nil, err
			}
			msg.Attributes = append(msg.Attributes, Attribute{
				Name:      "name_" + n.SelectAttr("cfLangCode"),
				Value:     strings.TrimSpace(n.InnerText()),
				StartDate: &t1,
				EndDate:   models.NormalizeEndOfTime(t2),
			})
		}

		for _, v := range ValuesByClassName(doc, "cfOrgUnit_Class", "/be.ugent/organisatie/type/campus", "") {
			startDate := v.StartDate
			endDate := v.EndDate
			msg.Attributes = append(msg.Attributes, Attribute{
				Name:      "type",
				Value:     "campus",
				StartDate: &startDate,
				EndDate:   models.NormalizeEndOfTime(endDate),
			})
		}
		for _, v := range ValuesByClassName(doc, "cfOrgUnit_Class", "/be.ugent/organisatie/type/vakgroep", "") {
			startDate := v.StartDate
			endDate := v.EndDate
			msg.Attributes = append(msg.Attributes, Attribute{
				Name:      "type",
				Value:     "department",
				StartDate: &startDate,
				EndDate:   models.NormalizeEndOfTime(endDate),
			})
		}
		for _, v := range ValuesByClassName(doc, "cfOrgUnit_Class", "/be.ugent/organisatie/type/faculteit", "") {
			startDate := v.StartDate
			endDate := v.EndDate
			msg.Attributes = append(msg.Attributes, Attribute{
				Name:      "type",
				Value:     "faculty",
				StartDate: &startDate,
				EndDate:   models.NormalizeEndOfTime(endDate),
			})
		}
		for _, v := range ValuesByClassName(doc, "cfOrgUnit_Class", "/be.ugent/organisatie/type/universiteit", "") {
			startDate := v.StartDate
			endDate := v.EndDate
			msg.Attributes = append(msg.Attributes, Attribute{
				Name:      "type",
				Value:     "university",
				StartDate: &startDate,
				EndDate:   models.NormalizeEndOfTime(endDate),
			})
		}
		for _, v := range ValuesByClassName(doc, "cfOrgUnit_OrgUnit", "/be.ugent/gismo/organisatie-organisatie/type/kind-van", "cfOrgUnitId1") {
			startDate := v.StartDate
			endDate := v.EndDate
			msg.Attributes = append(msg.Attributes, Attribute{
				Name:      "parent_id",
				Value:     v.Value,
				StartDate: &startDate,
				EndDate:   models.NormalizeEndOfTime(endDate),
			})
		}
		// e.g. 000006047
		for _, v := range ValuesByClassName(doc, "cfFedId", "/be.ugent/gismo/organisatie/federated-id/memorialis", "cfFedId") {
			startDate := v.StartDate
			endDate := v.EndDate
			msg.Attributes = append(msg.Attributes, Attribute{
				Name:      "ugent_memorialis_id",
				Value:     v.Value,
				StartDate: &startDate,
				EndDate:   models.NormalizeEndOfTime(endDate),
			})
		}
		// e.g. "WE03V"
		for _, v := range ValuesByClassName(doc, "cfFedId", "/be.ugent/gismo/organisatie/federated-id/org-code", "cfFedId") {
			startDate := v.StartDate
			endDate := v.EndDate
			msg.Attributes = append(msg.Attributes, Attribute{
				Name:      "code",
				Value:     v.Value,
				StartDate: &startDate,
				EndDate:   models.NormalizeEndOfTime(endDate),
			})
		}
		// e.g. WE03*
		for _, v := range ValuesByClassName(doc, "cfFedId", "/be.ugent/gismo/organisatie/federated-id/biblio-code", "cfFedId") {
			startDate := v.StartDate
			endDate := v.EndDate
			msg.Attributes = append(msg.Attributes, Attribute{
				Name:      "biblio_code",
				Value:     v.Value,
				StartDate: &startDate,
				EndDate:   models.NormalizeEndOfTime(endDate),
			})
		}

	}

	return msg, nil
}

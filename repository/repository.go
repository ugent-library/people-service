package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/samber/lo"
	"github.com/ugent-library/crypt"
	"github.com/ugent-library/people-service/models"
)

const (
	personPageLimit          = 200
	organizationPageLimit    = 200
	organizationSuggestLimit = 10
	personSuggestLimit       = 10
)

type repository struct {
	client *sql.DB
	secret []byte
}
type setCursor struct {
	// IMPORTANT: auto increment (of id) starts with 1, so default value 0 should never match
	LastID int `json:"l"`
}

type organizationParent struct {
	id                           int
	dateCreated                  *time.Time
	dateUpdated                  *time.Time
	organizationID               int
	parentOrganizationID         int
	parentOrganizationExternalID string
}

type organizationMember struct {
	id                     int
	dateCreated            *time.Time
	dateUpdated            *time.Time
	personID               int
	organizationID         int
	organizationExternalID string
}

type organizationIdentifier struct {
	id             int
	organizationID int
	value          string
}

type personIdentifier struct {
	id       int
	personID int
	value    string
}

func NewRepository(config *Config) (*repository, error) {
	client, err := openClient(config.DbUrl)
	if err != nil {
		return nil, err
	}
	return &repository{
		client: client,
		secret: []byte(config.AesKey),
	}, nil
}

func (repo *repository) getOrganizationIdentifiers(ctx context.Context, ids ...int) ([]*organizationIdentifier, error) {
	query := `SELECT "id", "organization_id", "value" FROM "organization_identifiers" WHERE "organization_id" = any($1)`
	rows, err := repo.client.QueryContext(
		ctx,
		query,
		pgIntArray(ids),
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	organizationIdentifiers := []*organizationIdentifier{}

	for rows.Next() {
		oid := &organizationIdentifier{}
		err := rows.Scan(&oid.id, &oid.organizationID, &oid.value)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		organizationIdentifiers = append(organizationIdentifiers, oid)
	}

	return organizationIdentifiers, nil
}

func (repo *repository) getPersonIdentifiers(ctx context.Context, ids ...int) ([]*personIdentifier, error) {
	query := `SELECT "id", "person_id", "value" FROM "person_identifiers" WHERE "person_id" = any($1)`
	rows, err := repo.client.QueryContext(
		ctx,
		query,
		pgIntArray(ids),
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	pids := []*personIdentifier{}

	for rows.Next() {
		pid := &personIdentifier{}
		err := rows.Scan(&pid.id, &pid.personID, &pid.value)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		pids = append(pids, pid)
	}

	return pids, nil
}

func (repo *repository) getOrganizationMembers(ctx context.Context, ids ...int) ([]organizationMember, error) {
	query := `
SELECT
	"id",
	"organization_id",
    "person_id",
	"date_created",
	"date_updated",
	(SELECT "external_id" FROM "organizations" WHERE "id" = op.organization_id) AS "organization_external_id"
FROM "organization_members" op
WHERE "person_id" = any($1)
ORDER by "organization_id" ASC
	`
	rows, err := repo.client.QueryContext(
		ctx,
		query,
		pgIntArray(ids),
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	organizationMembers := []organizationMember{}

	for rows.Next() {
		om := organizationMember{}
		err := rows.Scan(
			&om.id,
			&om.organizationID,
			&om.personID,
			&om.dateCreated,
			&om.dateUpdated,
			&om.organizationExternalID,
		)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		organizationMembers = append(organizationMembers, om)
	}

	return organizationMembers, nil
}

func (repo *repository) getOrganizationParents(ctx context.Context, ids ...int) ([]organizationParent, error) {
	query := `
SELECT
	"id",
	"organization_id",
    "parent_organization_id",
	"date_created",
	"date_updated",
	(SELECT "external_id" FROM "organizations" WHERE "id" = op.parent_organization_id) AS "parent_organization_external_id"
FROM "organization_parents" op
WHERE "organization_id" = any($1)
ORDER by "parent_organization_id" ASC
	`
	rows, err := repo.client.QueryContext(
		ctx,
		query,
		pgIntArray(ids),
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	organizationParents := []organizationParent{}

	for rows.Next() {
		op := organizationParent{}
		err := rows.Scan(
			&op.id,
			&op.organizationID,
			&op.parentOrganizationID,
			&op.dateCreated,
			&op.dateUpdated,
			&op.parentOrganizationExternalID,
		)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		organizationParents = append(organizationParents, op)
	}

	return organizationParents, nil
}

func (repo *repository) GetOrganization(ctx context.Context, id string) (*models.Organization, error) {
	query := `
SELECT 
	"id", 
	"external_id", 
	"date_created", 
	"date_updated", 
	"name_dut",
	"name_eng",
	"acronym",
	"type"
FROM organizations WHERE external_id = $1 LIMIT 1`

	org := &models.Organization{}
	var rowID int
	err := repo.client.QueryRowContext(ctx, query, id).Scan(
		&rowID,
		&org.ID,
		&org.DateCreated,
		&org.DateUpdated,
		&org.NameDut,
		&org.NameEng,
		&org.Acronym,
		&org.Type,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNotFound
	}

	organizationParents, err := repo.getOrganizationParents(ctx, rowID)
	if err != nil {
		return nil, err
	}

	for _, op := range organizationParents {
		org.Parent = append(org.Parent, &models.OrganizationParent{
			ID:          op.parentOrganizationExternalID,
			DateCreated: op.dateCreated,
			DateUpdated: op.dateUpdated,
		})
	}

	oids, err := repo.getOrganizationIdentifiers(ctx, rowID)
	if err != nil {
		return nil, err
	}

	for _, oid := range oids {
		urn, _ := models.ParseURN(oid.value)
		org.AddIdentifier(urn)
	}

	return org, nil
}

func (repo *repository) GetOrganizationsByIdentifier(ctx context.Context, urns ...models.URN) ([]*models.Organization, error) {
	urnValues := make([]string, 0, len(urns))
	for _, urn := range urns {
		urnValues = append(urnValues, urn.String())
	}

	orgs := []*models.Organization{}
	rowIDS := []int{}

	query := `
	SELECT "id", "external_id", "date_created", "date_updated", "type", "name_dut", "name_eng", "acronym"
	FROM "organizations" WHERE "id" IN (
		SELECT "organization_id" FROM "organization_identifiers" WHERE "value" = any($1)
	)
	`

	rows, err := repo.client.QueryContext(ctx, query, pgTextArray(urnValues))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rowID int
		org := &models.Organization{}
		err = rows.Scan(
			&rowID,
			&org.ID,
			&org.DateCreated,
			&org.DateUpdated,
			&org.Type,
			&org.NameDut,
			&org.NameEng,
			&org.Acronym,
		)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		rowIDS = append(rowIDS, rowID)
		orgs = append(orgs, org)
	}

	if len(orgs) == 0 {
		return orgs, nil
	}

	allOrganizationParents, err := repo.getOrganizationParents(ctx, rowIDS...)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(orgs); i++ {
		rowID := rowIDS[i]
		org := orgs[i]
		for _, op := range allOrganizationParents {
			if op.organizationID == rowID {
				org.AddParent(&models.OrganizationParent{
					ID:          op.parentOrganizationExternalID,
					DateCreated: op.dateCreated,
					DateUpdated: op.dateUpdated,
				})
			}
		}
	}

	allOrganizationIdentifiers, err := repo.getOrganizationIdentifiers(ctx, rowIDS...)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(orgs); i++ {
		rowID := rowIDS[i]
		org := orgs[i]
		for _, oid := range allOrganizationIdentifiers {
			if oid.organizationID == rowID {
				urn, _ := models.ParseURN(oid.value)
				org.AddIdentifier(urn)
			}
		}
	}

	// TODO: order by array_position cannot work on array itself. Find another way
	return orgs, nil
}

func (repo *repository) SaveOrganization(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	if org.IsStored() {
		return repo.UpdateOrganization(ctx, org)
	}
	return repo.CreateOrganization(ctx, org)
}

func (repo *repository) CreateOrganization(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	now := time.Now().UTC()
	org.DateCreated = &now
	org.DateUpdated = &now
	org.ID = ulid.Make().String()
	for _, parent := range org.Parent {
		parent.DateCreated = &now
		parent.DateUpdated = &now
	}

	tx, err := repo.client.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: %w", err)
	}
	defer tx.Rollback()

	// add organization
	query := `
	INSERT INTO "organizations" (
		"external_id",
		"date_created",
		"date_updated",
		"name_dut",
		"name_eng",
		"type",
		"acronym",
		"ts_vals"
	)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING "id"
	`
	tsVals := []string{}
	tsVals = append(tsVals, org.NameDut, org.NameEng)
	tsVals = append(tsVals, org.Acronym)
	tsVals = append(tsVals, org.GetIdentifierValues()...)
	var rowID int
	err = tx.QueryRowContext(
		ctx, query,
		org.ID,
		org.DateCreated,
		org.DateUpdated,
		org.NameDut,
		org.NameEng,
		org.Type,
		org.Acronym,
		toJSON(tsVals),
	).Scan(&rowID)
	if err != nil {
		return nil, err
	}

	// add parents
	parentOrganizationExternalIds := []string{}
	parentOrganizationIDs := []int{}
	for _, parent := range org.Parent {
		parentOrganizationExternalIds = append(parentOrganizationExternalIds, parent.ID)
	}
	parentOrganizationExternalIds = lo.Uniq(parentOrganizationExternalIds)
	query = `SELECT "id" FROM "organizations" WHERE "external_id" = any($1) ORDER BY array_position($1, external_id)`
	parentRows, err := tx.QueryContext(
		ctx,
		query,
		pgTextArray(parentOrganizationExternalIds),
	)
	if err != nil {
		return nil, err
	}
	defer parentRows.Close()
	for parentRows.Next() {
		var rowID int
		parentRows.Scan(&rowID)
		parentOrganizationIDs = append(parentOrganizationIDs, rowID)
	}
	if len(parentOrganizationExternalIds) != len(parentOrganizationIDs) {
		return nil, models.ErrInvalidReference
	}
	query = `
INSERT INTO "organization_parents"
	("organization_id", "parent_organization_id", "date_created", "date_updated")
	VALUES($1, $2, $3, $4);
`
	for i := 0; i < len(parentOrganizationIDs); i++ {
		parentOrganizationID := parentOrganizationIDs[i]
		orgParent := org.Parent[i]
		fmt.Fprintf(os.Stderr, "parent org: %+v\n", orgParent)
		_, err := tx.ExecContext(ctx, query, rowID, parentOrganizationID, orgParent.DateCreated, orgParent.DateUpdated)
		if err != nil {
			return nil, err
		}
	}

	// add identifiers
	query = `
INSERT INTO "organization_identifiers"("organization_id", "date_created", "date_updated", "value")
VALUES($1, $2, $3, $4)
	`
	for _, urn := range org.Identifier {
		tx.ExecContext(ctx, query, rowID, now, now, urn.String())
	}

	// commit
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("unable to commit transaction: %w", err)
	}

	return org, nil
}

func (repo *repository) UpdateOrganization(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	now := time.Now().UTC()
	org.DateUpdated = &now
	for _, parent := range org.Parent {
		if parent.DateCreated == nil {
			parent.DateCreated = &now
		}
		if parent.DateUpdated == nil {
			parent.DateUpdated = &now
		}
	}

	// start transaction
	tx, err := repo.client.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: %w", err)
	}
	defer tx.Rollback()

	// update organization
	tsVals := []string{}
	tsVals = append(tsVals, org.NameDut, org.NameEng)
	tsVals = append(tsVals, org.Acronym)
	tsVals = append(tsVals, org.GetIdentifierValues()...)

	query := `
UPDATE "organizations"
SET
	"date_updated" = $2,
	"name_dut" = $3,
	"name_eng" = $4,
	"type" = $5,
	"acronym" = $6,
	"ts_vals" = $7
WHERE "external_id" = $1
RETURNING "id"
	`
	var rowID int
	err = tx.QueryRowContext(
		ctx,
		query,
		org.ID,
		now,
		org.NameDut,
		org.NameEng,
		org.Type,
		org.Acronym,
		toJSON(tsVals),
	).Scan(&rowID)
	if err != nil {
		return nil, err
	}

	// update organization parents
	var newOrganizationParents []organizationParent
	if len(org.Parent) > 0 {
		parentOrganizationExternalIDs := []string{}
		organizationParents := []*organizationParent{}

		for _, parent := range org.Parent {
			parentOrganizationExternalIDs = append(parentOrganizationExternalIDs, parent.ID)
		}
		parentOrganizationExternalIDs = lo.Uniq(parentOrganizationExternalIDs)

		pgExternalIds := pgTextArray(parentOrganizationExternalIDs)
		rows, err := tx.QueryContext(ctx, "SELECT id, external_id FROM organizations WHERE external_id = any($1) ORDER BY array_position($1, external_id)", pgExternalIds)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var parentID int
			var parentExternalID string
			err = rows.Scan(&parentID, &parentExternalID)
			if errors.Is(err, sql.ErrNoRows) {
				return nil, models.ErrInvalidReference
			}
			if err != nil {
				return nil, err
			}
			organizationParents = append(organizationParents, &organizationParent{
				organizationID:               rowID,
				parentOrganizationID:         parentID,
				parentOrganizationExternalID: parentExternalID,
			})
		}
		if len(parentOrganizationExternalIDs) != len(organizationParents) {
			return nil, models.ErrInvalidReference
		}

		for _, parent := range org.Parent {
			newOrganizationParent := organizationParent{
				organizationID: rowID,
				dateCreated:    parent.DateCreated,
				dateUpdated:    parent.DateUpdated,
			}
			for _, op := range organizationParents {
				if op.parentOrganizationExternalID == parent.ID {
					newOrganizationParent.parentOrganizationID = op.parentOrganizationID
					break
				}
			}
			newOrganizationParents = append(newOrganizationParents, newOrganizationParent)
		}
	}

	updatedRelIds := []int{}
	query = `
INSERT INTO "organization_parents"("organization_id", "parent_organization_id", "date_created", "date_updated")
VALUES($1, $2, $3, $4)
ON CONFLICT("organization_id", "parent_organization_id")
DO UPDATE SET date_updated = EXCLUDED.date_updated
RETURNING "id"
	`

	if len(newOrganizationParents) > 0 {
		for _, newOrganizationParent := range newOrganizationParents {
			var relId int
			err = tx.QueryRowContext(ctx, query, rowID, newOrganizationParent.parentOrganizationID, newOrganizationParent.dateCreated, newOrganizationParent.dateUpdated).Scan(&relId)
			if err != nil {
				return nil, err
			}
			updatedRelIds = append(updatedRelIds, relId)
		}
	}

	query = `DELETE FROM "organization_parents" WHERE "organization_id" = $1 AND NOT "id" = any($2)`
	_, err = tx.ExecContext(ctx, query, rowID, pgIntArray(updatedRelIds))
	if err != nil {
		return nil, err
	}

	// update identifiers
	updatedRelIds = []int{}
	urnValues := org.GetIdentifierQualifiedValues()

	if len(urnValues) > 0 {
		query = `
INSERT INTO "organization_identifiers"("organization_id", "value", "date_created", "date_updated")
VALUES($1, $2, $3, $4)
ON CONFLICT("value")
DO UPDATE SET date_updated = EXCLUDED.date_updated
RETURNING "id"
`
		for _, urnValue := range urnValues {
			var relID int
			err = tx.QueryRowContext(ctx, query, rowID, urnValue, now, now).Scan(&relID)
			if err != nil {
				return nil, err
			}
			updatedRelIds = append(updatedRelIds, relID)
		}
	}

	query = `DELETE FROM "organization_identifiers" WHERE "organization_id" = $1 AND NOT "id" = any($2)`
	_, err = tx.ExecContext(ctx, query, rowID, pgIntArray(updatedRelIds))
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("unable to commit transaction: %w", err)
	}

	return org, nil
}

func (repo *repository) DeleteOrganization(ctx context.Context, id string) error {
	_, err := repo.client.ExecContext(ctx, "DELETE FROM organizations WHERE external_id = $1", id)
	return err
}

func (repo *repository) EachOrganization(ctx context.Context, cb func(*models.Organization) bool) error {
	cursor := setCursor{}

	for {
		organizations, newCursor, err := repo.getOrganizations(ctx, cursor)
		if err != nil {
			return err
		}

		for _, organization := range organizations {
			if !cb(organization) {
				return nil
			}
		}

		if len(organizations) == 0 {
			break
		}
		if newCursor.LastID <= 0 {
			break
		}
		cursor = newCursor
	}

	return nil
}

func (repo *repository) SuggestOrganizations(ctx context.Context, query string) ([]*models.Organization, error) {
	tsQuery, tsQueryArgs := toTSQuery(query)

	sqlQuery := fmt.Sprintf(
		`SELECT "id", "external_id", "date_created", "date_updated", "type", "name_dut", "name_eng", "acronym" FROM "organizations" WHERE ts @@ %s LIMIT %d`,
		tsQuery,
		organizationSuggestLimit)

	rows, err := repo.client.QueryContext(ctx, sqlQuery, tsQueryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orgRowIDs := []int{}
	orgs := []*models.Organization{}

	for rows.Next() {
		var rowID int
		org := &models.Organization{}
		err = rows.Scan(
			&rowID,
			&org.ID,
			&org.DateCreated,
			&org.DateUpdated,
			&org.Type,
			&org.NameDut,
			&org.NameEng,
			&org.Acronym,
		)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		orgRowIDs = append(orgRowIDs, rowID)
		orgs = append(orgs, org)
	}

	if len(orgs) == 0 {
		return nil, nil
	}

	allOrganizationParents, err := repo.getOrganizationParents(ctx, orgRowIDs...)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(orgs); i++ {
		org := orgs[i]
		orgRowID := orgRowIDs[i]
		for _, op := range allOrganizationParents {
			if op.organizationID == orgRowID {
				org.AddParent(&models.OrganizationParent{
					ID:          op.parentOrganizationExternalID,
					DateCreated: op.dateCreated,
					DateUpdated: op.dateUpdated,
				})
			}
		}
	}

	allOrganizationIdentifiers, err := repo.getOrganizationIdentifiers(ctx, orgRowIDs...)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(orgs); i++ {
		org := orgs[i]
		orgRowID := orgRowIDs[i]
		for _, orgIdentifier := range allOrganizationIdentifiers {
			if orgIdentifier.organizationID == orgRowID {
				urn, _ := models.ParseURN(orgIdentifier.value)
				org.AddIdentifier(urn)
			}
		}
	}

	return orgs, nil
}

func (repo *repository) GetOrganizations(ctx context.Context) ([]*models.Organization, string, error) {
	organizations, newCursor, err := repo.getOrganizations(ctx, setCursor{})
	if err != nil {
		return nil, "", err
	}

	var encodedCursor string
	if newCursor.LastID > 0 {
		encodedCursor, err = repo.encodeCursor(newCursor)
		if err != nil {
			return nil, "", err
		}
	}
	return organizations, encodedCursor, nil
}

func (repo *repository) GetMoreOrganizations(ctx context.Context, tokenValue string) ([]*models.Organization, string, error) {
	cursor := setCursor{}
	if err := repo.decodeCursor(tokenValue, &cursor); err != nil {
		return nil, "", err
	}
	organizations, newCursor, err := repo.getOrganizations(ctx, cursor)
	if err != nil {
		return nil, "", err
	}

	var encodedCursor string
	if newCursor.LastID > 0 {
		encodedCursor, err = repo.encodeCursor(newCursor)
		if err != nil {
			return nil, "", err
		}
	}

	return organizations, encodedCursor, nil
}

func (repo *repository) getOrganizations(ctx context.Context, cursor setCursor) ([]*models.Organization, setCursor, error) {
	newCursor := setCursor{}

	// get organizations
	query := `
SELECT 
	"id",
	"external_id",
	"date_created",
	"date_updated",
	"type",
	"name_dut",
	"name_eng", 
	"acronym"
FROM "organizations"
WHERE "id" > $1 ORDER BY "id" ASC LIMIT ` + fmt.Sprintf("%d", organizationPageLimit)

	rows, err := repo.client.QueryContext(ctx, query, cursor.LastID)
	if err != nil {
		return nil, newCursor, err
	}
	defer rows.Close()

	orgRowIDs := []int{}
	orgs := []*models.Organization{}
	for rows.Next() {
		var rowID int
		org := &models.Organization{}
		err = rows.Scan(
			&rowID,
			&org.ID,
			&org.DateCreated,
			&org.DateUpdated,
			&org.Type,
			&org.NameDut,
			&org.NameEng,
			&org.Acronym,
		)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, newCursor, nil
		}
		if err != nil {
			return nil, newCursor, err
		}
		orgRowIDs = append(orgRowIDs, rowID)
		orgs = append(orgs, org)
	}
	if len(orgs) == 0 {
		return nil, newCursor, nil
	}

	// get uncapped total
	var total int
	err = repo.client.QueryRowContext(ctx, `SELECT COUNT(*) "total" FROM "organizations"`).Scan(&total)
	if err != nil {
		return nil, newCursor, err
	}

	// attach organization parents
	allOrganizationParents, err := repo.getOrganizationParents(ctx, orgRowIDs...)
	if err != nil {
		return nil, newCursor, err
	}
	for i := 0; i < len(orgRowIDs); i++ {
		rowID := orgRowIDs[i]
		org := orgs[i]
		for _, op := range allOrganizationParents {
			if op.organizationID == rowID {
				org.AddParent(&models.OrganizationParent{
					ID:          op.parentOrganizationExternalID,
					DateCreated: op.dateCreated,
					DateUpdated: op.dateUpdated,
				})
			}
		}
	}

	// attach identifiers
	allOrganizationIdentifiers, err := repo.getOrganizationIdentifiers(ctx, orgRowIDs...)
	if err != nil {
		return nil, setCursor{}, err
	}
	for i := 0; i < len(orgRowIDs); i++ {
		rowID := orgRowIDs[i]
		org := orgs[i]
		for _, oid := range allOrganizationIdentifiers {
			if oid.organizationID == rowID {
				urn, _ := models.ParseURN(oid.value)
				org.AddIdentifier(urn)
			}
		}
	}

	// set next cursor
	if total > len(orgRowIDs) {
		newCursor = setCursor{
			LastID: orgRowIDs[len(orgRowIDs)-1],
		}
	}

	return orgs, newCursor, nil
}

func (repo *repository) SavePerson(ctx context.Context, p *models.Person) (*models.Person, error) {
	if p.IsStored() {
		return repo.UpdatePerson(ctx, p)
	}
	return repo.CreatePerson(ctx, p)
}

func (repo *repository) CreatePerson(ctx context.Context, p *models.Person) (*models.Person, error) {
	now := time.Now().UTC()
	p.DateCreated = &now
	p.DateUpdated = &now
	p.ID = ulid.Make().String()
	for _, orgMember := range p.Organization {
		orgMember.DateCreated = &now
		orgMember.DateUpdated = &now
	}

	tx, err := repo.client.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
INSERT INTO "people"
	(
		"external_id",
		"date_created",
		"date_updated",
		"active",
		"birth_date",
		"job_category",
		"email",
		"given_name",
		"preferred_given_name",
		"name",
		"family_name",
		"preferred_family_name",
		"honorific_prefix",
		"role",
		"settings",
		"object_class",
		"expiration_date",
		"token",
		"ts_vals"
	)
	VALUES
	(
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9,
		$10,
		$11,
		$12,
		$13,
		$14,
		$15,
		$16,
		$17,
		$18,
		$19
	)
	RETURNING "id"
	`
	var rowID int
	queryArgs := []any{}
	queryArgs = append(queryArgs, p.ID)
	queryArgs = append(queryArgs, p.DateCreated)
	queryArgs = append(queryArgs, p.DateUpdated)
	queryArgs = append(queryArgs, p.Active)
	queryArgs = append(queryArgs, p.BirthDate)
	queryArgs = append(queryArgs, p.JobCategory)
	queryArgs = append(queryArgs, p.Email)
	queryArgs = append(queryArgs, p.GivenName)
	queryArgs = append(queryArgs, p.PreferredGivenName)
	queryArgs = append(queryArgs, p.Name)
	queryArgs = append(queryArgs, p.FamilyName)
	queryArgs = append(queryArgs, p.PreferredFamilyName)
	queryArgs = append(queryArgs, p.HonorificPrefix)
	queryArgs = append(queryArgs, toJSON(p.Role))
	queryArgs = append(queryArgs, toJSON(p.Settings))
	queryArgs = append(queryArgs, toJSON(p.ObjectClass))
	queryArgs = append(queryArgs, p.ExpirationDate)
	tokens := make([]string, 0, len(p.Token))
	for _, token := range p.Token {
		eToken, err := encryptMessage(repo.secret, token.Value)
		if err != nil {
			return nil, fmt.Errorf("unable to encrypt %s: %w", token.Namespace, err)
		}
		eURN := models.NewURN(token.Namespace, eToken)
		tokens = append(tokens, eURN.String())
	}
	queryArgs = append(queryArgs, toJSON(tokens))
	tsVals := []string{}
	if p.Name != "" {
		tsVals = append(tsVals, p.Name)
	}
	queryArgs = append(queryArgs, toJSON(tsVals))

	err = tx.QueryRowContext(ctx, query, queryArgs...).Scan(&rowID)
	if err != nil {
		return nil, err
	}

	if len(p.Organization) > 0 {
		var organizationExternalIDs []string
		for _, orgMember := range p.Organization {
			organizationExternalIDs = append(organizationExternalIDs, orgMember.ID)
		}
		organizationExternalIDs = lo.Uniq(organizationExternalIDs)

		orgRowIDS := []int{}
		rows, err := tx.QueryContext(
			ctx,
			`SELECT "id" FROM "organizations" WHERE "external_id" = any($1)`,
			pgTextArray(organizationExternalIDs))
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var rowID int
			err = rows.Scan(&rowID)
			if errors.Is(err, sql.ErrNoRows) {
				return nil, models.ErrInvalidReference
			}
			if err != nil {
				return nil, err
			}
			orgRowIDS = append(orgRowIDS, rowID)
		}

		if len(organizationExternalIDs) != len(orgRowIDS) {
			return nil, fmt.Errorf("%w: person.organization_id contains invalid organization id's", models.ErrInvalidReference)
		}

		for i := 0; i < len(orgRowIDS); i++ {
			insertQuery := `
			INSERT INTO "organization_members"("date_created", "date_updated", "organization_id", "person_id")
			VALUES($1, $2, $3, $4)
			`
			_, err = tx.ExecContext(ctx, insertQuery, now, now, orgRowIDS[i], rowID)
			if err != nil {
				return nil, err
			}
		}

	}

	// add identifiers
	query = `
INSERT INTO "person_identifiers"("person_id", "date_created", "date_updated", "value")
VALUES($1, $2, $3, $4)
	`
	for _, urnValue := range p.GetIdentifierQualifiedValues() {
		tx.ExecContext(ctx, query, rowID, now, now, urnValue)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("unable to commit transaction: %w", err)
	}

	return p, nil
}

// TODO: make transaction safe
func (repo *repository) SetPersonOrcid(ctx context.Context, id string, orcid string) error {
	person, err := repo.GetPerson(ctx, id)
	if err != nil {
		return err
	}

	if orcid == "" {
		newIdentifiers := make([]*models.URN, 0, len(person.Identifier))
		for _, urn := range person.Identifier {
			if urn.Namespace != "orcid" {
				newIdentifiers = append(newIdentifiers, urn)
			}
		}
		person.ClearIdentifier()
		for _, urn := range newIdentifiers {
			person.AddIdentifier(urn)
		}
	} else {
		for _, urn := range person.Identifier {
			if urn.Namespace == "orcid" {
				urn.Value = orcid
			}
		}
	}

	_, err = repo.UpdatePerson(ctx, person)

	return err
}

// TODO: make transaction safe
func (repo *repository) SetPersonOrcidToken(ctx context.Context, id string, orcidToken string) error {
	person, err := repo.GetPerson(ctx, id)
	if err != nil {
		return err
	}

	if orcidToken == "" {
		person.Token = lo.Filter(person.Token, func(token *models.URN, idx int) bool {
			return token.Namespace != "orcid"
		})
	} else {
		person.ClearToken()
		person.AddToken("orcid", orcidToken)
	}

	_, err = repo.UpdatePerson(ctx, person)
	return err
}

func (repo *repository) UpdatePerson(ctx context.Context, p *models.Person) (*models.Person, error) {
	now := time.Now().UTC()
	p.DateUpdated = &now
	for _, orgMember := range p.Organization {
		if orgMember.DateCreated == nil {
			orgMember.DateCreated = &now
		}
		if orgMember.DateUpdated == nil {
			orgMember.DateUpdated = &now
		}
	}

	tx, err := repo.client.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// update person
	query := `
UPDATE "people"
SET "date_updated" = $1,
	"active" = $2,
	"birth_date" = $3,
	"job_category" = $4,
	"email" = $5,
	"given_name" = $6,
	"preferred_given_name" = $7,
	"name" = $8,
	"family_name" = $9,
	"preferred_family_name" = $10,
	"honorific_prefix" = $11,
	"role" = $12,
	"settings" = $13,
	"object_class" = $14,
	"expiration_date" = $15,
	"token" = $16,
	"ts_vals" = $17
WHERE "external_id" = $18 
RETURNING "id"
	`
	var rowID int
	queryArgs := []any{}
	queryArgs = append(queryArgs, p.DateUpdated)
	queryArgs = append(queryArgs, p.Active)
	queryArgs = append(queryArgs, p.BirthDate)
	queryArgs = append(queryArgs, p.JobCategory)
	queryArgs = append(queryArgs, p.Email)
	queryArgs = append(queryArgs, p.GivenName)
	queryArgs = append(queryArgs, p.PreferredGivenName)
	queryArgs = append(queryArgs, p.Name)
	queryArgs = append(queryArgs, p.FamilyName)
	queryArgs = append(queryArgs, p.PreferredFamilyName)
	queryArgs = append(queryArgs, p.HonorificPrefix)
	queryArgs = append(queryArgs, toJSON(p.Role))
	queryArgs = append(queryArgs, toJSON(p.Settings))
	queryArgs = append(queryArgs, toJSON(p.ObjectClass))
	queryArgs = append(queryArgs, p.ExpirationDate)
	tokens := make([]string, 0, len(p.Token))
	for _, token := range p.Token {
		eToken, err := encryptMessage(repo.secret, token.Value)
		if err != nil {
			return nil, fmt.Errorf("unable to encrypt %s: %w", token.Namespace, err)
		}
		eURN := models.NewURN(token.Namespace, eToken)
		tokens = append(tokens, eURN.String())
	}
	queryArgs = append(queryArgs, toJSON(tokens))
	tsVals := []string{}
	if p.Name != "" {
		tsVals = append(tsVals, p.Name)
	}
	queryArgs = append(queryArgs, toJSON(tsVals))
	queryArgs = append(queryArgs, p.ID)

	err = tx.QueryRowContext(ctx, query, queryArgs...).Scan(&rowID)
	if err != nil {
		return nil, err
	}

	// update "organization_members"
	updatedOrganizationMemberIds := []int{}
	if len(p.Organization) > 0 {
		orgRowIDs := make([]int, 0, len(p.Organization))
		orgExternalIDs := make([]string, 0, len(p.Organization))
		for _, orgMem := range p.Organization {
			orgExternalIDs = append(orgExternalIDs, orgMem.ID)
		}
		orgExternalIDs = lo.Uniq(orgExternalIDs)
		rows, err := tx.QueryContext(
			ctx,
			`SELECT "id" FROM "organizations" WHERE "external_id" = any($1) ORDER BY array_position($1, external_id)`,
			pgTextArray(orgExternalIDs),
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var rowID int
			if err = rows.Scan(&rowID); err != nil {
				return nil, err
			}
			orgRowIDs = append(orgRowIDs, rowID)
		}
		orgRowIDs = lo.Uniq(orgRowIDs)
		if len(orgExternalIDs) != len(orgRowIDs) {
			return nil, models.ErrInvalidReference
		}

		for i := 0; i < len(orgRowIDs); i++ {
			orgRowID := orgRowIDs[i]
			memberOrg := p.Organization[i]
			insertQuery := `
			INSERT INTO "organization_members"
				("date_created", "date_updated", "person_id", "organization_id")
			VALUES($1, $2, $3, $4)
			ON CONFLICT("person_id", "organization_id")
			DO UPDATE SET date_updated = EXCLUDED.date_updated
			RETURNING "id"
			`
			var relID int
			err = tx.QueryRowContext(ctx, insertQuery, memberOrg.DateCreated, memberOrg.DateUpdated, rowID, orgRowID).Scan(&relID)
			if err != nil {
				return nil, err
			}
			updatedOrganizationMemberIds = append(updatedOrganizationMemberIds, relID)
		}
	}
	query = `DELETE FROM "organization_members" WHERE "person_id" = $1`
	queryArgs = []any{rowID}
	if len(updatedOrganizationMemberIds) > 0 {
		query += ` AND NOT "id" = any($2)`
		queryArgs = append(queryArgs, pgIntArray(updatedOrganizationMemberIds))
	}
	_, err = tx.ExecContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}

	// "organization_identifiers"
	updatedOrganizationIdentifierIds := []int{}
	if len(p.Identifier) > 0 {
		for _, urnValue := range p.GetIdentifierQualifiedValues() {
			insertQuery := `
			INSERT INTO "organization_identifiers"
				("date_created", "date_updated", "person_id", "value")
			VALUES($1, $2, $3, $4)
			ON CONFLICT("value")
			DO UPDATE SET date_updated = EXCLUDED.date_updated
			RETURNING "id"
			`
			var relID int
			err = tx.QueryRowContext(ctx, insertQuery, now, now, rowID, urnValue).Scan(&relID)
			if err != nil {
				return nil, err
			}
			updatedOrganizationIdentifierIds = append(updatedOrganizationIdentifierIds, relID)
		}
	}
	query = `DELETE FROM "organization_identifiers" WHERE "person_id" = $1`
	queryArgs = []any{rowID}
	if len(updatedOrganizationMemberIds) > 0 {
		query += ` AND NOT "id" = any($2)`
		queryArgs = append(queryArgs, pgIntArray(updatedOrganizationIdentifierIds))
	}
	_, err = tx.ExecContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return p, nil
}

func (repo *repository) GetPerson(ctx context.Context, id string) (*models.Person, error) {
	query := `
SELECT
	"id",
	"date_created",
	"date_updated",
	"external_id",
	"active",
	"birth_date",
	"email",
	"given_name",
	"name",
	"family_name",
	"job_category",
	"preferred_given_name",
	"preferred_family_name",
	"honorific_prefix", 
	"role",
	"settings", 
	"object_class",
	"expiration_date",
	"token"
FROM "people" WHERE "external_id" = $1
LIMIT 1
	`

	var rowID int
	var encTokens []string
	p := &models.Person{}
	err := repo.client.QueryRowContext(ctx, query, id).Scan(
		&rowID,
		&p.DateCreated,
		&p.DateUpdated,
		&p.ID,
		&p.Active,
		&p.BirthDate,
		&p.Email,
		&p.GivenName,
		&p.Name,
		&p.FamilyName,
		&p.JobCategory,
		&p.PreferredGivenName,
		&p.PreferredFamilyName,
		&p.HonorificPrefix,
		&p.Role,
		&p.Settings,
		&p.ObjectClass,
		&p.ExpirationDate,
		encTokens,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	p.Token, err = repo.decryptTokens(encTokens)
	if err != nil {
		return nil, err
	}

	personIdentifiers, err := repo.getPersonIdentifiers(ctx, rowID)
	if err != nil {
		return nil, err
	}
	for _, pid := range personIdentifiers {
		urn, _ := models.ParseURN(pid.value)
		p.AddIdentifier(urn)
	}

	orgMembers, err := repo.getOrganizationMembers(ctx, rowID)
	if err != nil {
		return nil, err
	}
	for _, orgMember := range orgMembers {
		p.AddOrganizationMember(&models.OrganizationMember{
			ID:          orgMember.organizationExternalID,
			DateCreated: orgMember.dateCreated,
			DateUpdated: orgMember.dateUpdated,
		})
	}

	return p, nil
}

func (repo *repository) GetPeopleByIdentifier(ctx context.Context, urns ...models.URN) ([]*models.Person, error) {
	query := `
	SELECT
		"id",
		"date_created",
		"date_updated",
		"external_id",
		"active",
		"birth_date",
		"email",
		"given_name",
		"name",
		"family_name",
		"job_category",
		"preferred_given_name",
		"preferred_family_name",
		"honorific_prefix", 
		"role",
		"settings", 
		"object_class",
		"expiration_date",
		"token"
	FROM "people"
	WHERE "id" IN (SELECT "person_id" FROM "person_identifiers" WHERE "value" = any($1))
	`

	ids := make([]string, 0, len(urns))
	for _, urn := range urns {
		ids = append(ids, urn.String())
	}
	rows, err := repo.client.QueryContext(ctx, query, pgTextArray(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rowIDs := []int{}
	people := []*models.Person{}

	for rows.Next() {
		var rowID int
		var encTokens []string
		p := &models.Person{}
		err = rows.Scan(
			&rowID,
			&p.DateCreated,
			&p.DateUpdated,
			&p.ID,
			&p.Active,
			&p.BirthDate,
			&p.Email,
			&p.GivenName,
			&p.Name,
			&p.FamilyName,
			&p.JobCategory,
			&p.PreferredGivenName,
			&p.PreferredFamilyName,
			&p.HonorificPrefix,
			&p.Role,
			&p.Settings,
			&p.ObjectClass,
			&p.ExpirationDate,
			encTokens,
		)
		if err != nil {
			return nil, err
		}
		p.Token, err = repo.decryptTokens(encTokens)
		if err != nil {
			return nil, err
		}
		rowIDs = append(rowIDs, rowID)
		people = append(people, p)
	}

	if len(people) == 0 {
		return people, nil
	}

	allPersonOrganizationMembers, err := repo.getOrganizationMembers(ctx, rowIDs...)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(rowIDs); i++ {
		rowID := rowIDs[i]
		person := people[i]
		for _, orgMember := range allPersonOrganizationMembers {
			if orgMember.personID == rowID {
				person.AddOrganizationMember(&models.OrganizationMember{
					ID:          orgMember.organizationExternalID,
					DateCreated: orgMember.dateCreated,
					DateUpdated: orgMember.dateUpdated,
				})
			}
		}
	}

	allPersonIdentifiers, err := repo.getPersonIdentifiers(ctx, rowIDs...)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(rowIDs); i++ {
		rowID := rowIDs[i]
		person := people[i]
		for _, pid := range allPersonIdentifiers {
			if pid.personID == rowID {
				urn, _ := models.ParseURN(pid.value)
				person.AddIdentifier(urn)
			}
		}
	}

	return people, nil
}

func (repo *repository) DeletePerson(ctx context.Context, id string) error {
	_, err := repo.client.ExecContext(ctx, `DELETE FROM "people" WHERE "external_id" = $1`, id)
	return err
}

func (repo *repository) EachPerson(ctx context.Context, cb func(*models.Person) bool) error {
	cursor := setCursor{}

	for {
		people, newCursor, err := repo.getPeople(ctx, cursor)
		if err != nil {
			return err
		}

		for _, person := range people {
			if !cb(person) {
				return nil
			}
		}

		if len(people) == 0 {
			break
		}
		if newCursor.LastID <= 0 {
			break
		}
		cursor = newCursor
	}

	return nil

}

func (repo *repository) SuggestPeople(ctx context.Context, query string) ([]*models.Person, error) {
	tsQuery, tsQueryArgs := toTSQuery(query)
	sqlQuery := `
SELECT
	"id",
	"date_created",
	"date_updated",
	"external_id",
	"active",
	"birth_date",
	"email",
	"given_name",
	"name",
	"family_name",
	"job_category",
	"preferred_given_name",
	"preferred_family_name",
	"honorific_prefix", 
	"role",
	"settings", 
	"object_class",
	"expiration_date",
	"token",
	ts_rank(ts, %s) AS rank
FROM "people" WHERE ts @@ %s ORDER BY "rank" DESC LIMIT %d
`
	sqlQuery = fmt.Sprintf(
		sqlQuery,
		tsQuery,
		tsQuery,
		personSuggestLimit,
	)
	rows, err := repo.client.QueryContext(ctx, sqlQuery, tsQueryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rowIDs := []int{}
	people := []*models.Person{}

	for rows.Next() {
		var rowID int
		var encTokens []string
		p := &models.Person{}
		err = rows.Scan(
			&rowID,
			&p.DateCreated,
			&p.DateUpdated,
			&p.ID,
			&p.Active,
			&p.BirthDate,
			&p.Email,
			&p.GivenName,
			&p.Name,
			&p.FamilyName,
			&p.JobCategory,
			&p.PreferredGivenName,
			&p.PreferredFamilyName,
			&p.HonorificPrefix,
			&p.Role,
			&p.Settings,
			&p.ObjectClass,
			&p.ExpirationDate,
			encTokens,
		)
		if err != nil {
			return nil, err
		}
		p.Token, err = repo.decryptTokens(encTokens)
		if err != nil {
			return nil, err
		}
		rowIDs = append(rowIDs, rowID)
		people = append(people, p)
	}

	if len(people) == 0 {
		return people, nil
	}

	allPersonOrganizationMembers, err := repo.getOrganizationMembers(ctx, rowIDs...)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(rowIDs); i++ {
		rowID := rowIDs[i]
		person := people[i]
		for _, orgMember := range allPersonOrganizationMembers {
			if orgMember.personID == rowID {
				person.AddOrganizationMember(models.OrganizationMember{
					ID:          orgMember.organizationExternalID,
					DateCreated: orgMember.dateCreated,
					DateUpdated: orgMember.dateUpdated,
				})
			}
		}
	}

	allPersonIdentifiers, err := repo.getPersonIdentifiers(ctx, rowIDs...)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(rowIDs); i++ {
		rowID := rowIDs[i]
		person := people[i]
		for _, pid := range allPersonIdentifiers {
			if pid.personID == rowID {
				urn, _ := models.ParseURN(pid.value)
				person.AddIdentifier(*urn)
			}
		}
	}

	return people, nil
}

// TODO: make transaction safe
func (repo *repository) SetPersonRole(ctx context.Context, id string, roles []string) error {
	person, err := repo.GetPerson(ctx, id)
	if err != nil {
		return err
	}
	person.Role = roles
	_, err = repo.UpdatePerson(ctx, person)
	return err
}

// TODO: make transaction safe
func (repo *repository) SetPersonSettings(ctx context.Context, id string, settings map[string]string) error {
	person, err := repo.GetPerson(ctx, id)
	if err != nil {
		return err
	}

	person.Settings = settings

	_, err = repo.UpdatePerson(ctx, person)
	return err
}

func (repo *repository) AutoExpirePeople(ctx context.Context) (int64, error) {
	query := "UPDATE people SET active = false WHERE expiration_date <= $1 AND active = true"

	res, err := repo.client.ExecContext(ctx, query, time.Now().Local().Format("20060101"))
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (repo *repository) GetPeople(ctx context.Context) ([]*models.Person, string, error) {
	people, newCursor, err := repo.getPeople(ctx, setCursor{})
	if err != nil {
		return nil, "", err
	}

	var encodedCursor string
	if newCursor.LastID > 0 {
		encodedCursor, err = repo.encodeCursor(newCursor)
		if err != nil {
			return nil, "", err
		}
	}
	return people, encodedCursor, nil
}

func (repo *repository) GetMorePeople(ctx context.Context, tokenValue string) ([]*models.Person, string, error) {
	cursor := setCursor{}
	if err := repo.decodeCursor(tokenValue, &cursor); err != nil {
		return nil, "", err
	}

	people, newCursor, err := repo.getPeople(ctx, cursor)
	if err != nil {
		return nil, "", err
	}

	var encodedCursor string
	if newCursor.LastID > 0 {
		encodedCursor, err = repo.encodeCursor(newCursor)
		if err != nil {
			return nil, "", err
		}
	}
	return people, encodedCursor, nil
}

func (repo *repository) getPeople(ctx context.Context, cursor setCursor) ([]*models.Person, setCursor, error) {
	newCursor := setCursor{}

	query := `
SELECT
	"id",
	"date_created",
	"date_updated",
	"external_id",
	"active",
	"birth_date",
	"email",
	"given_name",
	"name",
	"family_name",
	"job_category",
	"preferred_given_name",
	"preferred_family_name",
	"honorific_prefix", 
	"role",
	"settings", 
	"object_class",
	"expiration_date",
	"token"
FROM "people" WHERE "id" > $1 ORDER BY "id" ASC LIMIT $2
	`

	rows, err := repo.client.QueryContext(ctx, query, cursor.LastID, personPageLimit)
	if err != nil {
		return nil, newCursor, err
	}
	defer rows.Close()

	rowIDs := []int{}
	people := []*models.Person{}

	for rows.Next() {
		var rowID int
		var encTokens []string
		p := &models.Person{}
		err = rows.Scan(
			&rowID,
			&p.DateCreated,
			&p.DateUpdated,
			&p.ID,
			&p.Active,
			&p.BirthDate,
			&p.Email,
			&p.GivenName,
			&p.Name,
			&p.FamilyName,
			&p.JobCategory,
			&p.PreferredGivenName,
			&p.PreferredFamilyName,
			&p.HonorificPrefix,
			&p.Role,
			&p.Settings,
			&p.ObjectClass,
			&p.ExpirationDate,
			encTokens,
		)
		if err != nil {
			return nil, newCursor, err
		}
		p.Token, err = repo.decryptTokens(encTokens)
		if err != nil {
			return nil, newCursor, err
		}
		rowIDs = append(rowIDs, rowID)
		people = append(people, p)
	}

	if len(people) == 0 {
		return people, newCursor, nil
	}

	allPersonOrganizationMembers, err := repo.getOrganizationMembers(ctx, rowIDs...)
	if err != nil {
		return nil, newCursor, err
	}

	for i := 0; i < len(rowIDs); i++ {
		rowID := rowIDs[i]
		person := people[i]
		for _, orgMember := range allPersonOrganizationMembers {
			if orgMember.personID == rowID {
				person.AddOrganizationMember(models.OrganizationMember{
					ID:          orgMember.organizationExternalID,
					DateCreated: orgMember.dateCreated,
					DateUpdated: orgMember.dateUpdated,
				})
			}
		}
	}

	allPersonIdentifiers, err := repo.getPersonIdentifiers(ctx, rowIDs...)
	if err != nil {
		return nil, newCursor, err
	}
	for i := 0; i < len(rowIDs); i++ {
		rowID := rowIDs[i]
		person := people[i]
		for _, pid := range allPersonIdentifiers {
			if pid.personID == rowID {
				urn, _ := models.ParseURN(pid.value)
				person.AddIdentifier(*urn)
			}
		}
	}

	var total int
	err = repo.client.QueryRowContext(ctx, `SELECT COUNT(*) AS "total" FROM "people"`).Scan(&total)
	if err != nil {
		return nil, newCursor, err
	}

	if total > len(people) {
		newCursor = setCursor{
			LastID: rowIDs[len(rowIDs)-1],
		}
	}

	return people, newCursor, nil
}

func (repo *repository) encodeCursor(c any) (string, error) {
	plaintext, _ := json.Marshal(c)
	ciphertext, err := crypt.Encrypt(repo.secret, plaintext)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (repo *repository) decodeCursor(encryptedCursor string, c any) error {
	ciphertext, err := base64.URLEncoding.DecodeString(encryptedCursor)
	if err != nil {
		return err
	}
	plaintext, err := crypt.Decrypt(repo.secret, ciphertext)
	if err != nil {
		return err
	}
	return json.Unmarshal(plaintext, c)
}

func (repo *repository) descryptToken(encToken string) (*models.URN, error) {
	eURN, _ := models.ParseURN(encToken)
	rawTokenVal, err := decryptMessage(repo.secret, eURN.Value)
	if err != nil {
		return nil, fmt.Errorf("unable to decrypt %s token: %w", eURN.Namespace, err)
	}
	return &models.URN{
		Namespace: eURN.Namespace,
		Value:     rawTokenVal,
	}, nil
}

func (repo *repository) decryptTokens(encTokens []string) ([]*models.URN, error) {
	urns := make([]*models.URN, 0, len(encTokens))
	for _, encToken := range encTokens {
		urn, err := repo.descryptToken(encToken)
		if err != nil {
			return nil, err
		}
		urns = append(urns, urn)
	}
	return urns, nil
}

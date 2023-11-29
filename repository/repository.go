package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgtype"
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
	from                         *time.Time
	until                        *time.Time
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

type organization struct {
	id          int
	externalID  pgtype.Text
	dateCreated pgtype.Timestamptz
	dateUpdated pgtype.Timestamptz
	Type        pgtype.Text
	nameDut     pgtype.Text
	nameEng     pgtype.Text
	acronym     pgtype.Text
}

type person struct {
	id                  int
	token               pgtype.JSONB
	externalID          pgtype.Text
	active              pgtype.Bool
	dateCreated         pgtype.Timestamptz
	dateUpdated         pgtype.Timestamptz
	name                pgtype.Text
	givenName           pgtype.Text
	familyName          pgtype.Text
	email               pgtype.Text
	preferredGivenName  pgtype.Text
	preferredFamilyName pgtype.Text
	birthDate           pgtype.Text
	honorificPrefix     pgtype.Text
	jobCategory         pgtype.JSONB
	role                pgtype.JSONB
	settings            pgtype.JSONB
	objectClass         pgtype.JSONB
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

func (repo *repository) getOrganizationIdentifiers(ctx context.Context, organizationIDs ...int) ([]*organizationIdentifier, error) {
	query := `
SELECT
	"id", "organization_id", "value" FROM "organization_identifiers"
WHERE "organization_id" = any($1) ORDER BY array_position($1, organization_id), "value" ASC`
	rows, err := repo.client.QueryContext(
		ctx,
		query,
		pgIntArray(organizationIDs),
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
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return organizationIdentifiers, nil
}

func (repo *repository) getPersonIdentifiers(ctx context.Context, personIDs ...int) ([]*personIdentifier, error) {
	query := `
SELECT
	"id", "person_id", "value" FROM "person_identifiers"
WHERE "person_id" = any($1) ORDER BY array_position($1, person_id), "value" ASC`
	rows, err := repo.client.QueryContext(
		ctx,
		query,
		pgIntArray(personIDs),
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
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pids, nil
}

func (repo *repository) getOrganizationMembers(ctx context.Context, personIDs ...int) ([]*organizationMember, error) {
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
ORDER BY array_position($1, person_id), "organization_id" ASC
	`
	rows, err := repo.client.QueryContext(
		ctx,
		query,
		pgIntArray(personIDs),
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	organizationMembers := []*organizationMember{}

	for rows.Next() {
		om := &organizationMember{}
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
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return organizationMembers, nil
}

func (repo *repository) getOrganizationParents(ctx context.Context, organizationIDs ...int) ([]*organizationParent, error) {
	query := `
SELECT
	"id",
	"organization_id",
    "parent_organization_id",
	"date_created",
	"date_updated",
	"from",
	"until",
	(SELECT "external_id" FROM "organizations" WHERE "id" = op.parent_organization_id) AS "parent_organization_external_id"
FROM "organization_parents" op
WHERE "organization_id" = any($1)
ORDER by array_position($1, organization_id), "parent_organization_id" ASC
	`
	rows, err := repo.client.QueryContext(
		ctx,
		query,
		pgIntArray(organizationIDs),
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	organizationParents := []*organizationParent{}

	for rows.Next() {
		op := &organizationParent{}
		err := rows.Scan(
			&op.id,
			&op.organizationID,
			&op.parentOrganizationID,
			&op.dateCreated,
			&op.dateUpdated,
			&op.from,
			&op.until,
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
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return organizationParents, nil
}

func (repo *repository) GetOrganization(ctx context.Context, externalId string) (*models.Organization, error) {
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

	orgRec := &organization{}
	err := repo.client.QueryRowContext(ctx, query, externalId).Scan(
		&orgRec.id,
		&orgRec.externalID,
		&orgRec.dateCreated,
		&orgRec.dateUpdated,
		&orgRec.nameDut,
		&orgRec.nameEng,
		&orgRec.acronym,
		&orgRec.Type,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNotFound
	}

	orgs, err := repo.unpackOrganizations(ctx, orgRec)
	if err != nil {
		return nil, err
	}

	return orgs[0], nil
}

func (repo *repository) GetOrganizationsByIdentifier(ctx context.Context, urns ...*models.URN) ([]*models.Organization, error) {
	urnValues := make([]string, 0, len(urns))
	for _, urn := range urns {
		urnValues = append(urnValues, urn.String())
	}

	orgRecs := []*organization{}

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
		orgRec := &organization{}
		err = rows.Scan(
			&orgRec.id,
			&orgRec.externalID,
			&orgRec.dateCreated,
			&orgRec.dateUpdated,
			&orgRec.Type,
			&orgRec.nameDut,
			&orgRec.nameEng,
			&orgRec.acronym,
		)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		orgRecs = append(orgRecs, orgRec)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(orgRecs) == 0 {
		return nil, nil
	}

	orgs, err := repo.unpackOrganizations(ctx, orgRecs...)
	if err != nil {
		return nil, err
	}

	return orgs, nil
}

func (repo *repository) unpackOrganizations(ctx context.Context, orgRecs ...*organization) ([]*models.Organization, error) {
	orgs := make([]*models.Organization, 0, len(orgRecs))

	if len(orgRecs) == 0 {
		return orgs, nil
	}

	rowIDs := make([]int, 0, len(orgRecs))
	for _, orgRec := range orgRecs {
		rowIDs = append(rowIDs, orgRec.id)
		orgs = append(orgs, &models.Organization{
			ID:          orgRec.externalID.String,
			DateCreated: &orgRec.dateCreated.Time,
			DateUpdated: &orgRec.dateUpdated.Time,
			Type:        orgRec.Type.String,
			NameDut:     orgRec.nameDut.String,
			NameEng:     orgRec.nameEng.String,
			Acronym:     orgRec.acronym.String,
		})
	}

	allOrganizationParents, err := repo.getOrganizationParents(ctx, rowIDs...)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(orgRecs); i++ {
		orgRec := orgRecs[i]
		org := orgs[i]
		for _, op := range allOrganizationParents {
			if op.organizationID == orgRec.id {
				org.AddParent(&models.OrganizationParent{
					ID:          op.parentOrganizationExternalID,
					DateCreated: op.dateCreated,
					DateUpdated: op.dateUpdated,
					From:        op.from,
					Until:       op.until,
				})
			}
		}
	}

	allOrganizationIdentifiers, err := repo.getOrganizationIdentifiers(ctx, rowIDs...)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(orgRecs); i++ {
		orgRec := orgRecs[i]
		org := orgs[i]
		for _, oid := range allOrganizationIdentifiers {
			if oid.organizationID == orgRec.id {
				urn, _ := models.ParseURN(oid.value)
				org.AddIdentifier(urn)
			}
		}
	}

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

	// start transaction
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
	tsVals := []string{org.NameDut, org.NameEng, org.Acronym}
	tsVals = append(tsVals, org.GetIdentifierValues()...)
	tsVals = vacuum(tsVals)
	var rowID int
	err = tx.QueryRowContext(
		ctx, query,
		org.ID,
		org.DateCreated,
		org.DateUpdated,
		nullString(org.NameDut),
		nullString(org.NameEng),
		org.Type,
		nullString(org.Acronym),
		nullJSON(tsVals),
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
	pRows, err := tx.QueryContext(
		ctx,
		query,
		pgTextArray(parentOrganizationExternalIds),
	)
	if err != nil {
		return nil, err
	}
	defer pRows.Close()
	for pRows.Next() {
		var rowID int
		pRows.Scan(&rowID)
		parentOrganizationIDs = append(parentOrganizationIDs, rowID)
	}
	parentOrganizationIDs = lo.Uniq(parentOrganizationIDs)
	if len(parentOrganizationExternalIds) != len(parentOrganizationIDs) {
		return nil, models.ErrInvalidReference
	}
	query = `
INSERT INTO "organization_parents"
	("organization_id", "parent_organization_id", "date_created", "date_updated", "from", "until")
VALUES($1, $2, $3, $4, $5, $6);
`
	for i := 0; i < len(parentOrganizationIDs); i++ {
		parentOrganizationID := parentOrganizationIDs[i]
		orgParent := org.Parent[i]
		_, err := tx.ExecContext(
			ctx,
			query,
			rowID,
			parentOrganizationID,
			orgParent.DateCreated,
			orgParent.DateUpdated,
			orgParent.From,
			orgParent.Until,
		)
		if err != nil {
			return nil, err
		}
	}

	// add identifiers
	query = `
INSERT INTO "organization_identifiers"
	("organization_id", "date_created", "date_updated", "value")
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
	tsVals := []string{org.NameDut, org.NameEng, org.Acronym}
	tsVals = append(tsVals, org.GetIdentifierValues()...)
	tsVals = vacuum(tsVals)

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
		nullString(org.NameDut),
		nullString(org.NameEng),
		org.Type,
		nullString(org.Acronym),
		nullJSON(tsVals),
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
		pRows, err := tx.QueryContext(ctx, "SELECT id, external_id FROM organizations WHERE external_id = any($1) ORDER BY array_position($1, external_id)", pgExternalIds)
		if err != nil {
			return nil, err
		}
		defer pRows.Close()

		for pRows.Next() {
			var parentID int
			var parentExternalID string
			err = pRows.Scan(&parentID, &parentExternalID)
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
		if err := pRows.Err(); err != nil {
			return nil, err
		}

		if len(parentOrganizationExternalIDs) != len(organizationParents) {
			return nil, models.ErrInvalidReference
		}

		for _, parent := range org.Parent {
			newOrganizationParent := organizationParent{
				organizationID: rowID,
				dateCreated:    parent.DateCreated,
				dateUpdated:    parent.DateUpdated,
				from:           parent.From,
				until:          parent.Until,
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
INSERT INTO "organization_parents"
	("organization_id", "parent_organization_id", "date_created", "date_updated", "from", "until")
VALUES($1, $2, $3, $4, $5, $6)
ON CONFLICT("organization_id", "parent_organization_id", "from")
DO UPDATE SET date_updated = EXCLUDED.date_updated, until = EXCLUDED.until
RETURNING "id"
	`

	if len(newOrganizationParents) > 0 {
		for _, newOrganizationParent := range newOrganizationParents {
			var relId int
			err = tx.QueryRowContext(
				ctx,
				query,
				rowID,
				newOrganizationParent.parentOrganizationID,
				newOrganizationParent.dateCreated,
				newOrganizationParent.dateUpdated,
				newOrganizationParent.from,
				newOrganizationParent.until,
			).Scan(&relId)
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
ON CONFLICT("organization_id", "value")
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
		`SELECT
		"id", 
		"external_id", 
		"date_created", 
		"date_updated", 
		"type", 
		"name_dut", 
		"name_eng", 
		"acronym",
		ts_rank(ts, %s) AS rank 
		FROM "organizations" WHERE ts @@ %s LIMIT %d`,
		tsQuery,
		tsQuery,
		organizationSuggestLimit)

	rows, err := repo.client.QueryContext(ctx, sqlQuery, tsQueryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orgRecs := []*organization{}

	for rows.Next() {
		orgRec := &organization{}
		var rank float64
		err = rows.Scan(
			&orgRec.id,
			&orgRec.externalID,
			&orgRec.dateCreated,
			&orgRec.dateUpdated,
			&orgRec.Type,
			&orgRec.nameDut,
			&orgRec.nameEng,
			&orgRec.acronym,
			&rank,
		)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		orgRecs = append(orgRecs, orgRec)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(orgRecs) == 0 {
		return nil, nil
	}

	orgs, err := repo.unpackOrganizations(ctx, orgRecs...)
	if err != nil {
		return nil, err
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
WHERE "id" > $1 ORDER BY "id" ASC LIMIT $2`

	rows, err := repo.client.QueryContext(ctx, query, cursor.LastID, organizationPageLimit)
	if err != nil {
		return nil, newCursor, err
	}
	defer rows.Close()

	orgRecs := []*organization{}
	for rows.Next() {
		orgRec := &organization{}
		err = rows.Scan(
			&orgRec.id,
			&orgRec.externalID,
			&orgRec.dateCreated,
			&orgRec.dateUpdated,
			&orgRec.Type,
			&orgRec.nameDut,
			&orgRec.nameEng,
			&orgRec.acronym,
		)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, newCursor, nil
		}
		if err != nil {
			return nil, newCursor, err
		}
		orgRecs = append(orgRecs, orgRec)
	}

	if err := rows.Err(); err != nil {
		return nil, newCursor, err
	}

	if len(orgRecs) == 0 {
		return nil, newCursor, nil
	}

	// get uncapped total
	var total int
	err = repo.client.QueryRowContext(ctx, `SELECT COUNT(*) "total" FROM "organizations"`).Scan(&total)
	if err != nil {
		return nil, newCursor, err
	}

	orgs, err := repo.unpackOrganizations(ctx, orgRecs...)
	if err != nil {
		return nil, newCursor, err
	}

	// set next cursor
	if total > len(orgRecs) {
		newCursor = setCursor{
			LastID: orgRecs[len(orgRecs)-1].id,
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
		$18
	)
	RETURNING "id"
	`
	var rowID int
	queryArgs := []any{
		p.ID,
		p.DateCreated,
		p.DateUpdated,
		p.Active,
		nullString(p.BirthDate),
		nullJSON(p.JobCategory),
		nullString(p.Email),
		nullString(p.GivenName),
		nullString(p.PreferredGivenName),
		nullString(p.Name),
		nullString(p.FamilyName),
		nullString(p.PreferredFamilyName),
		nullString(p.HonorificPrefix),
		nullJSON(p.Role),
		nullJSON(p.Settings),
		nullJSON(p.ObjectClass),
	}
	tokens := make([]string, 0, len(p.Token))
	for _, token := range p.Token {
		eToken, err := encryptMessage(repo.secret, token.Value)
		if err != nil {
			return nil, fmt.Errorf("unable to encrypt %s: %w", token.Namespace, err)
		}
		eURN := models.NewURN(token.Namespace, eToken)
		tokens = append(tokens, eURN.String())
	}
	queryArgs = append(queryArgs,
		nullJSON(tokens),
		nullJSON(vacuum([]string{p.Name})),
	)

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
		if err := rows.Err(); err != nil {
			return nil, err
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

	person.Identifier = lo.Filter(person.Identifier, func(identifier *models.URN, idx int) bool {
		return identifier.Namespace != "orcid"
	})
	if orcid != "" {
		person.AddIdentifier(models.NewURN("orcid", orcid))
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

	person.Token = lo.Filter(person.Token, func(token *models.URN, idx int) bool {
		return token.Namespace != "orcid"
	})
	if orcidToken != "" {
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
	"token" = $15,
	"ts_vals" = $16
WHERE "external_id" = $17
RETURNING "id"
	`
	var rowID int
	queryArgs := []any{
		p.DateUpdated,
		p.Active,
		nullString(p.BirthDate),
		nullJSON(p.JobCategory),
		nullString(p.Email),
		nullString(p.GivenName),
		nullString(p.PreferredGivenName),
		nullString(p.Name),
		nullString(p.FamilyName),
		nullString(p.PreferredFamilyName),
		nullString(p.HonorificPrefix),
		nullJSON(p.Role),
		nullJSON(p.Settings),
		nullJSON(p.ObjectClass),
	}
	tokens := make([]string, 0, len(p.Token))
	for _, token := range p.Token {
		eToken, err := encryptMessage(repo.secret, token.Value)
		if err != nil {
			return nil, fmt.Errorf("unable to encrypt %s: %w", token.Namespace, err)
		}
		eURN := models.NewURN(token.Namespace, eToken)
		tokens = append(tokens, eURN.String())
	}
	queryArgs = append(queryArgs,
		nullJSON(tokens),
		nullJSON(vacuum([]string{p.Name})),
		p.ID,
	)

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

	// "person_identifiers"
	updatedPersonIdentifierIds := []int{}
	if len(p.Identifier) > 0 {
		for _, urnValue := range p.GetIdentifierQualifiedValues() {
			insertQuery := `
			INSERT INTO "person_identifiers"
				("date_created", "date_updated", "person_id", "value")
			VALUES($1, $2, $3, $4)
			ON CONFLICT("person_id", "value")
			DO UPDATE SET date_updated = EXCLUDED.date_updated
			RETURNING "id"
			`
			var relID int
			err = tx.QueryRowContext(ctx, insertQuery, now, now, rowID, urnValue).Scan(&relID)
			if err != nil {
				return nil, err
			}
			updatedPersonIdentifierIds = append(updatedPersonIdentifierIds, relID)
		}
	}
	query = `DELETE FROM "person_identifiers" WHERE "person_id" = $1`
	queryArgs = []any{rowID}
	if len(updatedPersonIdentifierIds) > 0 {
		query += ` AND NOT "id" = any($2)`
		queryArgs = append(queryArgs, pgIntArray(updatedPersonIdentifierIds))
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

func (repo *repository) GetPerson(ctx context.Context, externalID string) (*models.Person, error) {
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
	"token"
FROM "people" WHERE "external_id" = $1
LIMIT 1
	`

	p := &person{}
	err := repo.client.QueryRowContext(ctx, query, externalID).Scan(
		&p.id,
		&p.dateCreated,
		&p.dateUpdated,
		&p.externalID,
		&p.active,
		&p.birthDate,
		&p.email,
		&p.givenName,
		&p.name,
		&p.familyName,
		&p.jobCategory,
		&p.preferredGivenName,
		&p.preferredFamilyName,
		&p.honorificPrefix,
		&p.role,
		&p.settings,
		&p.objectClass,
		&p.token,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	people, err := repo.unpackPeople(ctx, p)
	if err != nil {
		return nil, err
	}

	return people[0], nil
}

func (repo *repository) GetPeopleByIdentifier(ctx context.Context, urns ...*models.URN) ([]*models.Person, error) {
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

	personRecs := []*person{}

	for rows.Next() {
		p := &person{}
		err = rows.Scan(
			&p.id,
			&p.dateCreated,
			&p.dateUpdated,
			&p.externalID,
			&p.active,
			&p.birthDate,
			&p.email,
			&p.givenName,
			&p.name,
			&p.familyName,
			&p.jobCategory,
			&p.preferredGivenName,
			&p.preferredFamilyName,
			&p.honorificPrefix,
			&p.role,
			&p.settings,
			&p.objectClass,
			&p.token,
		)
		if err != nil {
			return nil, err
		}
		personRecs = append(personRecs, p)
	}

	if len(personRecs) == 0 {
		return nil, nil
	}

	people, err := repo.unpackPeople(ctx, personRecs...)
	if err != nil {
		return nil, err
	}

	return people, nil
}

func (repo *repository) unpackPeople(ctx context.Context, personRecs ...*person) ([]*models.Person, error) {
	people := make([]*models.Person, 0, len(personRecs))

	if len(personRecs) == 0 {
		return people, nil
	}

	rowIDs := make([]int, 0, len(personRecs))
	for _, personRec := range personRecs {
		rowIDs = append(rowIDs, personRec.id)

		person := &models.Person{
			ID:                  personRec.externalID.String,
			Active:              personRec.active.Bool,
			DateCreated:         &personRec.dateCreated.Time,
			DateUpdated:         &personRec.dateUpdated.Time,
			Name:                personRec.name.String,
			GivenName:           personRec.givenName.String,
			FamilyName:          personRec.familyName.String,
			Email:               personRec.email.String,
			PreferredGivenName:  personRec.preferredGivenName.String,
			PreferredFamilyName: personRec.preferredFamilyName.String,
			BirthDate:           personRec.birthDate.String,
			HonorificPrefix:     personRec.honorificPrefix.String,
		}
		if vals, err := fromNullStringArray(personRec.jobCategory.Bytes); err != nil {
			return nil, err
		} else {
			person.JobCategory = vals
		}
		if vals, err := fromNullStringArray(personRec.role.Bytes); err != nil {
			return nil, err
		} else {
			person.Role = vals
		}
		if vals, err := fromNullStringArray(personRec.objectClass.Bytes); err != nil {
			return nil, err
		} else {
			person.ObjectClass = vals
		}
		if vals, err := fromNullStringArray(personRec.token.Bytes); err != nil {
			return nil, err
		} else {
			urns := make([]*models.URN, 0, len(vals))
			for _, encToken := range vals {
				urn, err := repo.descryptToken(encToken)
				if err != nil {
					return nil, err
				}
				urns = append(urns, urn)
			}
			person.Token = urns
		}
		if m, err := fromNullMap(personRec.settings.Bytes); err != nil {
			return nil, err
		} else {
			person.Settings = m
		}
		people = append(people, person)
	}

	allPersonOrganizationMembers, err := repo.getOrganizationMembers(ctx, rowIDs...)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(personRecs); i++ {
		personRec := personRecs[i]
		person := people[i]
		for _, orgMember := range allPersonOrganizationMembers {
			if orgMember.personID == personRec.id {
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
	for i := 0; i < len(personRecs); i++ {
		personRec := personRecs[i]
		person := people[i]
		for _, pid := range allPersonIdentifiers {
			if pid.personID == personRec.id {
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

	personRecs := []*person{}

	for rows.Next() {
		personRec := &person{}
		var rank float64
		err = rows.Scan(
			&personRec.id,
			&personRec.dateCreated,
			&personRec.dateUpdated,
			&personRec.externalID,
			&personRec.active,
			&personRec.birthDate,
			&personRec.email,
			&personRec.givenName,
			&personRec.name,
			&personRec.familyName,
			&personRec.jobCategory,
			&personRec.preferredGivenName,
			&personRec.preferredFamilyName,
			&personRec.honorificPrefix,
			&personRec.role,
			&personRec.settings,
			&personRec.objectClass,
			&personRec.token,
			&rank,
		)
		if err != nil {
			return nil, err
		}
		personRecs = append(personRecs, personRec)
	}

	if len(personRecs) == 0 {
		return nil, nil
	}

	people, err := repo.unpackPeople(ctx, personRecs...)
	if err != nil {
		return nil, err
	}

	return people, nil
}

func (repo *repository) SetPersonRole(ctx context.Context, externalID string, roles []string) error {
	res, err := repo.client.ExecContext(
		ctx,
		`UPDATE "people" SET date_updated = now(), role = $1 WHERE external_id = $2`,
		nullJSON(roles),
		externalID,
	)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return models.ErrNotFound
	}

	return nil
}

func (repo *repository) SetPersonSettings(ctx context.Context, externalID string, settings map[string]string) error {
	res, err := repo.client.ExecContext(
		ctx,
		`UPDATE "people" SET date_updated = now(), settings = $1 WHERE external_id = $2`,
		nullJSON(settings),
		externalID,
	)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return models.ErrNotFound
	}

	return nil
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
	"token"
FROM "people" WHERE "id" > $1 ORDER BY "id" ASC LIMIT $2
	`

	rows, err := repo.client.QueryContext(ctx, query, cursor.LastID, personPageLimit)
	if err != nil {
		return nil, newCursor, err
	}
	defer rows.Close()

	personRecs := []*person{}

	for rows.Next() {
		personRec := &person{}
		err = rows.Scan(
			&personRec.id,
			&personRec.dateCreated,
			&personRec.dateUpdated,
			&personRec.externalID,
			&personRec.active,
			&personRec.birthDate,
			&personRec.email,
			&personRec.givenName,
			&personRec.name,
			&personRec.familyName,
			&personRec.jobCategory,
			&personRec.preferredGivenName,
			&personRec.preferredFamilyName,
			&personRec.honorificPrefix,
			&personRec.role,
			&personRec.settings,
			&personRec.objectClass,
			&personRec.token,
		)
		if err != nil {
			return nil, newCursor, err
		}
		personRecs = append(personRecs, personRec)
	}

	if len(personRecs) == 0 {
		return nil, newCursor, nil
	}

	people, err := repo.unpackPeople(ctx, personRecs...)
	if err != nil {
		return nil, newCursor, err
	}

	var total int
	err = repo.client.QueryRowContext(ctx, `SELECT COUNT(*) AS "total" FROM "people"`).Scan(&total)
	if err != nil {
		return nil, newCursor, err
	}
	if total > len(personRecs) {
		newCursor = setCursor{
			LastID: personRecs[len(personRecs)-1].id,
		}
	}

	return people, newCursor, nil
}

func (repo *repository) GetPersonIDActive(ctx context.Context, active bool) ([]string, error) {
	rows, err := repo.client.QueryContext(ctx, `SELECT "external_id" FROM "people" WHERE active = $1`, active)
	if err != nil {
		return nil, err
	}

	externalIDs := []string{}
	for rows.Next() {
		var externalID string
		err := rows.Scan(&externalID)
		if err != nil {
			return nil, err
		}
		externalIDs = append(externalIDs, externalID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return externalIDs, nil
}

func (repo *repository) SetPeopleActive(ctx context.Context, active bool, externalIDs ...string) error {
	_, err := repo.client.ExecContext(
		ctx,
		`UPDATE "people" SET date_updated = now(), active = $1 WHERE "external_id" = any($2)`,
		active,
		pgTextArray(externalIDs),
	)
	if err != nil {
		return err
	}
	return nil
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

package repository

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"github.com/samber/lo"
	"github.com/ugent-library/crypt"
	"github.com/ugent-library/people-service/old/models"
)

const (
	personPageLimit       = 200
	organizationPageLimit = 200
)

type repository struct {
	client *pgxpool.Pool
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

type organization struct {
	id          int
	externalID  pgtype.Text
	dateCreated pgtype.Timestamptz
	dateUpdated pgtype.Timestamptz
	Type        pgtype.Text
	nameDut     pgtype.Text
	nameEng     pgtype.Text
	acronym     pgtype.Text
	identifier  []byte
}

type person struct {
	id                  int
	token               []byte
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
	jobCategory         []byte
	role                []byte
	settings            []byte
	objectClass         []byte
	identifier          []byte
}

func NewRepository(config *Config) (*repository, error) {
	// cf. https://github.com/jackc/pgx/blob/master/pgxpool/pool.go#L286
	pool, err := pgxpool.New(context.TODO(), config.DbUrl)
	if err != nil {
		return nil, err
	}
	return &repository{
		client: pool,
		secret: []byte(config.AesKey),
	}, nil
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
	rows, err := repo.client.Query(
		ctx,
		query,
		personIDs,
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
	rows, err := repo.client.Query(
		ctx,
		query,
		organizationIDs,
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
	"type",
	"identifier"
FROM organizations WHERE external_id = $1 LIMIT 1`

	orgRec := &organization{}
	err := repo.client.QueryRow(ctx, query, externalId).Scan(
		&orgRec.id,
		&orgRec.externalID,
		&orgRec.dateCreated,
		&orgRec.dateUpdated,
		&orgRec.nameDut,
		&orgRec.nameEng,
		&orgRec.acronym,
		&orgRec.Type,
		&orgRec.identifier,
	)
	if errors.Is(err, pgx.ErrNoRows) {
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
	"acronym",
	"identifier"
FROM "organizations" WHERE "identifier" ?| $1`

	rows, err := repo.client.Query(ctx, query, urnValues)
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
			&orgRec.identifier,
		)
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
		dateCreated := orgRec.dateCreated.Time
		dateUpdated := orgRec.dateUpdated.Time
		org := &models.Organization{
			ID:          orgRec.externalID.String,
			DateCreated: &dateCreated,
			DateUpdated: &dateUpdated,
			Type:        orgRec.Type.String,
			NameDut:     orgRec.nameDut.String,
			NameEng:     orgRec.nameEng.String,
			Acronym:     orgRec.acronym.String,
		}
		orgs = append(orgs, org)
		urnValues := []string{}
		if err := json.Unmarshal(orgRec.identifier, &urnValues); err != nil {
			return nil, err
		}
		for _, urnVal := range urnValues {
			urn, err := models.ParseURN(urnVal)
			if err != nil {
				return nil, err
			}
			org.AddIdentifier(urn)
		}
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

	return orgs, nil
}

func (repo *repository) SaveOrganization(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	if org.IsStored() {
		return repo.UpdateOrganization(ctx, org)
	}
	return repo.CreateOrganization(ctx, org)
}

func (repo *repository) getTsValsForOrganization(org *models.Organization) []string {
	tsVals := []string{org.NameDut, org.NameEng, org.Acronym}
	tsVals = vacuum(append(tsVals, org.GetIdentifierValues()...))
	return tsVals
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
	tx, err := repo.client.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

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
		"identifier",
		"ts_vals"
	)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING "id"
	`
	var rowID int
	err = tx.QueryRow(
		ctx, query,
		org.ID,
		org.DateCreated,
		org.DateUpdated,
		pgtext(org.NameDut),
		pgtext(org.NameEng),
		org.Type,
		pgtext(org.Acronym),
		pgjson(org.GetIdentifierQualifiedValues()),
		pgjson(repo.getTsValsForOrganization(org)),
	).Scan(&rowID)
	if err != nil {
		return nil, err
	}

	// add parents
	parentOrganizationExternalIds := []string{}
	for _, parent := range org.Parent {
		parentOrganizationExternalIds = append(parentOrganizationExternalIds, parent.ID)
	}
	parentOrganizationExternalIds = lo.Uniq(parentOrganizationExternalIds)

	pOrgRows := []*organization{}
	query = `SELECT "id", "external_id" FROM "organizations" WHERE "external_id" = any($1)`
	pRows, err := tx.Query(
		ctx,
		query,
		parentOrganizationExternalIds,
	)
	if err != nil {
		return nil, err
	}
	defer pRows.Close()
	for pRows.Next() {
		o := &organization{}
		pRows.Scan(&o.id, &o.externalID)
		pOrgRows = append(pOrgRows, o)
	}

	if len(parentOrganizationExternalIds) != len(pOrgRows) {
		return nil, models.ErrInvalidReference
	}

	query = `
INSERT INTO "organization_parents"
	("organization_id", "parent_organization_id", "date_created", "date_updated", "from", "until")
VALUES($1, $2, $3, $4, $5, $6);
`
	for _, orgParent := range org.Parent {
		var pOrgRow *organization
		for _, po := range pOrgRows {
			if po.externalID.String == orgParent.ID {
				pOrgRow = po
				break
			}
		}
		_, err := tx.Exec(
			ctx,
			query,
			rowID,
			pOrgRow.id,
			orgParent.DateCreated,
			orgParent.DateUpdated,
			orgParent.From,
			orgParent.Until,
		)
		if err != nil {
			return nil, err
		}
	}

	// commit
	if err := tx.Commit(ctx); err != nil {
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
	tx, err := repo.client.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// update organization
	query := `
UPDATE "organizations"
SET
	"date_updated" = $2,
	"name_dut" = $3,
	"name_eng" = $4,
	"type" = $5,
	"acronym" = $6,
	"identifier" = $7,
	"ts_vals" = $8
WHERE "external_id" = $1
RETURNING "id"
	`
	var rowID int
	err = tx.QueryRow(
		ctx,
		query,
		org.ID,
		now,
		pgtext(org.NameDut),
		pgtext(org.NameEng),
		org.Type,
		pgtext(org.Acronym),
		pgjson(org.GetIdentifierQualifiedValues()),
		pgjson(repo.getTsValsForOrganization(org)),
	).Scan(&rowID)
	if err != nil {
		return nil, err
	}

	// update organization parents
	var newOrganizationParents []*organizationParent
	if len(org.Parent) > 0 {
		parentOrganizationExternalIDs := []string{}
		for _, parent := range org.Parent {
			parentOrganizationExternalIDs = append(parentOrganizationExternalIDs, parent.ID)
		}
		parentOrganizationExternalIDs = lo.Uniq(parentOrganizationExternalIDs)

		pRows, err := tx.Query(ctx, "SELECT id, external_id FROM organizations WHERE external_id = any($1)", parentOrganizationExternalIDs)
		if err != nil {
			return nil, err
		}
		defer pRows.Close()

		pOrgRows := []*organization{}

		for pRows.Next() {
			o := &organization{}
			err = pRows.Scan(&o.id, &o.externalID)
			if err != nil {
				return nil, err
			}
			pOrgRows = append(pOrgRows, o)
		}
		if err := pRows.Err(); err != nil {
			return nil, err
		}

		if len(parentOrganizationExternalIDs) != len(pOrgRows) {
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
			for _, pOrgRow := range pOrgRows {
				if pOrgRow.externalID.String == parent.ID {
					newOrganizationParent.parentOrganizationID = pOrgRow.id
					break
				}
			}
			newOrganizationParents = append(newOrganizationParents, &newOrganizationParent)
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
			err = tx.QueryRow(
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
	_, err = tx.Exec(ctx, query, rowID, updatedRelIds)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("unable to commit transaction: %w", err)
	}

	return org, nil
}

func (repo *repository) DeleteOrganization(ctx context.Context, id string) error {
	_, err := repo.client.Exec(ctx, "DELETE FROM organizations WHERE external_id = $1", id)
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

func (repo *repository) SuggestOrganizations(ctx context.Context, params models.OrganizationSuggestParams) ([]*models.Organization, error) {
	params = params.MergeDefault()
	tsQuery, tsQueryArgs := toTSQuery(params.Query)

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
		"identifier",
		ts_rank(ts, %s) AS rank 
		FROM "organizations" WHERE ts @@ %s LIMIT %d`,
		tsQuery,
		tsQuery,
		params.Limit)

	rows, err := repo.client.Query(ctx, sqlQuery, tsQueryArgs...)
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
			&orgRec.identifier,
			&rank,
		)
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

func (repo *repository) GetOrganizationsById(ctx context.Context, ids ...string) ([]*models.Organization, error) {
	query := `SELECT 
	"id",
	"external_id",
	"date_created",
	"date_updated",
	"type",
	"name_dut",
	"name_eng", 
	"acronym",
	"identifier"
FROM "organizations" WHERE "external_id" = any($1)`

	rows, err := repo.client.Query(ctx, query, ids)
	if err != nil {
		return nil, err
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
			&orgRec.identifier,
		)
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
	"acronym",
	"identifier"
FROM "organizations"
WHERE "id" > $1 ORDER BY "id" ASC LIMIT $2`

	rows, err := repo.client.Query(ctx, query, cursor.LastID, organizationPageLimit)
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
			&orgRec.identifier,
		)
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
	err = repo.client.QueryRow(ctx, `SELECT COUNT(*) "total" FROM "organizations"`).Scan(&total)
	if err != nil {
		return nil, newCursor, err
	}

	orgs, err := repo.unpackOrganizations(ctx, orgRecs...)
	if err != nil {
		return nil, newCursor, err
	}

	// set next cursor
	if len(orgRecs) >= organizationPageLimit {
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

func (repo *repository) getTsValsForPerson(p *models.Person) []string {
	tsVals := []string{p.Name}
	tsVals = vacuum(append(tsVals, p.GetIdentifierValues()...))
	return tsVals
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

	// ensure biblio_id
	p.EnsureBiblioID()

	tx, err := repo.client.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

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
		"identifier",
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
	queryArgs := []any{
		p.ID,
		p.DateCreated,
		p.DateUpdated,
		p.Active,
		pgtext(p.BirthDate),
		pgjson(p.JobCategory),
		pgtext(p.Email),
		pgtext(p.GivenName),
		pgtext(p.PreferredGivenName),
		pgtext(p.Name),
		pgtext(p.FamilyName),
		pgtext(p.PreferredFamilyName),
		pgtext(p.HonorificPrefix),
		pgjson(p.Role),
		pgjson(p.Settings),
		pgjson(p.ObjectClass),
	}
	eTokenMap := map[string]string{}
	for typ, val := range p.Token {
		eVal, err := encryptMessage(repo.secret, val)
		if err != nil {
			return nil, fmt.Errorf("unable to encrypt %s: %w", typ, err)
		}
		eTokenMap[typ] = eVal
	}
	queryArgs = append(queryArgs,
		pgjson(eTokenMap),
		pgjson(p.GetIdentifierQualifiedValues()),
		pgjson(repo.getTsValsForPerson(p)),
	)

	err = tx.QueryRow(ctx, query, queryArgs...).Scan(&rowID)
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
		rows, err := tx.Query(
			ctx,
			`SELECT "id" FROM "organizations" WHERE "external_id" = any($1)`,
			organizationExternalIDs)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var rowID int
			err = rows.Scan(&rowID)
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
			_, err = tx.Exec(ctx, insertQuery, now, now, orgRowIDS[i], rowID)
			if err != nil {
				return nil, err
			}
		}

	}

	if err := tx.Commit(ctx); err != nil {
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
func (repo *repository) SetPersonToken(ctx context.Context, id string, typ string, val string) error {
	person, err := repo.GetPerson(ctx, id)
	if err != nil {
		return err
	}

	if val == "" {
		delete(person.Token, typ)
	} else {
		person.Token[typ] = val
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
	// ensure biblio_id
	p.EnsureBiblioID()

	tx, err := repo.client.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

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
	"identifier" = $16,
	"ts_vals" = $17
WHERE "external_id" = $18
RETURNING "id"
	`
	var rowID int
	queryArgs := []any{
		p.DateUpdated,
		p.Active,
		pgtext(p.BirthDate),
		pgjson(p.JobCategory),
		pgtext(p.Email),
		pgtext(p.GivenName),
		pgtext(p.PreferredGivenName),
		pgtext(p.Name),
		pgtext(p.FamilyName),
		pgtext(p.PreferredFamilyName),
		pgtext(p.HonorificPrefix),
		pgjson(p.Role),
		pgjson(p.Settings),
		pgjson(p.ObjectClass),
	}
	eTokenMap := map[string]string{}
	for typ, val := range p.Token {
		eVal, err := encryptMessage(repo.secret, val)
		if err != nil {
			return nil, fmt.Errorf("unable to encrypt %s: %w", typ, err)
		}
		eTokenMap[typ] = eVal
	}
	queryArgs = append(queryArgs,
		pgjson(eTokenMap),
		pgjson(p.GetIdentifierQualifiedValues()),
		pgjson(repo.getTsValsForPerson(p)),
		p.ID,
	)

	err = tx.QueryRow(ctx, query, queryArgs...).Scan(&rowID)
	if err != nil {
		return nil, err
	}

	// update "organization_members"
	updatedOrganizationMemberIds := []int{}
	if len(p.Organization) > 0 {
		orgExternalIDs := make([]string, 0, len(p.Organization))
		for _, orgMem := range p.Organization {
			orgExternalIDs = append(orgExternalIDs, orgMem.ID)
		}
		orgExternalIDs = lo.Uniq(orgExternalIDs)

		rows, err := tx.Query(
			ctx,
			`SELECT "id", "external_id" FROM "organizations" WHERE "external_id" = any($1)`,
			orgExternalIDs,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		orgRows := make([]*organization, 0, len(p.Organization))
		for rows.Next() {
			o := &organization{}
			if err = rows.Scan(&o.id, &o.externalID); err != nil {
				return nil, err
			}
			orgRows = append(orgRows, o)
		}

		if len(orgExternalIDs) != len(orgRows) {
			return nil, models.ErrInvalidReference
		}

		for _, orgMember := range p.Organization {
			var orgId int
			for _, orgRow := range orgRows {
				if orgRow.externalID.String == orgMember.ID {
					orgId = orgRow.id
					break
				}
			}

			insertQuery := `
			INSERT INTO "organization_members"
				("date_created", "date_updated", "person_id", "organization_id")
			VALUES($1, $2, $3, $4)
			ON CONFLICT("person_id", "organization_id")
			DO UPDATE SET date_updated = EXCLUDED.date_updated
			RETURNING "id"
			`
			var relID int
			err = tx.QueryRow(ctx, insertQuery, orgMember.DateCreated, orgMember.DateUpdated, rowID, orgId).Scan(&relID)
			if err != nil {
				return nil, err
			}
			updatedOrganizationMemberIds = append(updatedOrganizationMemberIds, relID)
		}
	}

	query = `DELETE FROM "organization_members" WHERE "person_id" = $1 AND NOT "id" = any($2)`
	queryArgs = []any{rowID, updatedOrganizationMemberIds}
	_, err = tx.Exec(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
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
	"token",
	"identifier"
FROM "people" WHERE "external_id" = $1
LIMIT 1
	`

	p := &person{}
	err := repo.client.QueryRow(ctx, query, externalID).Scan(
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
		&p.identifier,
	)

	if errors.Is(err, pgx.ErrNoRows) {
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
		"token",
		"identifier"
	FROM "people" WHERE "identifier" ?| $1
	`

	ids := make([]string, 0, len(urns))
	for _, urn := range urns {
		ids = append(ids, urn.String())
	}
	rows, err := repo.client.Query(ctx, query, ids)
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
			&p.identifier,
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

func (repo *repository) GetPeopleById(ctx context.Context, ids ...string) ([]*models.Person, error) {
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
		"token",
		"identifier"
	FROM "people" WHERE "external_id" = any($1)
	`

	rows, err := repo.client.Query(ctx, query, ids)
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
			&p.identifier,
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
		dateCreated := personRec.dateCreated.Time
		dateUpdated := personRec.dateUpdated.Time
		person := &models.Person{
			ID:                  personRec.externalID.String,
			Active:              personRec.active.Bool,
			DateCreated:         &dateCreated,
			DateUpdated:         &dateUpdated,
			Name:                personRec.name.String,
			GivenName:           personRec.givenName.String,
			FamilyName:          personRec.familyName.String,
			Email:               personRec.email.String,
			PreferredGivenName:  personRec.preferredGivenName.String,
			PreferredFamilyName: personRec.preferredFamilyName.String,
			BirthDate:           personRec.birthDate.String,
			HonorificPrefix:     personRec.honorificPrefix.String,
		}
		if vals, err := fromPgTextArray(personRec.jobCategory); err != nil {
			return nil, err
		} else {
			person.JobCategory = vals
		}
		if vals, err := fromPgTextArray(personRec.role); err != nil {
			return nil, err
		} else {
			person.Role = vals
		}
		if vals, err := fromPgTextArray(personRec.objectClass); err != nil {
			return nil, err
		} else {
			person.ObjectClass = vals
		}
		if eTokenMap, err := fromPgMap(personRec.token); err != nil {
			return nil, err
		} else {
			tokenMap := map[string]string{}
			for typ, eVal := range eTokenMap {
				val, err := repo.decryptToken(eVal)
				if err != nil {
					return nil, err
				}
				tokenMap[typ] = val
			}
			person.Token = tokenMap
		}
		if vals, err := fromPgTextArray(personRec.identifier); err != nil {
			return nil, err
		} else {
			for _, val := range vals {
				urn, err := models.ParseURN(val)
				if err != nil {
					return nil, err
				}
				person.AddIdentifier(urn)
			}
		}

		if m, err := fromPgMap(personRec.settings); err != nil {
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

	return people, nil
}

func (repo *repository) DeletePerson(ctx context.Context, id string) error {
	_, err := repo.client.Exec(ctx, `DELETE FROM "people" WHERE "external_id" = $1`, id)
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

func (repo *repository) SuggestPeople(ctx context.Context, params models.PersonSuggestParams) ([]*models.Person, error) {
	params = params.MergeDefault()
	tsQuery, tsQueryArgs := toTSQuery(params.Query)
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
	"identifier",
	ts_rank(ts, %s) AS rank
FROM "people" WHERE ts @@ %s AND "active" = any(%s)  ORDER BY "rank" DESC LIMIT %d
`
	tsQueryArgs = append(tsQueryArgs, params.Active)
	sqlQuery = fmt.Sprintf(
		sqlQuery,
		tsQuery,
		tsQuery,
		fmt.Sprintf("$%d", len(tsQueryArgs)),
		params.Limit,
	)
	rows, err := repo.client.Query(ctx, sqlQuery, tsQueryArgs...)
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
			&personRec.identifier,
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
	res, err := repo.client.Exec(
		ctx,
		`UPDATE "people" SET date_updated = now(), role = $1 WHERE external_id = $2`,
		pgjson(roles),
		externalID,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}

func (repo *repository) SetPersonSettings(ctx context.Context, externalID string, settings map[string]string) error {
	res, err := repo.client.Exec(
		ctx,
		`UPDATE "people" SET date_updated = now(), settings = $1 WHERE external_id = $2`,
		pgjson(settings),
		externalID,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
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
	"token",
	"identifier"
FROM "people" WHERE "id" > $1 ORDER BY "id" ASC LIMIT $2
	`

	rows, err := repo.client.Query(ctx, query, cursor.LastID, personPageLimit)
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
			&personRec.identifier,
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
	err = repo.client.QueryRow(ctx, `SELECT COUNT(*) AS "total" FROM "people"`).Scan(&total)
	if err != nil {
		return nil, newCursor, err
	}

	if len(personRecs) >= personPageLimit {
		newCursor = setCursor{
			LastID: personRecs[len(personRecs)-1].id,
		}
	}

	return people, newCursor, nil
}

func (repo *repository) GetPersonIDActive(ctx context.Context, active bool) ([]string, error) {
	rows, err := repo.client.Query(ctx, `SELECT "external_id" FROM "people" WHERE active = $1`, active)
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

func (repo *repository) SetPersonActive(ctx context.Context, externalID string, active bool) error {
	_, err := repo.client.Exec(
		ctx,
		`UPDATE "people" SET date_updated = now(), active = $1 WHERE "external_id" = $2`,
		active,
		externalID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (repo *repository) SetPeopleActive(ctx context.Context, active bool, externalIDs ...string) error {
	_, err := repo.client.Exec(
		ctx,
		`UPDATE "people" SET date_updated = now(), active = $1 WHERE "external_id" = any($2)`,
		active,
		externalIDs,
	)
	if err != nil {
		return err
	}
	return nil
}

func (repo *repository) RebuildAutocompletePeople(ctx context.Context) error {
	tx, err := repo.client.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	rows, err := repo.client.Query(ctx, `SELECT "id", "name", "identifier" FROM "people"`)
	if err != nil {
		return err
	}
	for rows.Next() {
		pRec := &person{}
		if err = rows.Scan(&pRec.id, &pRec.name, &pRec.identifier); err != nil {
			return err
		}
		p := &models.Person{Name: pRec.name.String}
		if vals, err := fromPgTextArray(pRec.identifier); err != nil {
			return err
		} else {
			for _, val := range vals {
				urn, err := models.ParseURN(val)
				if err != nil {
					return err
				}
				p.AddIdentifier(urn)
			}
		}
		_, err = tx.Exec(ctx, `UPDATE "people" SET "ts_vals" = $1 WHERE "id" = $2`, pgjson(repo.getTsValsForPerson(p)), pRec.id)
		if err != nil {
			return err
		}
	}

	if err = rows.Err(); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (repo *repository) RebuildAutocompleteOrganizations(ctx context.Context) error {
	tx, err := repo.client.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	rows, err := repo.client.Query(ctx, `SELECT "id", "name_dut", "name_eng", "acronym", "identifier" FROM "organizations"`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		oRec := &organization{}
		if err = rows.Scan(&oRec.id, &oRec.nameDut, &oRec.nameEng, &oRec.acronym, &oRec.identifier); err != nil {
			return err
		}
		org := &models.Organization{
			NameDut: oRec.nameDut.String,
			NameEng: oRec.nameEng.String,
			Acronym: oRec.acronym.String,
		}
		if vals, err := fromPgTextArray(oRec.identifier); err != nil {
			return err
		} else {
			for _, val := range vals {
				urn, err := models.ParseURN(val)
				if err != nil {
					return err
				}
				org.AddIdentifier(urn)
			}
		}
		_, err = tx.Exec(ctx, `UPDATE "organizations" SET "ts_vals" = $1 WHERE "id" = $2`, pgjson(repo.getTsValsForOrganization(org)), oRec.id)
		if err != nil {
			return err
		}
	}

	if err = rows.Err(); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
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

func (repo *repository) decryptToken(eVal string) (string, error) {
	rawTokenVal, err := decryptMessage(repo.secret, eVal)
	if err != nil {
		return "", fmt.Errorf("unable to decrypt token: %w", err)
	}
	return rawTokenVal, nil
}

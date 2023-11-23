package repository

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgtype"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ugent-library/crypt"
)

var regexMultipleSpaces = regexp.MustCompile(`\s+`)
var regexNoBrackets = regexp.MustCompile(`[\[\]()\{\}]`)

func toTSQuery(query string) (string, []any) {
	// remove duplicate spaces
	query = regexMultipleSpaces.ReplaceAllString(query, " ")
	// trim
	query = strings.TrimSpace(query)

	queryParts := make([]string, 0)
	queryArgs := make([]any, 0)
	argCounter := 0

	for _, qp := range strings.Split(query, " ") {
		// remove terms that contain brackets
		if regexNoBrackets.MatchString(qp) {
			continue
		}
		argCounter++

		// $1 || ':*'
		queryParts = append(queryParts, fmt.Sprintf("$%d || ':*'", argCounter))
		queryArgs = append(queryArgs, qp)
	}

	// $1:* & $2:*
	tsQuery := fmt.Sprintf(
		"to_tsquery('usimple', %s)",
		strings.Join(queryParts, " || ' & ' || "),
	)

	return tsQuery, queryArgs
}

func openClient(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// TODO: make this configurable
	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Hour)
	db.SetConnMaxIdleTime(time.Minute)

	return db, nil
}

func encryptMessage(key []byte, message string) (string, error) {
	cryptedMsgInBytes, err := crypt.Encrypt(key, []byte(message))
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(cryptedMsgInBytes), nil
}

func decryptMessage(key []byte, cryptedMsg string) (string, error) {
	cryptedMsgInBytes, err := base64.URLEncoding.DecodeString(cryptedMsg)
	if err != nil {
		return "", err
	}

	msgInBytes, err := crypt.Decrypt(key, cryptedMsgInBytes)
	if err != nil {
		return "", err
	}

	return string(msgInBytes), nil
}

func pgTextArray(values []string) pgtype.TextArray {
	pgValues := pgtype.TextArray{}
	pgValues.Set(values)
	return pgValues
}

func pgIntArray(values []int) pgtype.Int4Array {
	pgValues := pgtype.Int4Array{}
	pgValues.Set(values)
	return pgValues
}

func toJSON(val any) []byte {
	bytes, _ := json.Marshal(val)
	return bytes
}

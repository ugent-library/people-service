package ugentldap

import "github.com/go-ldap/ldap/v3"

type Searcher interface {
	SearchPeople(string, func(*ldap.Entry) error) error
}

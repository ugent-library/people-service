package cli

import (
	"github.com/ugent-library/people-service/models"
	"github.com/ugent-library/people-service/repository"
	"github.com/ugent-library/people-service/ugentldap"
)

func newRepository() (models.Repository, error) {
	return repository.NewRepository(&repository.Config{
		DbUrl:  config.Db.Url,
		AesKey: config.Db.AesKey,
	})
}

func newUgentLdapClient() *ugentldap.Client {
	return ugentldap.NewClient(ugentldap.Config{
		Url:      config.Ldap.Url,
		Username: config.Ldap.Username,
		Password: config.Ldap.Password,
	})
}

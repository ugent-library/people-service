package cli

import (
	"github.com/ugent-library/people-service/old/models"
	"github.com/ugent-library/people-service/old/repository"
	"github.com/ugent-library/people-service/old/ugentldap"
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

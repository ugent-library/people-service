package ugentldap

import (
	"context"

	"github.com/go-ldap/ldap/v3"
)

type Client struct {
	url      string
	username string
	password string
}

type clientConn struct {
	conn *ldap.Conn
}

type Config struct {
	Url      string
	Username string
	Password string
}

const bufferSize = 2000

var ldapAttributes = []string{
	"objectClass",
	"uid",
	"ugentPreferredSn",
	"ugentPreferredGivenName",
	"ugentID",
	"ugentHistoricIDs",
	"ugentBirthDate",
	"mail",
	"ugentBarcode",
	"ugentJobCategory",
	"ugentAddressingTitle",
	"displayName",
	"departmentNumber",
	"ugentFaculty",
}

func NewClient(config Config) *Client {
	return &Client{
		url:      config.Url,
		username: config.Username,
		password: config.Password,
	}
}

func (cli *Client) newConn() (*clientConn, error) {
	conn, err := ldap.DialURL(cli.url)
	if err != nil {
		return nil, err
	}

	if err = conn.Bind(cli.username, cli.password); err != nil {
		defer conn.Close()
		return nil, err
	}

	return &clientConn{conn}, nil
}

func (conn *clientConn) close() error {
	return conn.conn.Close()
}

func (conn *clientConn) Search(ctx context.Context, req *ldap.SearchRequest, cb func(*ldap.Entry) error) error {
	res := conn.conn.SearchAsync(ctx, req, bufferSize)
	for res.Next() {
		if err := cb(res.Entry()); err != nil {
			break
		}
	}
	if err := res.Err(); err != nil {
		return err
	}
	return nil
}

func (cli *Client) SearchPeople(ctx context.Context, filter string, cb func(*ldap.Entry) error) error {
	uc, err := cli.newConn()
	if err != nil {
		return err
	}
	defer uc.close()

	searchReq := ldap.NewSearchRequest(
		"ou=people,dc=ugent,dc=be",
		ldap.ScopeSingleLevel,
		ldap.NeverDerefAliases,
		0, 0, false,
		filter,
		ldapAttributes,
		[]ldap.Control{},
	)

	return uc.Search(ctx, searchReq, cb)
}

# people-service
Person service

# Build process

```
go build
```

# Prepare environment variables

Declare the following environment variables,

or store them in file `.env` in the root of your folder (important: exclude `export` statement)
(see .env.example):

* `PEOPLE_PRODUCTION`:

  type: `bool`

  default: `false`

* `PEOPLE_API_KEY`

  type: `string`

  default: `people`

  description: api key. Used in authentication header `X-Api-Key` for server (openapi)

* `PEOPLE_API_PORT`

  type: `int`

  default: `3999`

  description: port address for server

* `PEOPLE_API_HOST`

  type : `string`

  default: `127.0.0.1`

  description: host address for server

* `PEOPLE_DB_URL`

  type: `string`

  default: `postgres://people:people@localhost:5432/authority?sslmode=disable`

  description: postgres database connection url

* `PEOPLE_DB_AES_KEY`

  type: `string`

  required: `true`

  description: [AES](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard) key. This is now used to encrypt attribute `orcid_token`. Note that an AES key must be 128 bits long (or 16 characters). Generate one with command `openssl rand -hex 16`

* `PEOPLE_LDAP_URL`

  type: `string`

  description: ldap connection url. e.g. `ldaps://ldaps.ugent.be:636`

  required: `true`

* `PEOPLE_LDAP_USERNAME`

  type: `string`

  description: ldap username

  required: `true`

  Note: internally we bind to scope `ou=people,dc=ugent,dc=be`, so make sure these
  credentials are valid for that scope.

* `PEOPLE_LDAP_PASSWORD`

  type: `string`

  description: ldap password

  required: `true`

# Run database migrations

Before starting the application you should run any pending database migrations.

In production we use [tern](https://github.com/jackc/tern). Make sure
that directories `ent/migrate/migrations` (for atlas) and `etc/migrations` (for tern) are kept in sync. And note that tern uses a different naming for sql
files (prefix is a padded number instead of a timestamp)

# Start the api server (openapi)

```
$ ./people-service server
```

# run in docker

Build base docker image `people-service`:

```
$ docker build -t ugentlib/people-service .
$ docker push ugentlib/people-service
```

If image `people-service` is already docker github,

then you may skip that step.

Start set of services using `docker compose`:

```
$ docker compose up
```

Docker compose uses that image `people-service`

# Update flow organizations

## nightly cron job that upserts person records

A nightly cron job reads person records from the ugent ldap

and inserts/updates person records. Matching on existing records

is done by matching on identifier `historic_ugent_id`.

The following attributes are overwritten:

* `identifier->'ugent_username'`
* `identifier->'ugent_id'`
* `identifier->'historic_ugent_id'`
* `identifier->'ugent_barcode'`
* `given_name`
* `family_name`
* `name`
* `birth_date`
* `email`
* `job_category`
* `honorific_prefix`
* `object_class`
* `organization`.

Note that if no organization be found based on `identifier->'ugent'` then no (dummy) organization record is made for it. In that case the attribute is ignored.

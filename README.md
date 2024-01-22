# people-service
Person service

## Build process

```
go build
```

## Prepare environment variables

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

## Run database migrations

We use [tern](https://github.com/jackc/tern) for database migrations.

Before starting the application you should run any pending database migrations:

```
. .env
tern status --conn-string $PEOPLE_DB_URL -m etc/migrations
tern migrate --conn-string $PEOPLE_DB_URL -m etc/migrations
```

## Start the api server (openapi)

```
$ ./people-service server
```

## Update flow organizations

### nightly cron job that upserts person records

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

## Dev Containers

This project supports [Dev Containers](https://code.visualstudio.com/docs/devcontainers/containers). Following these steps
will auto setup a containerized development environment for this project. In VS Code, you will be able to start a terminal
that logs into a Docker container. This will allow you to write and interact with the code inside a self-contained sandbox.

**Installing the Dev Containers extension**

1. Open VS Code.
2. Go to the [Dev Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) extension page.
3. Click the `install` button to install the extension in VS Code.

**Open in Dev Containers**

1. Open the project directory in VS Code.
2. Click on the "Open a remote window" button in the lower left window corner.
3. Choose "reopen in container" from the popup menu.
4. The green button should now read "Dev Container: App name" when successfully opened.
5. Open a new terminal in VS Code from the `Terminal` menu link.

You are now logged into the dev container and ready to develop code, write code, push to git or execute commands.

**Run the project**

1. Open a new terminal in VS Code from the `Terminal` menu link.
2. Execute this command `reflex -d none -c reflex.docker.conf`.
3. Once the application has started, VS Code will show a popup with a link that opens the project in your browser.

**Networking**

The application and its dependencies run on these ports:

| Application    | Port |
| -------------- | ---- |
| People Service | 3201 |
| DB Application | 3251 |
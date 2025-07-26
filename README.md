# Houston we have (no) problem

<p align="center">
    <img src="./houston.jpg">
</p>

**Houston** is a lightweight and efficient web application designed for monitoring the availability of other web
applications. Built with simplicity in mind, it provides an easy-to-use alternative to complex monitoring solutions like
Grafana.

**Key Features**:

- Periodic health checks of specified URLs
- HTTP status verification to detect downtime
- Email notifications in case of failures
- Visual representation of service uptime history through graphs

Houston ensures that you stay informed about the health of your applications without unnecessary complexity or resource
overhead.

**Contents**:

1. [Technology](#Technology)
2. [Migrations](#Migrations)
3. [Getting started](#Getting-started)
4. [Settings](#Settings)
5. [Tests](#Tests)

## Technology

1. **GoLang** - vanilla Go as core of the application
2. **SQLite** - relational database as a file
3. **Templ** - templates markup & HTML rendering
4. **Cobra** - CLI interface
5. **TailwindCSS** - templates styling engine
6. **FlyonUI** - library of ready to use TailwindCSS components
7. **HTMX** - JS magic without JS development
8. **Air** - application Hot Reload
9. **Docker** - universal running environment

## Migrations

Houston application stores data in relational database. For easier development, **GOOSE** migration tool is used
to automate database maintenance. Migrations lives in `./migrations`. New migration can be created with:

```bash
goose -dir=./migrations create [migration_name] sql
```

If there is no **GOOSE** installed, run:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Ready to use migration can be applied on existing database wih:

```bash
make up
```

Last migration rollback can be achieved with:

```bash
make down
```

Migrations status can be checked with:

```bash
make status
```

**WARNING**: all above is described for development purposes - application can migrate to the newest database version
by itself (migrations are running from application level on startup).

## Getting started

Application build & run is fully automated. Application can be run in one of two modes:

1. **Standard** - just build & run application
2. **Hot Reload** - build & run application, watch for files changes & apply them live

To run application in Hot Reload mode, use:

```bash
make
```

To run application in standard mode, use:

```bash
make run
```

By default, application will be ready under [localhost](http://127.0.0.1:8080).

### CLI

When application is built as binary, you can run it from CLI. Houston provides more than one CLI command, to see all
available options run:

```bash
houston --help
```

### Docker

Application is ready to use with Docker. Image can be built with:

```bash
docker build -t houston:dev .
```

Nextly, it can be run with `docker compose`:

```bash
docker compose up -d
```

By default, in compose file, all needed data are linked to the container as volumes to share with development environment:

1. `./configs/app.yml` - application configuration
2. `./data/app.db` - SQLite database file
3. `./logs` - application logs

Ready to use image for download is always available at [Docker Hub](https://hub.docker.com). It is recommended to link
above data to the running container for data persistence.

### Mailing

The main aim of Houston is to notify user about changes in web applications health checks. By default mailing client
is not configured - notifications are visible only in console output. More info about configuration can be found
in [Settings](#Settings) chapter. There is an option to test real life scenerio with usage of SMTP server. When running
project with usage of **docker compose**, `maildev` can be used to test notifications by SMTP. First of all,
application needs to be configured:

```yaml
smtp:
  enabled: true
  host: "smtp"
  port: 1025
  user: "no-reply@houston.com"
  from: "Houston <no-reply@houston.com>"
```

After running application with **docker compose** and this configuration, e-mail client will be accessible on
[local machine](http://127.0.0.1:1080). Here, sent e-mail notifications can be confirmed.

## Settings

Default application configuration is described below in form of YAML file. Application is searching for configuration in
3 ways:

1. Read `./app.yml`.
2. Read `./config/app.yml`.
3. If no configuration found, use defaults.

Additionally, if configuration file do not contain specific variables, default values will be used.


```yaml
# Logs configuration
logs:
  # Should log to file?
  file-enabled: false
  # Max single logs file size in MB
  max-size: 10
  # Max age of single logs file in days (logs rotation)
  max-age: 30

# Server configuration
server:
  # Host
  host: "0.0.0.0"
  # Port
  port: "8080"

# Database configuration
database:
  # Path to SQLite database file
  file: "./data/app.db"

# Authorization configuration
authorization:
  # Secret to generate JWT token
  secret: ""
  # JWT token expire after X hours
  expires-after: 12
  # Use secure cookies (required to work with HTTPS)
  over-https: true
  # Sign up configuration
  sign-up:
    # Public sign up enabled
    enabled: true

# Health checks data retention configuration
retention:
  # Retention enabled?
  enabled: false
  # Delete health checks data older than X days
  older-than: 30

# E-mail notifications configuration
smtp:
  # E-mail notifications enabled (by default log to console)?
  enabled: false
  # SMTP server address
  host: ""
  # SMTP server port
  port: -1
  # SMTP server username
  user: ""
  # SMTP server password
  password: ""
  # E-mail sender (by default username will be used)
  from: ""
```

## Tests

All unit tests are stored in `tests` directory. For unit testing data stores are mocked with in memory (maps)
implementation. Application is designed with commands execution approach, so only commands should be unit tested
(this covers whole business logic). Unit tests can be launched with:

```bash
make test
```

Above command will out full unit tests report. To run tests without detailed output use:

```bash
make test_no_output
```

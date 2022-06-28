# geolocation

## what is geolocation?

`geolocation` is a `CLI App` built with [cobra](https://github.com/spf13/cobra) and [viper](https://github.com/spf13/viper) and [go](https://go.dev/).
It parses location data from a `csv` file, sanitise and runs it through some validations and uploads it to a `postgres` table, 
after the loading is done the data can be accessed over an `api` call by passing in an `IP address`.
That's in nut shell what the app dose.

The `App` exposes 2 main commands:-

1. `ingest` which parses a `csv` file and persists it in database
2. `serve` which runs an `http server` exposing and endpoint named `resolve` which converts a given `IP Address` to `location` 
    using the database populated in step `#1`

## helps?

`./geolocation --help`

`./geolocation ingest --help`

`./geolocation serve --help`

## hmm interesting, are there any prerequisites?
The `App`, at bare minimim needs a `.env` at the project's root with following details
```
DB_HOST=postgres
BD_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_DATABASE=postgres
DB_DRIVER=postgres
DB_EXPLORER_EMAIL=someone@postgres.com
DB_EXPLORER_PASSWORD=postgres
``` 
and a `.csv` file to ingest location data from.

**For `docker-compose` to work, we've to mount the location where `.csv` file is present**
`line # 23 in docker-compose.yaml`

## architecture?

```
cli commands -> service -> database
```

## how to execute?
`docker-compose up` will bring up all services.
After all services have successfully come up visit
[localhost](http://localhost:8080/resolve?ip=125.159.20.54) and substitute the `ip`, either it will get resolved or you'll get an error.

To explore the postgres data visit [postgres](http://localhost:5050) and login with `DB_EXPLORER_EMAIL` and `DB_EXPLORER_PASSWORD` and configure the database in there.

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

## help?

`docker run bhardwaz007/geolocation:latest --help`

`docker run bhardwaz007/geolocation:latest ingest --help`

`docker run bhardwaz007/geolocation:latest serve --help`

## prerequisites?
The `App`, at bare minimim needs a `.env` at the project's root with postgres database details, following is the format
```ts
DB_HOST=<postgres host>
BD_PORT=<postgres port>
DB_USER=<postgres user>
DB_PASSWORD=<postgres password>
DB_DATABASE=<postgres database>
DB_DRIVER=postgres
``` 
and a `.csv` file to ingest location data from.

Sample `data_dump.csv`

```csv
ip_address,country_code,country,city,latitude,longitude,mystery_value
200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346
160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115
70.95.73.73,TL,Saudi Arabia,Gradymouth,-49.16675918861615,-86.05920084416894,2559997162
,PY,Falkland Islands (Malvinas),,75.41685191518815,-144.6943217219469,0
125.159.20.54,LI,Guyana,Port Karson,-78.2274228596799,-163.26218895343357,1337885276
```

**We expect you to mount these files when you run the image**

## architecture

```
cli commands -> service -> database
```

## how to execute?
you need to `ingest` the `csv` data first only then you'll be able to query it.
`ingest` internally takes care of `migration.`

### how to ingets? 
```docker
docker run -it \
--name geolocation-ingestor \
-v location/to/data_dump.csv:/root/data_dump.csv \ # csv file location
-v location/to/.env:/root/.env \ # .env file location
bhardwaz007/geolocation:latest ingest
```
and the `serve`.

### how to serve? 
```docker
docker run -it \
-p 8080:8080 \
--name geolocation-server \
-v location/to/.env:/root/.env \ # .env file location
bhardwaz007/geolocation:latest serve
```

you can visit [here](http://localhost:8080/resolve?ip=200.106.141.15) check.
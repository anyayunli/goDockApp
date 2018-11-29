goDockApp
------------
A RESTful application boilerplate in Go (golang),  taking packages and tools of Go echo + Postgres + Materialize + Docker.

### Requirements
```
Go1.11.1 or newer.

PostgreSQL 9.6.3 or newer.
```

### Installation
```
git clone git@github.com:anyayunli/goDockApp.git

```

or if you have the .zip files,  extract files

```
cd goDockApp
```

###  Build & Run
```
docker build -t go-dock-app .

docker-compose up
```

### HTTP server address and port

defaults to http://localhost:3344/

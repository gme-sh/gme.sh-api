<img align="right" src="https://forthebadge.com/images/badges/made-with-go.svg"></img>

[![Go Report Card](https://goreportcard.com/badge/github.com/gme-sh/gme.sh-api)](https://goreportcard.com/report/github.com/gme-sh/gme.sh-api)
![Build Go](https://github.com/gme-sh/gme.sh-api/workflows/Build%20Go/badge.svg)
![Build Docker](https://github.com/gme-sh/gme.sh-api/workflows/Build%20Docker/badge.svg)

# GMEshortener
$GME go brrrrr

> Allan, please add details!

## Run
### Docker
```bash
$ docker build . -t gmesh:latest
$ docker run -it --rm --name gmesh-api -p 80:80 gmesh:latest
```

#### BBolt, non-persistent
```bash
$ docker run -it --rm --name gmesh-api -e "GME_PERSISTENT_BACKEND=bbolt" -p 80:80 gmesh:latest
```

#### BBolt, persistent
```bash
$ docker run -it --rm --name gmesh-api -e "GME_PERSISTENT_BACKEND=bbolt" -v $PWD/data:/data  -p 80:80 gmesh:latest
```

### Docker-Compose
Copy `docker-compose-{preferred-option}.yml` and `docker-compose.env` from `docker/`

**Options:**
* redis-mongo
* redis
* scratch

```bash
$ cp docker/docker-compose-redis-mongo.yml ./docker-compose.yml
$ cp docker/docker-compose.env ./
```

Build
```bash
$ docker-compose build
```

Start it
```bash
$ docker-compose up [-d]
```
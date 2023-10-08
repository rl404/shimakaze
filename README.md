# Shimakaze

[![Go Report Card](https://goreportcard.com/badge/github.com/rl404/shimakaze)](https://goreportcard.com/report/github.com/rl404/shimakaze)
![License: MIT](https://img.shields.io/github/license/rl404/shimakaze)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/rl404/shimakaze)
[![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/rl404/shimakaze)](https://hub.docker.com/r/rl404/shimakaze)
[![publish & deploy](https://github.com/rl404/shimakaze/actions/workflows/publish-deploy.yml/badge.svg)](https://github.com/rl404/shimakaze/actions/workflows/publish-deploy.yml)

Shimakaze is [Vtuber Wikia](https://virtualyoutuber.fandom.com/wiki/Virtual_YouTuber_Wiki) scraper and REST API.

## Features

- Save vtuber details
  - Vtuber data
  - Vtuber channels data and videos
    - Youtube
    - Twitch
    - Bilibili
    - Niconico
- Save agency's vtuber list
- Auto update vtuber & agency data (cron)
- Interchangeable cache
  - no cache
  - inmemory
  - [Redis](https://redis.io/)
  - [Memcache](https://memcached.org/)
- Interchangeable pubsub
  - [NSQ](https://nsq.io/)
  - [RabbitMQ](https://www.rabbitmq.com/)
  - [Redis](https://redis.io/)
  - [Google PubSub](https://cloud.google.com/pubsub)
- [Swagger](https://github.com/swaggo/swag)
- [Docker](https://www.docker.com/)
- [Newrelic](https://newrelic.com/) monitoring
  - HTTP
  - Cron
  - Database
  - Cache
  - Pubsub
  - External API

_More will be coming soon..._

## Requirement

- [Go](https://go.dev/)
- [MongoDB](https://www.mongodb.com/)
- [Youtube API key](https://cloud.google.com/docs/authentication/api-keys)
- [Twitch Client id & secret](https://dev.twitch.tv/docs/authentication/register-app/)
- PubSub ([NSQ](https://nsq.io/)/[RabbitMQ](https://www.rabbitmq.com/)/[Redis](https://redis.io/)/[Google PubSub](https://cloud.google.com/pubsub))
- (optional) Cache ([Redis](https://redis.io/)/[Memcache](https://memcached.org/))
- (optional) [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- (optional) [Newrelic](https://newrelic.com/) license key

## Installation

### Without [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)

1. Clone the repository.

```sh
git clone github.com/rl404/shimakaze
```

2. Rename `.env.sample` to `.env` and modify the values according to your setup.
3. Run. You need at least 2 consoles/terminals.

```sh
# Run the API.
make

# Run the consumer.
make consumer
```

6. [localhost:45001](http://localhost:45001) is ready (port may varies depend on your `.env`).

#### Other commands

```sh
# Update old vtuber data.
make cron-update

# Fill missing vtuber data.
make cron-fill
```

### With [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)

1. Clone the repository.

```sh
git clone github.com/rl404/shimakaze
```

2. Rename `.env.sample` to `.env` and modify the values according to your setup.
3. Get docker image.

```sh
# Pull existing image.
docker pull rl404/shimakaze

# Or build your own.
make docker-build
```

4. Run the container. You need at least 2 consoles/terminals.

```sh
# Run the API.
make docker-api

# Run the consumer.
make docker-consumer
```

5. [localhost:45001](http://localhost:45001) is ready (port may varies depend on your `.env`).

#### Other commands

```sh
# Update old vtuber data.
make docker-cron-update

# Fill missing vtuber data.
make docker-cron-fill

# Stop running containers.
make docker-stop
```

## Environment Variables

| Env                               |           Default           | Description                                                                                                |
| --------------------------------- | :-------------------------: | ---------------------------------------------------------------------------------------------------------- |
| `SHIMAKAZE_APP_ENV`               |            `dev`            | Environment type (`dev`/`prod`).                                                                           |
| `SHIMAKAZE_HTTP_PORT`             |           `45001`           | HTTP server port.                                                                                          |
| `SHIMAKAZE_HTTP_READ_TIMEOUT`     |            `5s`             | HTTP read timeout.                                                                                         |
| `SHIMAKAZE_HTTP_WRITE_TIMEOUT`    |            `5s`             | HTTP write timeout.                                                                                        |
| `SHIMAKAZE_HTTP_GRACEFUL_TIMEOUT` |            `10s`            | HTTP graceful timeout.                                                                                     |
| `SHIMAKAZE_CACHE_DIALECT`         |         `inmemory`          | Cache type (`nocache`/`redis`/`inmemory`)                                                       |
| `SHIMAKAZE_CACHE_ADDRESS`         |                             | Cache address.                                                                                             |
| `SHIMAKAZE_CACHE_PASSWORD`        |                             | Cache password.                                                                                            |
| `SHIMAKAZE_CACHE_TIME`            |            `24h`            | Cache time.                                                                                                |
| `SHIMAKAZE_DB_ADDRESS`            | `mongodb://localhost:27017` | Database address with port.                                                                                |
| `SHIMAKAZE_DB_NAME`               |         `shimakaze`         | Database name.                                                                                             |
| `SHIMAKAZE_DB_USER`               |                             | Database username.                                                                                         |
| `SHIMAKAZE_DB_PASSWORD`           |                             | Database password.                                                                                         |
| `SHIMAKAZE_PUBSUB_DIALECT`        |         `rabbitmq`          | Pubsub type (`rabbitmq`/`redis`/`google`)                                                            |
| `SHIMAKAZE_PUBSUB_ADDRESS`        |                             | Pubsub address (if you are using `google`, this will be your google project id).                           |
| `SHIMAKAZE_PUBSUB_PASSWORD`       |                             | Pubsub password (if you are using `google`, this will be the content of your google service account json). |
| `SHIMAKAZE_CRON_UPDATE_LIMIT`     |            `10`             | Vtuber & agency count limit when updating old data.                                                        |
| `SHIMAKAZE_CRON_FILL_LIMIT`       |            `10`             | Vtuber & agency count limit when filling missing data.                                                     |
| `SHIMAKAZE_CRON_AGENCY_AGE`       |             `7`             | Age of old agency data (in days).                                                                          |
| `SHIMAKAZE_CRON_ACTIVE_AGE`       |             `1`             | Age of old active vtuber data (in days).                                                                   |
| `SHIMAKAZE_CRON_RETIRED_AGE`      |             `7`             | Age of old retired vtuber data (in days).                                                                  |
| `SHIMAKAZE_NEWRELIC_NAME`         |         `shimakaze`         | Newrelic application name.                                                                                 |
| `SHIMAKAZE_NEWRELIC_LICENSE_KEY`  |                             | Newrelic license key.                                                                                      |
| `SHIMAKAZE_YOUTUBE_KEY`           |                             | Youtube API key.                                                                                           |
| `SHIMAKAZE_YOUTUBE_MAX_AGE`       |            `60`             | Age limit of youtube videos (in days).                                                                     |
| `SHIMAKAZE_TWITCH_CLIENT_ID`      |                             | Twitch client id.                                                                                          |
| `SHIMAKAZE_TWITCH_CLIENT_SECRET`  |                             | Twitch client secret.                                                                                      |
| `SHIMAKAZE_TWITCH_MAX_AGE`        |            `60`             | Age limit of twitch videos (in days).                                                                      |
| `SHIMAKAZE_BILIBILI_MAX_AGE`      |            `60`             | Age limit of bilibili videos (in days).                                                                    |
| `SHIMAKAZE_NICONICO_MAX_AGE`      |            `60`             | Age limit of niconico videos (in days).                                                                    |

## Trivia

[Shimakaze](<https://en.wikipedia.org/wiki/Japanese_destroyer_Shimakaze_(1942)>)'s name is taken from one of the fastest japanese destroyer. Also, [exists](https://en.kancollewiki.net/Shimakaze) in Kantai Collection games and manga.

## Disclaimer

Shimakaze is meant for educational purpose and personal usage only. Please use it responsibly according to Wikia [License](https://www.fandom.com/licensing).

All data belong to their respective copyrights owners, shimakaze does not have any affiliation with content providers.

## License

MIT License

Copyright (c) 2023 Axel

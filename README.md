# go-new-relic-gin-zerolog

[![Lines Of Code](https://tokei.rs/b1/github/kevinmichaelchen/go-new-relic-gin-zerolog?category=code)](https://github.com/kevinmichaelchen/go-new-relic-gin-zerolog)

This demo uses
* [New Relic](https://newrelic.com/), an o11y platform
* [gin](https://github.com/gin-gonic/gin), an HTTP web framework
* [zap](https://github.com/uber-go/zap), a structured logger

## Getting started
### Start server
Start the server using your New Relic browser license key.
```shell
env \
  ENV=local \
  SERVICE_NAME=foobar \
  NEW_RELIC_KEY=YOUR-KEY \
  go run main.go
```

### Hit endpoint
```shell
curl localhost:8081/health
```

### Check New Relic UI
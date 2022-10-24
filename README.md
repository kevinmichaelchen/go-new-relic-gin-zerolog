# go-new-relic-gin-zap

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
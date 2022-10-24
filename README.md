# go-new-relic-gin-zerolog

[![Lines Of Code](https://tokei.rs/b1/github/kevinmichaelchen/go-new-relic-gin-zerolog?category=code)](https://github.com/kevinmichaelchen/go-new-relic-gin-zerolog)

This demo uses
* [New Relic](https://newrelic.com/), an o11y platform
* [gin](https://github.com/gin-gonic/gin), an HTTP web framework
* [zerolog]([https://github.com/uber-go/zap](https://github.com/rs/zerolog)), a structured logger

For this demo, I signed up for a free plan of New Relic.

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
Visit [New Relic One](https://one.newrelic.com) and navigate to `APM & services`,
and search for your _Trace group_. When you drill down into a trace, you should
be able to see all the logs belonging to that trace. This is known as _logs in context_.

module github.com/mintel/fluxcloud-filebeat

go 1.12

replace github.com/docker/distribution => github.com/2opremio/distribution v0.0.0-20190419185413-6c9727e5e5de

replace github.com/mintel/pkg/config => ./pkg/config

replace github.com/mintel/pkg/handler => /pkg/handler

replace github.com/mintel/pkg/server => ./pkg/server

require (
	github.com/fluxcd/flux v1.15.0
	github.com/gin-gonic/gin v1.4.0
	github.com/justinbarrick/fluxcloud v0.3.8
	github.com/kelseyhightower/envconfig v1.4.0
)

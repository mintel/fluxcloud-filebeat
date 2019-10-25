# Overview

A WebHook for [fluxcloud](https://github.com/justinbarrick/fluxcloud) to send events to filebeat (or any tcp-socket).

Once events are passed to filebeat, they can be forwarded on to ElasticSearch, and from there, queried by tools such as Grafana.

This provides a nice way to annotate Flux events (releases, commits, etc) and overlay them with other data-sources, such as those from Prometheus.

## Pre-requesites

- [flux](https://github.com/fluxcd/flux)
- [fluxcloud](https://github.com/justinbarrick/fluxcloud)

## Usage

```
# In another terminal, run netcat to listen on port 9000 (for debugging)
netcat -l -p 9000
```

```
# Export required vars
export VCS_ROOT_URL=http://github.com/whatever-repo-flux-monitors
export FILEBEAT_ADDRESS=127.0.0.1:9000
```

```
# Run it
make
./bin/fluxcloud-filebeat
```

```
# In another terminal, curl some example flux event (which would come from fluxcloud)
curl -d@examples/release.json http://localhost:8080/v1/event -H"Content-type:application/json"
```
## Envars

Supported environments vars

| Variable 	       | Description                               | Default        | Required |
|----------------------|-------------------------------------------|----------------|----------|
| PORT                 | Port that the webserver listens on        | 8080           | No       |
| FILEBEAT_ADDRESS     | TCP Socket address to forward events onto | 127.0.0.1:9000 | No       |
| VCS_ROOT_URL         | Root URL of your VCS                      | -              | Yes      |
| KEEP_FLUX_EVENTS     | Forward on orignial Flux events           | 0              | No       |

## Kubernetes Manifests

Typically used as sidecar alongside flux, fluxcloud and filebeat.

```yaml
  - name: fluxcloud
    image: justinbarrick/fluxcloud:master-b0312e82
    env:
    - name: EXPORTER_TYPE
      value: webhook
    - name: WEBHOOK_URL
      value: http://127.0.0.1:8080/v1/event
    - name: LISTEN_ADDRESS
      value: :3032
    ports:
    - containerPort: 3032
      protocol: TCP
  - name: fluxcloud-webhook
    image: mintel/fluxcloud-filebeat:latest
    env:
    - name: FILEBEAT_ADDRESS
      value: "127.0.0.1:9000"
    - name: VCS_ROOT_URL
      value: http://github.com/<your-vsc-root>
    ports:
    - containerPort: 8080
      protocol: TCP
  - image: docker.elastic.co/beats/filebeat-oss:6.5.4
    name: filebeat
    ports:
    - containerPort: 9000
      protocol: TCP
```

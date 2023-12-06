<a href="https://github.com/piplos/stream-tools/releases"><img src="https://img.shields.io/github/v/release/piplos/stream-tools?sort=semver&label=Release&color=651FFF"></a>
<a href="https://goreportcard.com/report/github.com/piplos/stream-tools"><img src="https://goreportcard.com/badge/github.com/piplos/stream-tools"></a>
<a href="https://www.codefactor.io/repository/github/piplos/stream-tools"><img src="https://www.codefactor.io/repository/github/piplos/stream-tools/badge" alt="CodeFactor" /></a>
<a href="https://github.com/piplos/stream-tools/actions/workflows/release.yml"><img src="https://github.com/piplos/stream-tools/actions/workflows/release.yml/badge.svg"></a>
<a href="https://github.com/piplos/stream-tools/blob/master/LICENSE"><img src="https://img.shields.io/badge/License-MIT-yellow.svg"></a>
<a href="https://hub.docker.com/r/piplosmedia/stream-tools/"><img src="https://img.shields.io/docker/pulls/piplosmedia/stream-tools.svg"></a>
<a href="https://hub.docker.com/r/piplosmedia/stream-tools/"><img src="https://img.shields.io/docker/image-size/piplosmedia/stream-tools/latest"></a>

# stream-tools

### Requests:
`/ping` - Pings the server for availability, useful for Docker health checks. On a successful response, it will print "ping".  

`/stream/play` - displays what is currently playing. Request parameters:  
**url** - stream link (mandatory)  

`/stream/status` - отображает статус стрима. Request parameters:  
**url** - stream link (mandatory)  
**duration** - stream idle time (optional, default 10)  
**volume** - stream sound level (optional, default -70.0)  

### Response:
The response will return with one of the HTTP statuses - 200, 400, 500 and JSON in the form:
```json
{
  "code": int,
  "message": string
}
```
**code** - duplicate http response status code  
**message** - message body
## Golang Healthcheck ##

**Background**

Dockerfile's support `HEALTHCHECK CMD`. Often the implementation of the healthcheck is of the form:
```
HEALTHCHECK \
--interval=5s \
--timeout=5s \
CMD curl --fail http://localhost:80/healthz
```

But, this assumes the container includes:
* a shell (the absence of [...] entails this being the shell form of the command. So-named because it runs in a shell in the container)
* `curl` to invoke the HTTP GET

Neither of these should be assumed. It is considered a good practice to build minimal container images perhaps an init (e.g. [dumb-init](https://github.com/Yelp/dumb-init)) system and your intended process; including a shell and extraneous tools such as curl increases complexity and broadens the container's attack surface.

With shell forms, signals go to the shell rather than to your process and this may cause undesired behavior.

**Description**

Building on the work of the Solutio folks ([link](https://github.com/Soluto/golang-docker-healthcheck-example)), here's a simple Golang-based healthcheck that you may build statically (no external dependencies) and drop into your containers (not just Golang-oriented ones)

The multi-stage Dockerfile uses a Google [distroless](https://github.com/GoogleContainerTools/distroless) image to provide a minimal runtime.

To mirror the `curl` command that it ofen replaces, `healtcheck` takes a single command-line parameter for the HTTP endpoint that it should check (GET):
```
./healthcheck http://localhost:8080/healthz
```
Or in a Dockerfile (using the preffered `exec` form):
```
HEALTHCHECK --interval=5s --timeout=5s CMD ["/healthcheck","http://localhost:8080/healthz"]
```

**Use**

`build` and -- if desired - `push` the image:
```
IMAGE=hellohenry

docker build \
--tag=${IMAGE} \
--file=Dockerfile \
.

docker push \
--tag=${IMAGE}
```
Run it:
```
docker run \
--name=${IMAGE} \
--interactive \
--rm \
--tty \
--publish=8080:8080 \
hellohenry
```
curl its endpoints:
```
curl http://localhost:8080
curl http://localhost:8080/argz
curl http://localhost:8080/healthz
curl http://localhost:8080/varz
```
It should have a '(healthy)' healthcheck status:
```
docker container ls \
--format="{{.ID}}\t{{.Image}}\t{{.Status}}"
9f098599ef77	hellohenry	Up 3 minutes (healthy)
```
If you have (and you should) `jq` installed:
```
docker container inspect ${IMAGE} \
| jq ".[] | select(.Name==\"/${IMAGE}\") | .State"

```
Which should yield something similar to:
```
{
  "Status": "running",
  ...
  "StartedAt": "2018-06-04T15:34:43.913977796Z",
  "FinishedAt": "0001-01-01T00:00:00Z",
  "Health": {
    "Status": "healthy",
    "FailingStreak": 0,
    "Log": [
      {
        "Start": "2018-06-04T08:41:31.090731863-07:00",
        "End": "2018-06-04T08:41:31.148874822-07:00",
        "ExitCode": 0,
        "Output": ""
      },
      {
        "Start": "2018-06-04T08:41:36.157541673-07:00",
        "End": "2018-06-04T08:41:36.240299555-07:00",
        "ExitCode": 0,
        "Output": ""
      },
      ...
    ]
  }
}
```

**Miscellany**
Includes:
* Dockerfile
* Kubernetes: deployment.yaml, ingress.yaml
* Docker Compose: docker-compose.yaml

**Caveat**

One downside in not having a shell is that you can't use environment variables to define the healthcheck's URL.

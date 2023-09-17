![build](https://github.com/Tanemahuta/aws-lambda-server/actions/workflows/verify.yml/badge.svg?branch=main)
[![go report](https://goreportcard.com/badge/github.com/Tanemahuta/aws-lambda-server)](https://goreportcard.com/report/github.com/Tanemahuta/aws-lambda-server)
[![codecov](https://codecov.io/gh/Tanemahuta/aws-lambda-server/branch/main/graph/badge.svg?token=FHO3AAZ41O)](https://codecov.io/gh/Tanemahuta/aws-lambda-server)
[![Go Reference](https://pkg.go.dev/badge/github.com/Tanemahuta/aws-lambda-server.svg)](https://pkg.go.dev/github.com/Tanemahuta/aws-lambda-server)
![GHCR](https://ghcr-badge.egpl.dev/tanemahuta/aws-lambda-server/tags?trim=major,minor&label=latest&ignore=sha256*,v*)

# aws-lambda-server

## description

A server which invokes AWS lambda functions from http requests, mapping the request to the payload.

### docker image

A docker image can be found at `ghcr.io/tanemahuta/aws-lambda-server:<tag>`.

### routing

Routing is achieved using [gorilla/mux](https://github.com/gorilla/mux).

When using a path in the request route, you may use [path variables](https://github.com/gorilla/mux#readme)
(e.g. `/test/{id}`), which will be parsed and propagated to the lambda invocation.

### function invocation

The AWS lambda function is invoked using the [aws-sdk](https://aws.amazon.com/de/sdk-for-go/).

If you need to attach an IAM role in an EKS cluster, check out
[this article](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html).

#### request

The parsed request is being adapted using [aws.LambdaRequest](pkg/aws/lambda_request.go).

To handle the request, you can use the first parameter in your handler:

```javascript
async function handler(req, ctx) {
    console.log("hostname", req.host);
    console.log("headers", req.headers);
    console.log("method", req.method);
    console.log("full uri", req.uri);
    console.log("parsed path variables", req.vars);
    console.log("read body", req.body);
}
```

#### response

The returned response is being adapted using [aws.LambdaResponse](pkg/aws/lambda_response.go).

An example response may look like this:

```javascript
async function handler(req, ctx) {
    return {
        statusCode: 200,
        headers: {
            "Content-Type": "text/plain"
        },
        body: "Hello world"
    }
}
```

Alternatively the body may be a JSON, which will be serialized by the server:

```javascript
async function handler(req, ctx) {
    return {
        statusCode: 200,
        headers: {
            "Content-Type": "application/json"
        },
        body: {
            "Hello": "World"
        }
    }
}
```

which will result in a `{"Hello":"World"}` in the server's HTTP response body.

## command-line args

When running the [binary](main.go), the following command line parameters can be used:

- `--devel=(true|false)`: run in development mode (logging)
- `--config-file=<path>`: use the provided config file (default: `/etc/aws-lambda-server/config.yaml`)
- `--listen=<addr>`: use the provided listen address (default: `:8080`) for serving the requests towards the lambda
- `--metrics-listen=<addr>`: use the provided listen address (default: `:8081`) for serving metrics/health/readiness
  checks

## configuration

The configuration adds request matchers to a function. For the schema, start [here](pkg/config/server.go)

An annotated example config can be found [here](pkg/config/testdata/config.yaml).

## metrics, healthz and readyz

The application provides health (`/healthz`) and readiness checks (`/readyz`) listening to the
configured `--metrics-listen` address.

Additionally, the following [metrics are available](pkg/metrics/global.go):

- `http_requests_total`: counter for total http requests served
- `http_request_duration_seconds`: histogram for http request duration
- `http_request_size_bytes`: histogram for http request size
- `http_response_size_bytes`: histogram for http response size
- `aws_lambda_invocation_total`: counter for AWS lambda invocations by function ARN
- `aws_lambda_invocation_errors_total`: gauge for AWS lambda invocation errors by function ARN
- `aws_lambda_invocation_duration_seconds`: histogram AWS lambda invocation duration by function ARN

## helm-chart

Helm charts are created from the [charts directory](charts) and published to [this repository](https://tanemahuta.github.io/aws-lambda-server).

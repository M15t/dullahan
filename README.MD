# Dullahan - Calculating personal finance

## Prerequisites

- [Go](https://golang.org/doc/install) ^1.12.4
- [Docker](https://docs.docker.com/install/) ^18.09.2
- [Docker Compose](https://docs.docker.com/compose/install/) ^1.23.2
- [Chamber](https://github.com/segmentio/chamber) ^2.2.0
- [Serverless](https://www.serverless.com/framework/docs/getting-started) ^3.25.1
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv1.html) ^1.16.14
- [jq](https://stedolan.github.io/jq/download/) ^1.6.0
- [fswatch](https://github.com/emcrisostomo/fswatch#getting-fswatch) ^1.14.0
- [go-swagger](https://goswagger.io/install.html#homebrewlinuxbrew) ^0.21.0

## Getting started

1. Initialize the app for the first time:
   ```
   make provision
   ```
2. Generate swagger API docs:
   ```
   make specs
   ```
3. Run the development server:
   ```
   make start
   ```

The application runs as an HTTP server at port 8081. You can log in as superadmin to the application by:

Then grab the `access_token` for authorization HTTP header:

```http
GET /v1/... HTTP/1.1
Host: localhost:8082
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ...
```

## Deployment

```
make deploy
```

### To AWS Lambda

Coming soon!

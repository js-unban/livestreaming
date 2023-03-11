# livestreaming

Project for sending a HLS stream

## How to run

The project can be run using docker and passing an environment variable
for the signing key for tokens (SIGNING_KEY).

Example:

```
docker run -e SIGNING_KEY=foo -ti -p 8080:8080 --rm livestreaming ./livestreaming
```

# livestreaming

Project for sending a HLS stream

## How to run

The project can be run using docker and passing an environment variable
for the signing key for tokens (SIGNING_KEY).

Also content should live under /content folder and include a folder named segments with all segments
in the following format segment?.ts where "?" is a number

Example:

```
docker run -e SIGNING_KEY=foo -ti -p 8080:8080 --rm livestreaming ./livestreaming
```

## Updating models

Models can be updated under the /models folder and then the `wire` command[https://github.com/google/wire] needs to be ran on that folder, this will generate a wire_gen.so file that must be added when you commit.

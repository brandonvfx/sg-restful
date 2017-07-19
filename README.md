# SG Restful - WIP [![Build Status](https://travis-ci.org/brandonvfx/sg-restful.svg?branch=master)](https://travis-ci.org/brandonvfx/sg-restful) [![Docker Build Status](https://img.shields.io/docker/build/brandonvfx/sg-restful.svg)](https://hub.docker.com/r/brandonvfx/sg-restful)


SG Restful is a restful interface for the [Shotgun](http://shotgunsoftware.com)  Api.

## Docker

### Tags

- `0.6.0-beta.1`, `latest`
- `0.5.2-beta.1`

### Deploy

```
docker run -e SG_HOST=<your shotgun server> -p 8000:8000 brandonvfx/sg-restful
```

## Endpoints

- Find by id
    - GET /[entity type]/[id]
- Find all
    - GET /[entity type]
- Summarize
    - GET /[entity type]/summarize
- Create
    - POST /[entity type]
- Update
    - PATCH /[entity type]/[id]
- Delete
    - DELETE /[entity type]/[id]


## Auth

SG Restful using basic auth for getting script and user credentials. This may change in the future.

Script `Authorization` header:
```
Basic <base64 script_name:script_key>
```

User `Authorization` header:
( This is now safer to use than it was but I still suggest using https or only with internal servers )
```
Basic-User <base64 user_name:user_password>
```

## Query Strings

### Read
- page (int): Page of results to return.
- limit (int): Number of results per page to return
- fields (comma separated listed of string): The fields/columns to return.
- q (string): The query to execute. Syntax below.

### Summarize 
- q (string): The query to execute. Syntax below.
- summaries (json): array of hashes. Each hash should have 2 key/value pairs:
    - field: the Shotgun field to summarize
    - type: the type of summary to do. (see the Shotgun documentation)
- groupings (json): array of hashes. Each hash should have 3 key/value pairs: 
    - field: the Shotgun field to group on
    - direction: the direction to sort
    - type: the grouping type. (see the Shotgun documentation)

## Query Syntax

There are 3 formats for the query but they all have the same basic structures for the filters themselves. Each filter is defined by an array of 3 values.

```
[<name>, <relation>, <values>]
```

- Name and relation are both string.
- Values can be either a value (string, int, bool, etc) or an array of values.

### Format 1

```
q=and([<name>, <relation>, <values>],...)
q=or([<name>, <relation>, <values>],...)
```

### Format 2

```
q={"logical_operator":"and", "conditions":[[<name>, <relation>, <values>],...]}
```

### Format 3

This format the logical_operator is always assumed to be "and".

```
q=[[<name>, <relation>, <values>],...]
```

## Testing

### Tags

In order to facilitate testing, sg-restful makes use of a mocks from the testify package. Because of this, tests should be run using `script/test`.

### Writing Out a Log File

Optionally, you can write out a log for a test run by setting `SG-RESTFUL_LOG_TO_FILE="yes"` or `SG-RESTFUL_LOG_TO_FILE=true`. This will result in a log file getting generated with with the following name:

    sg-restful-test.<time>.log

Where <time> is the current time using the following format string `"20060102150405"`

You can either persist this value in your shell or prefix your go test command like so:

```
env SG-RESTFUL_LOG_TO_FILE="yes" script/test
```

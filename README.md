# SG Restful - WIP [![Build Status](https://travis-ci.org/brandonvfx/sg-restful.svg?branch=master)](https://travis-ci.org/brandonvfx/sg-restful)

## NOT PRODUCTION READY!

SG Restful is a restful interface for the [Shotgun](http://shotgunsoftware.com)  Api.


## What works currently

### Entities

- Read
    - Get by id
    - Get all
    - Returning fields
    - Pagination
- Create
- Update
- Delete


## Auth

SG Restful using basic auth for getting script and user credentials. This may change in the future.

Script `Authorization` header:
```
Basic <base64 script_name:script_key>
```

User `Authorization` header:
( I don't suggest using this unless you have an internal Shotgun sever. )
```
Basic-user <base64 user_name:user_password>
```

## Query String

- page (int): Page of results to return.
- limit (int): Number of results per page to return
- fields (comma separated listed of string): The fields/columns to return.
- q (string): The query to execute. Syntax below.


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

In order to facilitate testing, sg-restful makes use of a mocks from the testify package. Because of this, tests should be run using `go test -tags test`.

### Writing Out a Log File

Optionally, you can write out a log for a test run by setting `SG-RESTFUL_LOG_TO_FILE="yes"` or `SG-RESTFUL_LOG_TO_FILE=true`. This will result in a log file getting generated with with the following name:

    sg-restful-test.<time>.log

Where <time> is the current time using the following format string `"20060102150405"`

You can either persist this value in your shell or prefix your go test command like so:

```
env SG-RESTFUL_LOG_TO_FILE="yes" go test -tag test 
```

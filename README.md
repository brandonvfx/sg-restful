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
- fields (comma seperated listed of string): The fields/columns to return.
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

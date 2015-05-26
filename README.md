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
Basic <base64 scipt_name:sript_key>
```

User `Authorization` header:
( I don't suggest using this unless you have an internal Shotgun sever. )
```
Basic-user <base64 user_name:user_password>
```
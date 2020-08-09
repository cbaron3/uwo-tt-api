# REST API example application

This is a list of basic usage of the <strong>Unofficial UWO TT API.</strong>

Review the README for setting up the API.

# REST API

The base URL can either be:

* http://localhost:8080/api/v1/

OR

* http://uwottapi.ca/api/v1/

## Get options

### Request

`GET /campuses/`

    curl -i -H 'Accept: application/json' http://localhost:8080/api/v1/campuses

### Response

    HTTP/1.1 200 OK
    Status: 200 OK
    Connection: close
    Content-Type: application/json
    Content-Length: 1173
    X-Ratelimit-Limit: 120
    X-Ratelimit-Remaining: 119

    [{...},]

## Get courses

## Get sections

## Get with sorting

## Get with pagination

## Get with filter

## Get pagination, sorting, filtering

## Get with exclusive filters
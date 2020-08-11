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

### Request

`GET /campuses/`

    curl -i -H 'Accept: application/json' http://localhost:8080/api/v1/courses

### Response

    HTTP/1.1 200 OK
    Status: 200 OK
    Connection: close
    Content-Type: application/json
    Transfer-Encoding: chunked
    X-Ratelimit-Limit: 120
    X-Ratelimit-Remaining: 118

    [{...},]

## Get sections


`GET /campuses/`

    curl -i -H 'Accept: application/json' http://localhost:8080/api/v1/sections

### Response

    HTTP/1.1 200 OK
    Status: 200 OK
    Connection: close
    Content-Type: application/json
    Transfer-Encoding: chunked
    X-Ratelimit-Limit: 120
    X-Ratelimit-Remaining: 117

    [{...},]

## Get with pagination

`GET /sections/`

    curl -i -H 'Accept: application/json' http://localhost:8080/api/v1/sections?offset=15&limit=5

### Response

    HTTP/1.1 200 OK
    Status: 200 OK
    Connection: close
    Content-Type: application/json
    Transfer-Encoding: chunked
    X-Ratelimit-Limit: 120
    X-Ratelimit-Remaining: 116

    [{...},]

## Get with sorting


`GET /campuses/`

    curl -i -H 'Accept: application/json' http://localhost:8080/api/v1/campuses?sortby=text&dec=true

### Response

    HTTP/1.1 200 OK
    Status: 200 OK
    Connection: close
    Content-Type: application/json
    Content-Length: 1173
    X-Ratelimit-Limit: 120
    X-Ratelimit-Remaining: 115

    [{...},]

## Get with filter

`GET /courses/`

    curl -i -H 'Accept: application/json' http://localhost:8080/api/v1/courses?course-faculty=exact:PSYCHOL&course-number=gte:2000&course-number=lt:3000

### Response

    HTTP/1.1 200 OK
    Status: 200 OK
    Connection: close
    Content-Type: application/json
    Transfer-Encoding: chunked
    X-Ratelimit-Limit: 120
    X-Ratelimit-Remaining: 114

    [{...},]

## Get pagination, sorting, filtering


`GET /sections/`

    curl -i -H 'Accept: application/json' http://localhost:8080/api/v1/sections?course-faculty=exact:ACTURSCI&sortby=course-number&dec=true&offset=1&limit=5

### Response

    HTTP/1.1 200 OK
    Status: 200 OK
    Connection: close
    Content-Type: application/json
    Transfer-Encoding: chunked
    X-Ratelimit-Limit: 120
    X-Ratelimit-Remaining: 113

    [{...},]


## Get with inclusive filters

`GET /campuses/`

    curl -i -H 'Accept: application/json' http://localhost:8080/api/v1/sections?course-faculty=exact:ACTURSCI&course-faculty=exact:MSE&inclusive=true

### Response

    HTTP/1.1 200 OK
    Status: 200 OK
    Connection: close
    Content-Type: application/json
    Transfer-Encoding: chunked
    X-Ratelimit-Limit: 120
    X-Ratelimit-Remaining: 112

    [{...},]

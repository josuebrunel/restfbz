[![restfbz](https://github.com/josuebrunel/restfbz/actions/workflows/tuma.yml/badge.svg)](https://github.com/josuebrunel/restfbz/actions/workflows/tuma.yml)

# Service

Write a simple fizz-buzz REST server.

"The original fizz-buzz consists in writing all numbers from 1 to 100, and just replacing all multiples of 3 by ""fizz"", all multiples of 5 by ""buzz"", and all multiples of 15 by ""fizzbuzz"".
The output would look like this: ""1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,16,...""."

Your goal is to implement a web server that will expose a REST API endpoint that:
- Accepts five parameters: three integers int1, int2 and limit, and two strings str1 and str2.
- Returns a list of strings with numbers from 1 to limit, where: all multiples of int1 are replaced by str1, all multiples of int2 are replaced by str2, all multiples of int1 and int2 are replaced by str1str2.

The server needs to be:
- Ready for production
- Easy to maintain by other developers

Bonus: add a statistics endpoint allowing users to know what the most frequent request has been. This endpoint should:
- Accept no parameter
- Return the parameters corresponding to the most used request, as well as the number of hits for this request


# Doc

## Development

I didn't go with something fancy with a lot of external libraries.
The application has an internal app which is the service `restfbz` and uses 2 packages:
* fizzbuzz: which compute a fizzbuzz series according the input
* stats: which records the query string received in a Sqlite database

I wrote 2 middlewares:
* logging middleware to log requests
* stats middleware which uses the *stats* package to record the query strings or the incoming request

Last but not least there is a *simple Makefile* provided

## Deployment

Since no instructin has been given there is a  *docker-compose* file provided.

Alternatively I would go with:
1. A versioned binary built through an CI/CD pipeline and push to a repository
2. An ansible roles which would install the binary from the repository, create the required user and rights then set up a systemd config for the service

## Endpoints

```
Path: /
Params:
    * int1 : smallest integer
    * ini2 : biggest integer
    * limit: limit of result
    * str1 : string replacing multiple of int1
    * str2 : string replacing multiple of int2
```

```
Path: /stats
```

## Build the service binary

```shell
$ make build
```

## Run test

```shell
$ make test
```

## Run service via Docker

```shell
$ make docker-run
```

## Run tests via Docker

```shell
$ make docker-test
```

## Example

```shell
$ curl -v "http://127.0.0.1:8999/?int1=3&int2=5&limit=30&str1=foo&str2=bar" | jq
{
  "count": 30,
  "error": null,
  "data": [
    "1",
    "2",
    "foo",
    "4",
    "bar",
    "foo",
    "7",
    "8",
    "foo",
    "bar",
    "11",
    "foo",
    "13",
    "14",
    "foobar",
    "16",
    "17",
    "foo",
    "19",
    "bar",
    "foo",
    "22",
    "23",
    "foo",
    "bar",
    "26",
    "foo",
    "28",
    "29",
    "foobar"
  ]
}
```

```shell
$ curl -v "http://127.0.0.1:8999/stats" | jq
{
  "count": 1,
  "error": null,
  "data": {
    "qs": "int1=3&int2=5&limit=30&str1=foo&str2=bar",
    "hit": 5
  }
}
```

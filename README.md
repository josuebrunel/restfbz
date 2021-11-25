# Problem

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

The */stats* endpoint is not implemented. Below is the descriptin of what I would have done:

1. Have a single table with columns corresponding to all the required params, plus a *hit* column. This give more flexibility in a way stats on a single parameter value can be returned.
   Alternatively, it's possible to use an encoded string of *ordered query strings params* ( base64.b64encode(b"int1=3&int2=5&limit=25&str1=fizz&str2=buzz") ) in a 2 columns table with *encoded* and *hits*, since it is deterministic.  
2. For each request, I would record those parameter in a `StatsMiddleware` in a *UPDATE OR INSERT* manner


## Run service via Docker

```shell
$ make docker-run
```

## Run tests

```shell
$ make docker-test
```

## Example

```shell
$ curl -v "http://127.0.0.1:8999/?int1=3&int2=5&limit=25&str1=fizz&str2=buzz" | jq
[
  "1",
  "2",
  "fizz",
  "4",
  "buzz",
  "fizz",
  "7",
  "8",
  "fizz",
  "buzz",
  "11",
  "fizz",
  "13",
  "14",
  "fizzbuzz",
  "16",
  "17",
  "fizz",
  "19",
  "buzz",
  "fizz",
  "22",
  "23",
  "fizz",
  "buzz"
]
```

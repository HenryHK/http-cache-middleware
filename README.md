# http-cache-middleware

http-cache-middleware is a high performance http cache using redis, ideal for RESFful API. To maximize simplicity, the implementation using built-in `net/http` package, though code can be simplified using frameworks like `Gin`, `beego`, etc.

Cache expires after 10 mins by default, and can be configured for whatever.

This middleware should used behined auth, a simple `auth+value` key is used for user-specified cache naively.

## How to run

### Autopilot Access Key
To run the middleware, make sure you have a plain text file named `access` under the project folder, containing autopilot api key. 

### Redis
Make sure redis running on port 6379. Or you can modify the default in code in `main.go`.

### Run Application
To run the middleware, simply run:
```bash
$ go run main.go
```

or 
```bash
$ go build
```
and run the built binary.

The service is running on port 8080 by default and the endpoint is `/contact` supporting `GET`, `POST`, `PUT`, `DELETE`.

To test the api, use tools like `postman` or use the scripts under `scripts` which use `cURL`. Remember to replace `ACCESSKEY` in scripts with your api key. You may need to change the method from `PUT` to `POST` in `add.sh`. 

## Go Test

To ultilize modularization, each module has its own test, simply go to each package and run:
```bash
$ go test -v
```

## Possible Improvements

* Server-side cache refresh from autopilot
* More robust test to improve coverage
* Light-weight framework like `Gin` can be used to simplify code

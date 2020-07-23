# URL shortener

URL shortener web service in Golang. It allows to generate an ID for URLs and use them to easily access them.

# Description

All URL IDs are base64-encoded timestamps in milliseconds mainly because timestamps are natural identifiers, with no possibility of clashing.

Data are saved using [bitcask](https://github.com/prologic/bitcask), a key-value data store.

There are two DBs, one for storing short URLs
where keys are the encoded timestamps and values are the original URLs, and the other to keep track of the redirections count, where the keys are always the encoded timestamps and values are indeed the number of redirections.

Base64-encoding is necessary because keys must be strings (or better, slices of bytes). Values as well must be strings.

# Installation

First of all, open a terminal and navigate in the source code folder.

To init the environment and install all dependencies, enter

``` 
go mod tidy
```

To build the executable, enter

``` 
go build
```

To run the executable, enter

* On Linux: `./url-shortner`
* On Windows: `url-shortner.exe`
* On Mac: `TODO`
##
To bind a custom port, enter

* On Linux: `./url-shortner 5555`
* On Windows: `url-shortner.exe 5555`
* On Mac: `TODO`
##
If no port is provided or it's invalid, the default one is `8080` .

To run unit tests and see their coverage, enter

``` 
go test -cover
```

# API commands example

The web service runs on by default on `localhost:8080` if no port is provided as argument. If you have any other service bound on that port, please close it, otherwise the service won't run.

Examples here use [curl](https://curl.haxx.se/).

**IMPORTANT**: all URLs **must be valid**, in the form of `https://www.mycoolwebsite.com` (query params and segments are fine too). 

You need to provide the full URL, complete with the protocol and www, otherwise the request is rejected. 

For `GET` requests we suggest to use a browser instead.

## Full API documentation

To access the API documentation, with response and errors description, from a browser go to `http://localhost:8080` (or the custom port you set).

## Shorten a URL

``` 
curl -v -X POST http://localhost:8080/api -d "https://duckduckgo.com/?q=very+long+query+for+a+long+url&t=ffnt&atb=v222-1&ia=web"
```

Response:

``` json
{
	"urlKey": "MTU5NTI3OTY5NjgwNw",
	"location": "http://localhost:8080/api/MTU5NTI3OTY5NjgwNw",
	"message": "OK"
}
```

We set the `Location` header as well and the redirections count is initialized to 0.

## Access the original URL

From a browser, using the URL provided in the `POST`/`PUT` response, you will be redirected to the original website.
Redirection counts is increased by 1.

## Update the URL using an ID

``` 
curl -v -X PUT http://localhost:8080/api/MTU5NTI3OTY5NjgwNw -d "https://duckduckgo.com/?q=another+very+long+query+for+another+long+url&t=ffnt&atb=v222-1&ia=web"
```

Response

``` json
{
	"urlKey": "MTU5NTI3OTY5NjgwNw",
	"location": "http://localhost:8080/api/MTU5NTI3OTY5NjgwNw",
	"message": "OK"
}
```

We set the `Location` header as well. Redirections count is reset to 0.

## Delete the URL by the ID

``` 
curl -v -X DELETE http://localhost:8080/api/MTU5NTI3OTY5NjgwNw
```

Response

``` json
{
	"message": "URL successfully deleted for key MTU5NTI3OTY5NjgwNw"
}
```
We delete the redirections count entry, as well.

## Get URL redirections count

``` 
curl -v http://localhost:8080/api/count/MTU5NTI4MDAwNzA2MQ
```

Response

``` json
{
	"redirectionsCount": 2
}
```

# 

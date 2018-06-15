# Shrink your url 
A Go lang and sql based project

## Repository Contains

```
public
    index.html
    404.html
main.go
main_test.go
Dockerfile
Readme.md
```

## How to run

### Build new docker image 
```
docker build -t ankushkawanpure/url-awesome .
```
OR 
### Pull using docker hub  
```
 Pull docker image
 docker pull ankushkawanpure/url-awesome
 
 Run docker image locally with localhost
 docker run --publish 8080:8080 -t ankushkawanpure/url-awesome
 
 OR
 
 Run with some domain
 docker run --publish 8080:8080 -e DOMAIN=http://example.com -t ankushkawanpure/url-awesome 
```

### Run in go
```
go run main.go
```

### Hop on to
```
http://localhost:8080/
```


## Testing for this project

To run all the test with coverage

```
go test -cover
``` 

Run with coverage profile 

```
go test -coverprofile=coverage.out   
go tool cover -html=coverage.out    
```


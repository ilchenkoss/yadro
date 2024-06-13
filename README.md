<h2 align="center">Service for searching and retrieving comics pictures from xkcd.com by keywords</h2>

<p align="center">
  <img src="https://github.com/ilchenkoss/yadro/blob/main/example.gif" alt="Example">
</p>

##  Service Structure

The service consists of three microservices:

1. **Xkcd Server**:
    - Responsible for multithreaded data parsing and keyword processing.
    - Includes a request limit to xkcd.com.
    - Uses a concurrency limiter to manage interactions with the service.
    - Uses a rate limiter to control the frequency of requests.
    - **Run**:
      ```sh
      make xkcd
      ```
    
2. **Auth Server (gRPC)**:
    - Responsible for user authorization and authentication.
    - **Run**:
      ```sh
      make auth
      ```
    
3. **Web Server**:
    - Responsible for displaying results graphically and facilitating easy interaction with the services..
    - **Run**:
      ```sh
      make web
      ```
    - **Connect**: [http://localhost:33333](http://localhost:33333)
    
All servers are running with Makefile:
```sh
make servers_start
```

## Authentication

```sh
curl -X POST \\
  -H "Content-Type: application/json" \\
  -d '{"login": "login", "password": "password"}' \\
  http://localhost:22222/login
```

## Privilege Elevation

```sh
token="superuser token"
curl -X POST \\
  -H "authorization: bearer $token" \\
  -d '{"login": "user1"}' \\
  http://localhost:22222/toadmin
```

## Retrieving Comic Pictures by Keywords

```sh
curl -X GET \\
  -H "authorization: bearer (any user's token)" \\
  "http://localhost:22222/pics?search=binary,christmas,tree"
```

## Update Comics

```sh
token="admin token"
curl -X POST \\
  -H "authorization: bearer $token" \\
  http://localhost:22222/update
```

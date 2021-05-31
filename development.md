# Development FYI

## Backend (Golang, Gin)

```shell
go run . serve
```

## Frontend (Vite, Vue)

```shell
npm i
vite
```

## Log Verbosity Conventions

- Exit: exceptions that cause the application unable to continue running;
- Error: exceptions with serious consequences and should be fixed ASAP;
- Warning: exceptions that can be handled, inspection is required;
- Info: messages showing the application flow, typically the bootstrap/shutdown process;
- V1: more detailed bootstrap/shutdown messages
- V2: verbose bootstrap/shutdown messages
- V3: undefined
- V4: goroutine (workers, cronjob) level messages
- V5: event level messages
- V6: print http request, like "GET https://spongeprojects.com/healthz 200 OK in 151 milliseconds"
- V7: print http request with request headers
- V8: print http request with request/response headers and truncated body
- V9: print http request with request/response headers and truncated body
- V10: print http request with request/response headers and complete body

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
- V2: undefined
- V3: undefined
- V4: worker/controller level messages
- V5: event level messages
- V6: print http request url
- V7: print http request url, request headers
- V8: print http request url, request/response headers and truncated body
- V9: print http request url, request/response headers and truncated body
- V10: print http request url, request/response headers and full body

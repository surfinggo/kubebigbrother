# kubebigbrother

Kube Big Brother is a platform that monitors and records everything happens in a cluster.

> Big brother is watching you.

## Usage

Start watching events lively without interacting with the backend:

```shell
./kbb watch --config=/path/to/config.yaml
```

Start the server:

```shell
./kbb serve --config=/path/to/config.yaml
```

Start recording events:

```shell
./kbb record --config=/path/to/config.yaml
```

## Development

### Backend (Golang, Gin)

```shell
go run . serve
```

### Frontend (Vite, Vue)

```shell
npm i
vite
```

## ToDo (PRs are welcomed)

- Avoid duplicated "ADDED" events and missed "DELETED" events, by storing the current state in database.
- Better UI.
- Tests.

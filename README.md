# kubebigbrother

Kube Big Brother is a platform that monitors and records everything happens in a cluster.

> Big brother is watching you.

## Usage

### Watch

You can start watching events lively without interacting with the backend:

```shell
./kbb watch
```

### Controller

Start the controller (only 1 controller should be running simultaneously):

```shell
./kbb controller
```

The controller is responsible for handling all events, including sending notifications and recording them into the
database.

### Serve

Start the frontend server:

```shell
./kbb serve
```

Instead of connecting to Kubernetes API server directly, the server is connected to the database.

### Query

You can query the database with query command:

```shell
./kbb query
```

## Config

### Channels

Currently, kbb supports these channel types:

## Development

[Development](./development.md)

## ToDo (PRs are welcomed)

- Avoid duplicated "ADDED" events and missed "DELETED" events, by storing the current state in database.
- Better UI.
- Tests.

# kubebigbrother

[![Go Report Card](https://goreportcard.com/badge/github.com/spongeprojects/kubebigbrother)](https://goreportcard.com/report/github.com/spongeprojects/kubebigbrother)
[![License](https://img.shields.io/github/license/spongeprojects/kubebigbrother?color=blue)](https://github.com/spongeprojects/kubebigbrother/blob/main/LICENSE)

Kubebigbrother is a platform that monitors and records everything happens in a cluster.

> Big brother is watching you.

## Usage

There are two interfaces of kubebigbrother: the GUI, and the CLI.

For the GUI, you need to start a controller to records events and a server to serving the frontend UI, for the CLI, you
can use the watch command to start watching events lively from the Kubernetes API server, without interacting with any
backend.

```text
Usage:
  kbb [command]

Available Commands:
  controller  Run controller, watch events and persistent into database (only one instance should be running)
  help        Help about any command
  query       Query event history
  serve       Run the server to serve backend APIs
  watch       Watch events lively
```

Global flags:

```text
  -c, --config string                    path to config file (klog flags are not loaded from file, like -v) (default "config/config.local.yaml")
      --env string                       environment (default "debug")
  -h, --help                             help for kbb
```

It's recommended to set flags in a config file, these file types are supported:

```text
"json", "toml", "yaml", "yml", "properties", "props", "prop", "hcl", "dotenv", "env", "ini"
```

Besides, all klog flags are registered as global flags:

```text
      --add_dir_header                   If true, adds the file directory to the header of the log messages
      --alsologtostderr                  log to standard error as well as files
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --log_file string                  If non-empty, use this log file
      --log_file_max_size uint           Defines the maximum size a log file can grow to. Unit is megabytes. If the value is 0, the maximum file size is unlimited. (default 1800)
      --logtostderr                      log to standard error instead of files (default true)
      --one_output                       If true, only write logs to their native severity level (vs also writing to each lower severity level)
      --skip_headers                     If true, avoid header prefixes in the log messages
      --skip_log_headers                 If true, avoid headers when opening log files
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
  -v, --v Level                          number for the log level verbosity
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```

Example:

```shell
./kbb

# You can always print more information if you encounter any problem,
# the log verbosity convention can be found in `./development.md`.
./kbb -v 5

# The default value for `env` is "debug",
# in debug mode, by default, config file "./config/config.local.yaml" will be used, 
# in production mode, it's recommended to specify a config file explicitly.
./kbb -e production -c /opt/config.yaml
```

### Watch

You can start watching events lively from the Kubernetes API server, without interacting with any backend:

```shell
./kbb watch
```

Supported flags:

```text
      --informers-config string   path to informers config file (default "config/informers-config.local.yaml")
      --kubeconfig string         path to kubeconfig file (default "/Users/wujunchao/.kube/config")
```

#### Controller

The controller is responsible for handling all events, including sending notifications and recording them into the
database.

Start the controller (only 1 controller should be running simultaneously):

```shell
./kbb controller
```

Supported flags:

```text
      --db-args string            database args
      --db-dialect string         database dialect [mysql, postgres, sqlite] (default "sqlite")
      --informers-config string   path to informers config file (default "config/informers-config.local.yaml")
      --kubeconfig string         path to kubeconfig file (default "~/.kube/config")
```

#### Serve

Start the frontend server:

```shell
./kbb serve
```

Supported flags:

```text
      --addr string         serving address (default "0.0.0.0:8984")
      --db-args string            database args
      --db-dialect string         database dialect [mysql, postgres, sqlite] (default "sqlite")
```

Instead of connecting to Kubernetes API server directly, the server is connected to the database.

#### Query

You can query the database with query command:

```shell
./kbb query
```

Supported flags:

```text
      --db-args string            database args
      --db-dialect string         database dialect [mysql, postgres, sqlite] (default "sqlite")
```

## Config

You can specify

### Channels

Currently, kbb supports these channel types:

## Development

[Development](./development.md)

## ToDo (PRs are welcomed)

- Avoid duplicated "ADDED" events and missed "DELETED" events, by storing the current state in database.
- Register watchers as CRD, maybe keep compatibility of config file.
- Better UI.
- Tests.

# Task Manager

This repository contains a sample daily task manager API, database and client application. It's a work in progress.

At present moment all elements are written in `Go` and will work in tandem when docker is installed.

## Prerequisites

- `Go` 1.14 or higher
- `Docker` 19.0 or higher

## Running

To run the server, use:

```
cd docker && ./run.sh
```

Then you can use a CLI client, a Web UI or access the API directly to view
and manipulate the tasks.

## Web UI

Go to `localhost:8690` for a rudimentary display of today's tasks.

## CLI Client

To build the client:

```
go build -mod=vendor ./client-go -o tasks
cp ./client-go/cfg.yml .
```

Then you can use the `task` binary to query the server, for example by typing:

```
./task t
```

To show tasks lined up for today. Or:

```
./task l
```

To list all available tasks.

Additional functionality has not yet been implemented.

## API

Feel free to play with the API using Postman or other tools. Documentation is
provided as a yaml OpenAPI file in the `doc` folder.

## Testing

Run `go test ./...` to perform unit and integration tests on all packages in the repository.

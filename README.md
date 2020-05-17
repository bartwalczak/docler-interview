# Docler-interview

This repository contains a sample daily task manager API, database and client application.

At present moment all elements are written in `Go` and will work in tandem when docker is installed.

## Prerequisites

- `Go` 1.14 or higher
- `Docker` 19.0 or higher

## Running

To run the server, use:

```
cd docker && ./run.sh
```

To run the client:

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

## Testing

Run `go test ./...` to perform unit and integration tests on all packages in the repository.

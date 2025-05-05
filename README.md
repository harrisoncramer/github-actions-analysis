# Github Analysis

This is a simple script for downloading Github action stats and performing an analyis on them.

## Usage

See the `.example.env` file for configuration options.

Build the binary:

```sh
go build .
```

First `collect` the Github Actions data:

```sh
go build . && ./gh-analysis collect
```

Then run an analysis on the data:

```sh
go build . && ./gh-analysis analyze
```

# Github Actions Analysis

This is a simple script for downloading Github Action runs and performing an analyis on them.

## Usage

See the `.example.env` file for configuration options.

Build the binary:

```sh
go build .
```

First `collect` the Github Actions data:

```sh
go build . && ./github-actions-analysis collect
```

Then run an analysis on the data:

```sh
go build . && ./github-actions-analysis analyze
```

The output of the analysis will look something like this:

```csv
Job Name,Count,Avg (seconds),Min (seconds),Max (seconds),P90 (seconds),P99 (seconds)
lint,66,244.30,0,823,554,615
build,25,84.00,56,209,104,187
release,1,258.00,258,258,258,258
```

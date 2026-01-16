# Benchmarks

A multi-language repo containing benchmarks I ran to test different options and approaches.

## Prerequisites

To run the benchmarks, you need to have `hyperfine` installed. Make sure to make the `run.sh` scripts executable by running `chmod +x ./**/run.sh` from the root folder.

## Running the benchmarks

### TUI (Recommended)

Use the included TUI. Run `./bench` (build with `cd benchmark_runner && go build -o ../bench .`) to start the TUI.

### Manual

Run `./<LANGUAGE_FOLDER>/run.sh <benchmark_name>`. For example:

```sh
./JS-TS/run.sh array_includes_vs_set_has
```
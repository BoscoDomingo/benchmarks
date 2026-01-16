#!/bin/bash
# Shared run script for JS-TS benchmarks
# Usage: ./run.sh <benchmark_path>
# Example: ./run.sh array_includes_vs_set_has

BENCHMARK_PATH="${1:-.}"

if [ ! -d "$BENCHMARK_PATH" ]; then
    echo "Error: Benchmark path '$BENCHMARK_PATH' does not exist"
    exit 1
fi

EXPORT_FLAG=""
if [ "${EXPORT_RESULTS:-0}" = "1" ]; then
    EXPORT_FLAG="--export-markdown $BENCHMARK_PATH/results.md"
fi

find "$BENCHMARK_PATH" -maxdepth 1 -name "*.ts" ! -name "setup.ts" -print0 | xargs -0 -I {} echo "bun {}" | xargs -d '\n' hyperfine --warmup 10 --min-runs "${MIN_RUNS:-300}" $EXPORT_FLAG

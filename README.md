# Overview

This is an attempt to solve the Vehicle Routing Problem

# Prerequisites

* Go version 1.23
* Python 3+

# Installation and Generating Schedules

run `make build` to compile the binary for your OS/Architecture. This should auto detect, if you need to build a binary for a system different than the compiling system, use `make build-cli`. The binary for your system should be at `build/vrp`.

Run `./build/vrp testdata/[pick a file].txt` to generate a schedule based on loads contained within the .txt file

Run `make test` to execute a benchmark and end to end test against all route/load files in `testdata`
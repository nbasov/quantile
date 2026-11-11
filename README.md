# Quantile calculation

This program reads a JSONL file line by line, takes a number from one field of each JSON object, and prints the current value of the chosen quantile after every line.

## How to build and run the code

You need Go 1.27. On macOS you can install it with following command:
```sh
make setup
```

wich is just the shortcut for Brew call:

```sh
brew install go@1.27
```

After that, add Go to your `PATH` as the command output says.
```
export PATH="$$(brew --prefix go@1.27)/bin:$$PATH"
```

Main commands:

```sh
make demo       # build and run the program on demo.jsonl (this is the default target)
make build      # build the binary file ./quantile
make test       # run all tests with the race detector
make load-gen   # create a test file 1M.jsonl with 1 000 000 lines
make load-run   # build and run the program on 1M.jsonl, show time and memory use
make clean      # remove the binary file and 1M.jsonl
```

You can change the number of lines in the test file with the `LOAD_N` variable, for example `make load-gen LOAD_N=5000000`.

To run the program by hand:

```sh
./quantile -q 0.75 -f latency_ms
```

Command line arguments:

* `-q` is the quantile, a number from 0.0 to 1.0. Required.
* `-f` is the name of the JSON field with the numbers, case sensitive. Required.
* `-v` turns on verbose output. In this mode the program also prints the arguments at the start and line statistics at the end. Off by default.

After start, the program asks for the absolute path to the input file. Type the absolute path and press Enter. You can also send the absolute path through a pipe:

```sh
echo /demo.jsonl | ./quantile -q 0.75 -f latency_ms
```

## Assumptions and choice of algorithm

There are two groups of algorithms for quantile calculation: exact and approximate.

Exact algorithms keep every value from the input. Because of this, they always give the exact quantile. They are also quite simple to write. The problem is memory: if you do not use a sliding window, memory use grows with every new value, and on large data sets it can become very big.

Approximate algorithms do not keep all the values. Instead, they keep a small set of numbers that describe the data. From these numbers they can estimate the quantile with a known error. They need much less memory, but the result is not exact and the code is more complex.

The task does not give a strict memory limit. It only says that the input can have "many millions of records". So I calculated how much memory an exact algorithm would need:

* Each value is stored as `float32`, which takes 4 bytes.
* 100 000 000 values take 400 000 000 bytes, which is about 381 MB of RAM.

This is acceptable for a normal computer, so I decided to use an exact algorithm.

I wrote two exact algorithms:

* A simple one that keeps the values in a slice (`quantile_simple.go`).
* One based on the Two Heaps method (`quantile_two_heaps.go`).

The final version uses the Two Heaps algorithm. Both algorithms use the same amount of memory, but adding a new value is about 10 times faster with two heaps.

How the Two Heaps algorithm works, in short:

* The values are split into two heaps. The left heap is a max-heap and keeps the smaller values. The right heap is a min-heap and keeps the bigger values.
* After each new value, the program moves one value between the heaps if needed, so that the left heap always has the right size for the chosen quantile.
* The quantile is then found from the top of the left heap and the top of the right heap, using linear interpolation between them.
* Adding a value takes O(log n) time. Reading the current quantile takes O(1) time. Memory use is O(n).

On my computer, 1 million records are processed in 2.8 seconds.

## Handling of JSON errors

If a line cannot be used, the program does not stop. Instead of the current quantile value, it prints an error message to `STDOUT` for this line and goes on to the next line. The messages are:

* `Malformed line` means the line is not valid JSON.
* `Non-numeric field value` means the field exists, but its value is not a number.
* `Missing field` means the JSON object does not have the field given in `-f`.

In verbose mode (`-v`) the program also prints at the end how many lines had each type of error, and the total number of lines.

## Implementation details

The solution has these parts:

* `main.go` is the entry point. It reads the arguments and connects all parts together.
* `reader.go` reads lines from the JSONL file.
* `parser.go` parses the JSON in each line and takes the value of the given field.
* `computer.go` is a common interface for the quantile algorithms, so they can be easily swapped.
* `quantile_two_heaps.go` is the Two Heaps algorithm. This is the one the program uses.
* `quantile_simple.go` is the simple algorithm. It is kept only for comparison and is not used.
* `quantile_dumb.go` is an empty algorithm that does nothing. It is kept only for comparison and is not used.

In `main.go` (around line 27) the work is organized as a pipeline. Each step runs in its own goroutine and sends its results to the next step through a channel: reading lines, then parsing, then calculating the quantile, then printing. Channels add a small cost compared to one big loop with many `if` branches. But the code is easier to read, and each step does only one job, so the logic of the program is easy to follow.

The field name argument (`-f`) is case sensitive. This is done on purpose. JSON keys are case sensitive themselves, so `latency_ms` and `Latency_ms` are different fields. Also, the program does not need to change the case of every key in every line, which saves work.

Numbers are stored as `float32`. Its range and precision are enough for this task, and it uses half the memory of `float64`. The smaller `float16` was not used because its range of values is too small.

JSON is parsed with `json.Unmarshal` into a map. The same map is cleared and used again for every line. For small JSON objects, like the ones in the example file, this makes fewer memory allocations than a solution based on `json.Decoder`.

## Use of AI

I wrote the program and the algorithms by hand.

I used AI for these tasks:

* Looking at existing ways to solve the problem
* Comparing algorithms
* Describing the steps and the general structure of specific algorithms
* Checking that all edge cases are covered
* Writing tests
* Improving and cleaning up the Makefile
* Improving and proofreading README.md

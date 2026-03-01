# t6t: T6-Tools

`t6t` is a program for analyzing any .gsc script written for T6 (Call of Duty: Black Ops II). It allows for fast, easy analysis of functions, variables, flags, dvars etc.

## Installation

### Install binary

1. Download the latest binary
2. (optional) move it into /usr/bin/

### Build from source (requires Go)

```bash
# Clone the repo
git clone https://maxvanasten/t6t
cd t6t
./build.sh
# Optional but recommended: move t6t binary into /usr/bin/ so its accessible from anywhere
sudo mv ./build/t6t /usr/bin/t6t
```

## Usage

`t6t` depends on [gscp](https://github.com/maxvanasten/gscp) to parse raw .gsc files into abstract syntax trees. The easiest way to use `t6t` is by piping the output of `gscp` into `t6t`.

```bash
# Reads ast from STDIN using gscp
gscp -p input_file.gsc | t6t -f
# Reads ast from input file
t6t -f ast.json

# I personally like to use this command to view the output in a pretty way using jq and bat
gscp -p input_file.gsc | t6t -f | jq | bat -l json
```

## Flags

### -f

Function signatures and function calls

## Demo

You can find a demo analysis output in [./demo/analysis.json](./demo/analysis.json)

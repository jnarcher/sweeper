# Sweeper 

Play Minesweeper inside your terminal.

[screenshot](github/screenshot.png)

## Getting Started

### Dependencies

Make sure to have go downloaded and installed ([download here](https://go.dev/doc/install))

### Installing

First clone the repo
```bash
git clone https://github.com/jnarcher/sweeper
```
Then enter the directory and build the executable
```bash
cd sweeper && make
```
The executable can be found in the `target` directory.

### Usage

The following flags can be set to configure the board:
- `w` : width  (columns)
- `h` : height (rows)
- `b` : number of bombs

If no value for the number of bombs is provided, the program will default to 1/6th of the total squares available.

## Acknowledgments

* [bubbletea](https://github.com/charmbracelet/bubbletea/tree/master)

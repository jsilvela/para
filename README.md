# para

A simple command-line utility to wrap text to a given number of columns.

*para* wraps text to a given column, filling each line as much as possible,
to form compact paragraphs.

However, it respects full stops and paragraph breaks.
It also respects Markdown headers and list items.
It would be easy to implement more of Markdown.

*para* is meant for wrapping of text in *Sam* and *Acme* , or any other
editors which break words or don’t wrap lines.

## Usage

``` sh
para {optional wrap column}
```

If the optional number is not provided, the default wrap column is 80.
*para* takes input from `stdin` and dumps it to `stdout`; this is so that it can be
used as a piped command in *Acme* or *Sam*.

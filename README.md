# para
A simple command-line util to wrap text to a given column.
*para* compresses the text to form compact paragraphs, and respects full stops
and paragraph breaks.
It also respects Markdown headers and list items (easy to implement more of it).

It is meant for display of text on the *Sam* and *Acme* editors, or any other
editors which break words or don’t wrap lines.

## Usage
```
> para {optional wrap column}
```
If the optional number is not provided, the default wrap column is 80.
*para* takes input from stdin and dumps it to stdout; this is so that it can be
used as a pipe command in *Acme* or *Sam*.

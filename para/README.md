# para
A simple command-line util to wrap text to a given column.
Useful for display of text on the Sam and Acme editors, or any other editors which break words or don’t wrap lines.

# Usage
```
> para {optional wrap column}
```
If the optional number is not provided, the default wrap column is 90.
This takes input from stdin and dumps it to stdout. Again, useful for Acme.
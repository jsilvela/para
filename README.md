# para

A simple command-line utility to wrap text to a given number of columns. \
Meant for wrapping of text in *Sam* or *Acme*, or any other
editors which don't do their own wrapping/compacting.

*para* wraps text to a given column, filling each line as much as possible,
to form compact paragraphs.
When compacting, it respects full stops and paragraph breaks.
It also respects Markdown headers and list items.

## Usage

``` sh
para [wrapColumn (optional)]
```

If the optional number is not provided, the default width is 80.

*para* takes input from `stdin`, and outputs to `stdout`; this is so that
it can be used as a piped command in *Acme* or *Sam*.

## Building / testing

Clone the repo, then.

``` sh
cd cloned-repo
go test
go build
go install
```

Note that depending on your Go environment, you may need to execute the
`go install` with privileges.

## How it works

`para` fills two goals: wrapping and compacting.

Wrapping on its own is easy. We keep a running count of the length since
the last newline, and when we exceed the max width, we change the last
whitespace into a newline:

``` text
word1 word2 word3 … wordOverflows wordNext …
```

Would become:

``` text
word1 word2 word3 …
wordOverflows wordNext …`
```

But if we want to compact paragraphs, we need some interaction between lines
of input.

Imagine that after wrapping a given input line, we get the following.

``` text
word word word word word word word word word word word word
word word word word word word word word word word word word
word
```

We would want the next input line to continue filling the same paragraph.

``` text
word word word word word word word word word word word word
word word word word word word word word word word word word
word newLineWord newLineWord newLineWord newLineWord …
```

This means that when wrapping a line of input:

1. we don't necessarily want to "close" the paragraph with a newline
1. we may want to append to an un-closed paragraph seen previously

To fill these needs, we keep track of the *carry* between lines of input.
c.f. the signature of the central function:

``` go
wrapLine(line string, carry int) (wrapped string, newCarry int)
```

With this signature in place, the code becomes clear. Special care needs to
be taken with:

- paragraph breaks (i.e. empty lines)
- full stops (i.e. lines of input ending with `".\n"`)
- Markdown sections
- Markdown lists

These should be respected. Compaction should be skipped for them.
We achieve this by *flushing*  the un-closed running text:

``` go
flushRunningText := func() {
    if carry > 0 {
        writer.WriteString("\n")
        carry = 0
    }
}
```

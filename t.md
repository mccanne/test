# Test the foo bar

These are some notes on the foo bar test.

# Stdin

```
here
is
some
input
data
```

# Command

The foo gets its input from stdin and should produce the right output...

```
foo -r - run-sort | bar -prefix "x: "
```

# Expected

```
x: data
x: here
x: input
x: is
x: some
```

# Testing

~~~
foo
bar
~~~


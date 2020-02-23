# foo bar test

These are some notes on the foo bar test.  This text here is just
some descriptive notes and would be ignored by the test harness.

## stdin

```
here
is
some
input
data
```

## script

The foo gets its input from stdin and should produce the right output...

```
foo -r - run-sort | bar -prefix "x: "
```

## expected

```
x: data
x: here
x: input
x: is
x: some
```

# simple test

This is a simple test showing how you can use files.

## Command

The foo gets its input from stdin and should produce the right output...

```
echo a > F
echo c > F
echo b > F
sort F > F2
cat F2
```

When the test is done, the files in this test environment would be deleted,
unless there was an error, in which case, they would be left behind so
the tester can debug.

## Expected

```
a
b
c
```

# pcap command test

Steve wishes he could write his pcaps test like this, or using
some similar methodology.

He would put a test input pcap
brimsec/zqd/tests/data/pcap/test.pcap.

When the script is run a shell variable (e.g., *$data*)
is set to the data input
root (how the data path is bound to the particular would have to
be worked out and could either be supplied here or implied by the
test name).

## script

The foo gets its input from stdin and should produce the right output...

```
pcap index -x test.index -r $data/test.pcap
pcap slice -x test.index -r $data/test.pcap -from X -to Y -w out.pcap <args>
pcap ts -r out.pcap
```

## Expected

```
1425567432.793221
1425567047.803929
1425567047.804906
```

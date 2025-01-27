When writing godoc, sentences should always end with a dot.

When writing golang unit-tests, always try to use a map to store test cases.
If should follow the same pattern as:
```go
tests := map[string]struct {
  input string
  result string
} {
  "empty string":  {
    input: "",
    result: "",
  },
  "one character": {
    input: "x",
    result: "x",
  },
  "one multi byte glyph": {
    input: "🎉",
    result: "🎉",
  },
  "string with multiple multi-byte glyphs": {
    input: "🥳🎉🐶",
    result: "🐶🎉🥳",
  },
}

for name, test := range tests {
  t.Run(name, func(t *testing.T) {
    t.Parallel()
    g := NewWithT(t)
    result := reverse(test.input)
    g.Expect(result).To(Equal(test.result))
  })
}
```
It should use github.com/onsi/gomega to run assertions. github.com/onsi/gomega should be imported using dot.
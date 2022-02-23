# Reference Gen

This is a utility for generating markdown references from Go structs.
It is used to generate part of the documentation of the [OAuth2 Proxy](https://github.com/oauth2-proxy/oauth2-proxy) project.

For example, a struct as below:

```golang
// MyStruct contains a collection of fields.
type MyStruct struct {
	// Name is the name of the MyStruct.
	Name string `json:"name"`

  // Data is a slice of raw byte data.
	Data []byte `json:"bytes"`
}
```

Would be turned into markdown such as:

```markdown
### MyStruct

MyStruct contains a collection of fields.

| Field | Type | Description |
| ----- | ---- | ----------- |
| `name` | _string_ | Name is the name of the MyStruct. |
| `data` | _[]byte_ | Data is a slice of raw byte data. |
```

Where a struct contains another struct, the other structs will also be included
in the api references. Check out the [test data](https://github.com/oauth2-proxy/tools/tree/master/reference-gen/pkg/generator/testdata)
for full examples of more complex struct documentation generation.

## Running tests

Tests can be executed using the `test` target from the Makefile.

```bash
make test
```

The tests will also run the `lint` target which requires [golangci-lint](https://golangci-lint.run/usage/install/).
You will be prompted to install it should it not already be installed.

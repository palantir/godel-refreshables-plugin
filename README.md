<p align=right>
<a href=https://autorelease.general.dmz.palantir.tech/palantir/godel-refreshables-plugin><img src=https://img.shields.io/badge/Perform%20an-Autorelease-success.svg alt=Autorelease></a>
</p>

# godel-refreshables-plugin

A godel plugin for generating strongly-typed [refreshable](http://pkg.go.dev/github.com/palantir/pkg/refreshable)
wrappers for arbitrary types.

For each target type, and the types required to construct it, the plugin will generate an interface definition (and
implementation) which provides typed variants of the `Current`, `Map`, and `Subscribe` methods from the
`refreshable.Refreshable` interface. Struct types additionally have generated methods to access each field as a Refreshable.

### Plugin Configuration

The plugin reads a configuration file at `godel/config/refreshables-plugin.yml` which specifies the types for which
Refreshable wrappers will be generated.

Example:

```yaml
refreshables:
  # Relative path to local package
  ./pkg/mypackage:
    types:
      - MyType
  # Absolute path to local package
  github.com/user/project/pkg/mypackage:
    types:
      - MyType
  # Absolute path to remote package
  # In this case, an output path is required.
  github.com/otheruser/otherpackage:
    output: ./generated/otherpackage/zz_generated_refreshables.go
    types:
      - MyType
```

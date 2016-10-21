godeps2bazel
============

`godeps2bazel` takes a [Godeps](https://github.com/tools/godep) JSON file and
turns it into the equivalent [bazel](https://www.bazel.io/) repository
commands.

It's handy to use alongside
[gazelle](https://github.com/bazelbuild/rules_go/tree/master/go/tools/gazelle/gazelle),
the Go BUILD file generator.

Install with `go get github.com/jmhodges/godeps2bazel`.

Run as `godeps2bazel /path/to/Godeps.json` and copy the output into your
WORKSPACE (or wherever you're storing your repositories).

Be sure to add the [rules_go](https://github.com/bazelbuild/rules_go) repository
to your build and to include `"new_go_repository"` in the `load` call for
`rules_go` where you store the output.

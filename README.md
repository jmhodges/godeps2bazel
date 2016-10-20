= godeps2bazel

`godeps2bazel` takes a [Godeps](https://github.com/tools/godep) JSON file and
turns it into the equivalent [bazel](https://www.bazel.io/) repository commands.

Install with `go get github.com/jmhodges/godeps2bazel`.

Run as `godeps2bazel /path/to/Godeps.json` and copy the output into your
WORKSPACE (or where ever you're storing your repositories).

Wherever you do copy that data, because sure to include the
[rules_go](https://github.com/bazelbuild/rules_go) rules and to include
`"new_go_repository"` in the `load` call for `rules_go`.

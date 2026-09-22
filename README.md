# wc-tools

WywyCode's tool monorepo: the 6 default tools, each a Go module producing static binaries
that share the `wc-tools/shared` runtime lib, plus per-tool manifests. wc-tools is a git
submodule of the wc meta repo; its deployment contract is `build-tools.sh`, which stages
the tool tree into the pod image at `/home/wc/tools/` (decision 112).

## Tools

| Tool      | Purpose                                                                 |
| --------- | ----------------------------------------------------------------------- |
| `read`    | Print a file's contents, optionally paged (offset / line / byte limits) |
| `write`   | Create or overwrite a file, then run configured formatters              |
| `edit`    | Exact search/replace edit, then run configured formatters               |
| `bash`    | Execute a shell command, surfacing exit code                            |
| `skill`   | Resolve skills from `$XDG_CONFIG_HOME/skills`                           |
| `compact` | Mark older session messages dead via the harness `POST /set-active`     |

## Deployment

The wc Dockerfile's build stage runs `wc-tools/build-tools.sh <dest>` to stage the tool tree. It creates the following inside `/home/wc/tools`:

- `<dest>/<name>` executables
- `<dest>/manifests/<name>.yaml`

```sh
./build-tools.sh /home/wc/tools
```

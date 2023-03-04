# dotlink

Dotfile symlinker

## Usage

```
dotlink [--rm] source target
```

The structure of `source` is mirrored in `target`. `source` is traversed recursively until either:

1. A leaf is reached, or
2. A directory with prefix `link-` is reached.

On 1, the leaf (file) is linked, with parent dirs created as necessary. On 2, the directory is linked, with the `link-` prefix stripped.

Filenames and dirs beginning with a literal `.` (such as `.git`, `.gitignore`) are ignored. To create a link beginning with `.`, use the prefix `dot-`. When combined with `link-`, the order is `link-dot-`.

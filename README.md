# kurz

kurz allows you to view markdown documents on the command-line in a feature-rich UI. 

## Features

- Expand/collapse sections
- Copy selected text to your clipboard
- Load remote or local files
- Cache remote files for offline access
- Automatically discover README of remote Git repositories on GitHub, BitBucket and GitLab

## Usage

There are three primary uses of `kurz`:


1. Load a local markdown file: 

```
$ kurz ./path/to/file.md
```

2. Or use a remote URL:

```
$ kurz https://example.com/markdown-file.md
```

3. Provide a Git repository to view its README:

```
$ kurz github.com/KyleBanks/depth
```

## Options

`-raw` prints the markdown document as plain text to the console.

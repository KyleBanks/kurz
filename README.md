# kurz

[![GoDoc](https://godoc.org/github.com/KyleBanks/kurz?status.svg)](https://godoc.org/github.com/KyleBanks/kurz)&nbsp; 
[![Build Status](https://travis-ci.org/KyleBanks/kurz.svg?branch=master)](https://travis-ci.org/KyleBanks/kurz)&nbsp;
[![Go Report Card](https://goreportcard.com/badge/github.com/KyleBanks/kurz)](https://goreportcard.com/report/github.com/KyleBanks/kurz)&nbsp;
[![Coverage Status](https://coveralls.io/repos/github/KyleBanks/kurz/badge.svg?branch=master)](https://coveralls.io/github/KyleBanks/kurz?branch=master)


**This project is very early and is in active development!**

`kurz` allows you to view markdown documents on the command-line in a feature-rich UI. 

!['kurz' Readme Example](./docs/screenshot.png)

## Features

- **TODO** Expand/collapse sections
- **TODO** Copy selected text to your clipboard
- **TODO** Load remote or local files
- **TODO** Cache remote files for offline access
- **TODO** Automatically discover README of remote Git repositories on GitHub, BitBucket and GitLab
- **TODO** Syntax highlighting for code snippets

## Usage

There are three primary ways to use `kurz`:

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
$ kurz github.com/KyleBanks/kurz
```

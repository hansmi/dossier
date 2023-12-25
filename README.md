# Extract information from PDF documents

[![Latest release](https://img.shields.io/github/v/release/hansmi/dossier)][releases]
[![CI workflow](https://github.com/hansmi/dossier/actions/workflows/ci.yaml/badge.svg)](https://github.com/hansmi/dossier/actions/workflows/ci.yaml)
[![Go reference](https://pkg.go.dev/badge/github.com/hansmi/dossier.svg)](https://pkg.go.dev/github.com/hansmi/dossier)

Dossier is a library for extracting textual information from PDF documents. It
is written using the Go programming language.

Currently PDF is the only supported format (using [MuPDF][mupdf]). Other
formats can be implemented using custom parsers or by amending the library.

[Sketches](#sketches) provide a declarative approach to locating information as
an alternative to imperative/procedural access.


## Sketches

[Protocol buffers][protobuf] are used to define a sketch. The [sketch protobuf
definition](proto/sketch.proto) documents available configuration options.
Usually [textproto][textproto] will be the format used for writing sketches.

A web-based viewer is included in the command line utility. Screenshot of the
viewer with an [example sketch for
invoices](/pkg/sketch/testdata/acme-invoice.textproto):

![Graphical viewer showing an example invoice analysis](/docs/viewer-acme-invoice.png)

Invocation:

```shell
$ dossiercli web ./invoice.pdf ./sketch.textproto
2023/12/31 00:00:00 HTTP server listening on http://[::1]:8080
```


## Installation

```shell
go get github.com/hansmi/dossier
```

Command line utility:

```shell
go install github.com/hansmi/dossier/cmd/dossiercli@latest
```


[releases]: https://github.com/hansmi/dossier/releases/latest
[mupdf]: https://mupdf.com/
[protobuf]: https://protobuf.dev/
[textproto]: https://protobuf.dev/reference/protobuf/textformat-spec/

<!-- vim: set sw=2 sts=2 et : -->

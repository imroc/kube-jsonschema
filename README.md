# kubeschema

## Install

```bash
go install github.com/imroc/kubeschema@latest
```

## Requirements

- Make sure `kubectl` is installed and can operate current cluster.

## Usage

Dump json schema:

```bash
kubeschema dump
```

Index json schema:

```bash
kubeschema index
```

Dump json schema and index:

```bash
kubeschema dump --index
```

Dump json schema and index with extra directory:

```bash
kubeschema dump --index --extra-dir=../other-dir
```

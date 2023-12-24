# kubeschema

## Install

```bash
go install github.com/imroc/kubeschema@latest
```

## Usage

Start `kubectl proxy`:

```bash
$ kubectl proxy
Starting to serve on 127.0.0.1:8001
```

Dump json schema:

```bash
kubeschema dump
```

Index json schema:

```bash
kubeschema index
```

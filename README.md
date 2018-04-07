# textile-go

Textile's REST API and daemon

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/textileio/textile-go)](https://goreportcard.com/report/github.com/textileio/textile-go) [![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](http://commitizen.github.io/cz-cli/) [![CircleCI](https://circleci.com/gh/textileio/textile-go/tree/master.svg?style=shield)](https://circleci.com/gh/textileio/textile-go/tree/master)

This repository contains Textile's API and daemon used to create a cross platform cli and mobile frameworks. The ```textile-go``` library is primarily used in the [Textile Photos](https://www.textile.photos) mobile application. 

Until [Textile Photos](https://www.textile.photos) is ready for public release, this library will be rapidly evolving.

## Contributing

```
go get github.com/textileio/textile-go
```

You'll need a few different tools here to get setup...

#### Install `dep`

Golang package manager:

```
brew install dep
```

#### Install `gx`

IPFS package manager:

```
go get -u github.com/whyrusleeping/gx
go get -u github.com/whyrusleeping/gx-go
```

#### Install `node`

NodeJS is used for git hooks and some build tooling:

```
brew install node
```

#### Install dependencies

Finally, download deps managed by `gx` and `dep`:

```
npm run setup
```

#### Commitizen

The easiest way to write a valid commit message is to use the `npm` script:

```
npm run cm
```

This will start the interactive commit prompt. 

## Building

Build the cli based daemon:

```
make build
```

Build the iOS Framework:

```
make ios_framework
``` 

### Desktop client

```
go get -u github.com/asticode/go-astitools
go get -u github.com/asticode/go-astilectron-bundler/...
go get -u github.com/asticode/go-astilectron-bootstrap/...
```

```
cd desktop/
astilectron-bundler -v
```

you can now open the desktop build, stored in `desktop/output/darwin-amd64` (for mac). or run it in dev mode

```
cd desktop/
go run *.go
```


## Acknowledgments

Thanks to @cpacia, @drwasho and the rest of the OpenBazaar contributors for their work on [openbazaar-go](https://github.com/OpenBazaar/openbazaar-go). 

## License

MIT

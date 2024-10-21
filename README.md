![Tiny Geth Mascot](./tinygeth.jpeg)

# Tiny Geth

[![API Reference](https://pkg.go.dev/badge/github.com/karalabe/tinygeth)](https://pkg.go.dev/github.com/karalabe/tinygeth?tab=doc) [![Build Status](https://github.com/karalabe/tinygeth/actions/workflows/tests.yaml/badge.svg)](https://github.com/karalabe/tinygeth/actions/workflows/tests.yaml) [![Code Coverage](https://codecov.io/github/karalabe/tinygeth/graph/badge.svg?token=E6GEXC2AU3)](https://codecov.io/github/karalabe/tinygeth)

Tiny Geth is Peter's personal, opinionated "fork" of go-ethereum. It's not a complete fork, rather combined bits and bobs from upstream with some local changes that were rejected upstream. I want to explore what Geth could be with hands untied, without having to care about breaking code for that one developer in Nebraska[¹](https://xkcd.com/2347/).

***Disclaimers:***

- ***This is a personal project. I take no responsibility for anyone using it and incurring arbitrary losses.***
- ***Unless promised personally, there is no API guarantee anywhere in this project, not even on the RPCs.***

## Binaries

Tiny Geth is only pre-built as a multi-arch docker image for `amd64` and `arm64` in [`tinygeth/tinygeth`](https://hub.docker.com/r/tinygeth/tinygeth).

For anything else, you can build with Go:

```sh
go install ./cmd/tinygeth
```

There are no plans to distribute to any other platforms.

## Changelog

These are the approximate changes compared upstream:

- Tiny Geth uses a number of upstream packages.
  - There is no point in maintaining the same utility functionality in multiple repositories / forks. Any Go library code that can be reused as a pre-canned chunk of code from upstream, will be done so. These include: `accounts/abi`, `common`, `crypto`, `log`, `metrics`, `p2p`, `rlp`, `rpc` at the moment.
- Tiny Geth does not support account management.
  - Ethereum nodes should not concern themselves with being wallets. Whilst having a wallet implementation in Go is useful (please use go-ethereum as a library), having a keystore inside a running node is not.

## Authorship

Tracking authorship information through git commits from upstream is not possible due to the different code and dependency structure and changes to fundamental types. I'll try to explicitly highlight commits that are merging upstream code downstream into this repo to try and explicitly highlight its origins.

Upstream code sync history:

- 2024.10.18: [v1.14.12-unstable-f32f868](https://github.com/ethereum/go-ethereum/tree/f32f8686cd35343b8f888e5607518314770f4661)

## License

This project reuses probably 99.999% code from [go-ethereum](https://github.com/ethereum/go-ethereum), some as a dependency, some verbatim. The license is the same as upstream:

- The library (i.e. code outside the cmd folder) is licensed under the GNU Lesser General Public License v3.0.
- The binaries (i.e. code inside the cmd folder) are licensed under the GNU General Public License v3.0.

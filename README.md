# Tiny Geth

Tiny Geth is Peter's personal, opinionated "fork" of go-ethereum. It's not a complete fork, rather combined bits and bobs from upstream with some local changes that were rejected upstream. I want to explore what Geth could be with hands untied, without having to care about breaking code for that one developer in Nebraska[¹](https://xkcd.com/2347/).

***Disclaimers:***

- ***This is a personal project. I take no responsibility for anyone using it and incurring arbitrary losses.***
- ***Unless promised personally, there is no API guarantee anywhere in this project, not even on the RPCs.***

## Commit history

Tracking authorship information through git commits from upstream is not possible due to the different code and dependency structure and changes to fundamental types. I'll try to explicitly highlight commits that are merging upstream code downstream into this repo to try and explicitly highlight its origins.

Upstream code sync history:

- 2024.10.18: [v1.14.12-unstable-f32f868](https://github.com/ethereum/go-ethereum/tree/f32f8686cd35343b8f888e5607518314770f4661)

## License

This project reuses probably 99.999% code from [go-ethereum](https://github.com/ethereum/go-ethereum), some as a dependency, some verbatim. The license is the same as upstream:

- The library (i.e. code outside the cmd folder) is licensed under the GNU Lesser General Public License v3.0.
- The binaries (i.e. code inside the cmd folder) are licensed under the GNU General Public License v3.0.

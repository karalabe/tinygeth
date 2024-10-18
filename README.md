# Tiny Geth

Tiny Geth is Peter's personal, opinionated "fork" of go-ethereum. It's not a complete fork, rather combined bits and bobs from upstream with some local changes that were rejected upstream.

## Disclaimers

- ***This is a personal project. I take no responsibility for anyone using it and incurring arbitrary losses due to it.***
- ***Unless promised by me personally, there is no API guarantee anywhere in this project, not even on the RPCs.***

## But why?

I want to explore what Geth could be with hands untied, without having to care about breaking some obscure code for that one developer in Nebraska[¹](https://xkcd.com/2347/).

## License

This project reuses probably 99.999% code from [go-ethereum](https://github.com/ethereum/go-ethereum), some as a dependency, some verbatim. The license is the same as upstream:

- The library (i.e. all code outside the cmd directory) is licensed under the GNU Lesser General Public License v3.0.

- The binaries (i.e. all code inside the cmd directory) are licensed under the GNU General Public License v3.0.

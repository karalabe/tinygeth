// Copyright 2024 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package console

import (
	_ "embed"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

//go:generate yarn install
//go:generate npx esbuild console.js --bundle --platform=node --minify --keep-names --outfile=bundle.js

// Bundle is the compiled together NodeJS module containing Ethers.js, a pretty
// printed JavaScript REPL interpreter and some Geth bridge code.
//
//go:embed bundle.js
var bundle string

func Run(url string) error {
	// Find the NodeJS executable
	path, err := exec.LookPath("node")
	if err != nil {
		return err
	}
	path, _ = filepath.Abs(path)

	// Dump out the console bundle to disk and remove it after
	repl, err := os.CreateTemp("", "")
	if err != nil {
		return err
	}
	defer os.Remove(repl.Name())

	if _, err = repl.Write([]byte(bundle)); err != nil {
		return err
	}
	if err = repl.Close(); err != nil {
		return err
	}
	// Start the NodeJS REPL with the console loaded
	return syscall.Exec(path, []string{path, repl.Name(), url}, syscall.Environ())
}

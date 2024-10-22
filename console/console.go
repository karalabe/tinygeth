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

	"github.com/karalabe/tinygeth/console/jsrepl"
	"github.com/karalabe/tinygeth/params"
)

// RunAsProc starts up the console and completely switches over the process to
// it, exiting when the console goes down.
func RunAsProc(url string, eval string, scripts []string) error {
	// Configure the console, cleaning up after
	nodejs, bundle, envs, err := configure(eval, scripts)
	if err != nil {
		return err
	}
	defer os.Remove(bundle)

	// Execute the console, taking over the process space
	return syscall.Exec(nodejs, []string{nodejs, bundle, url}, append(envs, syscall.Environ()...))
}

// RunInProc starts up the console inside the current process, multiplexing the
// standard outputs and consuming the standard input.
func RunInProc(url string, eval string, scripts []string) error {
	// Configure the console, cleaning up after
	nodejs, bundle, envs, err := configure(eval, scripts)
	if err != nil {
		return err
	}
	// Execute the console, keeping it in the current process space
	repl := exec.Command(nodejs, bundle, url)
	repl.Env = append(envs, os.Environ()...)
	repl.Stdin = os.Stdin
	repl.Stdout = os.Stdout
	repl.Stderr = os.Stderr
	return repl.Run()
}

// configure attempts to configure a nodejs based console with a JavaScript REPL
// bundle exported to disk, with some env vars specified.
func configure(eval string, scripts []string) (node string, bundle string, envs []string, err error) {
	// Find the NodeJS executable
	path, err := exec.LookPath("node")
	if err != nil {
		return "", "", nil, err
	}
	path, _ = filepath.Abs(path)

	// Assemble the js repl bundle script
	repl, err := os.CreateTemp("", "tinygeth-console-")
	if err != nil {
		return "", "", nil, err
	}
	defer repl.Close()

	if _, err = repl.WriteString(jsrepl.Bundle + "\n"); err != nil {
		return "", "", nil, err
	}
	if eval == "" {
		// Running in interactive REPL mode
		if _, err = repl.WriteString("module.exports.startup().then(async () => {\n"); err != nil {
			return "", "", nil, err
		}
		for _, script := range scripts {
			blob, err := os.ReadFile(script)
			if err != nil {
				return "", "", nil, err
			}
			if _, err = repl.WriteString(string(blob) + "\n"); err != nil {
				return "", "", nil, err
			}
		}
		if _, err = repl.WriteString("});"); err != nil {
			return "", "", nil, err
		}
	} else {
		// Running in single-shot evaluation mode
		if _, err = repl.WriteString("module.exports.evaluate().then(async () => {\n"); err != nil {
			return "", "", nil, err
		}
		if _, err = repl.WriteString("console.log(await " + eval + "); process.exit();\n"); err != nil {
			return "", "", nil, err
		}
		if _, err = repl.WriteString("});"); err != nil {
			return "", "", nil, err
		}
	}
	// Inject our specific APIs and ENV variables
	envs = []string{
		"TINYGETH_CONSOLE_VERSION=" + params.VersionWithMeta,
	}
	// Return everything for different startups
	return path, repl.Name(), envs, nil
}

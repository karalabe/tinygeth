// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"github.com/karalabe/tinygeth/cmd/utils"
	"github.com/karalabe/tinygeth/console"
	"github.com/karalabe/tinygeth/internal/flags"
	"github.com/urfave/cli/v2"
)

var (
	// Define the scripts to fine tune the console
	consoleEvalFlag = &cli.StringFlag{
		Name:  "eval",
		Usage: "Evaluate a JavaScript statement after starting the console",
	}
	consoleScriptFlag = &cli.StringSliceFlag{
		Name:  "scripts",
		Usage: "Evaluate a JavaScript script(s) after starting the console",
	}
	consoleFlags = []cli.Flag{consoleEvalFlag, consoleScriptFlag}

	// Define the commands to start the console
	consoleCommand = &cli.Command{
		Action: localConsole,
		Name:   "console",
		Usage:  "Start a node with an interactive JavaScript console connected",
		Flags:  flags.Merge(nodeFlags, rpcFlags, consoleFlags),
		Description: `
The Geth console is an interactive shell for the JavaScript runtime environment
which exposes an admin interface, as well as Ðapp scripting through Ethers.js.`,
	}
	attachCommand = &cli.Command{
		Action:    remoteConsole,
		Name:      "attach",
		Usage:     "Start an interactive JavaScript console connected to a remote node",
		ArgsUsage: "[endpoint]",
		Flags:     flags.Merge([]cli.Flag{utils.DataDirFlag, utils.HttpHeaderFlag}, consoleFlags),
		Description: `
The Geth console is an interactive shell for the JavaScript runtime environment
which exposes an admin interface, as well as Ðapp scripting through Ethers.js.
This command opens a console on a running, remote Ethereum node.`,
	}
)

// localConsole starts a new node, attaching a js console to it at the same time.
func localConsole(ctx *cli.Context) error {
	// Create and start the node based on the CLI flags
	prepare(ctx)
	stack := makeFullNode(ctx)
	startNode(ctx, stack, true)
	defer stack.Close()

	// Attach to the newly started node and create the console.
	eval := ctx.String(consoleEvalFlag.Name)
	incl := ctx.StringSlice(consoleScriptFlag.Name)

	return console.RunInProc(stack.IPCEndpoint(), eval, incl)
}

// remoteConsole will connect to a remote node, attaching a js console to it.
func remoteConsole(ctx *cli.Context) error {
	eval := ctx.String(consoleEvalFlag.Name)
	incl := ctx.StringSlice(consoleScriptFlag.Name)

	return console.RunAsProc(ctx.Args().First(), eval, incl)
}

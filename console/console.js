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

// Pull in a few packages needed for the console
const ethers = require('ethers');
const repl   = require('pretty-repl');
const util   = require('util');

import chalk from 'chalk';

// Connect to the requested node and inject into the repl (ethers)
function dial(url) {
    if (url.startsWith('ws://') || url.startsWith('wss://')) {
        return new ethers.WebSocketProvider(url);
    } else if (url.startsWith('http://') || url.startsWith('https://')) {
        return new ethers.JsonRpcProvider(url);
    } else if (url.startsWith('ipc://')) {
        return new ethers.IpcSocketProvider(url);
    } else {
        throw Error("unknown node type to dial");
    }
}
const api = dial(process.argv[2]);

// Collect all the custom methods we want to expose
const context = {
    ethers: ethers,
    client: api,

    eth: {
        getBlock: async (block, txs = false) => {
            return await api.getBlock(block, txs)
        }
    }
};

// Print some startup headers
const welcome = async () => {
    const client  = await api.send("web3_clientVersion");
    const block   = await api.getBlock();
    const modules = await api.send("rpc_modules");

    console.log(`Welcome to the Tiny Geth console!\n`);

    console.log(`REPL: ${chalk.green("nodejs")} ${chalk.yellow(process.version)}`);
    console.log(`Web3: ${chalk.green("ethers")} ${chalk.yellow("v" + ethers.version)}\n`);

    console.log(`Attached: ${client}`);
    console.log(`At block: ${block.number} (${new Date(1000 * Number(block.timestamp))})`);
    console.log(` Exposed: ${Object.keys(modules).join(' ')}\n`);

    console.log(chalk.grey("• The web3 library uses promises, you need to await appropriately."));
    console.log(chalk.grey("• The usual Geth API methods are exposed to the root namespace."));
    console.log(chalk.grey("• The web3 provider is exposed fully via the `ethers` field."));
    console.log(chalk.grey("• The web3 connection is exposed via the `client` field.\n"));
};

// Start the REPL server and inject all context into it
const startup = async () => {
    const server = repl.start({
        ignoreUndefined: true,
        useGlobal: true,

        prompt: chalk.green('→ '),
        writer: (output) => {
            if (output && typeof output.stack === 'string' && typeof output.message === 'string') {
                return chalk.red(output.stack || output.message);
            }
            return chalk.gray(util.inspect(output, {colors: true, depth: null}));
        }
    });
    Object.assign(server.context, context);
};

welcome().then(() => startup());

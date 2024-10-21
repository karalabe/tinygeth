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
import {ipc} from "viem/node";

const ethers = require('ethers');
const viem   = require('viem');
const repl   = require('pretty-repl');
const util   = require('util');
const web3   = require('web3');

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
const apiEthers = dial(process.argv[2]);

// Connect to the requested node and inject into the repl (viem)
function transport(url) {
    if (url.startsWith('ws://') || url.startsWith('wss://')) {
        return new viem.webSocket(url);
    } else if (url.startsWith('http://') || url.startsWith('https://')) {
        return new viem.http(url);
    } else if (url.startsWith('ipc://')) {
        return new ipc(url);
    } else {
        throw Error("unknown node type to dial");
    }
}
const apiViem = viem.createPublicClient({
    transport: transport(process.argv[2]),
    batch: {
        multicall: true,
    },
    pollingInterval: 4000,
});

const apiWeb3 = new web3.Web3(process.argv[2]);

// Collect all the custom methods we want to expose
const context = {
    ethers: ethers,
    viem:   viem,
    web3:   web3,

    web3Ethers: apiEthers,
    web3Viem:   apiViem,
    web3Web3:   apiWeb3,

    eth: {
        getBlock: async (block, txs = false) => {
            if (typeof block === 'number') {
                return await apiViem.getBlock({blockNumber: block, includeTransactions: txs});
            }
            return await apiViem.getBlock({blockHash: block, includeTransactions: txs});
        }
    }
};

// Print some startup headers
const welcome = async () => {
    const client  = await apiViem.request({method: "web3_clientVersion"});
    const block   = await apiViem.getBlock();
    const modules = await apiViem.request({method: "rpc_modules"});

    const ethersver = require('viem/package.json').version;
    const viemver = require('viem/package.json').version;
    //const web3ver = require('web3/package.json').version;

    console.log(`Welcome to the Tiny Geth console!\n`);

    console.log(`REPL: ${chalk.green("nodejs")} ${chalk.yellow(process.version)}`);
    console.log(`Web3: ${chalk.green("ethers")} ${chalk.yellow("v" + ethers.version)}`);
    console.log(`Web3: ${chalk.green("  viem")} ${chalk.yellow("v" + viemver)}\n`);
    console.log(`Web3: ${chalk.green("web3js")} ${chalk.yellow("v" + web3.version)}\n`);

    console.log(`Attached: ${client}`);
    console.log(`At block: ${block.number} (${new Date(1000 * Number(block.timestamp))})`);
    console.log(` Exposed: ${Object.keys(modules).join(' ')}\n`);

    console.log(chalk.grey("• The web3 library uses promises, you need to await appropriately."));
    console.log(chalk.grey("• The usual Geth API methods are exposed to the root namespace."));
    console.log(chalk.grey("• The web3 provider can be directly accessed via the `web3`.\n"));
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
            return util.inspect(output, {colors: true, depth: null});
        },
    });
    Object.assign(server.context, context);
};

welcome().then(() => startup());

/*
 * SPDX-License-Identifier: Apache-2.0
 */

"use strict";

const { FileSystemWallet, Gateway } = require("fabric-network");
const path = require("path");

const ccpPath = path.resolve(__dirname, "connection-BuyerBank.json");

async function query(fcn, args) {
  try {
    // Create a new file system based wallet for managing identities.
    const walletPath = path.join(process.cwd(), "wallet");
    const wallet = new FileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);

    // Check to see if we've already enrolled the user.
    const userExists = await wallet.exists("BuyerBankUser");
    if (!userExists) {
      console.log(
        'An identity for the user "BuyerBankUser" does not exist in the wallet'
      );
      console.log("Run the registerBuyerBankUser.js application before retrying");
      return;
    }

    // Create a new gateway for connecting to our peer node.
    const gateway = new Gateway();
    await gateway.connect(ccpPath, {
      wallet,
      identity: "BuyerBankUser",
      discovery: { enabled: false, asLocalhost: true }
    });

    // Get the network (channel) our contract is deployed to.
    const network = await gateway.getNetwork("tfbcchannel");

    // Get the contract from the network.
    const contract = network.getContract("tfbccc");

    // Evaluate the specified transaction.
    // queryCar transaction - requires 1 argument, ex: ('queryCar', 'CAR4')
    // queryAllCars transaction - requires no arguments, ex: ('queryAllCars')
    const result = await contract.evaluateTransaction(fcn, args);
    console.log(
      `Transaction has been evaluated, result is: ${result.toString()}`
    );
    return result;
  } catch (error) {
    console.error(`Failed to evaluate transaction: ${error}`);
    throw err;
  }
}

exports.query = query;

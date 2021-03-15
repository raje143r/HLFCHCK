/*
 * SPDX-License-Identifier: Apache-2.0
 */

"use strict";

const { json } = require("body-parser");
const { FileSystemWallet, Gateway } = require("fabric-network");
const path = require("path");

const ccpPath = path.resolve(__dirname, "connection-BuyerBank.json");

async function invoke(fcn, args) {
  try {
    // Create a new file system based wallet for managing identities.
    console.log(fcn, args);
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

    // Submit the specified transaction.
    // createCar transaction - requires 5 argument, ex: ('createCar', 'CAR12', 'Honda', 'Accord', 'Black', 'Tom')
    // changeCarOwner transaction - requires 2 args , ex: ('changeCarOwner', 'CAR10', 'Dave')

    //await contract.submitTransaction(fcn, args);
    let result = await contract.submitTransaction(fcn, args);

    console.log("Transaction has been submitted");

    // Disconnect from the gateway.
    await gateway.disconnect();
    //return JSON.parse(result);
    return result;
  } catch (error) {

    
    console.error(`Failed to submit transaction: ${error}`);
    throw error;
  }
}

exports.invoke = invoke;

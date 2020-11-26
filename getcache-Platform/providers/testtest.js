// Create a new gateway for connecting to our peer node.
//const gateway = new Gateway();
//await gateway.connect(ccp, { wallet, identity: 'xxxx' });

//const client = gateway.getClient();
'use strict';
const configProviders = require('./config.json');
let helper = require('./helper.js');
let createChannel = require('./create-channel');
let join = require('./join-channel');
let install = require('./install-chaincode');
let instantiate = require('./instantiate-chaincode');
let config = require('../services/config.json');
const fabricService = require('../services/fabricserver.service');

let util = require('util');
let hfc = require('fabric-client');
const path = require('path');
let logger = helper.getLogger('install-chaincode');

async function getClientForOrg (userorg, username) {
	logger.debug('getClientForOrg - ****** START %s %s', userorg, username)
	let config = '-connection-profile-path';

	let client = hfc.loadFromConfig(hfc.getConfigSetting('network'+config));
	client.loadFromConfig(hfc.getConfigSetting(userorg+config));
	await client.initCredentialStores();
	if(username) {
		let user = await client.getUserContext(username, true);
		if(!user) {
			throw new Error(util.format('User was not found :', username));
		} else {
			logger.debug('User %s was found to be registered and enrolled', username);
		}
	}
	logger.debug('getClientForOrg - ****** END %s %s \n\n', userorg, username)

	return client;
}

let upgradeChaincode = async function(peers, chaincodeName, chaincodePath,
	chaincodeVersion, chaincodeType, username, org_name) {
	logger.debug('\n\n============ Install chaincode on organizations ============\n',peers,chaincodeName,chaincodePath,chaincodeType,chaincodeVersion,username,org_name);
	//process.env.GOPATH = path.join(__dirname, hfc.getConfigSetting('CC_SRC_PATH'));
    //logger.debug('\n',path,'\n')
	let error_message = null;
    // first setup the client for this org
    helper.setupChaincodeDeploy();
	
    logger.debug('Successfully got the fabric client for the organization "%s"', org_name);

	try {
        let client = await getClientForOrg(org_name, username);
        let request = {
			targets: peers,
			chaincodePath: chaincodePath,
			chaincodeId: chaincodeName,
			chaincodeVersion: chaincodeVersion,
			chaincodeType: chaincodeType
		};
        let installResponse = await client.installChaincode(request);
        logger.debug(installResponse);

        let channel = await client.getChannel(channelName);

        let tx_id = client.newTransactionID(true);
        let proposalResponse = await channel.sendUpgradeProposal({
            targets: peers,
            chaincodeType: chaincodeType,
            chaincodeId: chaincodeName,
            chaincodeVersion: chaincodeVersion,
            txId: tx_id
        });

        console.log(proposalResponse);

        console.log('Sending the Transaction ..');
        const transactionResponse = channel.sendTransaction({
            proposalResponses: proposalResponse[0],
            proposal: proposalResponse[1]
        });

        console.log(transactionResponse);
		logger.info('Calling peers in organization "%s" to join the channel', org_name);

		
	} catch(error) {
		logger.error('Failed to install due to error: ' + error.stack ? error.stack : error);
		error_message = error.toString();
	}

};



let peers1 = ["peer0.users.com","peer1.users.com","peer2.users.com","peer3.users.com"];
let peers2 = ["peer0.companies.com","peer1.companies.com","peer2.companies.com","peer3.companies.com"];
let peers3 = ["peer0.exchanges.com","peer1.exchanges.com","peer2.exchanges.com","peer3.exchanges.com"];
//let peers = peers1.concat(peers2);
let chaincodeCoinBalance = "github.com/CoinBalance"
let chaincodeDataProtocol = "github.com/DataProtocol"
let chaincodeTraditionalLending = "github.com/TraditionalLending"
let chaincodePRIVIcredit = "github.com/PRIVIcredit"
let chaincodePodSwapping = "github.com/PodSwapping"
let chaincodePodNFT = "github.com/PodNFT"
let chaincodePodInsurance = "github.com/PodInsurance"
let chaincodeType= "golang"
let chaincodeVersion= "v2"
let channelName = "broadcast"
let chaincodeName = config.PRIVIcredit
let chaincodePath = chaincodePRIVIcredit
let username = "admin"
let org_name = "companies"
//process.env.GOPATH = path.join(__dirname, hfc.getConfigSetting('CC_SRC_PATH'));

//logger.debug(hfc.getConfigSetting('CC_SRC_PATH',true));
upgradeChaincode([peers2[0]], chaincodeName, chaincodePath,
	chaincodeVersion, chaincodeType, username, org_name);
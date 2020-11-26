'use strict';
let util = require('util');
let helper = require('./helper.js');
let logger = helper.getLogger('upgrade-chaincode');

let upgradeChaincode = async function(peers, channelName, chaincodeName, chaincodeVersion, functionName, chaincodeType, args, username, org_name) {
	logger.debug('\n\n============ Upgrade chaincode on channel ' + channelName +
		' ============\n');

	let error_message = null;

	try {
		let client = await helper.getClientForOrg(org_name, username);
		logger.debug('Successfully got the fabric client for the organization "%s"', org_name);
		let channel = client.getChannel(channelName);
		if(!channel) {
			let message = util.format('Channel %s was not defined in the connection profile', channelName);
			logger.error(message);
			throw new Error(message);
		}
		let tx_id = client.newTransactionID(true);
		let deployId = tx_id.getTransactionID();

		// send proposal to endorser
		let request = {
			targets : peers,
			chaincodeId: chaincodeName,
			chaincodeType: chaincodeType,
			chaincodeVersion: chaincodeVersion,
			args: args,
			txId: tx_id,
		};

			// request.fcn = functionName;
			logger.debug("INSIDE:::::::::::::::::::::::::::::::::::::::",request)

		let results = await channel.sendUpgradeProposal(request, 240000); //upgrade takes much longer

		let proposalResponses = results[0];
		let proposal = results[1];
		let all_good = true;
		for (var i in proposalResponses) {
			let one_good = false;
			if (proposalResponses && proposalResponses[i].response &&
				proposalResponses[i].response.status === 200) {
				one_good = true;
				logger.info('upgrade proposal was good');
			} else {
				logger.error('upgrade proposal was bad');
			}
			all_good = all_good & one_good;
		}

		if (all_good) {
			logger.info(util.format(
				'Successfully sent Proposal and received ProposalResponse: Status - %s, message - "%s", metadata - "%s", endorsement signature: %s',
				proposalResponses[0].response.status, proposalResponses[0].response.message,
				proposalResponses[0].response.payload, proposalResponses[0].endorsement.signature));

			let promises = [];
			let event_hubs = channel.getChannelEventHubsForOrg();
			logger.debug('found %s eventhubs for this organization %s',event_hubs.length, org_name);
			event_hubs.forEach((eh) => {
				let upgradeEventPromise = new Promise((resolve, reject) => {
					logger.debug('upgradeEventPromise - setting up event');
					let event_timeout = setTimeout(() => {
						let message = 'REQUEST_TIMEOUT:' + eh.getPeerAddr();
						logger.error(message);
						eh.disconnect();
					}, 120000);
					eh.registerTxEvent(deployId, (tx, code, block_num) => {
						logger.info('The chaincode upgrade transaction has been committed on peer %s',eh.getPeerAddr());
						logger.info('Transaction %s has status of %s in blocl %s', tx, code, block_num);
						clearTimeout(event_timeout);

						if (code !== 'VALID') {
							let message = util.format('The chaincode upgrade transaction was invalid, code:%s',code);
							logger.error(message);
							reject(new Error(message));
						} else {
							let message = 'The chaincode upgrade transaction was valid.';
							logger.info(message);
							resolve(message);
						}
					}, (err) => {
						clearTimeout(event_timeout);
						logger.error(err);
						reject(err);
					},
						{unregister: true, disconnect: true}
					);
					eh.connect();
				});
				promises.push(upgradeEventPromise);
			});

			let orderer_request = {
				txId: tx_id,
				proposalResponses: proposalResponses,
				proposal: proposal
			};
			let sendPromise = channel.sendTransaction(orderer_request);
			promises.push(sendPromise);
			let results = await Promise.all(promises);
			logger.debug(util.format('------->>> R E S P O N S E : %j', results));
			let response = results.pop(); //  orderer results are last in the results
			if (response.status === 'SUCCESS') {
				logger.info('Successfully sent transaction to the orderer.');
			} else {
				error_message = util.format('Failed to order the transaction. Error code: %s',response.status);
				logger.debug(error_message);
			}

			// now see what each of the event hubs reported
			for(let i in results) {
				let event_hub_result = results[i];
				let event_hub = event_hubs[i];
				logger.debug('Event results for event hub :%s',event_hub.getPeerAddr());
				if(typeof event_hub_result === 'string') {
					logger.debug(event_hub_result);
				} else {
					if(!error_message) error_message = event_hub_result.toString();
					logger.debug(event_hub_result.toString());
				}
			}
		} else {
			error_message = util.format('Failed to send Proposal and receive all good ProposalResponse');
			logger.debug(error_message);
		}
	} catch (error) {
		logger.error('Failed to send upgrade due to error: ' + error.stack ? error.stack : error);
		error_message = error.toString();
	}

	if (!error_message) {
		let message = util.format(
			'Successfully upgrade chaincode in organization %s to the channel \'%s\'',
			org_name, channelName);
		logger.info(message);
		
		let response = {
			success: true,
			message: message
		};
		return response;
	} else {
		let message = util.format('Failed to upgrade. cause:%s',error_message);
		logger.error(message);
		throw new Error(message);
	}
};
exports.upgradeChaincode = upgradeChaincode;

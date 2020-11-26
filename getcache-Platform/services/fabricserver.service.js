const hfc = require('fabric-client');
const util = require('util');
const log4js = require('log4js');
const logger = log4js.getLogger('service call');
const helper = require('../providers/helper');
const config= require('./config.json')
let services = {};
services.registerWithCertificates = registerWithCertificates;
services.GetInfo=GetInfo;
services.Register=Register;
services.Invoke=Invoke
services.richquery=richquery;
module.exports=services


/*---------------------------------------------------------------------------------------
FUNCTION TO REGISTER CACHE PARTICIPANTS WITH CERTIFICATES IN BLOCKCHAIN
---------------------------------------------------------------------------------------*/

async function registerWithCertificates( req ) {

    try {
		// Get user public Id key and organisation type //
        let publicIdKey = req.body.publicIdKey;
		let orgName = req.body.orgName;
        if (!publicIdKey) {
            return {error: getErrorMessage('\'publicIdKey\'') };
        } if (!orgName) {
            return {error:getErrorMessage('\'orgName\'')};
        }
	
		// Register user with certificates //
        let response = await helper.getRegisteredUser( publicIdKey, orgName, true );
        if (response && typeof response !== 'string') {
			logger.debug( 'PUBLIC ID KEY %s SUCCESFULLY REGISTERED ' +
			              'FOR ORGANISATION %s', publicIdKey, orgName );
            if (response.success==false) {
				logger.debug( 'PUBLIC ID KEY %s ALREADY REGISTERED ' +
			              'FOR ORGANISATION %s', publicIdKey, orgName );
                return {error: publicId +  " ALREADY REGISTERED" }
            }
			return response
		// Failed to register user //
        } else {
			logger.debug( 'FAILED TO REGISTER PUBLIC KEY ID %s FOR ORGANISATION '+
			              ' %s WITH ERROR:\n %s', publicIdKey, orgName, response );
            return ( {success: false, message: response} );
        }
    } catch(error){
        return {error:error}
    }
}

async function GetInfo(req,res){
    try{
		let orgName,peers
		let username=req.body.username
        if(req.body.chaincodeName){
			let chaincodeName = req.body.chaincodeName
		}else{
			let chaincodeName = config.chaincodeName
		}
        if(req.body.type=="user") {
            orgName="users"
            peers=config.UserPeers
        } else if(req.body.type="company") {
            orgName="companies"
            peers=config.CompanyPeers
        } else {
			return res.json(getErrorMessage("enter correct type"))
		}
		let client = await helper.getClientForOrg(orgName, username);
		let channel = client.getChannel(config.channelName);
		if(!channel) {
			let message = util.format('Channel %s was not defined in the connection profile', channelName);
			logger.error(message);
			return getErrorMessage(message);
		}
		let request = {
			targets : peers,
			chaincodeId: chaincodeName,
			fcn: req.body.fcn,
			args: req.body.args
		};
		console.log("request",request)
		let response_payloads = await channel.queryByChaincode(request);
		if (response_payloads) {
			
			let response=response_payloads[0]
			console.log(response,"#############")
			logger.debug(("THIS IS RESPONSE",response))
            if (response.status){
				logger.debug("THIS IS RESPONSE",response)
                return getErrorMessage(response_payloads[0].toString('utf8'))
            } else {
				response=getSuccessMessage(response_payloads[0].toString('utf8'))
				logger.debug("THIS IS RESPONSE",response)
                return response
            }
		} else {
			logger.error('response_payloads is null');
			return getErrorMessage("null value recieved")
		}
    } catch(error){
        console.log("error",error)
        return getErrorMessage(error)
    }
}

async function richquery(req, res) {
	try {
		let creq = {
			body: {
				type: req.body.type,
				username: req.body.TransactionObj.Id,
				fcn: "richquery",
				args: [JSON.stringify(req.body.TransactionObj.query)]
			}
		 }
		let result = await GetInfo(creq)

		if (result.message.length < 0) {
			return getErrorMessage("Result Not Found")
		}
		result.message.forEach(function(msg){
			if (msg.Record.TxnObj){delete msg.Record.TxnObj;}
			console.log(msg.Record)
			})
		return result
	} catch (error) {
		console.log("error", error)
		return getErrorMessage(error)
	}
}

async function Register(req,res){
	try{
       let response = await Invoke(req,res)
       return response
    } catch(error){
        console.log("error",error)
        return getErrorMessage(error)
    }
}

function getErrorMessage(field) {
	let response = {
		success: false,
		message: field 
	};
	return response;
}

function getSuccessMessage(field) {
	let response = {
		success: true,
		message: JSON.parse(field )
	};
	return response;
}


/*---------------------------------------------------------------------------------------
FUNCTION TO INVOKE FUNCTIONS OF SMART CONTRACTS AND HANDLE RESPONSES
---------------------------------------------------------------------------------------*/

async function Invoke( req ) {
	

	// Variables containing the info to invoke the smart contract functions //
	// let peerNames, username, org_name
	let channelName = config.channelName
	let chaincodeName = req.body.chaincodeName
	let username = req.body.username

	let error_message = null;
	let tx_id_string = null;
	let payload = null;

	logger.debug(
		util.format(
			 '\n======  INVOKE TRANSACTION ON CHANNEL %s======\n', 
	         channelName ));
	
	try {
		// Determine peers to call //
        if(req.body.type=="user") {
            org_name="users"
            peerNames=config.UserPeers
        } else if(req.body.type=="business") {
            org_name="companies"
            peerNames=config.CompanyPeers
        }

		// Determine channel for the invokation //
		let client = await helper.getClientForOrg(org_name, username);
		let channel = client.getChannel(channelName);
		if (!channel) {
			let message = util.format( 
				'CHANNEL %s WAS NOT DEFINED IN THE CONNECTION',
			     channelName );
			logger.error(message);
			return getErrorMessage(message);
		}

		// Generate proposal for the invokation //
		let tx_id = client.newTransactionID();
		tx_id_string = tx_id.getTransactionID();
		let request = {
			targets: peerNames,
			chaincodeId: chaincodeName,
			fcn: req.body.fcn,
			args: req.body.args,
			chainId: channelName,
			txId: tx_id
		};
        console.log("################", request)
		let results = await channel.sendTransactionProposal(request);
		
		// Process responses of the smart contract invokation //
		let proposalResponses = results[0];
		let proposal = results[1];
		let all_good = true;
		for (var i in proposalResponses) {
			let one_good = false;
			if (proposalResponses && proposalResponses[i].response &&
				proposalResponses[i].response.status === 200) {
				one_good = true;
				logger.info( 'INVOKE CHAINCODE PROPOSAL WAS GOOD' );
				payload=results[0][0].response.payload
        		console.log( results[0], "#############" )
			} else {
                logger.error( 'INVOKE CHAINCODE PROPOSAL WAS BAD' );
                logger.debug( "RESULT::::::::::::::::::::::::::::",
				              proposalResponses[i]);      
                return getErrorMessage(
					          proposalResponses[i].toString('utf8'))          
			}
			all_good = all_good & one_good;
		}

		// If successful invokation register blockchain responses //
		if (all_good) {
			logger.info(util.format(
				'SUCCESSFULLY SENT PROPOSAL AND RECEIVED PROPOSAL REPONSE\n' +
				'STATUS - %s\n   MESSAGE: %s\n  METADATA: %s\n' +
				'ENDOSERMENT SIGNATURE: %s\n',
				proposalResponses[0].response.status, 
				proposalResponses[0].response.message,
				proposalResponses[0].response.payload, 
				proposalResponses[0].endorsement.signature));

			let promises = [];
			let event_hubs = channel.getChannelEventHubsForOrg();
			event_hubs.forEach((eh) => {
				logger.debug('INVOKEEVENTPROMISE - SETTING UP EVENT');
				let invokeEventPromise = new Promise(  (resolve, reject) => {
					// Time out //
					let event_timeout = setTimeout(() => {
						let message = 'REQUEST_TIMEOUT:' + eh.getPeerAddr();
						logger.error(message);
						eh.disconnect();
					}, 240000);
					// Register Transaction //
					eh.registerTxEvent( tx_id_string, (tx, code, block_num) => {
						logger.info( 'THE CHAINCODE INVOKATION TRANSACTION' +           
						             'COMMITTED ON PEER %s' ,eh.getPeerAddr() );
						logger.info( 'TRANSACTION %s HAS STATUS OF %s ' +
									 'IN BLOCK %s', tx, code, block_num );
						clearTimeout( event_timeout );
						// Invalid smartcode invokation //
						if (code !== 'VALID') {
							let message = util.format( 'INVOKE CHAINCODE ' +
							                'TRANSACTION INVALED, CODE: %s', code);
							logger.error(message);
							reject(new Error(message));
					    // Valid smartcode invokation //
						} else {
							let message = 'THE INVOKE CHAINCODE TRANSACTION '+
							              'WAS VALID'
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
				promises.push(invokeEventPromise);
			});

			// Send succesful transactions to Orderers //
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
				error_message = util.format('Failed to order the transaction. Error code: %s',
				                            response.status);
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
		logger.error('Failed to invoke due to error: ' + error.stack ? error.stack : error);
		error_message = error.toString();
	}

	// Return the results of the invokation to the API //
	if (!error_message) {
		let message = util.format(
			'Successfully invoked the chaincode %s to the channel \'%s\' for transaction ID: %s',
			org_name, channelName, tx_id_string);
		logger.info(message);

        return {success:true, message: tx_id_string, payload: payload.toString("utf-8") };
	} else {
		let message = util.format('Failed to invoke chaincode. cause:%s',error_message);
		logger.error(message);
		return getErrorMessage(message)
	}
};
const hfc = require('fabric-client');
const util = require('util');
const log4js = require('log4js');
const logger = log4js.getLogger('service call');
const helper = require('../providers/helper');
const config= require('./config.json')
let services = {};
services.eventlisten = eventlisten;
module.exports=services
let msg = ""

async function eventlisten() {
	console.log("Starting EVENT HUB");
	client = await helper.getClientForOrg("exchanges", "admin");
	channel = client.getChannel("broadcast");
	let chaincodeName=config.ICOName;
	let event_hubs = channel.getChannelEventHubsForOrg();
	//event_hubs.forEach((eh) => {
	event_hubs.forEach((eh) => {
		//let eh=event_hubs[0]
		eh.connect(true);
		let regid = eh.registerChaincodeEvent(chaincodeName, 'Transfer', (event, block_num) => {
		console.log("*****************DONE************")
		let event_payload = event.payload.toString('utf8');
		if (event_payload != null&&event_payload!=msg){
			msg=event_payload
			let temp=JSON.parse(event_payload)
			console.log(temp)
			
		} else {
			logger.info('already sent')
		}

		}, (error) => {
			console.log('Failed to receive the chaincode event ::' + error);
		},
		{
			unregister: false,
			disconnect: false
		}

	);
	})
}
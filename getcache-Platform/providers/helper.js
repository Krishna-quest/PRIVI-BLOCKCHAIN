
'use strict';
let log4js = require('log4js');
let logger = log4js.getLogger('Helper');
logger.setLevel('DEBUG');

let path = require('path');
let util = require('util');
const cryptoRandomString = require('crypto-random-string');
let hfc = require('fabric-client');
hfc.setLogger(logger);

// will be registering keys here for each org user
// Will be called at both register user&company

/*---------------------------------------------------------------------------------------
---------------------------------------------------------------------------------------*/

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

/*---------------------------------------------------------------------------------------
---------------------------------------------------------------------------------------*/

let getRegisteredUser = async function (username, userOrg, isJson) {
    try {
        // Initialise organisation to register new public Id Key //
        let client = await getClientForOrg( userOrg );
        logger.debug( 'SUCCESFULLY INITIALIZED THE CREDENTIAL STORES');
        let user = await client.getUserContext(username, true);

        // Public Id Key already registered on blockchain //
        if (user && user.isEnrolled()) {
            logger.info('PUBLIC ID KEY ALREADY REGISTERED ON BLOCKCHAIN');
            let response = { success: false, };
            return response;
        // Public Id Key not register on blockchain. Need to do it //
        } else {
            logger.info('PUBLIC ID KEY %s NOT ENROLLED, REGISTERING NOW.', username);
            let admins = hfc.getConfigSetting('admins');
            let adminUserObj = await client.setUserContext({ username: admins[0].username, 
                                                             password: admins[0].secret });

            if (adminUserObj.getAffiliation() != userOrg.toLowerCase()) {
                logger.info('Admin affiliation not registered. Registering now.');
                adminUserObj.setAffiliation(userOrg.toLowerCase());
                adminUserObj.setRoles(['peer', 'orderer', 'client', 'user']);
                adminUserObj = await client.setUserContext(adminUserObj);
            }

            logger.info('Admin User: % s', adminUserObj);
            let affiliation = userOrg.toLowerCase() + '.department1';

            let caClient = client.getCertificateAuthority();

            // CHECK IF THE ORGANISATION EXISTS //
            const affiliationService = caClient.newAffiliationService();
            const registeredAffiliations = await affiliationService.getAll(adminUserObj);
            if (!registeredAffiliations.result.affiliations.some(x => x.name == userOrg.toLowerCase())) {
                logger.info('Register the new affiliation: % s ', affiliation);
                await affiliationService.create({ name: affiliation, force: true }, adminUserObj);
            }
            let seed = cryptoRandomString({length: 10});
            let secret = await caClient.register({
                enrollmentID: username,
                affiliation: affiliation,
                attrs: [{ name: "seed", value: seed, ecert: true}],
            }, adminUserObj);
            logger.debug('SUCESSFULLY GOT THE SECRET FOR PUBLIC ID KEY % s — % s', username);
            user = await client.setUserContext({ 
                username: username, 
                password: secret,
                attr_reqs:[{ name: "seed", optional: false }]});
            logger.debug('Successfully enrolled username % s and setUserContext on the client object', username);

        }
        if (user && user.isEnrolled) {
            if (isJson && isJson === true) {
                let response = {
                    success: true,
                    secret: seed,
                    message: username + ' ENROLLED SUCCESSFULLY',
                };
                return response;
            }
        } else {
            throw new Error('USER WAS NOT ENROLLED');
        }
    } catch (error) {
        let message = 'FAILED TO REGISTER USER ' + username + ' FOR ORGANISATION ' + userOrg + 
                      ' WITH ERROR: ' + error.toString();
        logger.error(message);
        let response = {
            success: false,
            message: message,
        };
        return response;
    }
};

/*---------------------------------------------------------------------------------------
---------------------------------------------------------------------------------------*/

let setupChaincodeDeploy = function() {
	process.env.GOPATH = path.join(__dirname, hfc.getConfigSetting('CC_SRC_PATH'));
};

/*---------------------------------------------------------------------------------------
---------------------------------------------------------------------------------------*/

let getLogger = function(moduleName) {
	let logger = log4js.getLogger(moduleName);
	logger.setLevel('DEBUG');
	return logger;
};

/*---------------------------------------------------------------------------------------
---------------------------------------------------------------------------------------*/

exports.getClientForOrg = getClientForOrg;
exports.getLogger = getLogger;
exports.setupChaincodeDeploy = setupChaincodeDeploy;
exports.getRegisteredUser = getRegisteredUser;
// exports.getRegisteredAdmin = getRegisteredAdmin;
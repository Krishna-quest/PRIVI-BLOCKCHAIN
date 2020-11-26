const fabricService = require('../services/fabricserver.service');
const express = require('express');
let router = express.Router();
let config = require('../services/config.json')
let _ = require('lodash')

router.post('/register', register)
router.post('/getUser', getUser)
router.post('/getPrivacy', getPrivacy)
router.post('/modifyPrivacy', modifyPrivacy)
router.post('/getBusinessList', getBusinessList)
router.post('/getRoleList', getRoleList)
router.post('/getUserList', getUserList)
router.post('/getTIDList', getTIDList)
router.post('/encryptData', encryptData)
router.post('/decryptData', decryptData)
router.post('/insightDiscovery', insightDiscovery)
router.post('/insightPurchase', insightPurchase)
router.post('/insightTarget', insightTarget)

/*--------------------------------------------------------
--------------------------------------------------------*/

const PRIVI_PUBLIC_ID = "k3Xpi5IB61fvG3xNM4POkjnCQnx1"
const PRIVI_PRIVATE_ID = "ERSASDFAW5IB61fvG3xNM4POkjnCQnx1ALSNFsf901ASI2139abo329741"

/*--------------------------------------------------------
--------------------------------------------------------*/

async function register( req, res ) {
    try {
        let input={body: {}}
        input.body.username= PRIVI_PRIVATE_ID
        input.body.fcn = "register"
        input.body.chaincodeName = config.DataProtocol
        input.body.type = "business"
        input.body.args = [ JSON.stringify(req.body) ]

        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        res.json( response);
    } catch(error) {
        console.log("error", error)
        res.json( {"error": error.message} )
    }
    return false
}


/*--------------------------------------------------------
--------------------------------------------------------*/

async function getUser( req, res ) {
    try {
        let input = {body: {}}
        input.body.username= PRIVI_PRIVATE_ID
        input.body.fcn = "getUser"
        input.body.chaincodeName = config.DataProtocol
        input.body.type = "business"
        input.body.args = [ req.body.PublicId ]

        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload)
        response.output =  response.output
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("error",error)
        res.json( {"error": error.message} )
    }
    return false
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function getPrivacy( req, res ) {
    try {
        let input={body: {}}
        input.body.username= PRIVI_PRIVATE_ID
        input.body.fcn = "getPrivacy"
        input.body.chaincodeName = config.DataProtocol
        input.body.type = "business"
        input.body.args = [ req.body.PublicId ]

        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload)
        response.output =  response.output.Privacy
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("error", error)
        res.json( {"error": error.message} )
    }
    return false
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function modifyPrivacy( req, res ) {
    try {
        let input={body:{}}
        input.body.username= PRIVI_PRIVATE_ID
        input.body.fcn = "modifyPrivacy"
        input.body.chaincodeName = config.DataProtocol
        input.body.type = "business"
        input.body.args = [ JSON.stringify(req.body) ]

        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload)
        response.output = response.output.Privacy
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("error", error)
        res.json( {"error": error.message} )
    }
    return false
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function getBusinessList( req, res ) {
    try {
        let input={body:{}}
        input.body.username= PRIVI_PRIVATE_ID
        input.body.fcn = "getBusinessList"
        input.body.chaincodeName = config.DataProtocol
        input.body.type = "business"
        input.body.args = []

        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload)
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("error", error)
        res.json( {"error": error.message} )
    }
    return false
}


/*--------------------------------------------------------
--------------------------------------------------------*/

async function getRoleList( req, res ) {
    try {
        let input={body:{}}
        input.body.username= PRIVI_PRIVATE_ID
        input.body.fcn = "getRoleList"
        input.body.chaincodeName = config.DataProtocol
        input.body.type = "business"
        input.body.args = [ req.body.Role ]

        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload)
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("error", error)
        res.json( {"error": error.message} )
    }
    return false
}


/*--------------------------------------------------------
--------------------------------------------------------*/

async function getUserList( req, res ) {
    try {
        let input={body: {}}
        input.body.username= PRIVI_PRIVATE_ID
        input.body.fcn = "getUserList"
        input.body.chaincodeName = config.DataProtocol
        input.body.type = "business"
        input.body.args = []

        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload)
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error){
        console.log("error", error)
        res.json( {"error": error.message} )
    }
    return false
}


/*--------------------------------------------------------
--------------------------------------------------------*/

async function getTIDList( req, res ) {
    try {
        let input={body:{}}
        input.body.username= PRIVI_PRIVATE_ID
        input.body.fcn = "getTIDList"
        input.body.chaincodeName = config.DataProtocol
        input.body.type = "business"
        input.body.args = [ req.body.Id ]

        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload)
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("error", error)
        res.json( {"error": error.message} )
    }
    return false
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function encryptData( req, res ) {
    try {
        let input={body: {}}
        input.body.username= PRIVI_PRIVATE_ID
        input.body.fcn = "encryptData"
        input.body.chaincodeName = config.DataProtocol
        input.body.type = "business"
        input.body.args = [ req.body.PublicId ]

        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = response.payload
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("error", error)
        res.json( {"error": error.message} )
    }
    return false
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function decryptData( req, res ) {
    try {
        let input={body: {}}
        input.body.username= PRIVI_PRIVATE_ID
        input.body.fcn = "decryptData"
        input.body.chaincodeName = config.DataProtocol
        input.body.type = "business"
        input.body.args = [ req.body.Encryption_DID ]

        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = response.payload
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("error", error)
        res.json( {"error": error.message} )
    }
    return false
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function insightDiscovery( req, res ) {
    try{
        let input={body: {}}
        input.body.username= PRIVI_PRIVATE_ID
        input.body.fcn = "insightDiscovery"
        input.body.chaincodeName = config.DataProtocol
        input.body.type = "business"
        input.body.args = [ JSON.stringify(req.body) ]

        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload)
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("error", error)
        res.json( {"error": error.message} )
    }
    return false
}

/*--------------------------------------------------------
--------------------------------------------------------*/


async function insightPurchase( req, res ) {
    try {
        let input={body: {}}
        input.body.username= PRIVI_PRIVATE_ID
        input.body.fcn = "insightPurchase"
        input.body.chaincodeName = config.DataProtocol
        input.body.type = "business"
        input.body.args = [ JSON.stringify(req.body) ]

        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload)
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("error",error)
        res.json( {"error": error.message} )
    }
    return false
}


/*--------------------------------------------------------
--------------------------------------------------------*/


async function insightTarget( req, res ) {
    try {
        let input={body: {}}
        input.body.username= PRIVI_PRIVATE_ID
        input.body.fcn = "insightTarget"
        input.body.chaincodeName = config.DataProtocol
        input.body.type = "business"
        input.body.args = [ JSON.stringify(req.body) ]

        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload)
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("error", error)
        res.json( {"error": error.message} )
    }
    return false
}

/*--------------------------------------------------------
--------------------------------------------------------*/

module.exports = router;
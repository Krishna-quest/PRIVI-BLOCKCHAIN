// CONSTANTS //
const fabricService = require('../services/fabricserver.service');
const express = require('express');
let router = express.Router();
let config = require('../services/config.json')

// ROUTES //

router.post('/getParameters', getParameters)
router.post('/updateParameters', updateParameters)
router.post('/retrieveInsuranceList', retrieveInsuranceList)
router.post('/retrieveInsuranceInfo', retrieveInsuranceInfo)
router.post('/initiateInsurancePool', initiateInsurancePool)
router.post('/investInsurancePool', investInsurancePool)
router.post('/subscribeInsurancePool', subscribeInsurancePool)
router.post('/withdrawInsurancePool', withdrawInsurancePool)

/*--------------------------------------------------------
--------------------------------------------------------*/

const PRIVI_PUBLIC_ID = "k3Xpi5IB61fvG3xNM4POkjnCQnx1"
const PRIVI_PRIVATE_ID = "ERSASDFAW5IB61fvG3xNM4POkjnCQnx1ALSNFsf901ASI2139abo329741"

/*--------------------------------------------------------
--------------------------------------------------------*/

async function getParameters( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "getParameters"
        input.body.type = "business"
        input.body.chaincodeName = config.PodInsurance
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = []
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING getParameters: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function updateParameters( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "updateParameters"
        input.body.type = "business"
        input.body.chaincodeName = config.PodInsurance
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body) ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING updateParameters: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function retrieveInsuranceList( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "retrieveInsuranceList"
        input.body.type = "business"
        input.body.chaincodeName = config.PodInsurance
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = []
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING retrieveInsuranceList: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function retrieveInsuranceInfo( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "retrieveInsuranceInfo"
        input.body.type = "business"
        input.body.chaincodeName = config.PodInsurance
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ req.body.Id ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING retrieveInsuranceInfo: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function initiateInsurancePool( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "initiateInsurancePool"
        input.body.type = "business"
        input.body.chaincodeName = config.PodInsurance
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body) ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING initiateInsurancePool: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function investInsurancePool( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "investInsurancePool"
        input.body.type = "business"
        input.body.chaincodeName = config.PodInsurance
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body) ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING investInsurancePool: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/


async function subscribeInsurancePool( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "subscribeInsurancePool"
        input.body.type = "business"
        input.body.chaincodeName = config.PodInsurance
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body) ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING subscribeInsurancePool: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/


async function withdrawInsurancePool( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "withdrawInsurancePool"
        input.body.type = "business"
        input.body.chaincodeName = config.PodInsurance
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body) ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING withdrawInsurancePool: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

module.exports= router;
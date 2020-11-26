// CONSTANTS //
const fabricService = require('../services/fabricserver.service');
const express = require('express');
let router = express.Router();
let config = require('../services/config.json')

// ROUTES //
router.post('/updateRiskParameters', updateRiskParameters)
router.post('/initiatePRIVIcredit', initiatePRIVIcredit)
router.post('/getPRIVIcreditList', getPRIVIcreditList)
router.post('/getPRIVIcredit', getPRIVIcredit)
router.post('/getRiskParameters', getRiskParameters)
router.post('/modifyPRIVIparameters', modifyPRIVIparameters)
router.post('/withdrawFunds', withdrawFunds)
router.post('/depositFunds', depositFunds)
router.post('/borrowFunds', borrowFunds)
router.post('/assumePRIVIrisk', assumePRIVIrisk)
router.post('/managePRIVIcredits', managePRIVIcredits)

/*--------------------------------------------------------
--------------------------------------------------------*/

const PRIVI_PUBLIC_ID = "k3Xpi5IB61fvG3xNM4POkjnCQnx1"
const PRIVI_PRIVATE_ID = "ERSASDFAW5IB61fvG3xNM4POkjnCQnx1ALSNFsf901ASI2139abo329741"

/*--------------------------------------------------------
--------------------------------------------------------*/

async function updateRiskParameters( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "updateRiskParameters"
        input.body.type = "business"
        input.body.chaincodeName = config.PRIVIcredit
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ req.body.Token,
                            JSON.stringify(req.body.RiskParameters) ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = ""
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING updateRiskParameters: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function initiatePRIVIcredit( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "initiatePRIVIcredit"
        input.body.type = "business"
        input.body.chaincodeName = config.PRIVIcredit
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ req.body.InitialDeposit,
                            JSON.stringify(req.body.Loan_conditions) ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING initiatePRIVIcredit: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function getPRIVIcreditList( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "getPRIVIcreditList"
        input.body.type = "business"
        input.body.chaincodeName = config.PRIVIcredit
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
        console.log("ERROR CALLING getPRIVIcreditList: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function getPRIVIcredit( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "getPRIVIcredit"
        input.body.type = "business"
        input.body.chaincodeName = config.PRIVIcredit
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ req.body.LoanId ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING getPRIVIcredit: ", error)
        res.json( {"error": error.message} )
    }
}


/*--------------------------------------------------------
--------------------------------------------------------*/

async function getRiskParameters( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "getRiskParameters"
        input.body.type = "business"
        input.body.chaincodeName = config.PRIVIcredit
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ req.body.Token ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING getRiskParameters: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function modifyPRIVIparameters( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "modifyPRIVIparameters"
        input.body.type = "business"
        input.body.chaincodeName = config.PRIVIcredit
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
        console.log("ERROR CALLING modifyPRIVIcredit: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function withdrawFunds( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "withdrawFunds"
        input.body.type = "business"
        input.body.chaincodeName = config.PRIVIcredit
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
        console.log("ERROR CALLING withdrawFunds: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function depositFunds( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "depositFunds"
        input.body.type = "business"
        input.body.chaincodeName = config.PRIVIcredit
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
        console.log("ERROR CALLING depositFunds: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function borrowFunds( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "borrowFunds"
        input.body.type = "business"
        input.body.chaincodeName = config.PRIVIcredit
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
        console.log("ERROR CALLING borrowFunds: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function assumePRIVIrisk( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "assumePRIVIrisk"
        input.body.type = "business"
        input.body.chaincodeName = config.PRIVIcredit
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
        console.log("ERROR CALLING assumePRIVIrisk: ", error)
        res.json( {"error": error.message} )
    }
}


/*--------------------------------------------------------
--------------------------------------------------------*/

async function managePRIVIcredits( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "managePRIVIcredits"
        input.body.type = "business"
        input.body.chaincodeName = config.PRIVIcredit
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = []
        let response = await fabricService.Invoke(input);
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        res.json(response);
        /*response.output = []
        payload.forEach((element, i) => {
            response.output.push(JSON.parse(element));
            if(payload.length === i + 1) {
                delete response.message
                delete response.payload
                console.log("BLOCKCHAIN RESPONSE: ", response)
                res.json(response);
            }
        });*/
    } catch(error) {
        console.log("ERROR CALLING managePRIVIcredits: ", error)
        res.json( {"error": error.message} )
    }
}


/*--------------------------------------------------------
--------------------------------------------------------*/

module.exports= router;

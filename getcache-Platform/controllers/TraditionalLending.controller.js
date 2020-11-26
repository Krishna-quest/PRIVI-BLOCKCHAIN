// CONSTANTS //
const fabricService = require('../services/fabricserver.service');
const express = require('express');
let router = express.Router();
let config = require('../services/config.json')

// ROUTES //
router.post('/registerToken', registerToken)
router.post('/getTokenList', getTokenList)
router.post('/getLendingPools', getLendingPools)
router.post('/getTokenPool', getTokenPool)
router.post('/getDemandRatios', getDemandRatios)
router.post('/borrowFunds', borrowFunds)
router.post('/stakeToken', stakeToken)
router.post('/unStakeToken', unStakeToken)
router.post('/payInterests', payInterests)
router.post('/depositCollateral', depositCollateral)
router.post('/withdrawCollateral', withdrawCollateral)
router.post('/repayFunds', repayFunds)
router.post('/checkLiquidation', checkLiquidation)
router.post('/updateRiskParameters', updateRiskParameters)

/*--------------------------------------------------------
--------------------------------------------------------*/

const PRIVI_PUBLIC_ID = "k3Xpi5IB61fvG3xNM4POkjnCQnx1"
const PRIVI_PRIVATE_ID = "ERSASDFAW5IB61fvG3xNM4POkjnCQnx1ALSNFsf901ASI2139abo329741"

/*--------------------------------------------------------
--------------------------------------------------------*/

async function registerToken( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "registerToken"
        input.body.type = "business"
        input.body.chaincodeName = config.TraditionalLending
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body) ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = ""
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING registerToken: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function getTokenList( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "getTokenList"
        input.body.type = "business"
        input.body.chaincodeName = config.TraditionalLending
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = []
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload

        console.log("BLOCKCHAIN RESPONSE: ", response)
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING getTokenList: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function getDemandRatios( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "getDemandRatios"
        input.body.type = "business"
        input.body.chaincodeName = config.TraditionalLending
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = []
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload

        console.log("BLOCKCHAIN RESPONSE: ", response)
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING getDemandRatios: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function getLendingPools( req, res ) {
    try {
        let input = { body:{} }
        input.body.fcn = "getLendingPools"
        input.body.type = "business"
        input.body.chaincodeName = config.TraditionalLending
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
        console.log("ERROR CALLING getLendingPools: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function getTokenPool( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "getTokenPool"
        input.body.type = "business"
        input.body.chaincodeName = config.TraditionalLending
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ req.body.Token ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload

        console.log("BLOCKCHAIN RESPONSE: ", response)
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING getTokenPool: ", error)
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
        input.body.chaincodeName = config.TraditionalLending
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body) ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        console.log("BLOCKCHAIN RESPONSE: ", response)
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING borrowFunds: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function stakeToken( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "stakeToken"
        input.body.type = "business"
        input.body.chaincodeName = config.TraditionalLending
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body) ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        console.log("BLOCKCHAIN RESPONSE: ", response)
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING stakeToken: ", error)
        res.json( {"error": error.message} )
    }
}


/*--------------------------------------------------------
--------------------------------------------------------*/

async function unStakeToken( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "unStakeToken"
        input.body.type = "business"
        input.body.chaincodeName = config.TraditionalLending
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body) ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        console.log("BLOCKCHAIN RESPONSE: ", response)
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING unStakeToken: ", error)
        res.json( {"error": error.message} )
    }
}



/*--------------------------------------------------------
--------------------------------------------------------*/

async function payInterests( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "payInterests"
        input.body.type = "business"
        input.body.chaincodeName = config.TraditionalLending
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body) ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        console.log("BLOCKCHAIN RESPONSE: ", response)
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING payInterests: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function depositCollateral( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "depositCollateral"
        input.body.type = "business"
        input.body.chaincodeName = config.TraditionalLending
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body) ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        console.log("BLOCKCHAIN RESPONSE: ", response)
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING depositCollateral: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function withdrawCollateral( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "withdrawCollateral"
        input.body.type = "business"
        input.body.chaincodeName = config.TraditionalLending
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body) ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        console.log("BLOCKCHAIN RESPONSE: ", response)
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING withdrawCollateral: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function repayFunds( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "repayFunds"
        input.body.type = "business"
        input.body.chaincodeName = config.TraditionalLending
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body) ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        console.log("BLOCKCHAIN RESPONSE: ", response)
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING repayFunds: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function checkLiquidation( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "checkLiquidation"
        input.body.type = "business"
        input.body.chaincodeName = config.TraditionalLending
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body) ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        console.log("BLOCKCHAIN RESPONSE: ", response)
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING checkLiquidation: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function updateRiskParameters( req, res) {
    try {
        let jsonReq = req.body;
        let input = { body: {} }
        input.body.fcn = "updateRiskParameters"
        input.body.type = "business"
        input.body.chaincodeName = config.TraditionalLending
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ jsonReq.Token, 
                            JSON.stringify(jsonReq.RiskParameters) ]
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = "";
        delete response.message
        delete response.payload
        console.log("BLOCKCHAIN RESPONSE: ", response)
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING updateRiskParameters: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

module.exports= router;

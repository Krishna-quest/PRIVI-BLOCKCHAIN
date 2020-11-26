// CONSTANTS //
const fabricService = require('../services/fabricserver.service');
const express = require('express');
let router = express.Router();
let config = require('../services/config.json')

// ROUTES //
router.post('/updateRiskParameters', updateRiskParameters)
router.post('/retrievePoolList', retrievePoolList)
router.post('/retrievePodList', retrievePodList)
router.post('/retrievePodInfo', retrievePodInfo)
router.post('/retrievePoolInfo', retrievePoolInfo)
router.post('/createLiquidityPool', createLiquidityPool)
router.post('/depositLiquidity', depositLiquidity)
router.post('/withdrawLiquidity', withdrawLiquidity)
router.post('/managerLiquidity', managerLiquidity)
router.post('/initiatePOD', initiatePOD)
router.post('/deletePOD', deletePOD)
router.post('/investPOD', investPOD)
router.post('/interestPOD', interestPOD)
router.post('/swapPOD', swapPOD)
router.post('/managerPOD', managerPOD)
router.post('/liquidatePOD', liquidatePOD)
router.post('/checkPODLiquidation', checkPODLiquidation)

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
        input.body.chaincodeName = config.PodSwapping
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ req.body.Token,
                            JSON.stringify(req.body.RiskParameters) ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = ""
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

async function createLiquidityPool( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "createLiquidityPool"
        input.body.type = "business"
        input.body.chaincodeName = config.PodSwapping
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body)]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING createLiquidityPool: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function depositLiquidity( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "depositLiquidity"
        input.body.type = "business"
        input.body.chaincodeName = config.PodSwapping
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body)]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING depositLiquidity: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function withdrawLiquidity( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "withdrawLiquidity"
        input.body.type = "business"
        input.body.chaincodeName = config.PodSwapping
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
        console.log("ERROR CALLING withdrawLiquidity: ", error)
        res.json( {"error": error.message} )
    }
}


/*--------------------------------------------------------
--------------------------------------------------------*/

async function managerLiquidity( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "managerLiquidity"
        input.body.type = "business"
        input.body.chaincodeName = config.PodSwapping
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body) ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) { return res.json(response) }
        let payload = JSON.parse(response.payload)
        response.output = []
        res.json(response);
        payload.forEach((element, i) => {
            response.output.push(JSON.parse(element));
            if(payload.length === i + 1) {
                delete response.message
                delete response.payload
                console.log("BLOCKCHAIN RESPONSE: ", response)
                res.json(response);
            }
        });
    } catch(error) {
        console.log("ERROR CALLING managerLiquidity: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function initiatePOD( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "initiatePOD"
        input.body.type = "business"
        input.body.chaincodeName = config.PodSwapping
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body.PodInfo),
                            JSON.stringify(req.body.Collaterals),
                            JSON.stringify(req.body.RateChange) ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING initiatePOD: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function retrievePoolList( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "retrievePoolList"
        input.body.type = "business"
        input.body.chaincodeName = config.PodSwapping
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
        console.log("ERROR CALLING retrievePoolList: ", error)
        res.json( {"error": error.message} )
    }
}


/*--------------------------------------------------------
--------------------------------------------------------*/

async function retrievePodList( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "retrievePodList"
        input.body.type = "business"
        input.body.chaincodeName = config.PodSwapping
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
        console.log("ERROR CALLING retrievePodList: ", error)
        res.json( {"error": error.message} )
    }
}


/*--------------------------------------------------------
--------------------------------------------------------*/


async function retrievePodInfo( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "retrievePodInfo"
        input.body.type = "business"
        input.body.chaincodeName = config.PodSwapping
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ req.body.PodId ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING retrievePodInfo: ", error)
        res.json( {"error": error.message} )
    }
}


/*--------------------------------------------------------
--------------------------------------------------------*/


async function retrievePoolInfo( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "retrievePoolInfo"
        input.body.type = "business"
        input.body.chaincodeName = config.PodSwapping
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ req.body.PoolId ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING retrievePoolInfo: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function deletePOD( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "deletePOD"
        input.body.type = "business"
        input.body.chaincodeName = config.PodSwapping
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
        console.log("ERROR CALLING deletePOD: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function investPOD( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "investPOD"
        input.body.type = "business"
        input.body.chaincodeName = config.PodSwapping
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
        console.log("ERROR CALLING investPOD: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function interestPOD( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "interestPOD"
        input.body.type = "business"
        input.body.chaincodeName = config.PodSwapping
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
        console.log("ERROR CALLING interestPOD: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function swapPOD( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "swapPOD"
        input.body.type = "business"
        input.body.chaincodeName = config.PodSwapping
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
        console.log("ERROR CALLING swapPOD: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function managerPOD( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "managerPOD"
        input.body.type = "business"
        input.body.chaincodeName = config.PodSwapping
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
        console.log("ERROR CALLING managerPOD: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function liquidatePOD( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "liquidatePOD"
        input.body.type = "business"
        input.body.chaincodeName = config.PodSwapping
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body),
                            req.body.Type ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING liquidatePOD: ", error)
        res.json( {"error": error.message} )
    }
}


/*--------------------------------------------------------
--------------------------------------------------------*/

async function checkPODLiquidation( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "checkPODLiquidation"
        input.body.type = "business"
        input.body.chaincodeName = config.PodSwapping
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
        console.log("ERROR CALLING checkPODLiquidation: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

module.exports= router;
// CONSTANTS //
const fabricService = require('../services/fabricserver.service');
const express = require('express');
let router = express.Router();
let config = require('../services/config.json')

// ROUTES //
router.post('/retrievePodList', retrievePodList)
router.post('/retrievePodInfo', retrievePodInfo)
router.post('/retrieveInsuranceList', retrieveInsuranceList)
router.post('/retrieveInsuranceInfo', retrieveInsuranceInfo)
router.post('/initiatePodNFT', initiatePodNFT)
router.post('/getOrderBook', getOrderBook)
router.post('/newBuyOrder', newBuyOrder)
router.post('/deleteBuyOrder', deleteBuyOrder)
router.post('/newSellOrder', newSellOrder)
router.post('/deleteSellOrder', deleteSellOrder)
router.post('/buyPodNFT', buyPodNFT)
router.post('/sellPodNFT', sellPodNFT)
router.post('/initiateClaimProposal', initiateClaimProposal)

/*--------------------------------------------------------
--------------------------------------------------------*/

const PRIVI_PUBLIC_ID = "k3Xpi5IB61fvG3xNM4POkjnCQnx1"
const PRIVI_PRIVATE_ID = "ERSASDFAW5IB61fvG3xNM4POkjnCQnx1ALSNFsf901ASI2139abo329741"

/*--------------------------------------------------------
--------------------------------------------------------*/

async function initiatePodNFT( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "initiatePodNFT"
        input.body.type = "business"
        input.body.chaincodeName = config.PodNFT
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body.PodInfo),
                            JSON.stringify(req.body.Offers) ]
        
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

async function retrievePodList( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "retrievePodList"
        input.body.type = "business"
        input.body.chaincodeName = config.PodNFT
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
        input.body.chaincodeName = config.PodNFT
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

async function retrieveInsuranceList( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "retrieveInsuranceList"
        input.body.type = "business"
        input.body.chaincodeName = config.PodNFT
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
        input.body.chaincodeName = config.PodNFT
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

async function getOrderBook( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "getOrderBook"
        input.body.type = "business"
        input.body.chaincodeName = config.PodNFT
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
        console.log("ERROR CALLING getOrderBook: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function newBuyOrder( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "newBuyOrder"
        input.body.type = "business"
        input.body.chaincodeName = config.PodNFT
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
        console.log("ERROR CALLING newBuyOrder: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function deleteBuyOrder( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "deleteBuyOrder"
        input.body.type = "business"
        input.body.chaincodeName = config.PodNFT
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
        console.log("ERROR CALLING deleteBuyOrder: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function newSellOrder( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "newSellOrder"
        input.body.type = "business"
        input.body.chaincodeName = config.PodNFT
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
        console.log("ERROR CALLING newSellOrder: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function deleteSellOrder( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "deleteSellOrder"
        input.body.type = "business"
        input.body.chaincodeName = config.PodNFT
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
        console.log("ERROR CALLING deleteSellOrder: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function buyPodNFT( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "buyPodNFT"
        input.body.type = "business"
        input.body.chaincodeName = config.PodNFT
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
        console.log("ERROR CALLING buyPodNFT: ", error)
        res.json( {"error": error.message} )
    }
}


/*--------------------------------------------------------
--------------------------------------------------------*/

async function sellPodNFT( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "sellPodNFT"
        input.body.type = "business"
        input.body.chaincodeName = config.PodNFT
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
        console.log("ERROR CALLING sellPodNFT: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function initiateClaimProposal( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "initiateClaimProposal"
        input.body.type = "business"
        input.body.chaincodeName = config.PodNFT
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
        console.log("ERROR CALLING initiateClaimProposal: ", error)
        res.json( {"error": error.message} )
    }
}

/*--------------------------------------------------------
--------------------------------------------------------*/

module.exports= router;
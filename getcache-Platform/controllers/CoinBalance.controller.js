// CONSTANTS //
const fabricService = require('../services/fabricserver.service');
const express = require('express');
let router = express.Router();
let config = require('../services/config.json')

// ROUTES //
router.post('/registerToken', registerToken)
router.post('/getTokenList', getTokenList)
router.post('/getUserList', getUserList)
router.post('/registerWallet', registerWallet)
router.post('/balanceOf', balanceOf)
router.post('/transfer', transfer)
router.post('/multitransfer', multitransfer)
router.post('/spendFunds', spendFunds)
router.post('/getWallets', getWallets)
router.post('/updateMultiwallet', updateMultiwallet)
router.post('/getHistory', getHistory)
router.post('/mint', mint)
router.post('/burn', burn)
router.post('/swap', swap)
router.post('/withdraw', withdraw)

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
        input.body.chaincodeName = config.CoinBalance
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body),
                            PRIVI_PUBLIC_ID  ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
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
        input.body.chaincodeName = config.CoinBalance
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
        console.log("ERROR CALLING getTokenList: ", error)
        res.json( {"error": error.message} )
    }
}


/*--------------------------------------------------------
--------------------------------------------------------*/

async function getUserList( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "getUserList"
        input.body.type = "business"
        input.body.chaincodeName = config.CoinBalance
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
        console.log("ERROR CALLING getUserList: ", error)
        res.json( {"error": error.message} )
    }
    return false
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function registerWallet( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "registerWallet"
        input.body.type = "business"
        input.body.chaincodeName = config.CoinBalance
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ req.body.PublicId ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING registerBalance: ", error)
        res.json( {"error": error.message} )
    }
    return false
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function balanceOf( req, res ) {
    try {
        let input = { body:{} }
        input.body.fcn = "balanceOf"
        input.body.type = "business"
        input.body.chaincodeName = config.CoinBalance
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ req.body.PublicId ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload).Balances;
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING balanceOf: ",error)
        res.json( {"error": error.message} )
    }
    return false
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function transfer( req, res ) {
    try {
        let input = { body:{} }
        input.body.fcn = "transfer"
        input.body.type = "business"
        input.body.chaincodeName = config.CoinBalance
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
        console.log("ERROR CALLING transfer: ",error)
        res.json( {"error": error.message} )
    }
    return false
}

/*--------------------------------------------------------
--------------------------------------------------------*/


async function multitransfer(req,res){
    try{
        let arr=[]
        for(i=0;i<req.body.Multitransfer.length;i++){
            arr.push(JSON.stringify(req.body.Multitransfer[i]))
        }
        let input = { body:{} }
        input.body.fcn = "multitransfer"
        input.body.type = "business"
        input.body.chaincodeName = config.CoinBalance
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = arr
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING multitransfer: ",error)
        res.json( {"error": error.message} )
    }
    return false
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function spendFunds(req,res){
    try{
        let input = { body:{} }
        input.body.fcn = "spendFunds"
        input.body.type = "business"
        input.body.chaincodeName = config.CoinBalance
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body) ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = JSON.parse(response.payload);
        delete response.message
        delete response.payload
        res.json(response);
        /*let payload = JSON.parse(response.payload)
        response.output = []
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
        console.log("ERROR CALLING spendFunds: ",error)
        res.json( {"error": error.message} )
    }
    return false
}


/*--------------------------------------------------------
--------------------------------------------------------*/


async function getWallets(req,res){
    try{
        let input = { body:{} }
        input.body.fcn = "getWallets"
        input.body.type = "business"
        input.body.chaincodeName = config.CoinBalance
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
        console.log("ERROR CALLING updateMultiwallet: ",error)
        res.json( {"error": error.message} )
    }
    return false
}

/*--------------------------------------------------------
--------------------------------------------------------*/


async function updateMultiwallet(req, res){
    try{
        let input = { body:{} }
        input.body.fcn = "updateMultiwallet"
        input.body.type = "business"
        input.body.chaincodeName = config.CoinBalance
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body) ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = "";
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING updateMultiwallet: ",error)
        res.json( {"error": error.message} )
    }
    return false
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function getHistory(req, res) {
    try {
        let input = { body:{} }
        input.body.fcn = "getHistory"
        input.body.type = "business"
        input.body.chaincodeName = config.CoinBalance
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body) ]

        
        let response = await fabricService.Invoke(input)
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = []
        iterator_history = JSON.parse(response.payload)

        for (let i=0; i < iterator_history.length; i++) {
            iterator_txn = iterator_history[i].Value
            if (iterator_txn === null) continue;
            for (let j=0; j < iterator_txn.length; j++) {
                trans = iterator_txn[j]
                //trans.Id = iterator_history[i].TxId + String(j)
                trans.Timestamp = iterator_history[i].Timestamp
                response.output.push(trans)
            }
        }
        delete response.message
        delete response.payload
        res.json(response);
    } catch (error) {
        console.log("error", error)
        res.json( {"error": error.message} )
    }
    return false
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function mint( req, res ) {
    try {
        let input = { body: {} }
        input.body.fcn = "mint"
        input.body.type = "business"
        input.body.chaincodeName = config.CoinBalance
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body), 
                            PRIVI_PUBLIC_ID ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = ""
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING mint: ",error)
        res.json( {"error": error.message} )
    }
    return false
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function burn( req, res ) {
    try {
        let input = { body:{} }
        input.body.fcn = "burn"
        input.body.type = "business"
        input.body.chaincodeName = config.CoinBalance
        input.body.username = PRIVI_PRIVATE_ID
        input.body.args = [ JSON.stringify(req.body), 
                            PRIVI_PUBLIC_ID ]
        
        let response = await fabricService.Invoke(input);
        console.log("BLOCKCHAIN RESPONSE: ", response)
        if (response.success == false) {return res.json(response) }
        response.output = ""
        delete response.message
        delete response.payload
        res.json(response);
    } catch(error) {
        console.log("ERROR CALLING burn: ",error)
        res.json( {"error": error.message} )
    }
    return false
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function swap( req, res ) {
    try {
        let input = { body:{} }
        input.body.fcn = "swap"
        input.body.type = "business"
        input.body.chaincodeName = config.CoinBalance
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
        console.log("ERROR CALLING swap: ",error)
        res.json( {"error": error.message} )
    }
    return false
}

/*--------------------------------------------------------
--------------------------------------------------------*/

async function withdraw( req, res ) {
    try {
        let input = { body:{} }
        input.body.fcn = "withdraw"
        input.body.type = "business"
        input.body.chaincodeName = config.CoinBalance
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
        console.log("ERROR CALLING withdraw: ", error)
        res.json( {"error": error.message} )
    }
    return false
}

/*--------------------------------------------------------
--------------------------------------------------------*/

module.exports = router;
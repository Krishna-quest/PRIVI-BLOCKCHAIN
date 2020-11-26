/*
* Automating the flow to start the app
* Will be called as soon as app is started to
* Create Users
* Create Channel
* Join Channel
* Install Chaincode
* Instantiate Chaincode
*/

/* -------------------------------------------------------------------------------------------------
CONSTANTS AND CONFIG PATHS
------------------------------------------------------------------------------------------------- */
const PRIVI_PUBLIC_ID = "k3Xpi5IB61fvG3xNM4POkjnCQnx1"
const PRIVI_PRIVATE_ID = "ERSASDFAW5IB61fvG3xNM4POkjnCQnx1ALSNFsf901ASI2139abo329741"

const configProviders = require('./config.json');
let helper = require('./helper.js');
let logger = helper.getLogger('StartCacheKit');
let createChannel = require('./create-channel');
let join = require('./join-channel');
let install = require('./install-chaincode');
let instantiate = require('./instantiate-chaincode');
let upgrade = require('./upgrade-chaincode');
let config = require('../services/config.json');
const fabricService = require('../services/fabricserver.service');

/* -------------------------------------------------------------------------------------------------
INVOKE DEPLOYMENT OF CACHE BLOCKCHAIN
------------------------------------------------------------------------------------------------- */

(async function invoke() {
    let username="admin"
    let orgName=["users","companies","exchanges"]
    let response, message
    let channelName = "broadcast"
    let channelConfigPath="../artifacts/channel/channel.tx"

    // Register Cache Admin to Blockchain //
    
    /*response1 = await helper.getRegisteredUser( username, orgName[0], true);
    response2= await helper.getRegisteredUser( PRIVI_PRIVATE_ID, orgName[1], true);
    response3 = await helper.getRegisteredUser( username, orgName[2], true);
    console.log( response1, response2, response3)
    logger.debug('-- returned from registering the username %s for organization %s',username,orgName);
    logger = helper.getLogger('Create-Channel');
    // //create channel
    message = await createChannel.createChannel(channelName, channelConfigPath, username, orgName[1]);
    logger.debug("message",message)
    await new Promise(done => setTimeout(done, 10000));*/

    //join channel

    let peers1 = ["peer0.users.com","peer1.users.com","peer2.users.com","peer3.users.com"];
    let peers2 = ["peer0.companies.com","peer1.companies.com","peer2.companies.com","peer3.companies.com"];
    let peers3 = ["peer0.exchanges.com","peer1.exchanges.com","peer2.exchanges.com","peer3.exchanges.com"];
    let peers = peers1.concat(peers2);


    let chaincodeCoinBalance = "github.com/CoinBalance"
    let chaincodeDataProtocol = "github.com/DataProtocol"
    let chaincodeTraditionalLending = "github.com/TraditionalLending"
    let chaincodePRIVIcredit = "github.com/PRIVIcredit"
    let chaincodePodSwapping = "github.com/PodSwapping"
    let chaincodePodNFT = "github.com/PodNFT"
    let chaincodePodInsurance = "github.com/PodInsurance"
    let chaincodeType= "golang"
    logger.debug('CHANNEL NAME : ' + channelName);
    logger.debug('PEERS : ' + peers);


    /* -------------------------------------------------------------------------------------------------
    JOIN PEERS TO CHANNEL
    ------------------------------------------------------------------------------------------------- */
    
    /*logger.info('<<<<<<<<<<<<<<<<< J O I N  C H A N N E L >>>>>>>>>>>>>>>>>');
    message =  await join.joinChannel( channelName, peers1, username, orgName[0] );
    logger.debug( "MESSAGE 1: ", message )
    message =  await join.joinChannel( channelName, peers2, username, orgName[1] );
    logger.debug( "MESSAGE 2: ", message )
    message =  await join.joinChannel( channelName, peers3, username, orgName[2] );
    logger.debug( "MESSAGE 3: ", message )
    await new Promise(done => setTimeout(done, 20000));*/

    /* -------------------------------------------------------------------------------------------------
    INSTANTIATE COIN BALANCE SMART CONTRACT
    ------------------------------------------------------------------------------------------------- */
    /*
    message = await install.installChaincode( peers1, config.CoinBalance, chaincodeCoinBalance, 
                                              config.CoinBalanceVersion, chaincodeType, username, orgName[0] )
    logger.debug("INSTALL: ",message)
    message = await install.installChaincode( peers2, config.CoinBalance, chaincodeCoinBalance, 
                                              config.CoinBalanceVersion, chaincodeType, username, orgName[1] )
    logger.debug("INSTALL: ",message)
    message = await install.installChaincode( peers3, config.CoinBalance, chaincodeCoinBalance, 
                                              config.CoinBalanceVersion, chaincodeType, username, orgName[2])
    logger.debug("INSTALL: ",message)
    await new Promise(done => setTimeout(done, 10000));*/
    /*
    message = await instantiate.instantiateChaincode( [peers2[0]], channelName, config.CoinBalance,
                     config.CoinBalanceVersion, "Init", chaincodeType, [ PRIVI_PRIVATE_ID, "INSTANTIATE" ],  
                     username, orgName[1] );
    logger.debug("INSTANTIATE COIN BALANCE CHAINCODE: SUCCESFUL.", message)*/
    /*
    message = await upgrade.upgradeChaincode( [peers2[0]], channelName, config.CoinBalance,
        config.CoinBalanceVersion, "Invoke", chaincodeType, [ PRIVI_PRIVATE_ID, "UPGRADE" ],  
        username, orgName[1] );
    logger.debug("UPGRADE PRIVI CHAINCODE: SUCCESFUL.", message)*/

    /* -------------------------------------------------------------------------------------------------
    INSTANTIATE DATA PROTOCOL SMART CONTRACT
    ------------------------------------------------------------------------------------------------- */
    /*
    message = await install.installChaincode( peers1, config.DataProtocol, chaincodeDataProtocol, 
                            config.DataProtocolVersion, chaincodeType, username, orgName[0] )
    logger.debug("INSTALL: ", message)
    message = await install.installChaincode(peers2, config.DataProtocol, chaincodeDataProtocol, 
                            config.DataProtocolVersion, chaincodeType, username, orgName[1] )
    logger.debug("INSTALL: ", message)
    message = await install.installChaincode(peers3, config.DataProtocol, chaincodeDataProtocol, 
                            config.DataProtocolVersion, chaincodeType, username, orgName[2] )
    logger.debug("INSTALL: ", message)
    await new Promise(done => setTimeout(done, 10000));*/
    /*
    message = await instantiate.instantiateChaincode([peers2[0]], channelName, config.DataProtocol, 
                      config.DataProtocolVersion, "Init", chaincodeType, [ PRIVI_PRIVATE_ID, "INSTANTIATE" ],  
                      username, orgName[1]);
    logger.debug("INSTANTIATE DATA PROTOCOL CHAINCODE: SUCCESFUL.", message)*/
    /*
    message = await upgrade.upgradeChaincode( [peers2[0]], channelName, config.DataProtocol,
        config.DataProtocolVersion, "Init", chaincodeType, [ PRIVI_PRIVATE_ID, "UPGRADE" ],  
        username, orgName[1] );
    logger.debug("UPGRADE DATA PROTOCOL CHAINCODE: SUCCESFUL.", message)*/

    /* -------------------------------------------------------------------------------------------------
    INSTANTIATE TRADITIONAL LENDING SMART CONTRACT
    ------------------------------------------------------------------------------------------------- */
    /*
    message = await install.installChaincode( peers1, config.TraditionalLending, chaincodeTraditionalLending, 
                            config.TraditionalLendingVersion, chaincodeType, username, orgName[0] )
    logger.debug("INSTALL: ",message)
    message = await install.installChaincode( peers2, config.TraditionalLending, chaincodeTraditionalLending, 
                             config.TraditionalLendingVersion, chaincodeType, username, orgName[1] )
    logger.debug("INSTALL: ",message)
    message = await install.installChaincode( peers3, config.TraditionalLending, chaincodeTraditionalLending, 
                             config.TraditionalLendingVersion, chaincodeType, username, orgName[2] ) 
    logger.debug("INSTALL: ",message)
    await new Promise(done => setTimeout(done, 10000));
    
    message = await instantiate.instantiateChaincode( [peers2[0]], channelName, config.TraditionalLending,
                         config.TraditionalLendingVersion, "Init", chaincodeType, [ PRIVI_PRIVATE_ID, "INSTATIATE" ],  
                        username, orgName[1] );
    logger.debug("INSTANTIATE TRADITIONAL LENDING CHAINCODE: SUCCESFUL.", message)*/
    /*
    message = await upgrade.upgradeChaincode( [peers2[0]], channelName, config.TraditionalLending,
        config.TraditionalLendingVersion, "Init", chaincodeType, [ PRIVI_PRIVATE_ID, "UPGRADE" ],  
        username, orgName[1] );
    logger.debug("UPGRADE TRADITIONAL LENDING CHAINCODE: SUCCESFUL.", message)*/

    /* -------------------------------------------------------------------------------------------------
    INSTANTIATE PRIVI CREDIT SMART CONTRACT
    ------------------------------------------------------------------------------------------------- */
    /*
    message = await install.installChaincode( peers1, config.PRIVIcredit, chaincodePRIVIcredit, 
            config.PRIVIcreditVersion, chaincodeType, username, orgName[0] )
    logger.debug("INSTALL: ",message)
    message = await install.installChaincode( peers2, config.PRIVIcredit, chaincodePRIVIcredit, 
            config.PRIVIcreditVersion, chaincodeType, username, orgName[1] )
    logger.debug("INSTALL: ",message)
    message = await install.installChaincode( peers3, config.PRIVIcredit, chaincodePRIVIcredit, 
            config.PRIVIcreditVersion, chaincodeType, username, orgName[2] ) 
    logger.debug("INSTALL: ",message)
    await new Promise(done => setTimeout(done, 10000));*/
    /*
    message = await instantiate.instantiateChaincode( [peers2[0]], channelName, config.PRIVIcredit,
        config.PRIVIcreditVersion, "Init", chaincodeType, [ PRIVI_PRIVATE_ID, "INSTANTIATE" ],  
        username, orgName[1] );
    logger.debug("INSTANTIATE PRIVI CREDIT CHAINCODE: SUCCESFUL.", message)*/
    /*message = await upgrade.upgradeChaincode( [peers2[0]], channelName, config.PRIVIcredit,
        config.PRIVIcreditVersion, "Init", chaincodeType, [ PRIVI_PRIVATE_ID, "UPGRADE" ],  
        username, orgName[1] );
    logger.debug("UPGRADE PRIVI CREDIT CHAINCODE: SUCCESFUL.", message)*/
    
    /* -------------------------------------------------------------------------------------------------
    INSTANTIATE POD SWAPPING SMART CONTRACT
    ------------------------------------------------------------------------------------------------- */
    /*
    message = await install.installChaincode( peers1, config.PodSwapping, chaincodePodSwapping, 
            config.PodSwappingVersion, chaincodeType, username, orgName[0] )
    logger.debug("INSTALL: ",message)
    message = await install.installChaincode( peers2, config.PodSwapping, chaincodePodSwapping, 
            config.PodSwappingVersion, chaincodeType, username, orgName[1] )
    logger.debug("INSTALL: ",message)
    message = await install.installChaincode( peers3, config.PodSwapping, chaincodePodSwapping, 
            config.PodSwappingVersion, chaincodeType, username, orgName[2] ) 
    logger.debug("INSTALL: ",message)
    await new Promise(done => setTimeout(done, 10000));*/
    /*
    message = await instantiate.instantiateChaincode( [peers2[0]], channelName, config.PodSwapping,
        config.PodSwappingVersion, "Init", chaincodeType, [ PRIVI_PRIVATE_ID, "INSTANTIATE" ],  
        username, orgName[1] );
    logger.debug("INSTANTIATE POD SWAPPING CHAINCODE: SUCCESFUL.", message)*/
    /*message = await upgrade.upgradeChaincode( [peers2[0]], channelName, config.PodSwapping,
        config.PodSwappingVersion, "Init", chaincodeType, [ PRIVI_PRIVATE_ID, "UPGRADE" ],  
        username, orgName[1] );
    logger.debug("UPGRADE POD SWAPPING CHAINCODE: SUCCESFUL.", message)*/

    /* -------------------------------------------------------------------------------------------------
    INSTANTIATE POD NFT SMART CONTRACT
    ------------------------------------------------------------------------------------------------- */
    
    message = await install.installChaincode( peers1, config.PodNFT, chaincodePodNFT, 
            config.PodNFTVersion, chaincodeType, username, orgName[0] )
    logger.debug("INSTALL: ",message)
    message = await install.installChaincode( peers2, config.PodNFT, chaincodePodNFT, 
            config.PodNFTVersion, chaincodeType, username, orgName[1] )
    logger.debug("INSTALL: ",message)
    message = await install.installChaincode( peers3, config.PodNFT, chaincodePodNFT, 
            config.PodNFTVersion, chaincodeType, username, orgName[2] ) 
    logger.debug("INSTALL: ",message)
    await new Promise(done => setTimeout(done, 10000));
    /*
    message = await instantiate.instantiateChaincode( [peers2[0]], channelName, config.PodNFT,
        config.PodNFTVersion, "Init", chaincodeType, [ PRIVI_PRIVATE_ID, "INSTANTIATE" ],  
        username, orgName[1] );
    logger.debug("INSTANTIATE POD NFT CHAINCODE: SUCCESFUL.", message)*/
    message = await upgrade.upgradeChaincode( [peers2[0]], channelName, config.PodNFT,
        config.PodNFTVersion, "Init", chaincodeType, [ PRIVI_PRIVATE_ID, "UPGRADE" ],  
        username, orgName[1] );
    logger.debug("UPGRADE POD NFT CHAINCODE: SUCCESFUL.", message)

    /* -------------------------------------------------------------------------------------------------
    INSTANTIATE POD INSURANCE SMART CONTRACT
    ------------------------------------------------------------------------------------------------- */
    /*
    message = await install.installChaincode( peers1, config.PodInsurance, chaincodePodInsurance, 
            config.PodInsuranceVersion, chaincodeType, username, orgName[0] )
    logger.debug("INSTALL: ",message)
    message = await install.installChaincode( peers2, config.PodInsurance, chaincodePodInsurance, 
             config.PodInsuranceVersion, chaincodeType, username, orgName[1] )
    logger.debug("INSTALL: ",message)
    message = await install.installChaincode( peers3, config.PodInsurance, chaincodePodInsurance, 
             config.PodInsuranceVersion, chaincodeType, username, orgName[2] ) 
    logger.debug("INSTALL: ",message)
    await new Promise(done => setTimeout(done, 10000));*/
    /*
    message = await instantiate.instantiateChaincode( [peers2[0]], channelName, config.PodInsurance,
         config.PodInsuranceVersion, "Init", chaincodeType, [ PRIVI_PRIVATE_ID, "INSTANTIATE"  ],  
        username, orgName[1] );
    logger.debug("INSTANTIATE POD INSURANCE CHAINCODE: SUCCESFUL.", message);*/

    /*
    message = await upgrade.upgradeChaincode( [peers2[0]], channelName, config.PodInsurance,
         config.PodInsuranceVersion, "Init", chaincodeType, [ PRIVI_PRIVATE_ID,  "UPGRADE"  ],  
        username, orgName[1] );
    logger.debug("INSTANTIATE POD INSURANCE CHAINCODE: SUCCESFUL.", message);*/
    

    /* -------------------------------------------------------------------------------------------------
    REGISTER PRIVI AS ADMIN OF THE NETWORK AND PRIVI COIN AND BASE COIN
    ------------------------------------------------------------------------------------------------- */
    
    // Register PRIVI with public ID (this instantiate also a wallet) //
    /*PRIVI_register = { "PublicId": PRIVI_PUBLIC_ID, "Role": "ADMIN" }
    let input1={body:{}}
    input1.body.username= PRIVI_PRIVATE_ID
    input1.body.fcn = "register"
    input1.body.chaincodeName = config.DataProtocol
    input1.body.type = "business"
    input1.body.args = [ JSON.stringify(PRIVI_register) ]
    message = await fabricService.Invoke(input1);
    logger.debug("PRIVI REGISTRATION SUCCESFUL.", message)*/

    /* -------------------------------------------------------------------------------------------------
    REGISTER TOKENS ON THE SYSTEM
    ------------------------------------------------------------------------------------------------- */
    let PRIVI_COIN = { "Name": "PRIVI Coin", "Symbol": "PC", "Supply": configProviders.INITIAL_SUPPLY_PC }
    let BASE_COIN = { "Name": "Base Coin", "Symbol": "BC", "Supply": configProviders.INITIAL_SUPPLY_BC }
    let PRIVI_DATA_TOKEN = { "Name": "PRIVI Data Token", "Symbol": "PDT", "Supply": configProviders.INITIAL_SUPPLY_PDT }
    let BAL_TOKEN = { "Name": "Balancer", "Symbol": "BAL", "Supply": configProviders.INITIAL_SUPPLY_BAL }
    let BAT_TOKEN = { "Name": "Basic Attention Token", "Symbol": "BAT", "Supply": configProviders.INITIAL_SUPPLY_BAT }
    let COMP_TOKEN = { "Name": "Compound", "Symbol": "COMP", "Supply": configProviders.INITIAL_SUPPLY_COMP }
    let DAI_COIN = { "Name": "Dai Stablecoin", "Symbol": "DAI", "Supply": configProviders.INITIAL_SUPPLY_DAI }
    let ETH_COIN = { "Name": "Ethereum", "Symbol": "ETH", "Supply": configProviders.INITIAL_SUPPLY_ETH }
    let LINK_TOKEN = { "Name": "Chainlink", "Symbol": "LINK", "Supply": configProviders.INITIAL_SUPPLY_LINK }
    let MKR_COIN = { "Name": "MakerDAO", "Symbol": "MKR", "Supply": configProviders.INITIAL_SUPPLY_MKR }
    let UNI_TOKEN = { "Name": "Uniswap", "Symbol": "UNI", "Supply": configProviders.INITIAL_SUPPLY_UNI }
    let USDT_COIN = { "Name": "Tether", "Symbol": "USDT", "Supply": configProviders.INITIAL_SUPPLY_USDT }
    let WBTC_TOKEN = { "Name": "Wrapped Bitcoin", "Symbol": "WBTC", "Supply": configProviders.INITIAL_SUPPLY_WBTC }
    let YFI_TOKEN = { "Name": "Yearn Finance", "Symbol": "YFI", "Supply": configProviders.INITIAL_SUPPLY_YFI }

    let PRIVI_COIN_TL = { "Token": "PC", "Reserve": configProviders.INITIAL_RESERVE_PC, "RiskParameters": { "Col_withdrawal": configProviders.RISK_COL_WITHDRAWAL_PC, "Col_required": configProviders.RISK_COL_REQUIRED_PC, "Col_liquidation": configProviders.RISK_COL_LIQUIDATION_PC, "P_limit": configProviders.RISK_P_LIMIT_PC } }
    let BASE_COIN_TL = { "Token": "BC", "Reserve": configProviders.INITIAL_RESERVE_BC, "RiskParameters": { "Col_withdrawal": configProviders.RISK_COL_WITHDRAWAL_BC, "Col_required": configProviders.RISK_COL_REQUIRED_BC, "Col_liquidation": configProviders.RISK_COL_LIQUIDATION_BC, "P_limit": configProviders.RISK_P_LIMIT_BC } }
    let PRIVI_DATA_TOKEN_TL = { "Token": "PDT", "Reserve": configProviders.INITIAL_RESERVE_PDT, "RiskParameters": { "Col_withdrawal": configProviders.RISK_COL_WITHDRAWAL_PDT, "Col_required": configProviders.RISK_COL_REQUIRED_PDT, "Col_liquidation": configProviders.RISK_COL_LIQUIDATION_PDT, "P_limit": configProviders.RISK_P_LIMIT_PDT } }
    let BAL_TOKEN_TL = { "Token": "BAL", "Reserve": configProviders.INITIAL_RESERVE_BAL, "RiskParameters": { "Col_withdrawal": configProviders.RISK_COL_WITHDRAWAL_BAL, "Col_required": configProviders.RISK_COL_REQUIRED_BAL, "Col_liquidation": configProviders.RISK_COL_LIQUIDATION_BAL, "P_limit": configProviders.RISK_P_LIMIT_BAL } }
    let BAT_TOKEN_TL = { "Token": "BAT", "Reserve": configProviders.INITIAL_RESERVE_BAT, "RiskParameters": { "Col_withdrawal": configProviders.RISK_COL_WITHDRAWAL_BAT, "Col_required": configProviders.RISK_COL_REQUIRED_BAT, "Col_liquidation": configProviders.RISK_COL_LIQUIDATION_BAT, "P_limit": configProviders.RISK_P_LIMIT_BAT } }
    let COMP_TOKEN_TL = { "Token": "COMP", "Reserve": configProviders.INITIAL_RESERVE_COMP, "RiskParameters": { "Col_withdrawal": configProviders.RISK_COL_WITHDRAWAL_COMP, "Col_required": configProviders.RISK_COL_REQUIRED_COMP, "Col_liquidation": configProviders.RISK_COL_LIQUIDATION_COMP, "P_limit": configProviders.RISK_P_LIMIT_COMP } }
    let DAI_COIN_TL = { "Token": "DAI", "Reserve": configProviders.INITIAL_RESERVE_DAI, "RiskParameters": { "Col_withdrawal": configProviders.RISK_COL_WITHDRAWAL_DAI, "Col_required": configProviders.RISK_COL_REQUIRED_DAI, "Col_liquidation": configProviders.RISK_COL_LIQUIDATION_DAI, "P_limit": configProviders.RISK_P_LIMIT_DAI } }
    let ETH_COIN_TL = { "Token": "ETH", "Reserve": configProviders.INITIAL_RESERVE_ETH, "RiskParameters": { "Col_withdrawal": configProviders.RISK_COL_WITHDRAWAL_ETH, "Col_required": configProviders.RISK_COL_REQUIRED_ETH, "Col_liquidation": configProviders.RISK_COL_LIQUIDATION_ETH, "P_limit": configProviders.RISK_P_LIMIT_ETH } }
    let LINK_TOKEN_TL = { "Token": "LINK", "Reserve": configProviders.INITIAL_RESERVE_LINK, "RiskParameters": { "Col_withdrawal": configProviders.RISK_COL_WITHDRAWAL_LINK, "Col_required": configProviders.RISK_COL_REQUIRED_LINK, "Col_liquidation": configProviders.RISK_COL_LIQUIDATION_LINK, "P_limit": configProviders.RISK_P_LIMIT_LINK } }
    let MKR_COIN_TL = { "Token": "MKR", "Reserve": configProviders.INITIAL_RESERVE_MKR, "RiskParameters": { "Col_withdrawal": configProviders.RISK_COL_WITHDRAWAL_MKR, "Col_required": configProviders.RISK_COL_REQUIRED_MKR, "Col_liquidation": configProviders.RISK_COL_LIQUIDATION_MKR, "P_limit": configProviders.RISK_P_LIMIT_MKR } }
    let UNI_TOKEN_TL = { "Token": "USDT", "Reserve": configProviders.INITIAL_RESERVE_UNI, "RiskParameters": { "Col_withdrawal": configProviders.RISK_COL_WITHDRAWAL_UNI, "Col_required": configProviders.RISK_COL_REQUIRED_UNI, "Col_liquidation": configProviders.RISK_COL_LIQUIDATION_UNI, "P_limit": configProviders.RISK_P_LIMIT_UNI } }
    let USDT_COIN_TL = { "Token": "BC", "Reserve": configProviders.INITIAL_RESERVE_USDT, "RiskParameters": { "Col_withdrawal": configProviders.RISK_COL_WITHDRAWAL_USDT, "Col_required": configProviders.RISK_COL_REQUIRED_USDT, "Col_liquidation": configProviders.RISK_COL_LIQUIDATION_USDT, "P_limit": configProviders.RISK_P_LIMIT_USDT } }
    let WBTC_TOKEN_TL = { "Token": "WBTC", "Reserve": configProviders.INITIAL_RESERVE_WBTC, "RiskParameters": { "Col_withdrawal": configProviders.RISK_COL_WITHDRAWAL_WBTC, "Col_required": configProviders.RISK_COL_REQUIRED_WBTC, "Col_liquidation": configProviders.RISK_COL_LIQUIDATION_WBTC, "P_limit": configProviders.RISK_P_LIMIT_WBTC } }
    let YFI_TOKEN_TL = { "Token": "YFI", "Reserve": configProviders.INITIAL_RESERVE_YFI, "RiskParameters": { "Col_withdrawal": configProviders.RISK_COL_WITHDRAWAL_YFI, "Col_required": configProviders.RISK_COL_REQUIRED_YFI, "Col_liquidation": configProviders.RISK_COL_LIQUIDATION_YFI, "P_limit": configProviders.RISK_P_LIMIT_YFI } }

    let PRIVI_COIN_PRIVICREDIT_RP = { "Token": "PC", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_PC, "Interest_max": configProviders.INTEREST_MAX_PC, "P_incentive_min": configProviders.P_INCENTIVE_MIN_PC, "P_incentive_max": configProviders.P_INCENTIVE_MAX_PC, "P_premium_min": configProviders.P_PREMIUM_MIN_PC, "P_premium_max": configProviders.P_PREMIUM_MAX_PC } }
    let BASE_COIN_PRIVICREDIT_RP = { "Token": "BC", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_BC, "Interest_max": configProviders.INTEREST_MAX_BC, "P_incentive_min": configProviders.P_INCENTIVE_MIN_BC, "P_incentive_max": configProviders.P_INCENTIVE_MAX_BC, "P_premium_min": configProviders.P_PREMIUM_MIN_BC, "P_premium_max": configProviders.P_PREMIUM_MAX_BC }}
    let PRIVI_DATA_TOKEN_PRIVICREDIT_RP = { "Token": "PDT", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_PDT, "Interest_max": configProviders.INTEREST_MAX_PDT, "P_incentive_min": configProviders.P_INCENTIVE_MIN_PDT, "P_incentive_max": configProviders.P_INCENTIVE_MAX_PDT, "P_premium_min": configProviders.P_PREMIUM_MIN_PDT, "P_premium_max": configProviders.P_PREMIUM_MAX_PDT } }
    let BAL_TOKEN_PRIVICREDIT_RP = { "Token": "BAL", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_BAL, "Interest_max": configProviders.INTEREST_MAX_BAL, "P_incentive_min": configProviders.P_INCENTIVE_MIN_BAL, "P_incentive_max": configProviders.P_INCENTIVE_MAX_BAL, "P_premium_min": configProviders.P_PREMIUM_MIN_BAL, "P_premium_max": configProviders.P_PREMIUM_MAX_BAL } }
    let BAT_TOKEN_PRIVICREDIT_RP = { "Token": "BAT", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_BAT, "Interest_max": configProviders.INTEREST_MAX_BAT, "P_incentive_min": configProviders.P_INCENTIVE_MIN_BAT, "P_incentive_max": configProviders.P_INCENTIVE_MAX_BAT, "P_premium_min": configProviders.P_PREMIUM_MIN_BAT, "P_premium_max": configProviders.P_PREMIUM_MAX_BAT }}
    let COMP_TOKEN_PRIVICREDIT_RP = { "Token": "COMP", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_COMP, "Interest_max": configProviders.INTEREST_MAX_COMP, "P_incentive_min": configProviders.P_INCENTIVE_MIN_COMP, "P_incentive_max": configProviders.P_INCENTIVE_MAX_COMP, "P_premium_min": configProviders.P_PREMIUM_MIN_COMP, "P_premium_max": configProviders.P_PREMIUM_MAX_COMP } }
    let DAI_COIN_PRIVICREDIT_RP = { "Token": "DAI", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_DAI, "Interest_max": configProviders.INTEREST_MAX_DAI, "P_incentive_min": configProviders.P_INCENTIVE_MIN_DAI, "P_incentive_max": configProviders.P_INCENTIVE_MAX_DAI, "P_premium_min": configProviders.P_PREMIUM_MIN_DAI, "P_premium_max": configProviders.P_PREMIUM_MAX_DAI } }
    let ETH_COIN_PRIVICREDIT_RP = { "Token": "ETH", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_ETH, "Interest_max": configProviders.INTEREST_MAX_ETH, "P_incentive_min": configProviders.P_INCENTIVE_MIN_ETH, "P_incentive_max": configProviders.P_INCENTIVE_MAX_ETH, "P_premium_min": configProviders.P_PREMIUM_MIN_ETH, "P_premium_max": configProviders.P_PREMIUM_MAX_ETH } }
    let LINK_TOKEN_PRIVICREDIT_RP = { "Token": "LINK", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_LINK, "Interest_max": configProviders.INTEREST_MAX_LINK, "P_incentive_min": configProviders.P_INCENTIVE_MIN_LINK, "P_incentive_max": configProviders.P_INCENTIVE_MAX_LINK, "P_premium_min": configProviders.P_PREMIUM_MIN_LINK, "P_premium_max": configProviders.P_PREMIUM_MAX_LINK } }
    let MKR_COIN_PRIVICREDIT_RP = { "Token": "MKR", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_MKR, "Interest_max": configProviders.INTEREST_MAX_MKR, "P_incentive_min": configProviders.P_INCENTIVE_MIN_MKR, "P_incentive_max": configProviders.P_INCENTIVE_MAX_MKR, "P_premium_min": configProviders.P_PREMIUM_MIN_MKR, "P_premium_max": configProviders.P_PREMIUM_MAX_MKR } }
    let UNI_TOKEN_PRIVICREDIT_RP = { "Token": "UNI", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_UNI, "Interest_max": configProviders.INTEREST_MAX_UNI, "P_incentive_min": configProviders.P_INCENTIVE_MIN_UNI, "P_incentive_max": configProviders.P_INCENTIVE_MAX_UNI, "P_premium_min": configProviders.P_PREMIUM_MIN_UNI, "P_premium_max": configProviders.P_PREMIUM_MAX_UNI } }
    let USDT_COIN_PRIVICREDIT_RP = { "Token": "USDT", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_USDT, "Interest_max": configProviders.INTEREST_MAX_USDT, "P_incentive_min": configProviders.P_INCENTIVE_MIN_USDT, "P_incentive_max": configProviders.P_INCENTIVE_MAX_USDT, "P_premium_min": configProviders.P_PREMIUM_MIN_USDT, "P_premium_max": configProviders.P_PREMIUM_MAX_USDT } }
    let WBTC_TOKEN_PRIVICREDIT_RP = { "Token": "WBTC", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_WBTC, "Interest_max": configProviders.INTEREST_MAX_WBTC, "P_incentive_min": configProviders.P_INCENTIVE_MIN_WBTC, "P_incentive_max": configProviders.P_INCENTIVE_MAX_WBTC, "P_premium_min": configProviders.P_PREMIUM_MIN_WBTC, "P_premium_max": configProviders.P_PREMIUM_MAX_WBTC } }
    let YFI_TOKEN_PRIVICREDIT_RP = { "Token": "YFI", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_YFI, "Interest_max": configProviders.INTEREST_MAX_YFI, "P_incentive_min": configProviders.P_INCENTIVE_MIN_YFI, "P_incentive_max": configProviders.P_INCENTIVE_MAX_YFI, "P_premium_min": configProviders.P_PREMIUM_MIN_YFI, "P_premium_max": configProviders.P_PREMIUM_MAX_YFI } }

    let PRIVI_COIN_PODSWAPPING_RP = { "Token": "PC", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_POD_SWAPPING_PC, "Interest_max": configProviders.INTEREST_MAX_POD_SWAPPING_PC, "Pct_supply_lower": configProviders.PCT_SUPPLY_LOWER_POD_SWAPPING_PC, "Pct_supply_upper": configProviders.PCT_SUPPLY_UPPER_POD_SWAPPING_PC, "Liquidation_min": configProviders.LIQUIDATION_MIN_POD_SWAPPING_PC } }
    let BASE_COIN_PODSWAPPING_RP = { "Token": "BC", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_POD_SWAPPING_BC, "Interest_max": configProviders.INTEREST_MAX_POD_SWAPPING_BC, "Pct_supply_lower": configProviders.PCT_SUPPLY_LOWER_POD_SWAPPING_BC, "Pct_supply_upper": configProviders.PCT_SUPPLY_UPPER_POD_SWAPPING_BC, "Liquidation_min": configProviders.LIQUIDATION_MIN_POD_SWAPPING_BC } }
    let PRIVI_DATA_TOKEN_PODSWAPPING_RP = { "Token": "PDT", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_POD_SWAPPING_PDT, "Interest_max": configProviders.INTEREST_MAX_POD_SWAPPING_PDT, "Pct_supply_lower": configProviders.PCT_SUPPLY_LOWER_POD_SWAPPING_PDT, "Pct_supply_upper": configProviders.PCT_SUPPLY_UPPER_POD_SWAPPING_PDT, "Liquidation_min": configProviders.LIQUIDATION_MIN_POD_SWAPPING_PDT } }
    let BAL_TOKEN_PODSWAPPING_RP = { "Token": "BAL", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_POD_SWAPPING_BAL, "Interest_max": configProviders.INTEREST_MAX_POD_SWAPPING_BAL, "Pct_supply_lower": configProviders.PCT_SUPPLY_LOWER_POD_SWAPPING_BAL, "Pct_supply_upper": configProviders.PCT_SUPPLY_UPPER_POD_SWAPPING_BAL, "Liquidation_min": configProviders.LIQUIDATION_MIN_POD_SWAPPING_BAL } }
    let BAT_TOKEN_PODSWAPPING_RP = { "Token": "BAT", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_POD_SWAPPING_BAT, "Interest_max": configProviders.INTEREST_MAX_POD_SWAPPING_BAT, "Pct_supply_lower": configProviders.PCT_SUPPLY_LOWER_POD_SWAPPING_BAT, "Pct_supply_upper": configProviders.PCT_SUPPLY_UPPER_POD_SWAPPING_BAT, "Liquidation_min": configProviders.LIQUIDATION_MIN_POD_SWAPPING_BAT } }
    let COMP_TOKEN_PODSWAPPING_RP = { "Token": "COMP", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_POD_SWAPPING_COMP, "Interest_max": configProviders.INTEREST_MAX_POD_SWAPPING_COMP, "Pct_supply_lower": configProviders.PCT_SUPPLY_LOWER_POD_SWAPPING_COMP, "Pct_supply_upper": configProviders.PCT_SUPPLY_UPPER_POD_SWAPPING_COMP, "Liquidation_min": configProviders.LIQUIDATION_MIN_POD_SWAPPING_COMP } }
    let DAI_COIN_PODSWAPPING_RP = { "Token": "DAI", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_POD_SWAPPING_DAI, "Interest_max": configProviders.INTEREST_MAX_POD_SWAPPING_DAI, "Pct_supply_lower": configProviders.PCT_SUPPLY_LOWER_POD_SWAPPING_DAI, "Pct_supply_upper": configProviders.PCT_SUPPLY_UPPER_POD_SWAPPING_DAI, "Liquidation_min": configProviders.LIQUIDATION_MIN_POD_SWAPPING_DAI } }
    let ETH_COIN_PODSWAPPING_RP = { "Token": "ETH", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_POD_SWAPPING_ETH, "Interest_max": configProviders.INTEREST_MAX_POD_SWAPPING_ETH, "Pct_supply_lower": configProviders.PCT_SUPPLY_LOWER_POD_SWAPPING_ETH, "Pct_supply_upper": configProviders.PCT_SUPPLY_UPPER_POD_SWAPPING_ETH, "Liquidation_min": configProviders.LIQUIDATION_MIN_POD_SWAPPING_ETH } }
    let LINK_TOKEN_PODSWAPPING_RP = { "Token": "LINK", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_POD_SWAPPING_LINK, "Interest_max": configProviders.INTEREST_MAX_POD_SWAPPING_LINK, "Pct_supply_lower": configProviders.PCT_SUPPLY_LOWER_POD_SWAPPING_LINK, "Pct_supply_upper": configProviders.PCT_SUPPLY_UPPER_POD_SWAPPING_LINK, "Liquidation_min": configProviders.LIQUIDATION_MIN_POD_SWAPPING_LINK } }
    let MKR_COIN_PODSWAPPING_RP = { "Token": "MKR", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_POD_SWAPPING_MKR, "Interest_max": configProviders.INTEREST_MAX_POD_SWAPPING_MKR, "Pct_supply_lower": configProviders.PCT_SUPPLY_LOWER_POD_SWAPPING_MKR, "Pct_supply_upper": configProviders.PCT_SUPPLY_UPPER_POD_SWAPPING_MKR, "Liquidation_min": configProviders.LIQUIDATION_MIN_POD_SWAPPING_MKR } }
    let UNI_TOKEN_PODSWAPPING_RP = { "Token": "UNI", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_POD_SWAPPING_UNI, "Interest_max": configProviders.INTEREST_MAX_POD_SWAPPING_UNI, "Pct_supply_lower": configProviders.PCT_SUPPLY_LOWER_POD_SWAPPING_UNI, "Pct_supply_upper": configProviders.PCT_SUPPLY_UPPER_POD_SWAPPING_UNI, "Liquidation_min": configProviders.LIQUIDATION_MIN_POD_SWAPPING_UNI } }
    let USDT_COIN_PODSWAPPING_RP = { "Token": "USDT", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_POD_SWAPPING_USDT, "Interest_max": configProviders.INTEREST_MAX_POD_SWAPPING_USDT, "Pct_supply_lower": configProviders.PCT_SUPPLY_LOWER_POD_SWAPPING_USDT, "Pct_supply_upper": configProviders.PCT_SUPPLY_UPPER_POD_SWAPPING_USDT, "Liquidation_min": configProviders.LIQUIDATION_MIN_POD_SWAPPING_USDT } }
    let WBTC_TOKEN_PODSWAPPING_RP = { "Token": "WBTC", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_POD_SWAPPING_WBTC, "Interest_max": configProviders.INTEREST_MAX_POD_SWAPPING_WBTC, "Pct_supply_lower": configProviders.PCT_SUPPLY_LOWER_POD_SWAPPING_WBTC, "Pct_supply_upper": configProviders.PCT_SUPPLY_UPPER_POD_SWAPPING_WBTC, "Liquidation_min": configProviders.LIQUIDATION_MIN_POD_SWAPPING_WBTC } }
    let YFI_TOKEN_PODSWAPPING_RP = { "Token": "YFI", "RiskParameters": { "Interest_min": configProviders.INTEREST_MIN_POD_SWAPPING_YFI, "Interest_max": configProviders.INTEREST_MAX_POD_SWAPPING_YFI, "Pct_supply_lower": configProviders.PCT_SUPPLY_LOWER_POD_SWAPPING_YFI, "Pct_supply_upper": configProviders.PCT_SUPPLY_UPPER_POD_SWAPPING_YFI, "Liquidation_min": configProviders.LIQUIDATION_MIN_POD_SWAPPING_YFI } }


    /*
    await coinBalanceInvokeCoins(PRIVI_COIN, "PRIVI COIN LISTING SUCCESSFUL.");
    await coinBalanceInvokeCoins(BASE_COIN, "BASE COIN LISTING SUCCESSFUL.");
    await coinBalanceInvokeCoins(PRIVI_DATA_TOKEN, "PRIVI DATA TOKEN LISTING SUCCESSFUL.");
    await coinBalanceInvokeCoins(BAL_TOKEN, "BALANCER TOKEN LISTING SUCCESSFUL.");
    await coinBalanceInvokeCoins(BAT_TOKEN, "BASIC ATTENTION TOKEN LISTING SUCCESSFUL.");
    await coinBalanceInvokeCoins(COMP_TOKEN, "COMPOUND TOKEN LISTING SUCCESSFUL.");
    await coinBalanceInvokeCoins(DAI_COIN, "DAI STABLECOIN LISTING SUCCESSFUL.");
    await coinBalanceInvokeCoins(ETH_COIN, "ETHEREUM LISTING SUCCESSFUL.");
    await coinBalanceInvokeCoins(LINK_TOKEN, "CHAINLINK TOKEN LISTING SUCCESSFUL.");
    await coinBalanceInvokeCoins(MKR_COIN, "MAKER DAO COIN LISTING SUCCESSFUL.");
    await coinBalanceInvokeCoins(UNI_TOKEN, "UNISWAP TOKEN LISTING SUCCESSFUL.");
    await coinBalanceInvokeCoins(USDT_COIN, "TETHER TOKEN LISTING SUCCESSFUL.");
    await coinBalanceInvokeCoins(WBTC_TOKEN, "WRAPPED BITCOIN TOKEN LISTING SUCCESSFUL.");
    await coinBalanceInvokeCoins(YFI_TOKEN, "YEARN FINANCE TOKEN LISTING SUCCESSFUL.");*/
    /*
    await traditionalLendingInvokeRegisterToken(PRIVI_COIN_TL, "PRIVI COIN REGISTER IN TRADITIONAL LENDING SUCCESSFUL.");
    await traditionalLendingInvokeRegisterToken(BASE_COIN_TL, "BASE COIN REGISTER IN TRADITIONAL LENDING SUCCESSFUL.");
    await traditionalLendingInvokeRegisterToken(PRIVI_DATA_TOKEN_TL, "PRIVI DATA TOKEN REGISTER IN TRADITIONAL LENDING SUCCESSFUL.");
    await traditionalLendingInvokeRegisterToken(BAL_TOKEN_TL, "BALANCER TOKEN REGISTER IN TRADITIONAL LENDING SUCCESSFUL.");
    await traditionalLendingInvokeRegisterToken(BAT_TOKEN_TL, "BASIC ATTENTION TOKEN REGISTER IN TRADITIONAL LENDING SUCCESSFUL.");
    await traditionalLendingInvokeRegisterToken(COMP_TOKEN_TL, "COMPOUND TOKEN REGISTER IN TRADITIONAL LENDING SUCCESSFUL.");
    await traditionalLendingInvokeRegisterToken(DAI_COIN_TL, "DAI STABLECOIN REGISTER IN TRADITIONAL LENDING SUCCESSFUL.");
    await traditionalLendingInvokeRegisterToken(ETH_COIN_TL, "ETHEREUM REGISTER IN TRADITIONAL LENDING SUCCESSFUL.");
    await traditionalLendingInvokeRegisterToken(LINK_TOKEN_TL, "CHAINLINK TOKEN REGISTER IN TRADITIONAL LENDING SUCCESSFUL.");
    await traditionalLendingInvokeRegisterToken(MKR_COIN_TL, "MAKER DAO COIN REGISTER IN TRADITIONAL LENDING SUCCESSFUL.");
    await traditionalLendingInvokeRegisterToken(UNI_TOKEN_TL, "UNISWAP TOKEN REGISTER IN TRADITIONAL LENDING SUCCESSFUL.");
    await traditionalLendingInvokeRegisterToken(USDT_COIN_TL, "TETHER TOKEN REGISTER IN TRADITIONAL LENDING SUCCESSFUL.");
    await traditionalLendingInvokeRegisterToken(WBTC_TOKEN_TL, "WRAPPED BITCOIN TOKEN REGISTER IN TRADITIONAL LENDING SUCCESSFUL.");
    await traditionalLendingInvokeRegisterToken(YFI_TOKEN_TL, "YEARN FINANCE TOKEN REGISTER IN TRADITIONAL LENDING SUCCESSFUL.");*/
    /*
    await priviCreditTokenRiskParameters(PRIVI_COIN_PRIVICREDIT_RP, "PRIVI COIN UPDATE RISK PARAMETERS SUCCESSFUL.");
    await priviCreditTokenRiskParameters(BASE_COIN_PRIVICREDIT_RP, "BASE COIN UPDATE RISK PARAMETERS SUCCESSFUL.");
    await priviCreditTokenRiskParameters(PRIVI_DATA_TOKEN_PRIVICREDIT_RP, "PRIVI DATA TOKEN UPDATE RISK PARAMETERS SUCCESSFUL.");
    await priviCreditTokenRiskParameters(BAL_TOKEN_PRIVICREDIT_RP, "BALANCER TOKEN UPDATE RISK PARAMETERS SUCCESSFUL.");
    await priviCreditTokenRiskParameters(BAT_TOKEN_PRIVICREDIT_RP, "BASIC ATTENTION TOKEN UPDATE RISK PARAMETERS SUCCESSFUL.");
    await priviCreditTokenRiskParameters(COMP_TOKEN_PRIVICREDIT_RP, "COMPOUND TOKEN UPDATE RISK PARAMETERS SUCCESSFUL.");
    await priviCreditTokenRiskParameters(DAI_COIN_PRIVICREDIT_RP, "DAI STABLECOIN UPDATE RISK PARAMETERS SUCCESSFUL.");
    await priviCreditTokenRiskParameters(ETH_COIN_PRIVICREDIT_RP, "ETHEREUM UPDATE RISK PARAMETERS SUCCESSFUL.");
    await priviCreditTokenRiskParameters(LINK_TOKEN_PRIVICREDIT_RP, "CHAINLINK TOKEN UPDATE RISK PARAMETERS SUCCESSFUL.");
    await priviCreditTokenRiskParameters(MKR_COIN_PRIVICREDIT_RP, "MAKER DAO COIN UPDATE RISK PARAMETERS SUCCESSFUL.");
    await priviCreditTokenRiskParameters(UNI_TOKEN_PRIVICREDIT_RP, "UNISWAPUPDATE RISK PARAMETER SUCCESSFUL.");
    await priviCreditTokenRiskParameters(USDT_COIN_PRIVICREDIT_RP, "TETHER TOKEN UPDATE RISK PARAMETER SUCCESSFUL.");
    await priviCreditTokenRiskParameters(WBTC_TOKEN_PRIVICREDIT_RP, "WRAPPED BITCOIN TOKEN UPDATE RISK SUCCESSFUL.");
    await priviCreditTokenRiskParameters(YFI_TOKEN_PRIVICREDIT_RP, "YEARN FINANCE TOKEN UPDATE RISK PARAMETER SUCCESSFUL.");*/
    /*
    await podSwappingRiskParameters(PRIVI_COIN_PODSWAPPING_RP, "PRIVI COIN POD SWAPPING UPDATE RISK PARAMETERS SUCCESSFUL.");
    await podSwappingRiskParameters(BASE_COIN_PODSWAPPING_RP, "BASE COIN POD SWAPPING UPDATE RISK PARAMETERS SUCCESSFUL.");
    await podSwappingRiskParameters(PRIVI_DATA_TOKEN_PODSWAPPING_RP, "PRIVI DATA TOKEN POD SWAPPING UPDATE RISK PARAMETERS SUCCESSFUL.");
    await podSwappingRiskParameters(BAL_TOKEN_PODSWAPPING_RP, "BALANCER TOKEN POD SWAPPING UPDATE RISK PARAMETERS SUCCESSFUL.");
    await podSwappingRiskParameters(BAT_TOKEN_PODSWAPPING_RP, "BASIC ATTENTION TOKEN POD SWAPPING UPDATE RISK PARAMETERS SUCCESSFUL.");
    await podSwappingRiskParameters(COMP_TOKEN_PODSWAPPING_RP, "COMPOUND TOKEN POD SWAPPING UPDATE RISK PARAMETERS SUCCESSFUL.");
    await podSwappingRiskParameters(DAI_COIN_PODSWAPPING_RP, "DAI STABLECOIN POD SWAPPING UPDATE RISK PARAMETERS SUCCESSFUL.");
    await podSwappingRiskParameters(ETH_COIN_PODSWAPPING_RP, "ETHEREUM POD SWAPPING UPDATE RISK PARAMETERS SUCCESSFUL.");
    await podSwappingRiskParameters(LINK_TOKEN_PODSWAPPING_RP, "CHAINLINK TOKEN POD SWAPPING UPDATE RISK PARAMETERS SUCCESSFUL.");
    await podSwappingRiskParameters(MKR_COIN_PODSWAPPING_RP, "MAKER DAO COIN POD SWAPPING UPDATE RISK PARAMETERS SUCCESSFUL.");
    await podSwappingRiskParameters(UNI_TOKEN_PODSWAPPING_RP, "UNISWAP POD SWAPPING UPDATE RISK PARAMETER SUCCESSFUL.");
    await podSwappingRiskParameters(USDT_COIN_PODSWAPPING_RP, "TETHER TOKEN POD SWAPPING UPDATE RISK PARAMETER SUCCESSFUL.");
    await podSwappingRiskParameters(WBTC_TOKEN_PODSWAPPING_RP, "WRAPPED BITCOIN TOKEN POD SWAPPING UPDATE RISK SUCCESSFUL.");
    await podSwappingRiskParameters(YFI_TOKEN_PODSWAPPING_RP, "YEARN FINANCE TOKEN POD SWAPPING UPDATE RISK PARAMETER SUCCESSFUL.");*/

})();

/* ---------------------------------------------------------------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------------------------------------------------------------- */

async function coinBalanceInvokeCoins(coin, stringMessage) {
    let input2 = { body: {} }
    input2.body.username = configProviders.PRIVI_PRIVATE_ID
    input2.body.fcn = "registerToken"
    input2.body.chaincodeName = config.CoinBalance
    input2.body.type = "business"

    input2.body.args = [ JSON.stringify(coin), configProviders.PRIVI_PUBLIC_ID ]
    message = await fabricService.Invoke(input2);
    logger.debug(stringMessage, message)
}

async function traditionalLendingInvokeRegisterToken(coin, stringMessage) {
    let input = { body: {} }
    input.body.username = configProviders.PRIVI_PRIVATE_ID
    input.body.fcn = "registerToken"
    input.body.chaincodeName = config.TraditionalLending
    input.body.type = "business"

    input.body.args = [ JSON.stringify(coin) ]
    message = await fabricService.Invoke(input);
    logger.debug(stringMessage, message)
}


async function priviCreditTokenRiskParameters(coin, stringMessage) {
    let input = { body: {} }
    input.body.username = configProviders.PRIVI_PRIVATE_ID
    input.body.fcn = "updateRiskParameters"
    input.body.chaincodeName = config.PRIVIcredit
    input.body.type = "business"

    input.body.args = [ coin.Token,
                        JSON.stringify(coin.RiskParameters) ]
    message = await fabricService.Invoke(input);
    logger.debug(stringMessage, message)
}

async function podSwappingRiskParameters(coin, stringMessage) {
    let input = { body: {} }
    input.body.username = configProviders.PRIVI_PRIVATE_ID
    input.body.fcn = "updateRiskParameters"
    input.body.chaincodeName = config.PodSwapping
    input.body.type = "business"

    input.body.args = [ coin.Token,
                        JSON.stringify(coin.RiskParameters) ]
    message = await fabricService.Invoke(input);
    logger.debug(stringMessage, message)
}

/* ---------------------------------------------------------------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------------------------------------------------------------- */
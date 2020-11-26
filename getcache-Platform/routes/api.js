'use strict';

const express = require('express')

let router = express.Router();

router.use('/CoinBalance', require('../controllers/CoinBalance.controller'));
router.use('/DataProtocol', require('../controllers/DataProtocol.controller'));
router.use('/TraditionalLending', require('../controllers/TraditionalLending.controller'));
router.use('/PRIVIcredit', require('../controllers/PRIVIcredit.controller'));
router.use('/PodSwapping', require('../controllers/PodSwapping.controller'));
router.use('/PodNFT', require('../controllers/PodNFT.controller'));
router.use('/PodInsurance', require('../controllers/PodInsurance.controller'));
router.use(function (_, res) {
    res.status(404).send("Sorry can't find that!")
});

module.exports = router; 

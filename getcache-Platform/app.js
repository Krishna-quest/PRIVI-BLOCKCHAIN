'use strict';
let log4js = require('log4js');
let logger = log4js.getLogger('CacheBlockchainApp');
let express = require('express');
let bodyParser = require('body-parser');

let app = express();

let cors = require('cors');
//"basecoin":"bcbeta46",
require('./config.js');
// uncomment the line below if fresh start
require('./providers/startcachekit.js')
//require('./providers/testtest')
app.options('*', cors());
app.use(cors());
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({
	extended: false
}));

//require routes
let api = require('./routes/api');

app.use('/api', api);
//app.listen(3000);
let server = app.listen(4000, function() {
	console.log('Cache server listening on port ' + server.address().port);
  });
server.timeout = 999000;

logger.debug("****************** SERVER STARTED ************************");
//###########################



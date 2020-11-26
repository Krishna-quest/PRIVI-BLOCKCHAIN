HYPERLEDGER FABRIC NETWORK SETUP

0) pm2 stop app.js
1) docker rm -f $(docker ps -aq)
2) rm -rf /tmp/fab*
3) cd artifacts 
4) docker-compose -f docker-compose.yaml up -d
5) cd ..
6) node app.js
7) pm2 start app.js

docker ps
docker logs artifacts_peer0_companies_1

FILTER CONTAINERS CONTAINING STRING
docker ps -f name=dev-peer0.companies.com-CoinBalance

DELETE CONTAINERS CONTAINING STRING
docker stop $(docker ps -q --filter name=dev-peer0.companies.com-CoinBalanceS$i)



STEPS TO GET THE NETWORK UP AND RUNNING

1. Go to the artifacts directory and start the docker compose file, Use this command to do the same:
cd artifacts && docker-compose -f docker-compose.yaml up -d
2. Go back to the Home directory
    cd ..
3. Install the node modules and start the application:
    npm i 
    node app.js

NOTE:: IF YOU ARE RUNNING THE APPLICATION FOR THE FIRST TIME PLEASE UNCOMMENT THIS ON THE "App.js" file:
//uncomment the line below if fresh start
//require('./providers/startcachekit')

Also, remove the certificates from tmp folder 
rm -rf /tmp/fab*

This is a script I created to automate the configuration process.
will create channel, join peers , install chaincode and instantiate

IMPORTANT:: Comment the line back if this is not a fresh network, if the channel and chaincodes already exists.


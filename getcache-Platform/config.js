let util = require('util');
let path = require('path');
let hfc = require('fabric-client');
    
let file = 'network-config%s.yaml';
//var file = 'network-config.yaml';
    
let env = process.env.TARGET_NETWORK;
    if (env)
        file = util.format(file, '-' + env);
    else
        file = util.format(file, '');
hfc.setConfigSetting('network-connection-profile-path',path.join(__dirname, 'artifacts' ,file));
hfc.setConfigSetting('users-connection-profile-path',path.join(__dirname, 'artifacts', 'users.yaml'));
hfc.setConfigSetting('companies-connection-profile-path',path.join(__dirname, 'artifacts', 'companies.yaml'));
hfc.setConfigSetting('exchanges-connection-profile-path',path.join(__dirname, 'artifacts', 'exchanges.yaml'));
hfc.addConfigFile(path.join(__dirname, 'config.json'));
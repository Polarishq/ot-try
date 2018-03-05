"use strict";

//The name of the service being tested

//Location and configuration of the service being tested
var TARGET = {
    COMPONENT_NAME: 'ot-try',
    SERVICE: 'ot-try',
    VERSION: process.env.TARGET_VERSION || 'v1'
};
TARGET.SERVER = process.env.TARGET_SERVER || (TARGET.SERVICE + '-' + TARGET.VERSION);
TARGET.BASE_URL = 'http://' + TARGET.SERVER;// + '/' + TARGET.VERSION;

var GLOBAL_SETTINGS = {
    //Determines if dependant services are being run in stubbed mode
    TEST_DOUBLES: process.env.TEST_DOUBLES === 'true',
    SKIP_TIMEOUTS: process.env.SKIP_TIMEOUTS === 'true',
    ENVIRONMENT_NAME: process.env.ENVIRONMENT_NAME
};

console.log('TARGET: ', TARGET);
console.log('GLOBAL_SETTINGS:', GLOBAL_SETTINGS);

module.exports = {
    TARGET: TARGET,
    // Add ALL dependant services here so their cleanup endpoints are called in suite_teardown.js
    CALLOUT_TARGETS: [ TARGET.BASE_URL  ],
    GLOBAL_SETTINGS: GLOBAL_SETTINGS
};
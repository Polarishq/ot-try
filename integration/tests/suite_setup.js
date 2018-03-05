"use strict";

var http = require("http");
var request = require("request");
var settings = require("./settings.js");

var valid_codes = [200];

if (require.main === module) {
    var tries = 0;

    var waitForHealth = function(url) {
        request(url,
            function(error, response, body) {
                if ( error || ! valid_codes.includes(response.statusCode) ) {
                    var msg = error || response.body;
                    console.error("Error waiting for " + url +
                        " to come up: " + msg);
                    if (tries++ > 10) {
                        process.exit(1)
                    }
                    else {
                        setTimeout(function() {waitForHealth(url)}, 1000);
                    }
                }
                else {
                    console.info(url+ " is up: " + body);
                }
            }
        );
    };

    var services = settings.CALLOUT_TARGETS;
    for (var i in services) {
        var service = services[i];
        console.info("Start waiting for " + service + " to come up");
        waitForHealth(service + "/health");
    }
}

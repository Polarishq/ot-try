"use strict";

var http = require("http");
var request = require("request");
var settings = require("./settings.js");

var valid_codes = [200, 204];

if (require.main === module) {
    var cleanup = function(service) {
        return new Promise(function(resolve, reject) {
            console.info("Cleaning " + service + "...");
            request.post(service + "/cleanup",
                function(error, response, body) {
                    if ( error || ! valid_codes.includes(response.statusCode) ) {
                        console.error("Error cleaning " + service);
                        if ( error !== null ) {
                            console.error(error);
                        }
                        console.error("statusCode= " + response.statusCode);
                        console.error("body=" + response.body);
                        reject(error);
                        process.exit(1);
                    }
                    else {
                        console.info("Cleaned " + service);
                        console.error("statusCode= " + response.statusCode);
                        console.error("body=" + response.body);
                        resolve(service);
                    }
                }
            );
        });
    };

    var services = settings.CALLOUT_TARGETS;
    for (var i in services) {
        var service = services[i];
        cleanup(service).then(
            function(svc) {
            },
            function(err) {
            }
        );
    }
}

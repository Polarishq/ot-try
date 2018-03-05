"use strict";

var chakram = require('chakram'),
    expect = chakram.expect;

var polarisResponse = function (respObj, status) {
    expect(respObj).to.have.status(status);
    expect(respObj).to.have.header("content-type", "application/json");
};

chakram.addProperty("polaris", function(){});
chakram.addMethod("response", polarisResponse);

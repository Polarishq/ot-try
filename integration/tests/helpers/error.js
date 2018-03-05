"use strict";

var chakram = require('chakram'),
    expect = chakram.expect;

const ERROR_SCHEMA = {
    properties: {
        code: {type: "integer"},
        message: {type: "string"}
    },
    required: ["code", "message"]
};

var validatePolarisError = function (respObj, status, code, message) {
    expect(respObj).to.have.schema(ERROR_SCHEMA);
    expect(respObj).to.have.status(status);
    var err = respObj.body;
    var re = new RegExp(message);
    expect(err.code).to.eql(code);
    expect(err.message).to.match(re);
};

chakram.addProperty("polaris", function(){});
chakram.addMethod("error", validatePolarisError);

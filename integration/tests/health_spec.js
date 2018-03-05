"use strict";

var settings = require("./settings.js");
var chakram = require('chakram'),
    expect = chakram.expect;

describe("Health", function () {
    it("should return basic information", function () {
        return chakram.get(settings.TARGET.BASE_URL + '/health')
            .then(function (response) {
                expect(response).to.have.status(200);

                var health = response.body;
                expect(health.healthy).to.equal(true);
                expect(health.stub).to.equal(false);

                expect(settings.TARGET.SERVICE).to.equal(health.service_info.name);
                expect(settings.TARGET.VERSION).to.equal(health.service_info.version);
                console.info(settings.TARGET.COMPONENT_NAME + " is healthy and its revision is " + health.revision);
            })
    });
});

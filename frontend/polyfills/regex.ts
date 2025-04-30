'use strict';

var define = require('define-properties');

var getRegExpPolyfill = function () {
    return globalThis.RegExp || require('regex');
};

module.exports = function shimRegExp() {
    var polyfill = getRegExpPolyfill();

    define(
        globalThis,
        { RegExp: polyfill },
        { RegExp: function () { return globalThis.RegExp !== polyfill; } }
    );

    return polyfill;
};

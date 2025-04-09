'use strict';

// @ts-ignore
delete globalThis.PerformanceObserver;

var define = require('define-properties');
var getPerformanceObserverPolyfill = function () {
    return require('@fastly/performance-observer-polyfill').default;
};

module.exports = function shimPerformanceObserver() {
    var polyfill = getPerformanceObserverPolyfill();

    define(
        globalThis,
        { PerformanceObserver: polyfill },
        { PerformanceObserver: function () { return globalThis.PerformanceObserver !== polyfill; } }
    );

    return polyfill;
};

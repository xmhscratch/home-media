if (typeof globalThis !== 'object') {
    // the below function is not CSP-compliant, but reliably gets the
    // global object in sloppy mode in every engine.
    var getGlobal = Function('return this');

    /* when `globalThis` is not present */
    var shimmedGlobal = require('globalthis').shim();

    var shimmedFlat = require('array.prototype.flat').shim();

    var shimmedQueueMicrotask = require('./polyfills/queue-microtask')();

    var shimmedReplaceAll = require('string.prototype.replaceall').shim();

    var shimmedPerformanceObserver = require('./polyfills/performance-observer')();

    var shimmedMatchAll = require('string.prototype.matchall').shim();

    var shimmedFromEntries = !Object.fromEntries ? require('object.fromentries').shim() : undefined;
}

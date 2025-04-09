if (typeof globalThis !== 'object') {
    // the below function is not CSP-compliant, but reliably gets the
    // global object in sloppy mode in every engine.
    const getGlobal = Function('return this');

    /* when `globalThis` is not present */
    const shimmedGlobal = require('globalthis').shim();

    const shimmedFlat = require('array.prototype.flat').shim();

    const shimmedRegExp = require('./polyfills/regex')();

    const shimmedQueueMicrotask = require('./polyfills/queue-microtask')();

    const shimmedReplaceAll = require('string.prototype.replaceall').shim();

    const shimmedPerformanceObserver = require('./polyfills/performance-observer')();

    const shimmedMatchAll = require('string.prototype.matchall').shim();

    const shimmedFromEntries = !Object.fromEntries ? require('object.fromentries').shim() : undefined;
}

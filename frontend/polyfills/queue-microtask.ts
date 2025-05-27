'use strict';

var define = require('define-properties');

var getQueueMicrotaskPolyfill = function () {
  return globalThis.queueMicrotask || require('queue-microtask');
};

module.exports = function shimQueueMicrotask() {
  var polyfill = getQueueMicrotaskPolyfill();

  define(
    globalThis,
    { queueMicrotask: polyfill },
    {
      queueMicrotask: function () {
        return globalThis.queueMicrotask !== polyfill;
      },
    },
  );

  return polyfill;
};

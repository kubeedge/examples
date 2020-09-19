var syslogParser = require('../lib/glossy/parse.js'),
          assert = require('assert');

assert.ok(syslogParser, 'parser loaded');

assert.deepEqual(
    syslogParser.parseRfc3339("1985-04-12T23:20:50.52Z"),
    new Date(482196050000)
);

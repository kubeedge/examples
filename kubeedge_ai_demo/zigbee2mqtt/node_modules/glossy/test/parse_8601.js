var syslogParser = require('../lib/glossy/parse.js'),
          assert = require('assert');

assert.ok(syslogParser, 'parser loaded');

assert.deepEqual(
    syslogParser.parse8601('2011-10-10T14:48:00'), 
    new Date(Date.parse('2011-10-10T14:48:00'))
);

assert.equal(
    syslogParser.parse8601('foo'), 
    undefined
);

var syslogParser = require('../lib/glossy/parse.js'),
          assert = require('assert');

assert.ok(syslogParser, 'parser loaded');
assert.equal(syslogParser.decideValue(1), "1");
assert.equal(syslogParser.decideValue('-'), null);
assert.equal(syslogParser.decideValue('ー'), 'ー');

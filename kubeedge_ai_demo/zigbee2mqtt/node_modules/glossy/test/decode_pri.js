var syslogParser = require('../lib/glossy/parse.js'),
          assert = require('assert');

assert.ok(syslogParser, 'parser loaded');
assert.deepEqual(syslogParser.decodePri('<16>'), {
    prival: 16,
    facilityID: 2,
    severityID: 0,
    facility: 'mail',
    severity: 'emerg'
});

assert.deepEqual(syslogParser.decodePri('<66>1'), {
    prival: 66,
    facilityID: 8,
    severityID: 2,
    facility: 'uucp',
    severity: 'crit'
});


assert.equal(syslogParser.decodePri('1<16>'), false);
assert.equal(syslogParser.decodePri('<200>'), false);


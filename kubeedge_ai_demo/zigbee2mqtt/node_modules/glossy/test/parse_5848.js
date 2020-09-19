var syslogParser = require('../lib/glossy/parse.js'),
          assert = require('assert');

assert.ok(syslogParser, 'parser loaded');

var fullySigned = '<110>1 2009-05-03T14:00:39.519307+02:00 host.example.org syslogd 2138 - [ssign-cert VER="0111" RSID="1" SG="0" SPRI="0" TPBL="587" INDEX="1" FLEN="587" FRAG="2009-05-03T14:00:39.519005+02:00 K BACsLMZ NCV2NUAwe4RAeAnSQuvv2KS51SnHFAaWJNU2XVDYvW1LjmJgg4vKvQPo3HEOD+2hEkt1zcXADe03u5pmHoWy5FGiyCbglYxJkUJJrQqlTSS6vID9yhsmEnh07w3pOsxmb4qYo0uWQrAAenBweVMlBgV3ZA5IMA8xq8l+i8wCgkWJjCjfLar7s+0X3HVrRroyARv8EAIYoxofh9m N8n821BTTuQnz5hp40d6Z3UudKePu2di5Mx3GFelwnV0Qh5mSs0YkuHJg0mcXyUAoeYry5X6482fUxbm+gOHVmYSDtBmZEB8PTEt8Os8aedWgKEt/E4dT+Hmod4omECLteLXxtScTMgDXyC+bSBMjRRCaeWhHrYYdYBACCWMdTc12hRLJTn8LX99kv1I7qwgieyna8GCJv/rEgC ssS9E1qARM+h19KovIUOhl4VzBw3rK7v8Dlw/CJyYDd5kwSvCwjhO21LiReeS90VPYuZFRC1B82Sub152zOqIcAWsgd4myCCiZbWBsuJ8P0gtarFIpleNacCc6OV3i2Rg==" SIGN="AKAQEUiQptgpd0lKcXbuggGXH/dCdQCgdysrTBLUlbeGAQ4vwrnLOqSL7+c="]';

var fullySignedSD = syslogParser.parseStructure(fullySigned);
assert.deepEqual(fullySignedSD, { 
    'ssign-cert': { 
        VER: '0111',
        RSID: '1',
        SG: '0',
        SPRI: '0',
        TPBL: '587',
        INDEX: '1',
        FLEN: '587',
        FRAG: '2009-05-03T14:00:39.519005+02:00 K BACsLMZ NCV2NUAwe4RAeAnSQuvv2KS51SnHFAaWJNU2XVDYvW1LjmJgg4vKvQPo3HEOD+2hEkt1zcXADe03u5pmHoWy5FGiyCbglYxJkUJJrQqlTSS6vID9yhsmEnh07w3pOsxmb4qYo0uWQrAAenBweVMlBgV3ZA5IMA8xq8l+i8wCgkWJjCjfLar7s+0X3HVrRroyARv8EAIYoxofh9m N8n821BTTuQnz5hp40d6Z3UudKePu2di5Mx3GFelwnV0Qh5mSs0YkuHJg0mcXyUAoeYry5X6482fUxbm+gOHVmYSDtBmZEB8PTEt8Os8aedWgKEt/E4dT+Hmod4omECLteLXxtScTMgDXyC+bSBMjRRCaeWhHrYYdYBACCWMdTc12hRLJTn8LX99kv1I7qwgieyna8GCJv/rEgC ssS9E1qARM+h19KovIUOhl4VzBw3rK7v8Dlw/CJyYDd5kwSvCwjhO21LiReeS90VPYuZFRC1B82Sub152zOqIcAWsgd4myCCiZbWBsuJ8P0gtarFIpleNacCc6OV3i2Rg==',
        SIGN: 'AKAQEUiQptgpd0lKcXbuggGXH/dCdQCgdysrTBLUlbeGAQ4vwrnLOqSL7+c=' 
    } 
});

assert.deepEqual(syslogParser.parseSignedCertificate(fullySignedSD['ssign-cert']), {
    version: '01',
    hashAlgorithm: 1,
    hashAlgoString: 'SHA1',
    sigScheme: 1,
    rebootSessionID: 1,
    signatureGroup: 0,
    signaturePriority: 0,
    totalPayloadLength: 587,
    payloadIndex: 1,
    fragmentLength: 587,
    payloadTimestamp: new Date('2009-05-03T14:00:39.519005+02:00'),
    payloadType: 'K',
    payloadName: 'Public Key',
    keyBlob: new Buffer(fullySignedSD['ssign-cert']['FRAG'].split(/\s/)[2], encoding='base64'),
    thisSignature:  new Buffer(fullySignedSD['ssign-cert']['SIGN'], encoding='base64')
});


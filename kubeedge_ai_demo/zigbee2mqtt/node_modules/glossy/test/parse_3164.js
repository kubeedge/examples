var syslogParser = require('../lib/glossy/parse.js'),
 syslogGenerator = require('../lib/glossy/produce.js'),
          assert = require('assert');

assert.ok(syslogParser, 'parser loaded');

var gen = new syslogGenerator({type: 'bsd'});

var doubleSpaced = "<13>Feb  5 17:32:18 10.0.0.99 Use the BFG!";
syslogParser.parse(doubleSpaced, function(parsedMessage){
    var msg = gen.produce(parsedMessage);
    assert.equal(doubleSpaced, msg);    

    var expectedData = { 
        originalMessage: doubleSpaced,
        prival: 13,
        facilityID: 1,
        severityID: 5,
        facility: 'user',
        severity: 'notice',
        type: 'RFC3164',
        host: '10.0.0.99',
        message: 'Use the BFG!' };

    delete parsedMessage.date;
    delete parsedMessage.time;
    delete parsedMessage.timestamp;

    assert.deepEqual(parsedMessage, expectedData);
});

var withCommand = "<34>Oct 11 22:14:15 mymachine su: 'su root' failed for lonvick on /dev/pts/8";
syslogParser.parse(withCommand, function(parsedMessage){
    var expectedData = {
        originalMessage: withCommand,
        prival: 34,
        facilityID: 4,
        severityID: 2,
        facility: 'auth',
        severity: 'crit',
        type: 'RFC3164',
        host: 'mymachine',
        message: "su: 'su root' failed for lonvick on /dev/pts/8" };

    var parsedDate = parsedMessage.time;
    delete parsedMessage.time;

    assert.equal(parsedDate.getUTCMonth(), 9);
    assert.equal(parsedDate.getUTCHours(), 20);
    assert.deepEqual(parsedMessage, expectedData);

});

var withDifficultTime = "<191>94103: 51w2d: DHCPD: assigned IP address 10.10.1.94 to client 0100.01c4.21d3.b3";
syslogParser.parse(withDifficultTime, function(parsedMessage){
    var expectedData = { 
        originalMessage: withDifficultTime,
        prival: 191,
        facilityID: 23,
        severityID: 7,
        facility: 'local7',
        severity: 'debug',
        type: 'RFC3164',
        time: undefined,
        message: '51w2d: DHCPD: assigned IP address 10.10.1.94 to client 0100.01c4.21d3.b3'};

    assert.deepEqual(parsedMessage, expectedData);
});

var withYear = "<32>Mar 05 2011 22:21:02: %ASA-6-302013: Built inbound TCP connection 401 for outside:123.123.123.123/4413 (123.123.123.123/4413) to net:BOX/25 (BOX/25)";
syslogParser.parse(withYear, function(parsedMessage){
    var expectedData = { 
        originalMessage: withYear,
        prival: 32,
        facilityID: 4,
        severityID: 0,
        facility: 'auth',
        severity: 'emerg',
        type: 'RFC3164',
        time: undefined,
        host: '22:21:02:',
        message: '%ASA-6-302013: Built inbound TCP connection 401 for outside:123.123.123.123/4413 (123.123.123.123/4413) to net:BOX/25 (BOX/25)' };

   assert.deepEqual(parsedMessage, expectedData);
});

var withSpaces = "<13>Mar 15 11:22:40 myhost.com     0    11,03/15/12,11:22:38,§ó·s,10.10.10.171,,40C6A91373B6,";
syslogParser.parse(withSpaces, function(parsedMessage){
    var expectedData = { 
        originalMessage: withSpaces,
        prival: 13,
        facilityID: 1,
        severityID: 5,
        facility: 'user',
        severity: 'notice',
        type: 'RFC3164',
        host: 'myhost.com',
        message: '    0    11,03/15/12,11:22:38,§ó·s,10.10.10.171,,40C6A91373B6,' };

    delete parsedMessage.time;
    assert.deepEqual(parsedMessage, expectedData);

});


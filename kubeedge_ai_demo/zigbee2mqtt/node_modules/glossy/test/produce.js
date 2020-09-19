var assert       = require('assert');
var producer     = require('../lib/glossy/produce.js');

assert.ok(producer, 'producer loaded');

var syslogProducer = new producer();
assert.ok(syslogProducer, 'new SyslogProducer object created');
assert.equal(syslogProducer.type, 'RFC5424', 'Syslog Producer set correctly');

var BSDProducer = new producer({ type: 'BSD'});
assert.ok(BSDProducer, 'new BSDProducer object created');
assert.equal(BSDProducer.type, 'RFC3164', 'BSD Producer set correctly');

var presetProducer = new producer({
    type:     'bsd',
    facility: 'ntp',
    host:     'localhost',
    appName:   'kill'
});

var invalidProducer = new producer({
    type: 'invalid',
    facility: 'invalid',
});
assert.notEqual(invalidProducer, 'invalid producer is null');

var msg = syslogProducer.produce({
    facility: 'local4',
    severity: 'error',
    host: 'localhost',
    appName: 'sudo',
    pid: '123',
    date: new Date(1234567890000),
    message: 'Test Message'
});
assert.equal(msg, "<163>1 2009-02-13T23:31:30.00+01:00 localhost sudo 123 - - Test Message",'Valid message returned');

syslogProducer.produce({
    facility: 'audit',
    severity: 'error',
    host: '127.0.0.1',
    appName: 'sudo',
    pid: '419',
    date: new Date(1234567890000),
    message: 'Test Message'
}, function(cbMsg) {
    assert.equal(cbMsg, '<107>1 2009-02-13T23:31:30.00+01:00 127.0.0.1 sudo 419 - - Test Message', 'Valid message in callback returned');
});

BSDProducer.produce({
    facility: 'audit',
    severity: 'error',
    host: '127.0.0.1',
    appName: 'sudo',
    pid: '419',
    date: new Date(1234567890000),
    message: 'Test Message'
}, function(cbMsg){
    assert.equal(cbMsg, '<107>Feb 14 00:31:30 127.0.0.1 sudo[419]: Test Message');
});

var debugMsg = presetProducer.debug({
    facility: 'local2',
    message: 'Debug Message',
    date: new Date(1234567890000),
    pid: 91
});
assert.ok(debugMsg);
assert.equal(debugMsg, '<151>Feb 14 00:31:30 localhost kill[91]: Debug Message');

var infoMsg = presetProducer.info({
    facility: 'ntp',
    message: 'Info Message',
    pid: 42,
    date: new Date(1234567890000)
});
assert.ok(infoMsg);
assert.equal(infoMsg, '<102>Feb 14 00:31:30 localhost kill[42]: Info Message');

var noticeMsg = presetProducer.debug({
    facility: 'local2',
    message: 'Notice Message',
    pid: 16,
    date: new Date(1234567890000)
});
assert.ok(noticeMsg);
assert.equal(noticeMsg, '<151>Feb 14 00:31:30 localhost kill[16]: Notice Message');

var warnMsg = presetProducer.debug({
    facility: 'local4',
    message: 'Warning Message',
    pid: 91,
    date: new Date(1234567890000)
});
assert.ok(warnMsg);
assert.equal(warnMsg, '<167>Feb 14 00:31:30 localhost kill[91]: Warning Message');

var errorMsg = presetProducer.debug({
    facility: 'clock',
    message: 'Error Message',
    pid: 91,
    date: new Date(1234567890000)
});
assert.ok(errorMsg);
assert.equal(errorMsg, '<79>Feb 14 00:31:30 localhost kill[91]: Error Message');

var criticalMsg = presetProducer.crit({
    facility: 'local0',
    message: 'Critical Message',
    pid: 91,
    date: new Date(1234567890000)
});
assert.ok(criticalMsg);
assert.equal(criticalMsg, '<130>Feb 14 00:31:30 localhost kill[91]: Critical Message');

var alertMsg = presetProducer.alert({
    facility: 'clock',
    message: 'Alert Message',
    pid: 91,
    date: new Date(1234567890000)
});
assert.ok(alertMsg);
assert.equal(alertMsg, '<73>Feb 14 00:31:30 localhost kill[91]: Alert Message');

var emergencyMsg = presetProducer.emergency({
    facility: 'news',
    message: 'Emergency Message',
    pid: 91,
    date: new Date(1234567890000)
});
assert.ok(emergencyMsg);
assert.equal(emergencyMsg, '<56>Feb 14 00:31:30 localhost kill[91]: Emergency Message');

var structuredMsg = syslogProducer.produce({
    facility: 'local4',
    severity: 'error',
    host: 'mymachine.example.com',
    appName: 'evntslog',
    msgID: 'ID47',
    date: new Date(1234567890000),
    structuredData: {
        'exampleSDID@32473': {
            'iut':         "3",
            'eventSource': "Application",
            'eventID':     "1011",
            'seqNo':       "1"
        }
    },
    message: 'BOMAn application event log entry...'
});

assert.ok(structuredMsg);
assert.equal(structuredMsg, '<163>1 2009-02-13T23:31:30.00+01:00 mymachine.example.com evntslog - ID47 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011" seqNo="1"] BOMAn application event log entry...');

var messageWithOneDigitDate = presetProducer.emergency({
    facility: 'news',
    message: 'Emergency Message',
    pid: 91,
    date: new Date(1233531090000)
});
assert.ok(messageWithOneDigitDate);
assert.equal(messageWithOneDigitDate, '<56>Feb  2 00:31:30 localhost kill[91]: Emergency Message');

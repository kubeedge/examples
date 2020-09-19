var syslogParser = require('../lib/glossy/parse.js'),
          assert = require('assert');

assert.ok(syslogParser, 'parser loaded');
var withPrecisionTime = "<165>1 2003-08-24T05:14:15.000003-07:00 192.0.2.1 myproc 8710 - - %% It's time to make the do-nuts.";
syslogParser.parse(withPrecisionTime, function(parsedMessage){
    var expectedData = { 
        originalMessage: withPrecisionTime,
        prival: 165,
        facilityID: 20,
        severityID: 5,
        facility: 'local4',
        severity: 'notice',
        type: 'RFC5424',
        host: '192.0.2.1',
        appName: 'myproc',
        pid: '8710',
        msgID: null,
        message: "%% It's time to make the do-nuts." };
    
    delete parsedMessage.time;
    assert.deepEqual(parsedMessage, expectedData);

});

// FIXME 3 minute offset from UTC?!
var with8601 = "<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su - ID47 - BOM'su root' failed for lonvick on /dev/pts/8";
syslogParser.parse(with8601, function(parsedMessage){
    var expectedData = { 
        originalMessage: with8601,
        prival: 34,
        facilityID: 4,
        severityID: 2,
        facility: 'auth',
        severity: 'crit',
        type: 'RFC5424',
        time:  new Date('2003-10-11T22:14:15.003Z'),
        host: 'mymachine.example.com',
        appName: 'su',
        pid: null,
        msgID: 'ID47',
        message: "BOM'su root' failed for lonvick on /dev/pts/8" };
    
    assert.deepEqual(parsedMessage, expectedData);
});

// FIXME 3 minute offset from UTC?!
var withSD = '<165>1 2003-10-11T22:14:15.003Z mymachine.example.com evntslog - ID47 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] BOMAn application event log entry...';
syslogParser.parse(withSD, function(parsedMessage){
    var expectedData = { 
        originalMessage: withSD,
        prival: 165,
        facilityID: 20,
        severityID: 5,
        facility: 'local4',
        severity: 'notice',
        type: 'RFC5424',
        time:  new Date('2003-10-11T22:14:15.003Z'),
        host: 'mymachine.example.com',
        appName: 'evntslog',
        pid: null,
        msgID: 'ID47',
        structuredData: { 
            'exampleSDID@32473': { 
                iut: '3',
                eventSource: 'Application',
                eventID: '1011' 
            } 
        }, 
        message: 'BOMAn application event log entry...' };

    assert.deepEqual(parsedMessage, expectedData);
});

// FIXME 3 minute offset from UTC?!
var withDoubleSD =  '<165>1 2003-10-11T22:14:15.003Z mymachine.example.com evntslog - ID47 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"][examplePriority@32473 class="high"]';
syslogParser.parse(withDoubleSD, function(parsedMessage){
    var expectedStructuredData = { 
        'exampleSDID@32473': { 
            iut: '3', 
            eventSource: 'Application',
            eventID: '1011' 
         },
         'examplePriority@32473': {
             'class': 'high'
         }
    };

    var expectedData = { 
        originalMessage: withDoubleSD,
        prival: 165,
        facilityID: 20,
        severityID: 5,
        facility: 'local4',
        severity: 'notice',
        type: 'RFC5424',
        time:  new Date('2003-10-11T22:14:15.003Z'),
        host: 'mymachine.example.com',
        appName: 'evntslog',
        pid: null,
        msgID: 'ID47',
        structuredData: expectedStructuredData,  //FIXME Both sets should be there
        message: '' };
    assert.deepEqual(parsedMessage, expectedData);
});


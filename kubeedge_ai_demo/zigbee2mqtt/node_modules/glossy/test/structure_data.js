var syslogParser = require('../lib/glossy/parse.js'),
          assert = require('assert');

assert.ok(syslogParser, 'parser loaded');
var singleStructure = '[exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"]';
assert.deepEqual(syslogParser.parseStructure(singleStructure), {
    'exampleSDID@32473': {
        iut: '3', 
        eventSource: 'Application',
        eventID: '1011' 
    } 
});

var doubleStructure = '[exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"][examplePriority@32473 class="high"] ';
assert.deepEqual(syslogParser.parseStructure(doubleStructure), {
    'exampleSDID@32473': {
            iut: '3',
        eventID: '1011',
        eventSource: 'Application'
     },
     'examplePriority@32473': {
         'class': 'high'
     }
});

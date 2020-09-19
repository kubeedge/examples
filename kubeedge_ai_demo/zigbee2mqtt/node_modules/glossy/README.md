glossy
===========

glossy aims to be a very generic yet powerful library for both producing and
also parsing raw syslog messages. The library aims to be capable of adhearing to
RFC 3164, RFC 5424 and RFC 5848 and by itself does no network interactions, it's
up to you to use this library as a syslog producer, a consumer, relay or
something else entirely. In addition, glossy has no dependencies and can be
bootstrapped to operate in browser or other non-node.js environments.


Parsing
-------

    var syslogParser = require('glossy').Parse; // or wherever your glossy libs are
    
    parsedMessage = syslogParser.parse(message);

parsedMessage will return an object containing as many parsed values as
possible, as well as the original message. The date value will be a Date object.
Structured data will return as an object. Alternatively, you can give it a
callback as your second argument:

    syslogParser.parse(message, function(parsedMessage){
        console.log(parsedMessage);
    });


Producing
-------
Unless you stipulate for BSD/RFC 3164 style messages, it will default to
generating all messages as newer, RFC 5424 format. This might break consumers or
relays not expecting it.

    var syslogProducer = require('glossy').Produce; // or wherever glossy lives

    var msg = syslogProducer.produce({
        facility: 'local4', // these can either be a valid integer, 
        severity: 'error',  // or a relevant string
        host: 'localhost',
        appName: 'sudo',
        pid: '123',
        date: new Date(Date()),
        message: 'Nice, Neat, New, Oh Wow'
    });

Again, you can specify a callback for the second argument.

    var msg = syslogProducer.produce({
        facility: 'ntp', 
        severity: 'info',
        host: 'localhost',
        date: new Date(Date()),
        message: 'Lunch Time!'
    }, function(syslogMsg){
        console.log(syslogMsg);
    });

In addition, you can also predefined most of the values when you create the
object, to save having to repeat yourself:

    var syslogProducer = new require('glossy').Produce({
        type: 'BSD',
        facility: 'ftp',
        pid: 42,
        host: '::1'        
    });

For RFC5424 messages, you can also include structured data. Keys should comply
with the definition in [Section 7, RFC5424](http://tools.ietf.org/html/rfc5424#section-7) 
regarding names - keep them unique and your own custom keys should have at least
an @ sign.

    var msg = syslogProducer.produce({
        facility: 'local4', 
        severity: 'error',
        host: 'localhost',
        appName: 'starman',
        pid: '123',
        date: new Date(Date()),
        message: 'ACHTUNG!',
        structuredData: {
            'plack@host': {
                status: 'broken',
                hasTried: 'not really'
            }
        }
    });

Finally, we expose all the severities as functions themselves:

    var infoMsg = glossy.info({
           message: 'Info Message',
    });

Function names facilitating this are named debug, info, notice, warn, crit,
alert and emergency.

Parsing Example
-------
Handle incoming syslog messages coming in on UDP port 514:

    var syslogParser = require('glossy').Parse; // or wherever your glossy libs are
    var dgram  = require("dgram");
    var server = dgram.createSocket("udp4");
    
    server.on("message", function(rawMessage) {
        syslogParser.parse(rawMessage.toString('utf8', 0), function(parsedMessage){
            console.log(parsedMessage.host + ' - ' + parsedMessage.message);
        });
    });
    
    server.on("listening", function() {
        var address = server.address();
        console.log("Server now listening at " + 
            address.address + ":" + address.port);
    });
    
    server.bind(514); // Remember ports < 1024 need suid


Author
-------
Squeeks - privacymyass@gmail.com

License
-------
This is free software licensed under the MIT License - see the LICENSE file that
should be included with this package.

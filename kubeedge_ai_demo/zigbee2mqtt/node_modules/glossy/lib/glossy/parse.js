/*
 *    Glossy Parser - Parse incoming syslog messages
 *
 *    Copyright Squeeks <privacymyass@gmail.com>.
 *    This is free software licensed under the MIT License - 
 *    see the LICENSE file that should be included with this package.
 */

/*
 *    These values replace the integers in message that define the facility.
 */
var FacilityIndex = [
    'kern',     // kernel messages
    'user',     // user-level messages
    'mail',     // mail system
    'daemon',   // system daemons
    'auth',     // security/authorization messages
    'syslog',   // messages generated internally by syslogd
    'lpr',      // line printer subsystem
    'news',     // network news subsystem
    'uucp',     // UUCP subsystem
    'clock',    // clock daemon
    'sec',      // security/authorization messages
    'ftp',      // FTP daemon
    'ntp',      // NTP subsystem
    'audit',    // log audit
    'alert',    // log alert
    'clock',    // clock daemon (note 2)
    'local0',   // local use 0  (local0)
    'local1',   // local use 1  (local1)
    'local2',   // local use 2  (local2)
    'local3',   // local use 3  (local3)
    'local4',   // local use 4  (local4)
    'local5',   // local use 5  (local5)
    'local6',   // local use 6  (local6)
    'local7'    // local use 7  (local7)
];

// Note 1 - Various operating systems have been found to utilize
//           Facilities 4, 10, 13 and 14 for security/authorization,
//           audit, and alert messages which seem to be similar. 

// Note 2 - Various operating systems have been found to utilize
//           both Facilities 9 and 15 for clock (cron/at) messages.

/*
 *    These values replace the integers in message that define the severity.
 */
var SeverityIndex = [
    'emerg',    // Emergency: system is unusable
    'alert',    // Alert: action must be taken immediately
    'crit',     // Critical: critical conditions
    'err',      // Error: error conditions
    'warn',     // Warning: warning conditions
    'notice',   // Notice: normal but significant condition
    'info',     // Informational: informational messages
    'debug'     // Debug: debug-level messages
];

/*
 *    Defines the range matching BSD style months to integers.
 */
var BSDDateIndex = {
    'Jan': 0,
    'Feb': 1,
    'Mar': 2,
    'Apr': 3,
    'May': 4,
    'Jun': 5,
    'Jul': 6,
    'Aug': 7,
    'Sep': 8,
    'Oct': 9,
    'Nov': 10,
    'Dec': 11
};

// These values match the hasing algorithm values as defined in RFC 5848
var signedBlockValues = {

    // Section 4.2.1
    hashAlgorithm: [
        null,
        'SHA1',
        'SHA256'
    ],

    // Section 5.2.1
    keyBlobType: {
        'C': 'PKIX Certificate',
        'P': 'OpenPGP KeyID',
        'K': 'Public Key',
        'N': 'No key information',
        'U': 'Unknown'
    }

};

var GlossyParser = function() {};

/*  
 *  Parse the raw message received.
 *
 *  @param {String/Buffer} rawMessage Raw message received from socket
 *  @param {Function} callback Callback to run after parse is complete
 *  @return {Object} map containing all successfully parsed data.
 */
GlossyParser.prototype.parse = function(rawMessage, callback) {

    // Are you node.js? Is this a Buffer?
    if(typeof Buffer == 'function' && Buffer.isBuffer(rawMessage)) {
        rawMessage = rawMessage.toString('utf8', 0);
    } else if(typeof rawMessage != 'string') {
        return rawMessage;
    }

    // Always return the original message
    var parsedMessage = {
        originalMessage: rawMessage
    };
    
    var segments = rawMessage.split(' ');
    if(segments.length < 2) return parsedMessage;
    var priKeys = this.decodePri(segments[0]);
    if(priKeys) {
        for (var key in priKeys) parsedMessage[key] = priKeys[key];
    }

    var timeStamp;
    //TODO Could our detection between 3164/5424 be improved?
    if(segments[0].match(/^(<\d+>\d)$/))  {
        segments.shift(); // Shift the prival off
        timeStamp             = segments.shift();
        parsedMessage.type    = 'RFC5424';
        parsedMessage.time    = this.parseTimeStamp(timeStamp);
        parsedMessage.host    = this.decideValue(segments.shift());
        parsedMessage.appName = this.decideValue(segments.shift());
        parsedMessage.pid     = this.decideValue(segments.shift());
        parsedMessage.msgID   = this.decideValue(segments.shift());

        if(segments[0] !== '-') {
            var spliceMarker = 0;
            for (i = segments.length -1; i > -1; i--) {
                if(segments[i].substr(-1) === ']'){
                    spliceMarker = i;
                    spliceMarker++;
                    break;
                }
            }
            if(spliceMarker !== 0) {
                var sd = segments.splice(0, spliceMarker).join(' ');
                parsedMessage.structuredData = this.parseStructure(sd);

                if(parsedMessage.structuredData.ssign) {
                    parsedMessage.structuredData.signedBlock = 
                        this.parseSignedBlock(parsedMessage.structuredData);
                } else if(parsedMessage.structuredData['ssign-cert']) {
                    parsedMessage.structuredData.signedBlock = 
                        this.parseSignedCertificate(parsedMessage.structuredData);
                }

            }
        } else {
            segments.shift(); // Shift the SD marker off
        }
        parsedMessage.message = segments.join(' ');

    } else if (segments[0].match(/^(<\d+>\d+:)$/)) {
        parsedMessage.type    = 'RFC3164';
        timeStamp             = segments.splice(0,1).join(' ').replace(/^(<\d+>)/,'');
        parsedMessage.time    = this.parseBsdTime(timeStamp);
        parsedMessage.message = segments.join(' ');

    } else if(segments[0].match(/^(<\d+>\w+)/)) {
        parsedMessage.type    = 'RFC3164';
        if (segments[1] === '') segments.splice(1,1);
        timeStamp             = segments.splice(0,3).join(' ').replace(/^(<\d+>)/,'');
        parsedMessage.time    = this.parseBsdTime(timeStamp);
        parsedMessage.host    = segments.shift();
        parsedMessage.message = segments.join(' ');
    }

    if(callback) {
        callback(parsedMessage);
    } else {
        return parsedMessage;
    }

};

/*
 *  RFC5424 messages are supposed to specify '-' as the null value
 *  @param {String} a section from an RFC5424 message
 *  @return {Boolean/String} null if string is entirely '-', or the original value
 */
GlossyParser.prototype.decideValue = function(value) {
    return value === '-' ? null : value;
};

/*
 *  Parses the PRI value from the start of message
 *
 *  @param {String} message Supplied raw primary value and version
 *  @return {Object} Returns object containing Facility, Severity and Version
 *      if correctly parsed, empty values on failure.
 */
GlossyParser.prototype.decodePri = function(message) {
    if(typeof message != 'string') return;

    var privalMatch = message.match(/^<(\d+)>/);
    if(!privalMatch) return false;

    var returnVal = {
        prival: parseInt(privalMatch[1], 10)
    };

    if(privalMatch[2]) returnVal.versio = parseInt(privalMatch[2], 10);

    if(returnVal.prival && returnVal.prival >= 0 && returnVal.prival <= 191) {
    
        returnVal.facilityID = parseInt(returnVal.prival / 8, 10);
        returnVal.severityID = returnVal.prival - (returnVal.facilityID * 8);

        if(returnVal.facilityID < 24 && returnVal.severityID < 8) {
            returnVal.facility = FacilityIndex[returnVal.facilityID];
            returnVal.severity = SeverityIndex[returnVal.severityID];
        }
    } else if(returnVal.prival >= 191) {
        return false;
    }

    return returnVal;
};


/*
 *  Attempts to parse a given timestamp
 *  @param {String} timeStamp Supplied timestamp, should only be the timestamp, 
 *      not the entire message
 *  @return {Object} Date object on success
 */
GlossyParser.prototype.parseTimeStamp = function(timeStamp) {
    
    if(typeof timeStamp != 'string') return;
    var parsedTime;

    parsedTime = this.parse8601(timeStamp);
    if(parsedTime) return parsedTime;

    parsedTime = this.parseRfc3339(timeStamp);
    if(parsedTime) return parsedTime;

    parsedTime = this.parseBsdTime(timeStamp);
    if(parsedTime) return parsedTime;

    return parsedTime;

};

/*
 *  Parse RFC3339 style timestamps
 *  @param {String} timeStamp
 *  @return {Date/false} Timestamp, if parsed correctly
 *  @see http://blog.toppingdesign.com/2009/08/13/fast-rfc-3339-date-processing-in-javascript/
 */
GlossyParser.prototype.parseRfc3339 = function(timeStamp){
    var utcOffset, offsetSplitChar, offsetString,
        offsetMultiplier = 1,
        dateTime = timeStamp.split("T");
        if(dateTime.length < 2) return false;

        var date    = dateTime[0].split("-"),
        time        = dateTime[1].split(":"),
        offsetField = time[time.length - 1];

    offsetFieldIdentifier = offsetField.charAt(offsetField.length - 1);
    if (offsetFieldIdentifier === "Z") {
        utcOffset = 0;
        time[time.length - 1] = offsetField.substr(0, offsetField.length - 2);
    } else {
        if (offsetField[offsetField.length - 1].indexOf("+") != -1) {
            offsetSplitChar = "+";
            offsetMultiplier = 1;
        } else {
            offsetSplitChar = "-";
            offsetMultiplier = -1;
        }

        offsetString = offsetField.split(offsetSplitChar);
        if(offsetString.length < 2) return false;
        time[(time.length - 1)] = offsetString[0];
        offsetString = offsetString[1].split(":");
        utcOffset    = (offsetString[0] * 60) + offsetString[1];
        utcOffset    = utcOffset * 60 * 1000;
    }
               
    var parsedTime = new Date(Date.UTC(date[0], date[1] - 1, date[2], time[0], time[1], time[2]) + (utcOffset * offsetMultiplier ));
    return parsedTime;
};

/*
 *  Parse "BSD style" timestamps, as defined in RFC3164
 *  @param {String} timeStamp
 *  @return {Date/false} Timestamp, if parsed correctly
 */
GlossyParser.prototype.parseBsdTime = function(timeStamp) {
    var parsedTime;
    var d = timeStamp.match(/(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)\s+(\d{1,2})\s+(\d{2}):(\d{2}):(\d{2})/);
    if(d) {
        // Years are absent from the specification, use this year
        currDate   = new Date();
        parsedTime = new Date(
            currDate.getUTCFullYear(), 
            BSDDateIndex[ d[1] ], 
            d[2], 
            d[3], 
            d[4], 
            d[5]);
    }

    return parsedTime;
};

/*
 *  Parse ISO 8601 timestamps
 *  @param {String} timeStamp
 *  @return {Object/false} Timestamp, if successfully parsed
 */
GlossyParser.prototype.parse8601 = function(timeStamp) {
    var parsedTime = new Date(Date.parse(timeStamp));
    if(parsedTime.toString() === 'Invalid Date') return; //FIXME not the best
    return parsedTime;
};


/*
 *  Parse the structured data out of RFC5424 messages
 *  @param {String} msg The STRUCTURED-DATA section
 *  @return {Object} sdStructure parsed structure
 */
GlossyParser.prototype.parseStructure = function(msg) {
    var sdStructure = { };

    var state   = 0,
        ignore  = false,
        sdId    = '',
        sdParam = '',
        sdValue = '';

    /*
     * Build the structure using a horrible FSM.
     * The states we cycle are as following:
     *   0 1    2       34       20
     *     [sdID sdParam="sdValue"]
     */
    for(var i = 0; i < msg.length; i++) {
        var c = msg[i];
        switch(state) {
            case 0:  // SD-ELEMENT
                state = (c === '[') ? 1 : 0;
                break;
            case 1: // SD-ID
                if(c != ' ') {
                    sdId += c;
                } else {
                    sdStructure[sdId] = {};
                    state = 2;
                }
                break;
            case 2: // SD-PARAM
                if(c === '=') {
                    sdStructure[sdId][sdParam] = '';
                    state = 3;
                } else if(c === ']') {
                    sdId  = '';
                    state = 0;
                } else if(c != ' '){
                    sdParam += c;
                }
                break;
            case 3: // SD-PARAM/SD-VALUE
                state = c === '"' ? 4 : null; // FIXME Handle rubbish better
                break;
            case 4: // SD-VALUE
                if(c === '\\' && !ignore) {
                    ignore = true;
                } else if(c === '"' && !ignore) {
                    sdStructure[sdId][sdParam] = sdValue;
                    sdParam = '', sdValue = '';
                    state = 2;
                } else {
                    sdValue += c;
                    ignore = false;
                }
                break;
            default:
                break;
        }
    }
    return sdStructure;
};


/*
 *  Make sense of signed block messages
 *  @param {Object} block the parsed structured data containing signed data
 *  @return {Object} validatedBlock translated and named values, binary
 *      elements will be Buffer objects, if available
 */
GlossyParser.prototype.parseSignedBlock = function(block) {

    if(typeof block != 'object') return false;

    var signedBlock    = { };
    var validatedBlock = { };
    // Figure out where in the object the keys live...
    if(block.structuredData && block.structuredData.ssign) {
        signedBlock = block.structuredData.ssign;
    } else if(block.ssign) {
        signedBlock = block.ssign;
    } else if(block.VER) {
        signedBlock = block;
    } else {
        return false;
    }

    var versionMatch = signedBlock.VER.match(/^(\d{2})(\d|\w)(\d)$/);
    if(versionMatch !== null) {
        validatedBlock.version        = versionMatch[1];
        validatedBlock.hashAlgorithm  = parseInt(versionMatch[2], 10);
        validatedBlock.hashAlgoString = signedBlockValues.hashAlgorithm[validatedBlock.hashAlgorithm];
        validatedBlock.sigScheme      = parseInt(versionMatch[3], 10);
    }

    validatedBlock.rebootSessionID   = parseInt(signedBlock.RSID, 10);
    validatedBlock.signatureGroup    = parseInt(signedBlock.SG, 10);
    validatedBlock.signaturePriority = parseInt(signedBlock.SPRI, 10);
    validatedBlock.globalBlockCount  = parseInt(signedBlock.GBC, 10);
    validatedBlock.firstMsgNumber    = parseInt(signedBlock.FMN, 10);
    validatedBlock.msgCount          = parseInt(signedBlock.CNT, 10);
    validatedBlock.hashBlock         = signedBlock.HB.split(/\s/);

    // Check to see if we're in node or have a Buffer type
    if(typeof Buffer == 'function') {
        for(var hash in validatedBlock.hashBlock) {
            validatedBlock.hashBlock[hash] = new Buffer(
                validatedBlock.hashBlock[hash], encoding='base64'); 
        }
        validatedBlock.thisSignature = new Buffer(
            signedBlock.SIGN, encoding='base64');
    } else {
        validatedBlock.thisSignature = signedBlock.SIGN;
    }

    return validatedBlock;
    
};


/*
 *  Make sense of signed certificate messages
 *  @param {Object} block the parsed structured data containing signed data
 *  @return {Object} validatedBlock translated and named values, binary
 *      elements will be Buffer objects, if available
 */
GlossyParser.prototype.parseSignedCertificate = function(block) {

    if(typeof block != 'object') return false;

    var signedBlock    = { };
    var validatedBlock = { };
    // Figure out where in the object the keys live...
    if(block.structuredData && block.structuredData['ssign-cert']) {
        signedBlock = block.structuredData['ssign-cert'];
    } else if(block['ssign-cert']) {
        signedBlock = block['ssign-cert'];
    } else if(block.VER) {
        signedBlock = block;
    } else {
        return false;
    }

    var versionMatch = signedBlock.VER.match(/^(\d{2})(\d|\w)(\d)$/);
    if(versionMatch !== null) {
        validatedBlock.version        = versionMatch[1];
        validatedBlock.hashAlgorithm  = parseInt(versionMatch[2], 10);
        validatedBlock.hashAlgoString = signedBlockValues.hashAlgorithm[validatedBlock.hashAlgorithm];
        validatedBlock.sigScheme      = parseInt(versionMatch[3], 10);
    }

    validatedBlock.rebootSessionID     = parseInt(signedBlock.RSID, 10);
    validatedBlock.signatureGroup      = parseInt(signedBlock.SG, 10);
    validatedBlock.signaturePriority   = parseInt(signedBlock.SPRI, 10);
    validatedBlock.totalPayloadLength  = parseInt(signedBlock.TPBL, 10);
    validatedBlock.payloadIndex        = parseInt(signedBlock.INDEX, 10);
    validatedBlock.fragmentLength      = parseInt(signedBlock.FLEN, 10);

    var payloadFragment             = signedBlock.FRAG.split(/\s/);
    validatedBlock.payloadTimestamp = this.parseTimeStamp(payloadFragment[0]);
    validatedBlock.payloadType      = payloadFragment[1];
    validatedBlock.payloadName      = signedBlockValues.keyBlobType[payloadFragment[1]];

    if(typeof Buffer == 'function') {
        validatedBlock.keyBlob = new Buffer(
            payloadFragment[2], encoding='base64');
        validatedBlock.thisSignature = new Buffer(
            signedBlock.SIGN, encoding='base64');
    } else {
        validatedBlock.keyBlob       = payloadFragment[2];
        validatedBlock.thisSignature = signedBlock.SIGN;
    }

    return validatedBlock;

};


if(typeof module == 'object') {
    module.exports = new GlossyParser();
}

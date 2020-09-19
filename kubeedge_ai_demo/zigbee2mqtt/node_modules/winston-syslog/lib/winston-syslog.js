/*
 * syslog.js: Transport for logging to a remote syslog consumer
 *
 * (C) 2011 Squeeks and Charlie Robbins
 * MIT LICENCE
 *
 */

const dgram = require('dgram');
const net = require('net');
const secNet = require('tls');
const utils = require('./utils');
const glossy = require('glossy');
const winston = require('winston');
const Transport = require('winston-transport');
const { MESSAGE, LEVEL } = require('triple-beam');

const _noop = () => {};

// Ensure we have the correct winston here.
if (Number(winston.version.split('.')[0]) < 3) {
  throw new Error('Winston-syslog requires winston >= 3.0.0');
}

const levels = Object.keys({
  debug: 0,
  info: 1,
  notice: 2,
  warning: 3,
  warn: 3,
  error: 4,
  err: 4,
  crit: 5,
  alert: 6,
  emerg: 7
});

//
// ### function Syslog (options)
// #### @options {Object} Options for this instance.
// Constructor function for the Syslog Transport capable of sending
// RFC 3164 and RFC 5424 compliant messages.
//
class Syslog extends Transport {
  //
  // Expose the name of this Transport on the prototype
  //
  get name() {
    return 'syslog';
  }

  constructor(options = {}) {
    //
    // Inherit from `winston-transport`.
    //
    super(options);

    //
    // Setup connection state
    //
    this.connected = false;
    this.congested = false;
    this.retries = 0;
    this.queue = [];
    this.inFlight = 0;

    //
    // Merge the options for the target Syslog server.
    //
    this.setOptions(options);

    //
    // Setup our Syslog and network members for later use.
    //
    this.socket = null;
    var Producer = options.customProducer || glossy.Produce;
    this.producer = new Producer({
      type: this.type,
      appName: this.appName,
      pid: this.pid,
      facility: this.facility
    });
  }

  setOptions(options) {
    this.host = options.host || 'localhost';
    this.port = options.port || 514;
    this.path = options.path || null;
    this.protocol = options.protocol || 'udp4';
    this.protocolOptions = options.protocolOptions || {};
    this.endOfLine = options.eol;

    this.parseProtocol(this.protocol);

    //
    // Merge the default message options.
    //
    this.localhost =
      typeof options.localhost !== 'undefined'
        ? options.localhost
        : 'localhost';
    this.type = options.type || 'BSD';
    this.facility = options.facility || 'local0';
    this.pid = options.pid || process.pid;
    this.appName = options.appName || options.app_name || process.title;
  }

  parseProtocol(protocol = this.protocol) {
    const parsedProtocol = utils.parseProtocol(protocol);

    this.protocolType = parsedProtocol.type;
    this.protocolFamily = parsedProtocol.family;
    this.isDgram = parsedProtocol.isDgram;

    if (this.protocolType === 'unix' && !this.path) {
      throw new Error('`options.path` is required on unix dgram sockets.');
    }
  }

  //
  // ### function chunkMessage (buffer, callback)
  // #### @buffer {Buffer} Syslog message buffer.
  // #### @callback {function} Continuation to respond to when complete.
  // Syslog messages sent over the UDP transport must be 64KB bytes or less. In
  // order to avoid silent failures messages should be chunked when the buffer
  // to write is larger than the maximum message size.
  //
  chunkMessage(buffer, callback = _noop) {
    if (!this.connected) {
      this.queue.push(buffer);
    } else {
      // Maximum payload for the UDP transport
      // 65535 − 8 bytes UDP header − 20 bytes IP header
      // https://tools.ietf.org/html/rfc5426#section-3.2
      // This makes sense on loopback interfaces with MTU = 65536
      // For non-loopback messages, it's impossible to know in advance
      // the MTU of each interface through which a packet might be routed
      // https://nodejs.org/api/dgram.html
      const MAX_UDP_PAYLOAD = 65507;
      let offset = 0;

      while (offset < buffer.length) {
        this.inFlight++;
        const length =
          offset + MAX_UDP_PAYLOAD > buffer.length
            ? buffer.length - offset
            : MAX_UDP_PAYLOAD;
        this._sendChunk(
          buffer,
          { offset: offset, length: length, port: this.port, host: this.host },
          callback
        );
        offset += length;
      }
    }
  }

  //
  // ### function _sendChunk (buffer, options, callback)
  // #### @buffer {Buffer} Syslog message buffer.
  // #### @options {object} Options for the message send method.
  // #### @callback {function} Continuation to respond to when complete.
  // Sends a single chunk from an oversize UDP buffer.
  //
  _sendChunk(buffer, options, callback) {
    this.socket.send(
      buffer,
      options.offset,
      options.length,
      options.port,
      options.host,
      callback
    );
  }

  //
  // ### function log (info, callback)
  // #### @info {object} All relevant log information
  // #### @callback {function} Continuation to respond to when complete.
  // Core logging method exposed to Winston. Logs the `msg` and optional
  // metadata, `meta`, to the specified `level`.
  //
  log(info, callback) {
    let level = info[LEVEL];
    if (!~levels.indexOf(level)) {
      return callback(
        new Error('Cannot log unknown syslog level: ' + info[LEVEL])
      );
    }
    level = level === 'warn' ? 'warning' : level;
    const output = info[MESSAGE];

    const syslogMsg = this.producer.produce({
      severity: level,
      host: this.localhost,
      date: new Date(),
      message: this.endOfLine ? output + this.endOfLine : output
    });

    //
    // Attempt to connect to the socket
    //
    this.connect(err => {
      if (err) {
        //
        // If there was an error enqueue the message
        //
        this.queue.push(syslogMsg);

        return callback();
      }

      //
      // On any error writing to the socket, enqueue the message
      //
      const onError = logErr => {
        if (logErr) {
          this.queue.push(syslogMsg);
          this.emit('error', logErr);
        }
        this.emit('logged', info);
        this.inFlight--;
      };

      const onCongestion = () => {
        onError(new Error('Congestion Error'));
      };

      const sendDgram = () => {
        const buffer = Buffer.from(syslogMsg);

        if (this.protocolType === 'udp') {
          this.chunkMessage(buffer, onError);
        } else if (this.protocol === 'unix') {
          this.inFlight++;
          this.socket.send(buffer, 0, buffer.length, this.path, onError);
        } else if (this.congested) {
          this.queue.push(syslogMsg);
        } else {
          this.socket.once('congestion', onCongestion);
          this.inFlight++;
          this.socket.send(buffer, e => {
            this.socket.removeListener('congestion', onCongestion);
            onError(e);
          });
        }
      };

      //
      // Write to the `tcp*`, `udp*`, or `unix` socket.
      //
      if (this.isDgram) {
        sendDgram();
      } else {
        this.socket.write(syslogMsg, 'utf8', onError);
      }

      callback(null, true);
    });
  }

  //
  // ### function close ()
  // Closes the socket used by this transport freeing the resource.
  //
  close() {
    const max = 6;
    let attempt = 0;

    const _close = () => {
      if (attempt >= max || (this.queue.length === 0 && this.inFlight <= 0)) {
        if (this.socket) {
          if (this.socket.destroy) {
            // https://nodejs.org/api/net.html#net_socket_destroy_exception
            this.socket.destroy();
          } else if (this.socket.close) {
            // https://nodejs.org/api/dgram.html#dgram_socket_close_callback
            // https://www.npmjs.com/package/unix-dgram#socketclose
            this.socket.close();
          }
        }

        this.emit('closed', this.socket);
      } else {
        attempt++;
        setTimeout(_close, 200 * attempt);
      }
    };
    _close();
  }

  connectDgram(callback) {
    if (this.protocol === 'unix-connect') {
      return this._unixDgramConnect(callback);
    } else if (this.protocol === 'unix') {
      this.socket = require('unix-dgram').createSocket('unix_dgram');
    } else if (!this.socket) {
      // UDP protocol
      const proto = this.protocol === 'udp' ? 'udp4' : this.protocol;
      // https://nodejs.org/api/all.html#dgram_class_dgram_socket
      this.socket = dgram.createSocket({
        type: proto
      })
        .on('listening', () => {
          this.connected = true;
          let msg = this.queue.shift();
          while (msg) {
            this.chunkMessage(msg);
            msg = this.queue.shift();
          }
        });
      this.socket.bind();
    }

    return callback(null);
  }
  //
  // ### function connect (callback)
  // #### @callback {function} Continuation to respond to when complete.
  // Connects to the remote syslog server using `dgram` or `net` depending
  // on the `protocol` for this instance.
  //
  connect(callback) {
    //
    // If the socket already exists then respond
    //
    if (this.socket) {
      return !this.socket.readyState ||
        this.socket.readyState === 'open' ||
        this.socket.connected
        ? callback(null)
        : callback(true);
    }

    //
    // Create the appropriate socket type.
    //
    if (this.isDgram) {
      return this.connectDgram(callback);
    }

    this.socket = /^tls[4|6]?$/.test(this.protocol)
      ? new secNet.TLSSocket()
      : new net.Socket();
    this.socket.setKeepAlive(true);
    this.socket.setNoDelay();

    this.setupEvents();

    const connectConfig = Object.assign({}, this.protocolOptions, {
      host: this.host,
      port: this.port
    });

    if (this.protocolFamily) {
      connectConfig.family = this.protocolFamily;
    }

    this.socket.connect(connectConfig);

    //
    // Indicate to the callee that the socket is not ready. This
    // will enqueue the current message for later.
    //
    callback(true);
  }

  setupEvents() {
    const readyEvent = 'connect';
    //
    // On any error writing to the socket, emit the `logged` event
    // and the `error` event.
    //
    const onError = logErr => {
      if (logErr) {
        this.emit('error', logErr);
      }
      this.emit('logged');
      this.inFlight--;
    };

    //
    // Listen to the appropriate events on the socket that
    // was just created.
    //
    this.socket
      .on(readyEvent, () => {
        //
        // When the socket is ready, write the current queue
        // to it.
        //
        this.socket.write(this.queue.join(''), 'utf8', onError);

        this.emit('logged');
        this.queue = [];
        this.retries = 0;
        this.connected = true;
      })
      .on('error', function () {
        //
        // TODO: Pass this error back up
        //
      })
      .on('close', () => {
        //
        // Attempt to reconnect on lost connection(s), progressively
        // increasing the amount of time between each try.
        //
        const interval = Math.pow(2, this.retries);
        this.connected = false;

        setTimeout(() => {
          this.retries++;
          this.socket.connect(this.port, this.host);
        }, interval * 1000);
      })
      .on('timeout', () => {
        if (this.socket.destroy) {
          // https://nodejs.org/api/net.html#net_socket_settimeout_timeout_callback
          this.socket.destroy();
        } else if (this.socket.close) {
          // https://nodejs.org/api/dgram.html#dgram_socket_close_callback
          // https://www.npmjs.com/package/unix-dgram#socketclose
          this.socket.close();
        }
      });
  }

  _unixDgramConnect(callback) {
    const self = this;

    const flushQueue = () => {
      let sentMsgs = 0;
      this.queue.forEach(msg => {
        const buffer = Buffer.from(msg);

        if (!this.congested) {
          if (this.protocolType === 'udp') {
            this.chunkMessage(buffer, () => {
              ++sentMsgs;
            });
          } else {
            this.socket.send(buffer, function () {
              ++sentMsgs;
            });
          }
        }
      });

      this.queue.splice(0, sentMsgs);
    };

    this.socket = require('unix-dgram').createSocket('unix_dgram');
    this.socket.on('error', err => {
      this.emit('error', err);

      if (err.syscall === 'connect') {
        this.socket.close();
        this.socket = null;
        return callback(err);
      }
      if (err.syscall === 'send') {
        this.socket.close();
        this.socket = null;
      }
    });

    this.socket.on('connect', function () {
      this.on('congestion', () => {
        self.congested = true;
      });

      this.on('writable', () => {
        self.congested = false;
        flushQueue();
      });

      flushQueue();
      callback();
    });

    this.socket.connect(this.path);
  }
}

//
// Define a getter so that `winston.transports.Syslog`
// is available and thus backwards compatible.
//
winston.transports.Syslog = Syslog;

module.exports = {
  Syslog
};

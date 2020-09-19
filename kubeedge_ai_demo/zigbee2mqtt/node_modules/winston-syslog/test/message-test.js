'use strict';

const vows = require('vows');
const assert = require('assert');
const sinon = require('sinon');
const Syslog = require('../lib/winston-syslog.js').Syslog;
const dgram = require('dgram');

const PORT = 11229;
let server;
let transport;
let maxUdpPayload;
let message;
let sentMessage;
let numChunks;

const { MESSAGE, LEVEL } = require('triple-beam');

vows
  .describe('syslog message')
  .addBatch({
    'Opening fake syslog UDP server': {
      'topic': function () {
        const self = this;
        server = dgram.createSocket('udp4');
        server.on('listening', function () {
          // Maximum payload for the UDP transport
          // 65535 − 8 bytes UDP header − 20 bytes IP header
          // https://tools.ietf.org/html/rfc5426#section-3.2
          // This makes sense on loopback interfaces with MTU = 65536
          // For non-loopback messages, it's impossible to know in advance
          // the MTU of each interface through which a packet might be routed
          // https://nodejs.org/api/dgram.html
          maxUdpPayload = 65507;
          self.callback();
        });

        server.bind(PORT);
      },
      'logging an oversize message': {
        'topic': function () {
          // Generate message larger than max UDP message size.
          message = '#'.repeat(maxUdpPayload + 1000);
          transport = new Syslog({
            port: PORT
          });

          sinon.spy(transport, 'chunkMessage');
          sinon.spy(transport, '_sendChunk');

          transport.log({ [LEVEL]: 'debug', [MESSAGE]: message }, function (
            err
          ) {
            assert.ifError(err);
          });

          return null;
        },
        'correct number of chunks sent': function () {
          assert(transport.chunkMessage.calledTwice);

          sentMessage = transport.chunkMessage.getCall(0).args[0];
          numChunks = Math.ceil(sentMessage.length / maxUdpPayload);
          assert.equal(numChunks, transport._sendChunk.callCount);
        },
        'correct chunks sent': function () {
          let offset = 0;
          let i = 0;

          sentMessage = transport.chunkMessage.getCall(0).args[0];
          while (offset < sentMessage.length) {
            const length =
              offset + maxUdpPayload > sentMessage.length
                ? sentMessage.length - offset
                : maxUdpPayload;
            const buffer = Buffer.from(sentMessage);
            const options = {
              offset: offset,
              length: length,
              port: transport.port,
              host: transport.host
            };

            assert(transport._sendChunk.getCall(i).calledWith(buffer, options));

            offset += length;
            i++;
          }

          transport.close();
        }
      },
      'teardown': function () {
        transport.chunkMessage.restore();
        transport._sendChunk.restore();
        server.close();
      }
    }
  })
  .export(module);

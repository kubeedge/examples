/* eslint new-cap: ["error", { "newIsCapExceptions": ["createLogger"] }] */
/*
 * syslog-test.js: Tests for instances of the Syslog transport
 *
 * (C) 2010 Charlie Robbins
 * MIT LICENSE
 *
 */

const vows = require('vows');
const assert = require('assert');
const winston = require('winston');
const Syslog = require('../lib/winston-syslog').Syslog;

function assertSyslog(transport) {
  assert.instanceOf(transport, Syslog);
  assert.isFunction(transport.log);
  assert.isFunction(transport.connect);
}

function closeTopicInfo() {
  const transport = new winston.transports.Syslog();
  const logger = new winston.createLogger({ transports: [transport] });

  logger.log('info', 'Test message to actually use socket');
  logger.remove(transport);

  return transport;
}

function closeTopicDebug() {
  const transport = new winston.transports.Syslog();
  const logger = new winston.createLogger({ transports: [transport] });

  logger.log('debug', 'Test message to actually use socket');
  logger.remove(transport);

  return transport;
}

const transport = new Syslog();

vows.describe('winston-syslog').addBatch({
  'An instance of the Syslog Transport': {
    'should have the proper methods defined': function () {
      assertSyslog(transport);
    },
    'teardown': function () {
      transport.close();
    },
    'on close after not really writing': {
      topic: closeTopicDebug,
      on: {
        closed: {
          'closes the socket': function (socket) {
            assert.isNull(socket);
          }
        }
      }
    },
    'on close after really writing': {
      topic: closeTopicInfo,
      on: {
        closed: {
          'closes the socket': function (socket) {
            assert.isNull(socket._handle);
          }
        }
      }
    },
    'localhost option': {
      'should default to localhost': function () {
        const transportLocal = new winston.transports.Syslog();
        assert.equal(transportLocal.localhost, 'localhost');
        transportLocal.close();
      },
      'should accept other falsy entries as valid': function () {
        let transportNotLocal = new winston.transports.Syslog({ localhost: null });
        assert.isNull(transportNotLocal.localhost);
        transportNotLocal.close();
        transportNotLocal = new winston.transports.Syslog({ localhost: false });
        assert.equal(transportNotLocal.localhost, false);
        transportNotLocal.close();
      }
    },
    'adding / removing transport to syslog': {
      'should just work': function () {
        winston.add(new winston.transports.Syslog());
        winston.remove(new winston.transports.Syslog());
        winston.add(new winston.transports.Syslog());
        winston.remove(new winston.transports.Syslog());
      }
    }
  }
}).export(module);

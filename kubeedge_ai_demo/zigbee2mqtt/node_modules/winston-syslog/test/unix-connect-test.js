/* eslint no-sync: "off" */

const fs = require('fs');
const vows = require('vows');
const assert = require('assert');
const unix = require('unix-dgram');
const parser = require('glossy').Parse;
const Syslog = require('../lib/winston-syslog').Syslog;

const { MESSAGE, LEVEL } = require('triple-beam');

const SOCKNAME = '/tmp/unix_dgram.sock';

const transport = new Syslog({
  protocol: 'unix-connect',
  path: SOCKNAME
});

try {
  fs.unlinkSync(SOCKNAME);
} catch (e) {
  /* swallow */
}

let times = 0;
let server;

vows.describe('unix-connect').addBatch({
  'Trying to log to a non-existant log server': {
    'topic': function () {
      const self = this;
      transport.once('error', function (err) {
        self.callback(null, err);
      });

      transport.log({ [LEVEL]: 'debug', [MESSAGE]: `data${++times}` }, function (err) {
        assert.equal(null, err);
        assert.equal(transport.queue.length, 1);
      });
    },
    'should enqueue the log message': function (err) {
      assert(err);
      assert.equal(err.syscall, 'connect');
    }
  }
}).addBatch({
  'Logging when log server is up': {
    'topic': function () {
      const self = this;
      let n = 0;
      server = unix.createSocket('unix_dgram', function (buf) {
        parser.parse(buf, function (d) {
          ++n;
          assert(n <= 2);
          assert.equal(d.message, 'node[' + process.pid + ']: data' + n);
          assert.equal(d.severity, 'debug');
          if (n === 2) {
            self.callback();
          }
        });
      });

      server.bind(SOCKNAME);
      transport.log({ [LEVEL]: 'debug', [MESSAGE]: `data${++times}` }, function (err) {
        assert.ifError(err);
      });
    },
    'should print both the enqueed and the new msg': function (err) {
      assert.ifError(err);
    }
  }
}).addBatch({
  'Logging if server goes down again': {
    'topic': function () {
      const self = this;
      transport.once('error', function (err) {
        self.callback(null, err);
      });

      server.close();

      transport.log({ [LEVEL]: 'debug', [MESSAGE]: `data${++times}` }, function (err) {
        assert.ifError(err);
        assert.equal(transport.queue.length, 1);
      });
    },
    'should enqueue the log message': function (err) {
      assert(err);
      assert.equal(err.syscall, 'send');
      transport.close();
    }
  }
}).addBatch({
  'Logging works if server comes up again': {
    'topic': function () {
      const self = this;
      transport.once('error', function (err) {
        // Ignore error -- server hasn't come up yet, that's fine/expected
        assert(err);
        assert.equal(err.syscall, 'send');
      });
      let n = 2;
      try {
        fs.unlinkSync(SOCKNAME);
      } catch (e) {
        /* swallow */
      }
      server = unix.createSocket('unix_dgram', function (buf) {
        parser.parse(buf, function (d) {
          ++n;
          assert(n <= 4);
          assert.equal(d.message, 'node[' + process.pid + ']: data' + n);
          if (n === 4) {
            self.callback();
          }
        });
      });

      server.bind(SOCKNAME);
      transport.log({ [LEVEL]: 'debug', [MESSAGE]: `data${++times}` }, function (err) {
        assert.ifError(err);
      });
      return null;
    },
    'should print both the enqueed and the new msg': function (err) {
      assert.ifError(err);
      server.close();
      return null;
    }
  }

}).export(module);

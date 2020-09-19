'use strict';
/**
 * @file mkdir-recursive main
 * @module mkdir-recursive
 * @version 0.3.0
 * @author hex7c0 <hex7c0@gmail.com>
 * @copyright hex7c0 2015
 * @license GPLv3
 */

/*
 * initialize module
 */
var fs = require('fs');
var path = require('path');

/*
 * exports
 */
/**
 * make main. Check README.md
 * 
 * @exports mkdir
 * @function mkdir
 * @param {String} root - pathname
 * @param {Number} mode - directories mode, see Node documentation
 * @param {Function} callback - next callback
 */
function mkdir(root, mode, callback) {

  if (typeof mode === 'function') {
    var callback = mode;
    var mode = null;
  }
  if (typeof root !== 'string') {
    throw new Error('missing root');
  } else if (typeof callback !== 'function') {
    throw new Error('missing callback');
  }

  var chunks = root.split(path.sep); // split in chunks
  var chunk;
  if (path.isAbsolute(root) === true) { // build from absolute path
    chunk = chunks.shift(); // remove "/" or C:/
    if (!chunk) { // add "/"
      chunk = path.sep;
    }
  } else {
    chunk = path.resolve(); // build with relative path
  }

  return mkdirRecursive(chunk, chunks, mode, callback);
}
module.exports.mkdir = mkdir;

/**
 * makeSync main. Check README.md
 * 
 * @exports mkdirSync
 * @function mkdirSync
 * @param {String} root - pathname
 * @param {Number} mode - directories mode, see Node documentation
 * @return [{Object}]
 */
function mkdirSync(root, mode) {

  if (typeof root !== 'string') {
    throw new Error('missing root');
  }

  var chunks = root.split(path.sep); // split in chunks
  var chunk;
  if (path.isAbsolute(root) === true) { // build from absolute path
    chunk = chunks.shift(); // remove "/" or C:/
    if (!chunk) { // add "/"
      chunk = path.sep;
    }
  } else {
    chunk = path.resolve(); // build with relative path
  }

  return mkdirSyncRecursive(chunk, chunks, mode);
}
module.exports.mkdirSync = mkdirSync;

/**
 * remove main. Check README.md
 * 
 * @exports rmdir
 * @function rmdir
 * @param {String} root - pathname
 * @param {Function} callback - next callback
 */
function rmdir(root, callback) {

  if (typeof root !== 'string') {
    throw new Error('missing root');
  } else if (typeof callback !== 'function') {
    throw new Error('missing callback');
  }

  var chunks = root.split(path.sep); // split in chunks
  var chunk = path.resolve(root); // build absolute path
  // remove "/" from head and tail
  if (chunks[0] === '') {
    chunks.shift();
  }
  if (chunks[chunks.length - 1] === '') {
    chunks.pop();
  }

  return rmdirRecursive(chunk, chunks, callback);
}
module.exports.rmdir = rmdir;

/**
 * removeSync main. Check README.md
 * 
 * @exports rmdirSync
 * @function rmdirSync
 * @param {String} root - pathname
 * @return [{Object}]
 */
function rmdirSync(root) {

  if (typeof root !== 'string') {
    throw new Error('missing root');
  }

  var chunks = root.split(path.sep); // split in chunks
  var chunk = path.resolve(root); // build absolute path
  // remove "/" from head and tail
  if (chunks[0] === '') {
    chunks.shift();
  }
  if (chunks[chunks.length - 1] === '') {
    chunks.pop();
  }

  return rmdirSyncRecursive(chunk, chunks);
}
module.exports.rmdirSync = rmdirSync;

/*
 * functions
 */
/**
 * make directory recursively
 * 
 * @function mkdirRecursive
 * @param {String} root - absolute root where append chunks
 * @param {Array} chunks - directories chunks
 * @param {Number} mode - directories mode, see Node documentation
 * @param {Function} callback - next callback
 */
function mkdirRecursive(root, chunks, mode, callback) {

  var chunk = chunks.shift();
  if (!chunk) {
    return callback(null);
  }
  var root = path.join(root, chunk);

  return fs.exists(root, function(exists) {

    if (exists === true) { // already done
      return mkdirRecursive(root, chunks, mode, callback);
    }
    return fs.mkdir(root, mode, function(err) {

      if (err) {
        return callback(err);
      }
      return mkdirRecursive(root, chunks, mode, callback); // let's magic
    });
  });
}

/**
 * make directory recursively. Sync version
 * 
 * @function mkdirSyncRecursive
 * @param {String} root - absolute root where append chunks
 * @param {Array} chunks - directories chunks
 * @param {Number} mode - directories mode, see Node documentation
 * @return [{Object}]
 */
function mkdirSyncRecursive(root, chunks, mode) {

  var chunk = chunks.shift();
  if (!chunk) {
    return;
  }
  var root = path.join(root, chunk);

  if (fs.existsSync(root) === true) { // already done
    return mkdirSyncRecursive(root, chunks, mode);
  }
  var err = fs.mkdirSync(root, mode);
  return err ? err : mkdirSyncRecursive(root, chunks, mode); // let's magic
}

/**
 * remove directory recursively
 * 
 * @function rmdirRecursive
 * @param {String} root - absolute root where take chunks
 * @param {Array} chunks - directories chunks
 * @param {Function} callback - next callback
 */
function rmdirRecursive(root, chunks, callback) {

  var chunk = chunks.pop();
  if (!chunk) {
    return callback(null);
  }
  var pathname = path.join(root, '..'); // backtrack

  return fs.exists(root, function(exists) {

    if (exists === false) { // already done
      return rmdirRecursive(root, chunks, callback);
    }
    return fs.rmdir(root, function(err) {

      if (err) {
        return callback(err);
      }
      return rmdirRecursive(pathname, chunks, callback); // let's magic
    });
  });
}

/**
 * remove directory recursively. Sync version
 * 
 * @function rmdirRecursive
 * @param {String} root - absolute root where take chunks
 * @param {Array} chunks - directories chunks
 * @return [{Object}]
 */
function rmdirSyncRecursive(root, chunks) {

  var chunk = chunks.pop();
  if (!chunk) {
    return;
  }
  var pathname = path.join(root, '..'); // backtrack

  if (fs.existsSync(root) === false) { // already done
    return rmdirSyncRecursive(root, chunks);
  }
  var err = fs.rmdirSync(root);
  return err ? err : rmdirSyncRecursive(pathname, chunks); // let's magic
}

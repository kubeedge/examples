# [mkdir-recursive](https://github.com/hex7c0/mkdir-recursive)

[![NPM version](https://img.shields.io/npm/v/mkdir-recursive.svg)](https://www.npmjs.com/package/mkdir-recursive)
[![Linux Status](https://img.shields.io/travis/hex7c0/mkdir-recursive.svg?label=linux)](https://travis-ci.org/hex7c0/mkdir-recursive)
[![Windows Status](https://img.shields.io/appveyor/ci/hex7c0/mkdir-recursive.svg?label=windows)](https://ci.appveyor.com/project/hex7c0/mkdir-recursive)
[![Dependency Status](https://img.shields.io/david/hex7c0/mkdir-recursive.svg)](https://david-dm.org/hex7c0/mkdir-recursive)
[![Coveralls](https://img.shields.io/coveralls/hex7c0/mkdir-recursive.svg)](https://coveralls.io/r/hex7c0/mkdir-recursive)

make/remove (asynchronous/synchronous) directories recursively

## Installation

Install through NPM

```bash
npm install mkdir-recursive
```
or
```bash
git clone git://github.com/hex7c0/mkdir-recursive.git
```

## API

make 3 directories asynchronous recursively
```js
var fx = require('mkdir-recursive');

fx.mkdir('foo/bar/1', function(err) {

  console.log('done');
});
```

### mkdir(path [, mode], callback) [Node Doc](http://nodejs.org/api/fs.html#fs_fs_mkdir_path_mode_callback)

 - `path` - **String** Dir pathname *(default "required")*
 - `[mode]`- **Number** Scrivi *(default "0777")*
 - `callback` - **Function** Next callback when task is complete *(default "required")*

### mkdirSync(path [, mode]) [Node Doc](http://nodejs.org/api/fs.html#fs_fs_mkdirsync_path_mode)

 - `path` - **String** Dir pathname *(default "required")*
 - `[mode]`- **Number** Scrivi *(default "0777")*

### rmdir(path, callback) [Node Doc](http://nodejs.org/api/fs.html#fs_fs_rmdir_path_callback)

 - `path` - **String** Dir pathname *(default "required")*
 - `callback` - **Function** Next callback when task is complete *(default "required")*

### rmdirSync(path) [Node Doc](http://nodejs.org/api/fs.html#fs_fs_rmdirsync_path)

 - `path` - **String** Dir pathname *(default "required")*

## Examples

Take a look at my [examples](examples)

### [License GPLv3](LICENSE)

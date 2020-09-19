[typedoc](https://github.com/TypeStrong/typedoc) plugin to set custom source file URL links.

typedoc prints a *Defined in* statement showing the source file and line for all definitions. For projects hosted on GitHub this statement will automatically link to the source file.

This plugin allows to create links to files hosted on other platforms and sites like Bitbucket, GitLab or any custom site. It adds a `#L` anchor to the URL, linking to any specific line.

# Installation

    npm install --save-dev typedoc-plugin-sourcefile-url
    
typedoc will automatically detect and load the plugin from `node_modules`.

# Usage

#### Simple Prefix

    typedoc --sourcefile-url-prefix "https://www.your-repository.org/"
    
The `--sourcefile-url-prefix` option will create URLs by prefixing the given parameter in front of each source file.

*Defined in* `src/testfile.ts` will link to `https://www.your-repository.org/src/testfile.ts`.


#### Advanced Mappings

Sometimes more complex URL rules may be required. For example when grouping documentation of multiple repositories into one documentation page.

Advanced mappings are described in a JSON file.

    typedoc --sourcefile-url-map your-sourcefile-map.json
    
The `your-sourcefile-map.json` structure is: 

  
    [
      {
        "pattern": "^modules/module-one",
        "replace": "https://www.your-repository.org/module-one/"
      },     
      {
        "pattern": "^",
        "replace": "https://www.your-repository.org/main-project/"
      }
    ]

`pattern` is a regular expression (without enclosing slashes). Each *Defined in* statement is matched against the `pattern`. On match the `pattern` is replaced with the string from `replace` to create the URL.

There can be one or more mapping rules. For each *Defined in* only the first rule that matches is applied. In the above example the last rule would match all source files that did not start with `modules/module-one`. This compares to the *Simple Prefix* option.

---

The options are mutually exclusive. It is not possible to use `--sourcefile-url-prefix` and `--sourcefile-url-map` at the same time.

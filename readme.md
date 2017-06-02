## Envreplace

A small, statically linkable cli utility written in GoLang that replaces strings in files with the value of env vars.

The primary usage is for dockerizing existing applications. It's useful when you have an application that is looking somewhere in the filesystem for its configuration, and you'd prefer to pass that configuration as an environment variable.

It's roughly similar in functionality to [my node.js module with the same name](https://www.npmjs.com/package/envreplace), but since it can be statically linked and precompiled with no dependencies, it's much better suited for use in a docker container where you might not want a whole node.js install.

### A simple example

If you have a config file `/usr/src/conf/myconfig.properties`
```
my.property.foobar=#FOOBAR#
```
That you want to end up in `/usr/local/tomcat/conf` with its value `#FOOBAR#` taken from the environment variable `FOOBAR`

Then you'd call `envreplace` like this:
```bash
FOOBAR="Some value"
envreplace /usr/src/conf/myconfig.properties /usr/local/tomcat/conf/myconfig.properties
```
And you'd end up with the `/usr/local/tomcat/conf/myconfig.properties`
```
my.property.foobar=Some value
```

### Multiple files at once

If you have two files in the `/usr/src/conf/` folder that you want replaced, you can do multiple at once by creating an envreplace configuration json file.

The `files` entry in the JSON object a map of source files to destination files to be processed by `envreplace`.
You could create a file `/usr/src/conf/envreplace.json`
```json
{
  "files": {
    "./myconfig.properties": "/usr/local/tomcat/conf/myconfig.properties",
    "./somethingelse.file": "/usr/local/tomcat/conf/whatever.name.file"
  }
}
```

Then you'd run:
```bash
FOOBAR="Some value"
envreplace /usr/src/conf/envreplace.json
```
And both the files declared would have their `#FOOBAR#` replaced with `Some value`

## Configuration file options
The .json file that contains the config has two keys:
* `variableRegex` a string version of a RegEx.  It must contain a single capture group that will be used to pull out the environment variable name.
* `files` a json object.  This is used as a key/value map of file src to file destination. The value can be a string, or an array of strings.  If it is an array of strings, the src file is written to the multiple destinations specified in the array.
```
files: {
  "/path/to/source.file": "/path/to/destination.file",
  "/path/to/some/other/source.file": "/another/destination.file",
  "/path/to/a/file/that/should/go/two/places.file": ["/somewhere/place1.file", "/somewhereElse/place2.file"]
}
```

File paths listed in the files map are relative to the location of the configuration file.

## CLI options:

### Usage:
  * envreplace <srcFile> <destFile>
      Calls envreplace on <srcFile> and put the result into <destFile>

  * envreplace <configfile>
      Call envreplace with configuration stored in json formated file <configfile>

### Flags:
  * -help
      Prints help text
  * -regex string
      Regular expression that parses variable name into capture group 1. Remember to quote and escape to pass it through your shell. If no value is provided the default is '#([A-Z0-9_]+)#'
  * -silent
      Skips printing count of files and substitutions made
  * -verbose
      Prints info about each substitution made

## Fails on missing env vars
This tool will exit with a non-zero exit code and prints a helpful message if any of the files specified in the files option have match the variableRegex but do not have a matching env var.  This prevents you from accidentally using config that did not get substituted properly.

## Building
There's a very simple included makefile

* `make` build the executable
* `make test` build the executable and run the tests
* `make install` build the executable, run the tests, and copy the file to /usr/local/bin

## License (MIT)

Copyright (c) 2017 Ian Taylor

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

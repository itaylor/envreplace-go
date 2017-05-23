## Envreplace

A small, statically linkable cli utility written in GoLang that replaces strings in files with the value of env vars.

The primary usage is for docker-izing existing applications that don't look at the environment for their configuration. For when you need something to write the info from docker's set env vars into pre-existing configuration.

It's roughly similar in functionality to [my node.js module with the same name](https://www.npmjs.com/package/envreplace), but since it can be statically linked, it's much better suited for use in a docker container.

### A simple example

If you have a config file `/usr/src/conf/myconfig.properties`
```
my.property.foobar=#FOOBAR#
```
That you want to end up in `/usr/local/lib/tomcat/conf` with its value `#FOOBAR#` taken from the environment variable `FOOBAR`

You'd write an envreplace descriptor at `/usr/src/conf/envreplace.json`
```json
{
  "variableRegex": "#([A-Z0-9_]+)#",
  "files": {
    "/usr/src/conf/myconfig.properties": "/usr/local/lib/tomcat/conf/myconfig.properties",
  }
}
```

Then you'd run:
```bash
FOOBAR="This is foobar value"
envreplace /usr/src/conf/envreplace.json
```

And you'd end up with `/usr/local/lib/tomcat/conf/myconfig.properties`
```
my.property.foobar=This is foobar value
```

You'll probably want to wire this up so that it runs as part of your docker entrypoint, so that every time your container starts, you get the env vars copied into the config.

## Configuration Options
The .json file that contains the config has two keys:
* `variableRegex` a string version of a RegEx.  It must contain a single capture group that will be used to pull out the environment variable name.
* `files` a json object.  This is used as a key/value map of file src to file destination.
``` files: {
  "/path/to/source.file": "/path/to/destination.file",
  "/path/to/some/other/source.file": "/another/destination.file"
}
```

## Fails on missing env var
This tool will exit with a non-zero exit code and prints a helpful message if any of the files specified in the files option have match the variableRegex but do not have a matching env var.  This prevents you from accidentally using config that did not get substituted properly.

## Building
There's a very simple included makefile

* `make` build the executable
* `make test` build the executable and run the tests
* `make install` build the executable, run the tests, and copy the file to /usr/local/bin

// Go offers built-in support for [regular expressions](http://en.wikipedia.org/wiki/Regular_expression).
// Here are some examples of  common regexp-related tasks
// in Go.

package main

import "os"
//import "bytes"
import "fmt"
import "regexp"
import "encoding/json"
import "io/ioutil"
import "flag"
import "path/filepath"

func main() {

  var regex string
  var verbose bool
  var silent bool
  var help bool
  flag.StringVar(&regex, "regex", "", "Regular expression that parses variable name into capture group 1. Remember to quote and escape to pass it through your shell. If no value is provided the default is '#([A-Z0-9_]+)#'")
  flag.BoolVar(&verbose, "verbose", false, "Prints info about each substitution made")
  flag.BoolVar(&silent, "silent", false, "Skips printing count of files and substitutions made")
  flag.BoolVar(&help, "help", false, "Prints help text")
  flag.Parse()

  var argLen = len(flag.Args())
  if argLen == 0 || argLen > 2 || help {
    fmt.Println(`
Envreplace
  Replaces strings in files with their matching environment variable value.
  https://github.com/itaylor/envreplace-go

Usage:
  envreplace <srcFile> <destFile>
    Calls envreplace on <srcFile> and put the result into <destFile>

  envreplace <configfile>
    Call envreplace with configuration stored in json formated file <configfile>

Flags:`)
    flag.PrintDefaults()
    os.Exit(1)
  }

  var config Config
  if (argLen == 1) {
    config = loadConfig(flag.Args()[0])
  } else {
    config = Config{}
    config.Files = make(map[string]interface{})
    config.Files[flag.Args()[0]] = flag.Args()[1]
    config.BasePath = GetCwd();
  }

  if regex != "" {
    config.VariableRegex = regex;
  }
  if config.VariableRegex == "" {
    config.VariableRegex = "#([A-Z0-9_]+)#"
  }
  r := regexp.MustCompile(config.VariableRegex)

  fileCount := 0
  matchCount := 0

  for src, dest := range config.Files {

    dests := CoerceToDests(dest)

    var err interface{}
    fileCount++
    if verbose {
      fmt.Println(fmt.Sprintf("Processing file passed in as '%v'", src))
    }
    if !filepath.IsAbs(src) {
      src = filepath.Join(config.BasePath, src)
      src, err = filepath.Abs(src)
    }
    if verbose {
      fmt.Println(fmt.Sprintf("  Reading file '%v'", src))
    }

    srcData, err := ioutil.ReadFile(src)
    handleError(err)
    destData, replaceCount := doReplace(srcData, r, verbose, src)
    WriteOutput(destData, dests, config.BasePath, replaceCount, verbose)
    matchCount += replaceCount
  }
  if !silent {
    fmt.Println(fmt.Sprintf("Successfully processed %v files and made %v replacements", fileCount, matchCount))
  }
}

func WriteOutput(destData []byte, dests []string, basePath string, replaceCount int, verbose bool) {
  var err interface{}
  for _, d := range dests {
    if !filepath.IsAbs(d) {
      d = filepath.Join(basePath, d)
      d, err = filepath.Abs(d)
    }
    err = ioutil.WriteFile(d, destData, 0644)
    if verbose {
      fmt.Println(fmt.Sprintf("  Wrote %v replaced variables to destination '%v'", replaceCount, d))
    }
    handleError(err)
  }
}

func loadConfig (fileName string) Config {
  configFileName, err := filepath.Abs(fileName)
  handleError(err)
  configFileDir := filepath.Dir(configFileName);
  data, err := ioutil.ReadFile(configFileName)
  handleError(err)
  var config Config
  err = json.Unmarshal(data, &config)
  handleError(err)
  if config.BasePath == "" {
    config.BasePath = configFileDir
  }
  return config;
}

func CoerceToDests(dest interface{}) []string {
  var dests []string;
  switch dest := dest.(type) {
  case []interface{}:
    for _, d := range dest {
      dests = append(dests, fmt.Sprintf("%v", d))
    }
    case string:
      dests = append(dests, dest)
    default:
      handleError(fmt.Sprintf("Unexpected type '%T' of value '%v'", dest, dest))
  }
  return dests;
}

func doReplace(data []byte, regex *regexp.Regexp, verbose bool, fileName string) ([]byte, int) {
  str := string(data[:])
  count := 0
  replaced := ReplaceAllGroupFunc(regex, str, func(groups []string) string {
    var varName = groups[1]
    var envVar, wasSet = os.LookupEnv(varName)
    if !wasSet {
      handleError(fmt.Sprintf("Error: Env var '%v' found as '%v' in file '%v' is not set in the environment.", varName, groups[0], fileName))
    }
    count++
    if verbose {
      fmt.Println(fmt.Sprintf("  replaced '%v' with '%v'", groups[0], envVar))
    }
    return envVar
  })
  return []byte(replaced), count
}

type Config struct {
  VariableRegex string `json:"variableRegex,omitempty"`
  Files map[string]interface{} `json:"files"`
  BasePath string `json:"basePath,omitempty"`
}

func handleError(err interface{}) {
  if err != nil {
    fmt.Println(err)
    os.Exit(1);
  }
}

func GetCwd() string {
  dir, err := os.Getwd()
  if err != nil {
    panic(err)
  }
  return dir
}


func ReplaceAllGroupFunc(re *regexp.Regexp, str string, repl func([]string) string) string {
  result := ""
  lastIndex := 0
  for _, v := range re.FindAllSubmatchIndex([]byte(str), -1) {
    groups := []string{}
    for i := 0; i < len(v); i += 2 {
      groups = append(groups, str[v[i]:v[i+1]])
    }
    result += str[lastIndex:v[0]] + repl(groups)
    lastIndex = v[1]
  }
  return result + str[lastIndex:]
}

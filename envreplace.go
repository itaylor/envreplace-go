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

  flag.Parse();
  if len(flag.Args()) == 0 {
    handleError("Usage: envreplace <configfile>")
  }
  var configFileName = flag.Args()[0];

  configFileName, err := filepath.Abs(configFileName)
  handleError(err)
  configFileDir := filepath.Dir(configFileName);
  handleError(err)
  data, err := ioutil.ReadFile(configFileName)
  handleError(err)
  var config Config
  err = json.Unmarshal(data, &config)
  handleError(err)
  r := regexp.MustCompile(config.VariableRegex)

  for src, dest := range config.Files {

    if !filepath.IsAbs(src) {
      src = filepath.Join(configFileDir, src)
      src, err = filepath.Abs(src)
      handleError(err)
    }
    if !filepath.IsAbs(dest) {
      dest = filepath.Join(configFileDir, dest)
      dest, err = filepath.Abs(dest)
      handleError(err)
    }
    handleError(err)
    srcData, err := ioutil.ReadFile(src)
    handleError(err)
    destData := doReplace(srcData, r)
    err = ioutil.WriteFile(dest, destData, 0644)
    handleError(err)
  }
}

func doReplace(data []byte, regex *regexp.Regexp) []byte {
  str := string(data[:])
  replaced := ReplaceAllGroupFunc(regex, str, func(groups []string) string {
    var varName = groups[1]
    var envVar, wasSet = os.LookupEnv(varName)
    if !wasSet {
      handleError(fmt.Sprintf("The env var '%v' is not found in the environment.", varName))
    }
    return envVar
  })
  return []byte(replaced)
}

type Config struct {
  VariableRegex string `json:"variableRegex,omitempty"`//= "${env\\.(.*?)}"
  Files map[string]string `json:"files"`
}

func handleError(err interface{}) {
  if err != nil {
    fmt.Println(err)
    os.Exit(1);
  }
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

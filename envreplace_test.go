package main
import "testing"
import "os/exec"
import "os"
import "io/ioutil"
import "path"
import "encoding/json"
import "fmt"
import "reflect"
import "strings"

func TestMain(m *testing.M) {
  os.Exit(m.Run())
}

func TestPositiveCaseWithJsonConfig(t *testing.T) {
  t.Run("Test1 positive test", func(t *testing.T) {
    var cliArgs []string = []string {
        "./fixtures/test1.json",
      }
    cmd := getCommand(loadTestEnv("test1"), cliArgs)
    output, err := cmd.CombinedOutput()
    if err != nil {
      t.Fatal(err, string(output[:]))
    }
    actual1a := getOutput("test1a")
    expected1a := getExpected("test1a")
    compare(actual1a, expected1a, t)
    actual1b := getOutput("test1b")
    expected1b := getExpected("test1b")
    compare(actual1b, expected1b, t)
  })
}

func TestNegativeCaseWithJsonConfig(t *testing.T) {
  t.Run("Test2 negative test", func(t *testing.T) {
    var cliArgs []string = []string {
        "./fixtures/test2.json",
      }
      cmd := getCommand(loadTestEnv("test1"), cliArgs)
    output, err := cmd.CombinedOutput()
    if err == nil {
      t.Fatal("Expected a non-zero error code")
    }
    expectedStdOut := "Error: Env var 'FOOBAR' found as '#FOOBAR#' in file '/usr/src/myapp/fixtures/test2.txt' is not set in the environment.\n"
    expectedStdOut = strings.Replace(expectedStdOut, "/usr/src/myapp", getCwd(), -1)
    compare(string(output[:]), expectedStdOut, t)
  })
}

func TestPositiveCaseWithCliArgs(t *testing.T) {
  t.Run("Test3 cli args positive test", func(t *testing.T) {
    var envVars []string = []string {"FOO=test1", "BAR=test2"}
    var cliArgs []string = []string {"./fixtures/test3.txt", "./output/test3.txt.out"}
    cmd := getCommand(envVars, cliArgs)
    output, err := cmd.CombinedOutput()
    if err != nil {
      t.Fatal(err, string(output[:]))
    }
    actual := getOutput("test3")
    expected := getExpected("test3")
    compare(actual, expected, t)
  })
}

func TestRegexFlag(t *testing.T) {
  t.Run("Test4 cli args regex flag", func(t *testing.T) {
    var envVars []string = []string {"FOO=test4a", "BAR=test4b"}
    var cliArgs []string = []string {
      "-regex=\\$\\{([A-Z0-9_-]+)\\}",
      "./fixtures/test4.txt",
      "./output/test4.txt.out",
    }
    cmd := getCommand(envVars, cliArgs)
    output, err := cmd.CombinedOutput()
    if err != nil {
      t.Fatal(err, string(output[:]))
    }
    actual := getOutput("test4")
    expected := getExpected("test4")
    compare(actual, expected, t)
  })
}

func TestVerboseFlag(t *testing.T) {
  t.Run("Test verbose flag", func(t *testing.T) {
    var cliArgs []string = []string {
      "-verbose",
      "./fixtures/test1.json",
    }
    cmd := getCommand(loadTestEnv("test1"), cliArgs)
    output, err := cmd.CombinedOutput()
    if err != nil {
      t.Fatal(err, string(output[:]))
    }
    expectedStdOut :=
`Processing file passed in as './test1a.txt'
  Reading file '/usr/src/myapp/fixtures/test1a.txt'
  replaced '#FOO#' with 'test1'
  replaced '#BAR#' with 'test2'
  replaced '#FOO#' with 'test1'
  replaced '#FOO#' with 'test1'
  replaced '#BAR#' with 'test2'
  replaced '#FOO#' with 'test1'
  Wrote 6 replaced variables to destination '/usr/src/myapp/output/test1a.txt.out'
Processing file passed in as './test1b.txt'
  Reading file '/usr/src/myapp/fixtures/test1b.txt'
  replaced '#FOO#' with 'test1'
  Wrote 1 replaced variables to destination '/usr/src/myapp/output/test1b.txt.out'
Successfully processed 2 files and made 7 replacements
`
    expectedStdOut = strings.Replace(expectedStdOut, "/usr/src/myapp", getCwd(), -1)
    compare(string(output), expectedStdOut, t)
  })
}

func TestSilentFlag(t *testing.T) {
  t.Run("Test silent flag", func(t *testing.T) {
    var cliArgs []string = []string {
      "-silent",
      "./fixtures/test1.json",
    }
    cmd := getCommand(loadTestEnv("test1"), cliArgs)
    output, err := cmd.CombinedOutput()
    if err != nil {
      t.Fatal(err, string(output[:]))
    }
    compare(string(output), "", t)
  })
}

func TestNormalOutput(t *testing.T) {
  t.Run("Test normal output", func(t *testing.T) {
    var cliArgs []string = []string {
      "./fixtures/test1.json",
    }
    cmd := getCommand(loadTestEnv("test1"), cliArgs)
    output, err := cmd.CombinedOutput()
    if err != nil {
      t.Fatal(err, string(output[:]))
    }
    compare(string(output[:]), "Successfully processed 2 files and made 7 replacements\n", t)
  })
}

func TestConfigFilesWithArrays(t *testing.T) {
  t.Run("Test array in file config has multiple outputs for one input", func(t *testing.T) {
    var cliArgs []string = []string {
      "./fixtures/test5.json",
    }
    cmd := getCommand(loadTestEnv("test5"), cliArgs)
    output, err := cmd.CombinedOutput()
    if err != nil {
      t.Fatal(err, string(output[:]))
    }
    compare(string(output[:]), "Successfully processed 1 files and made 2 replacements\n", t)
    actual5a := getOutput("test5a")
    expected5 := getExpected("test5")
    compare(actual5a, expected5, t)
    actual5b := getOutput("test5b")
    compare(actual5b, expected5, t)
  })
}

func loadTestEnv(name string) []string {
  data, err := ioutil.ReadFile(path.Join("./", "fixtures", name + ".env.json"))
  handleErrorPanic(err)
  var testEnv map[string]string
  err2 := json.Unmarshal(data, &testEnv)
  handleErrorPanic(err2)
  env := os.Environ()
  for key, value := range testEnv {
    env = append(env, fmt.Sprintf("%s=%s", key, value))
  }
  return env
}

func getCwd() string {
  dir, err := os.Getwd()
  if err != nil {
    panic(err)
  }
  return dir
}

func getCommand(envVars []string, cliArgs []string) *exec.Cmd {
  dir := getCwd()
  cmd := exec.Command(path.Join(dir, "envreplace"), cliArgs...)
  cmd.Env = envVars;
  return cmd
}

func getExpected(name string) string {
  dir := getCwd()
  data, err := ioutil.ReadFile(path.Join(dir, "fixtures", name + ".expected.txt"))
  handleErrorPanic(err)
  return string(data[:])
}

func getOutput(name string) string {
  dir := getCwd()
  data, err := ioutil.ReadFile(path.Join(dir, "output", name + ".txt.out"))
  handleErrorPanic(err)
  return string(data[:])
}

func compare(actual string, expected string, t *testing.T) {
  if !reflect.DeepEqual(actual, expected) {
    t.Fatalf("actual = %s, expected = %s", actual, expected)
  }
}

func handleErrorPanic(err interface{}) {
  if err != nil {
    fmt.Println(err)
    panic(err)
  }
}

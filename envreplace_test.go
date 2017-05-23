package main
import "testing"
import "os/exec"
import "os"
import "io/ioutil"
import "path"
import "encoding/json"
import "fmt"
import "reflect"


func TestMain(m *testing.M) {
  os.Exit(m.Run())
}

func TestPositiveCase(t *testing.T) {
  t.Run("Test1 positive test", func(t *testing.T) {
    cmd := getCommand("test1")
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

func TestNegativeCase(t *testing.T) {
  t.Run("Test2 negative test", func(t *testing.T) {
    cmd := getCommand("test2")
    output, err := cmd.CombinedOutput()
    if err == nil {
      t.Fatal("Expected a non-zero error code")
    }
    compare(string(output[:]), "The env var 'FOOBAR' is not found in the environment.\n", t)
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

func getCommand(name string) *exec.Cmd {
  dir := getCwd()
  args :=  path.Join(dir, "fixtures", name + ".json")
  cmd := exec.Command(path.Join(dir, "envreplace"), args)
  cmd.Env = loadTestEnv(name)

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

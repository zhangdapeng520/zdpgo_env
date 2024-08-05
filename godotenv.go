// package zdpgo_env is a go port of the ruby dotenv library (https://github.com/bkeepers/dotenv)
//
// Examples/readme can be found on the GitHub page at https://github.com/joho/godotenv
//
// The TL;DR is that you make a .env file that looks something like
//
//	SOME_ENV_VAR=somevalue
//
// and then in your go code you can call
//
//	godotenv.Load()
//
// and all the env vars declared in .env will be available through os.Getenv("SOME_ENV_VAR")
package zdpgo_env

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const doubleQuoteSpecialChars = "\\\n\r\"!$`"

// Parse reads an env file from io.Reader, returning a map of keys and values.
func Parse(r io.Reader) (map[string]string, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		return nil, err
	}

	return UnmarshalBytes(buf.Bytes())
}

// Load 将读取您的env文件并将其加载到env中。
// 在程序开始时尽可能地调用这个函数(最好是在main中)。
// 如果你不带任何参数调用Load，它将默认加载当前路径中的.env。
// 你也可以告诉它加载哪些文件(可以有多个)，比如:
//
//	zdpgo_env.Load("fileone", "filetwo")
//
// 重要的是要注意，它不会覆盖已经存在的env变量-考虑使用.env文件设置dev变量或合理的默认值。
func Load(filenames ...string) (err error) {
	filenames = filenamesOrDefault(filenames)

	for _, filename := range filenames {
		err = loadFile(filename, false)
		if err != nil {
			return // return early on a spazout
		}
	}
	return
}

// Overload will read your env file(s) and load them into ENV for this process.
//
// Call this function as close as possible to the start of your program (ideally in main).
//
// If you call Overload without any args it will default to loading .env in the current path.
//
// You can otherwise tell it which files to load (there can be more than one) like:
//
//	godotenv.Overload("fileone", "filetwo")
//
// It's important to note this WILL OVERRIDE an env variable that already exists - consider the .env file to forcefully set all vars.
func Overload(filenames ...string) (err error) {
	filenames = filenamesOrDefault(filenames)

	for _, filename := range filenames {
		err = loadFile(filename, true)
		if err != nil {
			return // return early on a spazout
		}
	}
	return
}

// Read all env (with same file loading semantics as Load) but return values as
// a map rather than automatically writing values into env
func Read(filenames ...string) (envMap map[string]string, err error) {
	filenames = filenamesOrDefault(filenames)
	envMap = make(map[string]string)

	for _, filename := range filenames {
		individualEnvMap, individualErr := readFile(filename)

		if individualErr != nil {
			err = individualErr
			return // return early on a spazout
		}

		for key, value := range individualEnvMap {
			envMap[key] = value
		}
	}

	return
}

// Unmarshal reads an env file from a string, returning a map of keys and values.
func Unmarshal(str string) (envMap map[string]string, err error) {
	return UnmarshalBytes([]byte(str))
}

// UnmarshalBytes parses env file from byte slice of chars, returning a map of keys and values.
func UnmarshalBytes(src []byte) (map[string]string, error) {
	out := make(map[string]string)
	err := parseBytes(src, out)

	return out, err
}

// Exec loads env vars from the specified filenames (empty map falls back to default)
// then executes the cmd specified.
//
// Simply hooks up os.Stdin/err/out to the command and calls Run().
//
// If you want more fine grained control over your command it's recommended
// that you use `Load()`, `Overload()` or `Read()` and the `os/exec` package yourself.
func Exec(filenames []string, cmd string, cmdArgs []string, overload bool) error {
	op := Load
	if overload {
		op = Overload
	}
	if err := op(filenames...); err != nil {
		return err
	}

	command := exec.Command(cmd, cmdArgs...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

// Write serializes the given environment and writes it to a file.
func Write(envMap map[string]string, filename string) error {
	content, err := Marshal(envMap)
	if err != nil {
		return err
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content + "\n")
	if err != nil {
		return err
	}
	return file.Sync()
}

// Marshal outputs the given environment as a dotenv-formatted environment file.
// Each line is in the format: KEY="VALUE" where VALUE is backslash-escaped.
func Marshal(envMap map[string]string) (string, error) {
	lines := make([]string, 0, len(envMap))
	for k, v := range envMap {
		if d, err := strconv.Atoi(v); err == nil {
			lines = append(lines, fmt.Sprintf(`%s=%d`, k, d))
		} else {
			lines = append(lines, fmt.Sprintf(`%s="%s"`, k, doubleQuoteEscape(v)))
		}
	}
	sort.Strings(lines)
	return strings.Join(lines, "\n"), nil
}

func filenamesOrDefault(filenames []string) []string {
	if len(filenames) == 0 {
		return []string{".env"}
	}
	return filenames
}

func loadFile(filename string, overload bool) error {
	envMap, err := readFile(filename)
	if err != nil {
		return err
	}

	currentEnv := map[string]bool{}
	rawEnv := os.Environ()
	for _, rawEnvLine := range rawEnv {
		key := strings.Split(rawEnvLine, "=")[0]
		currentEnv[key] = true
	}

	for key, value := range envMap {
		if !currentEnv[key] || overload {
			_ = os.Setenv(key, value)
		}
	}

	return nil
}

func readFile(filename string) (envMap map[string]string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	return Parse(file)
}

func doubleQuoteEscape(line string) string {
	for _, c := range doubleQuoteSpecialChars {
		toReplace := "\\" + string(c)
		if c == '\n' {
			toReplace = `\n`
		}
		if c == '\r' {
			toReplace = `\r`
		}
		line = strings.Replace(line, string(c), toReplace, -1)
	}
	return line
}

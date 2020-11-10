package env

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	environmentFile = ".terraform/environment"
)

var (
	stokenvRegex = regexp.MustCompile("^[a-z0-9-]+/[a-z0-9-]+$")
)

// A string identifying a namespaced workspace, according to the format $namespace/$workspace, with
// helper functions to read and write the string to the file .terraform/environment
type StokEnv string

func ValidateAndParse(stokenv string) (namespace string, workspace string, err error) {
	if err = Validate(stokenv); err != nil {
		return workspace, namespace, err
	}
	parts := strings.Split(stokenv, "/")
	return parts[0], parts[1], nil
}

func Validate(stokenv string) error {
	if !stokenvRegex.MatchString(stokenv) {
		return fmt.Errorf("workspace must match pattern %s", stokenvRegex.String())
	}
	return nil
}

func WithOptionalNamespace(stokenv string) StokEnv {
	parts := strings.Split(stokenv, "/")
	if len(parts) == 1 {
		return StokEnv(fmt.Sprintf("default/%s", parts[0]))
	}
	return StokEnv(stokenv)
}

func NewStokEnv(namespace, workspace string) StokEnv {
	return StokEnv(fmt.Sprintf("%s/%s", namespace, workspace))
}

func (env StokEnv) Namespace() string {
	return strings.Split(string(env), "/")[0]
}

func (env StokEnv) Workspace() string {
	return strings.Split(string(env), "/")[1]
}

func ReadStokEnv(path string) (StokEnv, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadFile(filepath.Join(absPath, environmentFile))
	if err != nil {
		return "", err
	}

	return StokEnv(string(bytes)), nil
}

func (env StokEnv) Write(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	envPath := filepath.Join(absPath, environmentFile)
	if err := os.MkdirAll(filepath.Dir(envPath), 0755); err != nil {
		return err
	}

	return ioutil.WriteFile(envPath, []byte(string(env)), 0644)
}

func WriteEnvFile(path, content string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	envPath := filepath.Join(absPath, environmentFile)
	if err := os.MkdirAll(filepath.Dir(envPath), 0755); err != nil {
		return err
	}

	return ioutil.WriteFile(envPath, []byte(content), 0644)
}

// +build mage

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/mcandre/gmake-shim"
)

// ArtifactsPath describes where artifacts are produced.
var ArtifactsPath = "bin"

// Default references the default build task.
var Default = Port

// UnitTest executes the unit test suite.
func UnitTest() error {
	cmd := exec.Command("go", "test")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// IntegrationTest executes the integration test suite.
func IntegrationTest() error {
	cmd := exec.Command("gmake")
	cmd.Dir = "test"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Test executes unit and integration tests.
func Test() error {
	mg.Deps(UnitTest)
	mg.Deps(IntegrationTest)
	return nil
}

// PathSeparatorString models the os.PathSeparator as a string.
var PathSeparatorString = string(os.PathSeparator)

// AllPackagesPath denotes all Go packages in this project.
var AllPackagesPath = strings.Join([]string{".", "..."}, PathSeparatorString)

// GoVet runs go tool vet.
func GoVet() error {
	cmd := exec.Command("go", "vet", "-shadow", AllPackagesPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GoSourceFiles represents the set of Go source files in a project.
// Populdate with CollectGoFiles().
var GoSourceFiles = make(map[string]bool)

// CollectGoFiles queries Go source file paths.
func CollectGoFiles() error {
	var out bytes.Buffer

	cmd := exec.Command("go", "list", "-f", "{{$p := .}}{{range $f := .GoFiles}}{{$p.Dir}}/{{$f}}{{end}}{{range $f := .TestGoFiles}}{{$p.Dir}}/{{$f}}{{end}}", AllPackagesPath)
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(&out)

	for scanner.Scan() {
		pth := scanner.Text()
		GoSourceFiles[pth] = true
	}

	return nil
}

// GoLint runs golint.
func GoLint() error {
	mg.Deps(CollectGoFiles)

	for pth, _ := range GoSourceFiles {
		cmd := exec.Command("golint", pth)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

// Gofmt runs gofmt.
func GoFmt() error {
	mg.Deps(CollectGoFiles)

	for pth, _ := range GoSourceFiles {
		cmd := exec.Command("gofmt", "-s", "-w", pth)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

// GoImports runs goimports.
func GoImports() error {
	mg.Deps(CollectGoFiles)

	for pth, _ := range GoSourceFiles {
		cmd := exec.Command("goimports", "-w", pth)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

// GoErrcheck runs errcheck.
func Errcheck() error {
	cmd := exec.Command("errcheck", "-blank")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Nakedret runs nakedret.
func Nakedret() error {
	cmd := exec.Command("nakedret", "-l", "0", AllPackagesPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Lint runs the lint suite.
func Lint() error {
	mg.Deps(GoVet)
	mg.Deps(GoLint)
	mg.Deps(GoFmt)
	mg.Deps(GoImports)
	mg.Deps(Errcheck)
	mg.Deps(Nakedret)
	return nil
}

// Port archives build artifacts.
func Port() error {
	mg.Deps(Artifacts)

	archiveBasename := fmt.Sprintf("gmake-shim-%s", gmake.Version)
	archiveFilename := fmt.Sprintf("%s.zip", archiveBasename)

	cmd := exec.Command("zipc", "-chdir", ArtifactsPath, archiveFilename, archiveBasename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// AllCommandsPath denotes all Go application packages in this project.
var AllCommandsPath = strings.Join([]string{".", "cmd", "..."}, PathSeparatorString)

// Artifacts cross-compiles Go binaries.
func Artifacts() error {
	if err := os.MkdirAll(ArtifactsPath, os.ModeDir|0775); err != nil {
		return err
	}

	parts := []string{fmt.Sprintf("-output=%s", ArtifactsPath), fmt.Sprintf("gmake-shim-%s", gmake.Version), "{{.OS}}", "{{.Arch}}", "{{.Dir}}"}
	artifactStructure := strings.Join(parts, PathSeparatorString)

	cmd := exec.Command("gox", artifactStructure, AllCommandsPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Applications names which applications this project builds.
var Applications = []string{"gmake"}

// Install builds and installs Go applications.
func Install() error {
	cmd := exec.Command("go", "install", AllPackagesPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GoBin denotes the path to the Go binaries directory.
// Populated with LoadGoBinariesPath.
var GoBin = ""

// LoadBinariesPath queries the Go binaries directory.
func LoadGoBinariesPath() error {
	goPath := os.Getenv("GOPATH")

	if goPath == "" {
		user, err := user.Current()

		if err != nil {
			return err
		}

		goPath = path.Join(user.HomeDir, "go")
	}

	GoBin = path.Join(goPath, "bin")

	return nil
}

// Uninstall deletes installed Go applications.
func Uninstall() error {
	mg.Deps(LoadGoBinariesPath)

	for _, application := range Applications {
		if err := os.Remove(path.Join(GoBin, application)); err != nil {
			return err
		}
	}

	return nil
}

// Clean deletes artifacts.
func Clean() error {
	return os.RemoveAll(ArtifactsPath)
}

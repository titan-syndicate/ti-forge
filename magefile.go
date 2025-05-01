//go:build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Build builds the plugin binary
func Build() error {
	fmt.Println("Building plugin...")
	return sh.Run("go", "build", "-o", "bin/ti-scaffold")
}

// Run runs the plugin
func Run() error {
	mg.Deps(Build)
	fmt.Println("Running plugin...")
	return sh.Run("./bin/ti-scaffold")
}

// Dev runs the plugin in development mode
func Dev() error {
	mg.Deps(Build)
	fmt.Println("Running plugin in development mode...")
	return sh.Run("./bin/ti-scaffold", "--dev")
}

// Clean removes build artifacts
func Clean() error {
	fmt.Println("Cleaning...")
	// Remove bin directory
	if err := os.RemoveAll("bin"); err != nil {
		return err
	}
	// Remove output directory
	if err := os.RemoveAll("output"); err != nil {
		return err
	}
	return nil
}

// Install installs the plugin
func Install() error {
	mg.Deps(Build)
	fmt.Println("Installing plugin...")
	return sh.Run("cp", "bin/ti-scaffold", "/usr/local/bin/")
}

// Scaffold runs the scaffold command with test inputs
func Scaffold() error {
	mg.Deps(Build)
	fmt.Println("Running scaffold command...")
	return sh.Run("./bin/ti-scaffold", "--dev", "--config", "testdata/scaffold.yaml")
}

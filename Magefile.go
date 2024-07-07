//go:build mage
// +build mage

package main

import (
	"fmt"
	// mage:import
	build "github.com/grafana/grafana-plugin-sdk-go/build"
	"github.com/magefile/mage/mg"
)

// Hello prints a message (shows that you can define custom Mage targets).
func Hello() {
	fmt.Println("hello plugin developer!")
}

func LinuxARM64() {
	b := build.Build{}
	mg.Deps(b.LinuxARM64)
}

func Linux() {
	b := build.Build{}
	mg.Deps(b.Linux())
}

func LinuxARM() {
	b := build.Build{}
	mg.Deps(b.LinuxARM())
}

func Windows() { //revive:disable-line
	b := build.Build{}
	mg.Deps(b.Windows)
}

// Default configures the default target.
var Default = LinuxARM64

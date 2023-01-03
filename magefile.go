//go:build mage
// +build mage

package main

import (
	"github.com/carolynvs/magex/shx"
)

var must = shx.CommandBuilder{StopOnError: true}

func ci() {
	must.RunV("go", "run", "./ci/main.go")
}

func Build() {
	must.RunV("go", "build", "./...")
}

func Test() {

	must.RunV("go", "test", "./...")
}

func Upgrade() {
	must.RunV("go", "get", "-d", "-u", "./...")
	must.RunV("go", "mod", "tidy")
	// must.RunV("go", "test", "vendor")
}

package main

import (
	syngo "github.com/synesissoftware/syngo"
	ver2go "github.com/synesissoftware/ver2go"

	"fmt"
)

func main() {
	fmt.Printf("syngo v%s\n", ver2go.CalcVersionString(syngo.VersionMajor, syngo.VersionMinor, syngo.VersionPatch, syngo.VersionAB))
	fmt.Printf("ver2go v%s\n", ver2go.CalcVersionString(ver2go.VersionMajor, ver2go.VersionMinor, ver2go.VersionPatch, ver2go.VersionAB))
}

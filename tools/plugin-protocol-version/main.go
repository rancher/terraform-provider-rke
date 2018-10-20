package main

import (
	"fmt"

	"github.com/hashicorp/terraform/plugin/discovery"
)

func main() {
	fmt.Print(discovery.PluginInstallProtocolVersion)
}

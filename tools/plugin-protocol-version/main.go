package main

import (
	"fmt"

	"github.com/hashicorp/terraform/plugin"
)

func main() {
	fmt.Print(plugin.Handshake.ProtocolVersion)
}

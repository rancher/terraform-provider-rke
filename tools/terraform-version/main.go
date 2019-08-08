package main

import (
	"fmt"

	"github.com/hashicorp/terraform/version"
)

func main() {
	fmt.Print(version.Version)
}

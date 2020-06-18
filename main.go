package main

import (
	"log"

	"github.com/fenrirunbound/kubeconfig-merge/internal/commands"
)

func main() {
	cmd := commands.Root("kubeconfig-merge")

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

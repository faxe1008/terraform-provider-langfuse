package main

import (
	"context"
	"flag"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/faxe1008/terraform-provider-langfuse/langfuse"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "run in debug mode")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/langfuse/langfuse",
		Debug:   debug,
	}

	providerConstructor := func() provider.Provider {
		return langfuse.NewProvider("0.1.0")
	}

	if err := providerserver.Serve(context.Background(), providerConstructor, opts); err != nil {
		os.Exit(1)
	}
}

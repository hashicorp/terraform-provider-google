package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-provider-google/google"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	// concat with sdkv2 provider
	providers := []func() tfprotov5.ProviderServer{
		google.Provider().GRPCProvider, // sdk provider
	}

	// use the muxer
	muxServer, err := tf5muxserver.NewMuxServer(context.Background(), providers...)
	if err != nil {
		log.Fatalf(err.Error())
	}

	var serveOpts []tf5server.ServeOpt

	if debug {
		serveOpts = append(serveOpts, tf5server.WithManagedDebug())
	}

	err = tf5server.Serve(
		"registry.terraform.io/hashicorp/google",
		muxServer.ProviderServer,
		serveOpts...,
	)

	if err != nil {
		log.Fatal(err)
	}
}

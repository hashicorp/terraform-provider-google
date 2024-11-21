// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"

	"github.com/hashicorp/terraform-provider-google/google/fwprovider"
	"github.com/hashicorp/terraform-provider-google/google/provider"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	// primary is the SDKv2 implementation of the provider
	primary := provider.Provider()

	providers := []func() tfprotov5.ProviderServer{
		primary.GRPCProvider, // sdk provider
		providerserver.NewProtocol5(fwprovider.New(primary)), // framework provider
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

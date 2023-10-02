package main

import (
	"github.com/abagayev/terraform-provider-forwardemail/forwardemail"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: forwardemail.Provider,
	})
}

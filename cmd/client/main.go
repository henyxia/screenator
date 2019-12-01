package client

import (
	"github.com/henyxia/screenator/internal/config"
	"github.com/henyxia/screenator/internal/controlclient"
	"log"
)

func usage() {
	log.Fatalln("usage ./screenator-client CONFIG_FILE")
}

// Client starts the screenator client
func Client(configFile string) {
	log.Println("start screenator client")

	log.Println("read configuration file")
	conf := config.ReadConfigClient(configFile)
	log.Println("mac is set to:", conf.Mac)

	log.Println("control url is:", conf.RemoteProxyControl)
	var control = controlclient.ControlClient{
		Endpoint:            conf.RemoteProxyControl,
		Mac:                 conf.Mac,
		Proxy:               conf.RemoteProxyData,
		PortMin:             conf.LocalBindMin,
		PortMax:             conf.LocalBindMax,
		BrowserCommand:      conf.Browser,
		ProxyScheme:         conf.ProxyScheme,
		ProxyTLSCertificate: conf.ProxyTLSCertificate,
		ProxyTLSKey:         conf.ProxyTLSKey,
	}

	control.Start()
}

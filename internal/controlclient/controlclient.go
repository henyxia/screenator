package controlclient

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/henyxia/screenator/internal/browser"
	"github.com/henyxia/screenator/internal/model"
	"github.com/henyxia/screenator/internal/proxylocal"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"syscall"
	"time"
)

// ControlClient handler
type ControlClient struct {
	PortMin             uint16
	PortMax             uint16
	Proxy               string
	Endpoint            string
	Mac                 string
	Binds               map[int]*proxylocal.LocalProxy
	PortMap             map[int]int
	BrowserCommand      string
	Browser             *browser.Browser
	reverseDisplays     map[int]*model.Display
	ProxyScheme         string
	ProxyTLSCertificate string
	ProxyTLSKey         string
}

func (c ControlClient) getAPI(url string) (*http.Response, error) {
	var httpclient *http.Client

	/* create http client corresponding to the scheme */
	if c.ProxyScheme == "https" {
		log.Println("will use tls client")
		log.Println("tls certificate path: " + c.ProxyTLSCertificate)
		log.Println("tls key path: " + c.ProxyTLSKey)
		cert, err := tls.LoadX509KeyPair(c.ProxyTLSCertificate, c.ProxyTLSKey)
		if err != nil {
			log.Fatalf("server: loadkeys: %s", err)
		}
		config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
		transport := &http.Transport{TLSClientConfig: &config}
		httpclient = &http.Client{Transport: transport}
	} else if c.ProxyScheme == "http" {
		httpclient = &http.Client{}
	} else {
		log.Fatalln("retard alert! proxy scheme must be http or https!")
	}

	requestURL := c.Endpoint + url
	log.Println("get on:", requestURL)
	resp, err := httpclient.Get(requestURL)
	if err != nil {
		log.Printf("unable to get url '%s': %s", url, err)
		return nil, err
	}

	return resp, nil
}

func (c ControlClient) getDeviceID() int {
	resp, err := c.getAPI("/device/searchByMac/" + c.Mac)
	if err != nil {
		log.Fatalln("cannot retrieve device from mac")
	}
	if resp.StatusCode != 200 {
		log.Fatalln("device unknown for remote control")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var device model.Device
	err = json.Unmarshal(body, &device)
	if err != nil {
		log.Println("cannot unmarshal device")
	}

	return device.ID
}

func (c ControlClient) getDisplays(deviceID int) []model.Display {
	url := fmt.Sprintf("/device/%d/display", deviceID)
	resp, err := c.getAPI(url)
	if err != nil {
		log.Fatalln("cannot retrieve displays")
	}
	if resp.StatusCode != 200 {
		log.Fatalln("no displays")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var displays []model.Display
	err = json.Unmarshal(body, &displays)
	if err != nil {
		log.Println("cannot unmarshal displays")
	}

	return displays
}

func (c ControlClient) getNewPort(displayID int) int {
	numElem := c.PortMax - c.PortMin
	for i := 0; i < int(numElem); i++ {
		if _, ok := c.PortMap[i]; !ok {
			c.PortMap[i] = displayID
			return int(c.PortMin) + i
		}
	}

	log.Fatalln("no more port available")
	return 0
}

func (c ControlClient) getUrls() []string {
	var urls []string
	for port := range c.PortMap {
		url := "127.0.0.1:" + strconv.Itoa(int(c.PortMin)+port)
		urls = append(urls, url)
	}

	return urls
}

func (c ControlClient) removeBind(displayID int) {
	/* remove bind */
	c.Binds[displayID].Server.Shutdown(context.TODO())
	delete(c.Binds, displayID)

	/* remove from the reverse list */
	delete(c.reverseDisplays, displayID)

	/* remove from port list */
	for i, did := range c.PortMap {
		if did == displayID {
			delete(c.PortMap, i)
			break
		}
	}
}

// Start the controlclient
func (c ControlClient) Start() {
	/* init the port map */
	if c.PortMax <= c.PortMin {
		log.Fatalln("retard alert: port max must greater than port min!")
	}

	numElem := c.PortMax - c.PortMin
	c.PortMap = make(map[int]int, numElem)

	/* init browser */
	c.Browser = &browser.Browser{
		Command: c.BrowserCommand,
	}

	/* init var
	 * represents if the browser has been already started
	 */
	init := false

	/* reset var
	 * represents if the browser needs to be restarted
	 * it only happens if a display is removed
	 */
	needReset := false

	/* get assets */
	log.Println("get device id")
	deviceID := c.getDeviceID()

	log.Printf("hello, device#%d", deviceID)

	c.Binds = make(map[int]*proxylocal.LocalProxy)
	for {
		log.Println("get content to display")
		displays := c.getDisplays(deviceID)
		log.Printf("got %d content(s) to display", len(displays))

		/* create a reverse array for fast display removal detection */
		c.reverseDisplays = make(map[int]*model.Display)
		for _, display := range displays {
			c.reverseDisplays[display.ID] = &display
		}

		/* detect removed displays */
		for displayID := range c.Binds {
			if _, ok := c.reverseDisplays[displayID]; !ok {
				log.Println("remove old bind")
				c.removeBind(displayID)
				needReset = true
			}
		}

		if needReset {
			log.Println("reset browser")
			c.Browser.Cmd.Process.Signal(syscall.SIGTERM)
			c.Browser.Cmd.Wait()
			needReset = false
			init = false
		}

		/* bind required ports */
		for _, display := range displays {
			if _, ok := c.Binds[display.ID]; !ok {
				/* create missing binding */
				displayID := display.ID
				log.Printf("create new bind for display %d", displayID)
				/* get port */
				port := c.getNewPort(displayID)
				c.Binds[display.ID] = &proxylocal.LocalProxy{
					Host:                strconv.Itoa(display.Content),
					Bind:                "127.0.0.1:" + strconv.Itoa(port),
					Proxy:               c.Proxy,
					ProxyScheme:         c.ProxyScheme,
					ProxyTLSCertificate: c.ProxyTLSCertificate,
					ProxyTLSKey:         c.ProxyTLSKey,
				}
				c.Binds[displayID].Start()

				/* add it to the browser if initialized and not under restart */
				if init && !needReset {
					log.Println("start browser with the new url")
					urls := []string{"127.0.0.1:" + strconv.Itoa(port)}
					c.Browser.Run(urls)
				}
			}
		}

		if !init {
			log.Println("start browser for the first time")
			c.Browser.Run(c.getUrls())
			init = true
		}

		time.Sleep(10 * time.Second)
	}
}

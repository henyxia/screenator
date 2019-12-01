package browser

import (
	"log"
	"os/exec"
	"strings"
)

// Browser handler
type Browser struct {
	Command string
	Cmd     *exec.Cmd
}

// Run entrypoint
func (b *Browser) Run(urls []string) {
	browserSplit := strings.Split(b.Command, " ")
	browserName := browserSplit[0]
	browserArgs := browserSplit[1:]
	for _, url := range urls {
		browserArgs = append(browserArgs, url)
	}

	log.Println("start browser:", browserName)
	cmd := exec.Command(browserName, browserArgs...)
	if b.Cmd == nil {
		b.Cmd = cmd
	}
	err := cmd.Start()
	if err != nil {
		log.Fatalf("get current asset: %v", err)
	}
}

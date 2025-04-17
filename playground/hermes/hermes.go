package hermes

import (
	"fmt"
	"os"

	"github.com/hanchon/hanchond/playground/filesmanager"
)

type Hermes struct{}

func NewHermes() *Hermes {
	h := &Hermes{}

	if _, err := os.Stat(h.GetHermesBinary()); os.IsNotExist(err) {
		fmt.Printf("Hermes binary at %s does not exist; build it first\n", h.GetHermesBinary())
		os.Exit(1)
	}

	_ = filesmanager.CreateHermesFolder()
	h.initHermesConfig()
	return h
}

package hermes

import (
	"fmt"
	"os"

	"github.com/hanchon/hanchond/lib/utils"
	"github.com/hanchon/hanchond/playground/filesmanager"
)

type Hermes struct{}

func NewHermes() *Hermes {
	h := &Hermes{}

	if _, err := os.Stat(h.GetHermesBinary()); os.IsNotExist(err) {
		utils.ExitError(fmt.Errorf("Hermes binary at %s does not exist; build it first", h.GetHermesBinary()))
	}

	_ = filesmanager.CreateHermesFolder()

	h.initHermesConfig()

	return h
}

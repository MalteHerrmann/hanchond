package orbiter

import "os"

func (o *Orbiter) UpdateAppFile() error {
	appFile, err := o.OpenAppFile()
	if err != nil {
		return err
	}

	appFile = o.setMintDenom(appFile)
	return o.SaveAppFile(appFile)
}

func (o *Orbiter) setMintDenom(config []byte) []byte {
	
}

package manifestgen

import (
	"github.com/iLLeniumStudios/FiveMCarsMerger/pkg/dft"
	"github.com/iLLeniumStudios/FiveMCarsMerger/pkg/flags"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

type Manifest struct {
	HasCarcols          bool
	HasCarvariations    bool
	HasContentUnlocks   bool
	HasHandling         bool
	HasVehicleLayouts   bool
	HasVehicleModelsets bool
	HasVehicles         bool
	HasWeaponsFile      bool
}

type Generator interface {
	Generate() error
}

type generator struct {
	Flags flags.Flags
}

func New(_flags flags.Flags) Generator {
	return &generator{Flags: _flags}
}

func (g *generator) Generate() error {
	tmpl, err := template.New("manifestTemplate").Parse(manifestTemplate)
	if err != nil {
		return err
	}

	var manifest Manifest

	folders, err := ioutil.ReadDir(g.Flags.OutputPath + "/data")
	if err != nil {
		return err
	}

	for _, folder := range folders {
		if folder.Name() == strings.ToLower(dft.CARCOLS.String()) {
			manifest.HasCarcols = true
		} else if folder.Name() == strings.ToLower(dft.CARVARIATIONS.String()) {
			manifest.HasCarvariations = true
		} else if folder.Name() == strings.ToLower(dft.CONTENTUNLOCKS.String()) {
			manifest.HasContentUnlocks = true
		} else if folder.Name() == strings.ToLower(dft.HANDLING.String()) {
			manifest.HasHandling = true
		} else if folder.Name() == strings.ToLower(dft.VEHICLELAYOUTS.String()) {
			manifest.HasVehicleLayouts = true
		} else if folder.Name() == strings.ToLower(dft.VEHICLEMODELSETS.String()) {
			manifest.HasVehicleModelsets = true
		} else if folder.Name() == strings.ToLower(dft.VEHICLES.String()) {
			manifest.HasVehicles = true
		} else if folder.Name() == strings.ToLower(dft.WEAPONSFILE.String()) {
			manifest.HasWeaponsFile = true
		} else {
			log.Warn("Invalid folder name", folder.Name())
		}
	}

	fxManifest, err := os.Create(g.Flags.OutputPath + "/fxmanifest.lua")
	if err != nil {
		return err
	}

	defer fxManifest.Close()

	err = tmpl.Execute(fxManifest, manifest)
	if err != nil {
		return err
	}

	return nil
}

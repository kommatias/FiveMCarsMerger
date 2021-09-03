package main

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

const manifestTemplate = `fx_version 'cerulean'
game 'gta5'

files {
	{{ if .HasCarcols -}}
	'data/carcols/**.meta',
	{{ end -}}
	{{ if .HasCarvariations -}}
	'data/carvariations/**.meta',
	{{ end -}}
	{{ if .HasContentUnlocks -}}
	'data/contentunlocks/**.meta',
	{{ end -}}
	{{ if .HasHandling -}}
	'data/handling/**.meta',
	{{ end -}}
	{{ if .HasVehicleLayouts -}}
	'data/vehiclelayouts/**.meta',
	{{ end -}}
	{{ if .HasVehicleModelsets -}}
	'data/vehiclemodelsets/**.meta',
	{{ end -}}
	{{ if .HasVehicles -}}
	'data/vehicles/**.meta',
	{{ end -}}
	{{ if .HasWeaponsFile -}}
	'data/weaponsfile/**.meta',
	{{ end }}
}

{{ if .HasCarcols -}}
data_file 'CARCOLS_FILE' 'data/carcols/**.meta'
{{ end -}}
{{ if .HasCarvariations -}}
data_file 'VEHICLE_VARIATION_FILE' 'data/vehiclevariations/**.meta'
{{ end -}}
{{ if .HasContentUnlocks -}}
data_file 'CONTENT_UNLOCKING_META_FILE' 'data/contentunlocks/**.meta'
{{ end -}}
{{ if .HasHandling -}}
data_file 'HANDLING_FILE' 'data/handling/**.meta'
{{ end -}}
{{ if .HasVehicleLayouts -}}
data_file 'VEHICLE_LAYOUTS_FILE' 'data/vehiclelayouts/**.meta'
{{ end -}}
{{ if .HasVehicleModelsets -}}
data_file 'AMBIENT_VEHICLE_MODEL_SET_FILE' 'data/vehiclemodelsets/**.meta'
{{ end -}}
{{ if .HasVehicles -}}
data_file 'VEHICLE_METADATA_FILE' 'data/vehicles/**.meta'
{{ end -}}
{{ if .HasWeaponsFile -}}
data_file 'WEAPONINFO_FILE' 'data/weaponsfile/**.meta'
{{ end -}}
`

type DataFileManifest struct {
	Path string
	Type string
}

/*type Manifest struct {
	DataFiles []DataFileManifest
}*/

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

func GenerateFXManifest() error {
	tmpl, err := template.New("manifestTemplate").Parse(manifestTemplate)
	if err != nil {
		return err
	}

	var manifest Manifest

	folders, err := ioutil.ReadDir(flags.OutputPath + "/data")
	if err != nil {
		return err
	}

	for _, folder := range folders {
		if folder.Name() == strings.ToLower(CARCOLS.String()) {
			manifest.HasCarcols = true
		} else if folder.Name() == strings.ToLower(CARVARIATIONS.String()) {
			manifest.HasCarvariations = true
		} else if folder.Name() == strings.ToLower(CONTENTUNLOCKS.String()) {
			manifest.HasContentUnlocks = true
		} else if folder.Name() == strings.ToLower(HANDLING.String()) {
			manifest.HasHandling = true
		} else if folder.Name() == strings.ToLower(VEHICLELAYOUTS.String()) {
			manifest.HasVehicleLayouts = true
		} else if folder.Name() == strings.ToLower(VEHICLEMODELSETS.String()) {
			manifest.HasVehicleModelsets = true
		} else if folder.Name() == strings.ToLower(VEHICLES.String()) {
			manifest.HasVehicles = true
		} else if folder.Name() == strings.ToLower(WEAPONSFILE.String()) {
			manifest.HasWeaponsFile = true
		} else {
			log.Warn("Invalid folder name", folder.Name())
		}
	}

	fxManifest, err := os.Create(flags.OutputPath + "/fxmanifest.lua")
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

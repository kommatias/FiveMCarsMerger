package merger

import (
	"github.com/iLLeniumStudios/FiveMCarsMerger/pkg/carfinder"
	"github.com/iLLeniumStudios/FiveMCarsMerger/pkg/copier"
	"github.com/iLLeniumStudios/FiveMCarsMerger/pkg/dft"
	"github.com/iLLeniumStudios/FiveMCarsMerger/pkg/flags"
	"github.com/iLLeniumStudios/FiveMCarsMerger/pkg/manifestgen"
	"github.com/iLLeniumStudios/FiveMCarsMerger/pkg/typeidentifier"
	sliceutils "github.com/iLLeniumStudios/FiveMCarsMerger/pkg/utils/slice"
	"github.com/iLLeniumStudios/FiveMCarsMerger/pkg/validator"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type Merger interface {
	Merge() error
}

type merger struct {
	Flags          flags.Flags
	Generator      manifestgen.Generator
	Validator      validator.Validator
	TypeIdentifier typeidentifier.TypeIdentifier
	CarFinder      carfinder.CarFinder
	Copier         copier.Copier
}

func New(_flags flags.Flags) Merger {
	return &merger{
		Flags:          _flags,
		Generator:      manifestgen.New(_flags),
		Validator:      validator.New(),
		TypeIdentifier: typeidentifier.New(),
		CarFinder:      carfinder.New(_flags),
		Copier:         copier.New(_flags),
	}
}

func (m *merger) Merge() error {
	var streamFiles []dft.StreamFile
	var dataFiles []dft.DataFile

	log.Info("Creating Output Directory...")
	err := m.CreateOutputDirectory()
	if err != nil {
		return err
	}

	log.Info("Identifying cars in `" + m.Flags.InputPath + "` folder")

	err = filepath.Walk(m.Flags.InputPath, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return err
		}
		if m.Validator.IsValidStreamFile(f.Name()) {
			streamFiles = append(streamFiles, dft.StreamFile{
				Path: path,
				Name: f.Name(),
			})
			return err
		}
		if m.Validator.IsValidDataFile(f.Name()) {
			_type, err := m.TypeIdentifier.IdentifyDataFileType(path)
			if err != nil {
				return err
			}
			dataFile := dft.DataFile{
				Path: path,
				Name: f.Name(),
				Type: _type,
			}

			if dataFile.Type != dft.INVALID {
				dataFiles = append(dataFiles, dataFile)
			}
		}
		return err
	})
	if err != nil {
		return err
	}

	if len(dataFiles) == 0 || len(streamFiles) == 0 {
		log.Error("Cannot find any cars in the specified folder")

		if err := m.Cleanup(); err != nil {
			return err
		}
		return nil
	}

	log.Info("Copying Stream files..")
	err = m.Copier.CopyStreamFilesToOutputDirectory(streamFiles)
	if err != nil {
		return err
	}

	log.Info("Copying Data files...")
	err = m.Copier.CopyDataFilesToOutputDirectory(dataFiles)
	if err != nil {
		return err
	}

	log.Info("Generating fxmanifest.lua")
	err = m.Generator.Generate()
	if err != nil {
		return err
	}

	dataFileCars, err := m.CarFinder.FindDataFileCars()
	if err != nil {
		return err
	}
	streamFileCars, err := m.CarFinder.FindStreamFileCars()
	if err != nil {
		return err
	}

	validCars := sliceutils.RemoveDuplicates(m.CarFinder.FindValidCars(dataFileCars, streamFileCars))

	log.Info("Following cars are valid in the car pack: ", validCars)

	log.Info("Success. Copy the folder `" + m.Flags.OutputPath + "` and paste it into your resources.")

	return nil
}

func (m *merger) CreateOutputDirectory() error {
	if m.Flags.Clean {
		if err := m.Cleanup(); err != nil {
			return err
		}
	}
	err := os.Mkdir(m.Flags.OutputPath, 0755)
	if err != nil {
		return err
	}
	return nil
}

func (m *merger) Cleanup() error {
	if _, err := os.Stat(m.Flags.OutputPath); !os.IsNotExist(err) {
		err = os.RemoveAll(m.Flags.OutputPath)
		if err != nil {
			return err
		}
	}
	return nil
}

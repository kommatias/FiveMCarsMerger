package carfinder

import (
	"github.com/iLLeniumStudios/FiveMCarsMerger/pkg/flags"
	sliceutils "github.com/iLLeniumStudios/FiveMCarsMerger/pkg/utils/slice"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type CarFinder interface {
	FindValidCars(dataFileCars []string, streamFileCars []string) []string
	FindStreamFileCars() ([]string, error)
	FindDataFileCars() ([]string, error)
}

type carFinder struct {
	Flags flags.Flags
}

func New(_flags flags.Flags) CarFinder {
	return &carFinder{Flags: _flags}
}

func (cf *carFinder) FindValidCars(dataFileCars []string, streamFileCars []string) []string {
	var validCars []string
	var noStreamCars, noDataCars []string
	for _, dataFileCar := range dataFileCars {
		if sliceutils.ContainsElement(streamFileCars, dataFileCar) {
			validCars = append(validCars, dataFileCar)
		} else {
			noStreamCars = append(noStreamCars, dataFileCar)
		}
	}

	for _, streamFileCar := range streamFileCars {
		if sliceutils.ContainsElement(dataFileCars, streamFileCar) {
			validCars = append(validCars, streamFileCar)
		} else {
			noDataCars = append(noDataCars, streamFileCar)
		}
	}

	if len(noStreamCars) > 0 {
		log.Warn("Following cars have no stream files: ", noStreamCars)
	}
	if len(noDataCars) > 0 {
		log.Warn("Following cars have no data files: ", noDataCars)
	}

	return validCars
}

func (cf *carFinder) FindStreamFileCars() ([]string, error) {
	var streamFileCars []string
	outputStreamPath := cf.Flags.OutputPath + "/stream"

	files, err := ioutil.ReadDir(outputStreamPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".yft") && !strings.Contains(file.Name(), "_") {
			streamFileCars = append(streamFileCars, strings.ToLower(file.Name()[:len(file.Name())-4]))
		}
	}

	return streamFileCars, nil
}

func (cf *carFinder) FindDataFileCars() ([]string, error) {
	var dataFileCars []string
	outputDataPath := cf.Flags.OutputPath + "/data/vehicles/"

	files, err := ioutil.ReadDir(outputDataPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "vehicles_") {
			xmlFile, err := os.Open(outputDataPath + "/" + file.Name())
			if err != nil {
				return nil, err
			}
			byteValue, err := ioutil.ReadAll(xmlFile)
			if err != nil {
				return nil, err
			}
			re1 := regexp.MustCompile(`<modelName.*?>(.*)</modelName>`)

			matches := re1.FindAllStringSubmatch(string(byteValue), -1)

			for _, v := range matches {
				if !sliceutils.ContainsElement(dataFileCars, strings.ToLower(v[1])) {
					dataFileCars = append(dataFileCars, strings.ToLower(v[1]))
				}
			}

			err = xmlFile.Close()
			if err != nil {
				return nil, err
			}
		}
	}

	return dataFileCars, nil
}

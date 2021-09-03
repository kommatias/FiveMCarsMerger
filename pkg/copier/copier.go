package copier

import (
	"github.com/iLLeniumStudios/FiveMCarsMerger/pkg/dft"
	"github.com/iLLeniumStudios/FiveMCarsMerger/pkg/flags"
	fileutils "github.com/iLLeniumStudios/FiveMCarsMerger/pkg/utils/file"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

type Copier interface {
	CopyStreamFilesToOutputDirectory(streamFiles []dft.StreamFile) error
	CopyDataFilesToOutputDirectory(dataFiles []dft.DataFile) error
}

type copier struct {
	Flags flags.Flags
}

func New(_flags flags.Flags) Copier {
	return &copier{Flags: _flags}
}

func (c *copier) CopyDataFilesToOutputDirectory(dataFiles []dft.DataFile) error {
	err := c.CreateDirectoryInOutput("data")
	if err != nil {
		return err
	}

	dataPath := c.Flags.OutputPath + "/data/"
	dataFileFrequencies := make(map[dft.DataFileType]int)

	for _, dataFile := range dataFiles {
		dataFileFolder := dataPath + strings.ToLower(dataFile.Type.String()) + "/"
		log.Debug("Copying " + dataFile.Name)
		if _, ok := dataFileFrequencies[dataFile.Type]; !ok {
			_ = c.CreateDirectoryInOutput("data/" + strings.ToLower(dataFile.Type.String()))
			dataFileFrequencies[dataFile.Type] = 0
		}
		dataFileFrequencies[dataFile.Type] += 1
		destination := dataFileFolder + strings.ToLower(dataFile.Type.String()) + "_" + strconv.Itoa(dataFileFrequencies[dataFile.Type]) + ".meta"
		_, err = fileutils.CopyFile(dataFile.Path, destination)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *copier) CopyStreamFilesToOutputDirectory(streamFiles []dft.StreamFile) error {
	err := c.CreateDirectoryInOutput("stream")
	if err != nil {
		return err
	}

	streamPath := c.Flags.OutputPath + "/stream/"

	for _, streamFile := range streamFiles {
		log.Debug("Copying " + streamFile.Name)
		_, err := fileutils.CopyFile(streamFile.Path, streamPath+streamFile.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *copier) CreateDirectoryInOutput(name string) error {
	return os.MkdirAll(c.Flags.OutputPath+"/"+name, 0755)
}

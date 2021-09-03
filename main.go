package main

import (
	"encoding/xml"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var flags Flags

func init() {
	flag.BoolVar(&flags.Verbose, "verbose", false, "Enable verbose logging")
	flag.StringVar(&flags.InputPath, "input-path", ".", "Path to all cars")
	flag.StringVar(&flags.OutputPath, "output-path", "out", "Output path")
	flag.BoolVar(&flags.Clean, "clean", false, "Clear output directory")
	flag.Parse()

	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})

	if flags.Verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func main() {
	var streamFiles []StreamFile
	var dataFiles []DataFile

	log.Info("Creating Output Directory...")
	err := CreateOutputDirectory()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Identifying cars in `" + flags.InputPath + "` folder")

	err = filepath.Walk(flags.InputPath, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return err
		}
		if IsValidStreamFile(f.Name()) {
			streamFiles = append(streamFiles, StreamFile{
				Path: path,
				Name: f.Name(),
			})
			return err
		}
		if IsValidDataFile(f.Name()) {
			dataFile := DataFile{
				Path: path,
				Name: f.Name(),
				Type: IdentifyDataFileType(path),
			}

			if dataFile.Type != INVALID {
				dataFiles = append(dataFiles, dataFile)
			}
		}
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Copying Stream files..")
	err = CopyStreamFilesToOutputDirectory(streamFiles)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Copying Data files...")
	err = CopyDataFilesToOutputDirectory(dataFiles)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Generating fxmanifest.lua")
	err = GenerateFXManifest()
	if err != nil {
		log.Fatal(err)
	}

	dataFileCars, err := GetDataFileCars()
	if err != nil {
		log.Fatal(err)
	}
	streamFileCars, err := GetStreamFileCars()
	if err != nil {
		log.Fatal(err)
	}

	validCars := RemoveDuplicates(FindValidCars(dataFileCars, streamFileCars))

	log.Info("Following cars are valid in the car pack: ", validCars)

	log.Info("Success. Copy the folder `" + flags.OutputPath + "` and paste it into your resources.")
}

func FindValidCars(dataFileCars []string, streamFileCars []string) []string {
	var validCars []string
	var noStreamCars, noDataCars []string
	for _, dataFileCar := range dataFileCars {
		if ContainsElement(streamFileCars, dataFileCar) {
			validCars = append(validCars, dataFileCar)
		} else {
			noStreamCars = append(noStreamCars, dataFileCar)
		}
	}

	for _, streamFileCar := range streamFileCars {
		if ContainsElement(dataFileCars, streamFileCar) {
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

func GetStreamFileCars() ([]string, error) {
	var streamFileCars []string
	outputStreamPath := flags.OutputPath + "/stream"

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

func GetDataFileCars() ([]string, error) {
	var dataFileCars []string
	outputDataPath := flags.OutputPath + "/data/vehicles/"

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
				if !ContainsElement(dataFileCars, strings.ToLower(v[1])) {
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

func IdentifyDataFileType(path string) DataFileType {
	xmlFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	byteValue, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		log.Fatal(err)
	}

	defer xmlFile.Close()

	startTag, err := GetStartTag(byteValue)
	if err != nil {
		log.Fatal(err)
	}

	dataFileType := INVALID

	switch startTag {
	case "CVehicleModelInfo__InitDataList":
		dataFileType = VEHICLES
		break
	case "CAmbientModelSets":
		dataFileType = VEHICLEMODELSETS
		break
	case "CVehicleMetadataMgr":
		dataFileType = VEHICLELAYOUTS
		break
	case "SContentUnlocks":
		dataFileType = CONTENTUNLOCKS
		break
	case "CHandlingDataMgr":
		dataFileType = HANDLING
		break
	case "CVehicleModelInfoVariation":
		dataFileType = CARVARIATIONS
		break
	case "CVehicleModelInfoVarGlobal":
		dataFileType = CARCOLS
		break
	case "CWeaponInfoBlob":
		dataFileType = WEAPONSFILE
	case "":
		log.WithField("file", path).Debug("Invalid XML file")
		break
	default:
		log.WithFields(log.Fields{"file": path, "tag": startTag}).Debug("Unknown tag")
		break
	}

	return dataFileType
}

func CopyDataFilesToOutputDirectory(dataFiles []DataFile) error {
	err := CreateDirectoryInOutput("data")
	if err != nil {
		return err
	}

	dataPath := flags.OutputPath + "/data/"
	dataFileFrequencies := make(map[DataFileType]int)

	for _, dataFile := range dataFiles {
		dataFileFolder := dataPath + strings.ToLower(dataFile.Type.String()) + "/"
		log.Debug("Copying " + dataFile.Name)
		if _, ok := dataFileFrequencies[dataFile.Type]; !ok {
			CreateDirectoryInOutput("data/" + strings.ToLower(dataFile.Type.String()))
			dataFileFrequencies[dataFile.Type] = 0
		}
		dataFileFrequencies[dataFile.Type] += 1
		_, err = CopyFile(dataFile.Path, dataFileFolder+strings.ToLower(dataFile.Type.String())+"_"+strconv.Itoa(dataFileFrequencies[dataFile.Type])+".meta")
		if err != nil {
			return err
		}
	}

	return nil
}

func CopyStreamFilesToOutputDirectory(streamFiles []StreamFile) error {
	err := CreateDirectoryInOutput("stream")
	if err != nil {
		return err
	}

	streamPath := flags.OutputPath + "/stream/"

	for _, streamFile := range streamFiles {
		log.Debug("Copying " + streamFile.Name)
		_, err := CopyFile(streamFile.Path, streamPath+streamFile.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateDirectoryInOutput(name string) error {
	return os.MkdirAll(flags.OutputPath+"/"+name, 0755)
}

func CreateOutputDirectory() error {
	if flags.Clean {
		if _, err := os.Stat(flags.OutputPath); !os.IsNotExist(err) {
			err = os.RemoveAll(flags.OutputPath)
			if err != nil {
				return err
			}
		}
	}
	err := os.Mkdir(flags.OutputPath, 0755)
	if err != nil {
		return err
	}
	return nil
}

func HasAnyFileExtension(file string, extensions []string) bool {
	for _, extension := range extensions {
		if strings.HasSuffix(file, extension) {
			return true
		}
	}
	return false
}

func IsValidDataFile(file string) bool {
	validDataExtensions := []string{".meta", ".dat"}
	return HasAnyFileExtension(file, validDataExtensions)
}

func IsValidStreamFile(file string) bool {
	validStreamExtensions := []string{".yft", ".ytd"}
	return HasAnyFileExtension(file, validStreamExtensions)
}

func GetStartTag(bytes []byte) (string, error) {
	decoder := xml.NewDecoder(strings.NewReader(string(bytes)))

	for {
		t, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				return "", nil
			}
			return "", err
		}
		if se, ok := t.(xml.StartElement); ok {
			return se.Name.Local, nil
		}
	}
}

package main

type DataFileType int

const (
	CARCOLS DataFileType = iota + 1
	CARVARIATIONS
	CONTENTUNLOCKS
	HANDLING
	VEHICLELAYOUTS
	VEHICLEMODELSETS
	VEHICLES
	WEAPONSFILE
	INVALID
)

func (d DataFileType) String() string {
	return [...]string{"CARCOLS", "CARVARIATIONS", "CONTENTUNLOCKS", "HANDLING", "VEHICLELAYOUTS", "VEHICLEMODELSETS", "VEHICLES", "WEAPONSFILE", "INVALID"}[d-1]
}

func (d DataFileType) EnumIndex() int {
	return int(d)
}

type DataFile struct {
	Path string
	Name string
	Type DataFileType
}

type StreamFile struct {
	Path string
	Name string
}

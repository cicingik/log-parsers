package config

const (
	JsonFile = `json`
	CsvFile  = `csv`
	YamlFile = `yaml`
	YmlFile  = `yml`

	LevelNameKey  = `level_name`
	ValueKey      = `value`
	TimestampKey  = `timestamp`
	TotalValueKey = `total_value`

	YamlHeader = `---`

	DefaultOutputName = `out`

	MainMenuLevel = `main_menu`
)

var (
	AvailableInputFileList = []string{
		JsonFile,
		CsvFile,
	}

	AvailableOutputFileList = []string{
		JsonFile,
		YmlFile,
		YamlFile,
	}
)

package analytics

type OS uint8

const (
	Windows OS = iota + 1
	MacOS
	Linux
	Android
)

func (o OS) String() string {
	switch o {
	case Windows:
		return "Windows"
	case MacOS:
		return "MacOS"
	case Linux:
		return "Linux"
	case Android:
		return "Android"

	default:
		return "Unknown"
	}
}

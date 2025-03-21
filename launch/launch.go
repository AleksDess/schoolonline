package launch

import "os"

var Launch = "home"

func IsLaunch() {
	_, err := os.Stat("development.dv")

	if os.IsNotExist(err) {
		Launch = "server"
	}
}

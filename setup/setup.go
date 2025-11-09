package setup

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/mrechtien/mixgo/config"
)

func ConnectWifi(wifi *config.Wifi) {
	execWifiSetup(wifi.SSID, wifi.Passwd, wifi.Hidden)
}

func execWifiSetup(ssid string, passwd string, hidden bool) {
	// Options:
	// #> sudo raspi-config nonint do_wifi_ssid_passphrase myssid 'my passphrase'
	// #> sudo nmcli device wifi connect MySSID name MySSID password MyPassword
	// Change priority of network
	// #> nmcli c mod "xfinitywifi" conn.autoconnect-p -10
	command := exec.Command("bash", "-c", fmt.Sprintf("sudo nmcli device wifi connect '%s' name '%s' password '%s' hidden %s", ssid, ssid, passwd, strconv.FormatBool(hidden)))
	cmdOut, _ := command.StdoutPipe()
	err := command.Start()
	if err != nil {
		log.Fatalf("Fatal error calling wifi setup command: %s", command.String())
	}
	err = command.Wait()
	if err != nil {
		log.Fatalf("Fatal error running wifi setup command: %s\n%s", command.String(), err)
	}
	cmdOutput, _ := io.ReadAll(cmdOut)
	log.Printf("Wifi setup command findished with output: %s", cmdOutput)
}

func main() {

	if len(os.Args) != 3 {
		log.Fatalln("Missing arguments: call with wifi SSID and PWD")
	}

	ssid := os.Args[1]
	pwd := os.Args[2]

	execWifiSetup(ssid, pwd, false)
}

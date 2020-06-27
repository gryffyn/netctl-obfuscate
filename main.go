package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkUser() {
	currUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	if currUser.Uid != "0" {
		fmt.Println("Sorry, this must be run as root or with sudo.")
		os.Exit(126)
	}
}

func getPath() (fullpath string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the name of the netctl config file (eg. wlp2s0-ESSID): ")
	netctlFile, _ := reader.ReadString('\n')
	fullpath = filepath.Join("/etc/netctl/", netctlFile)
	fileValid := false
	for fileValid == false {
		if _, err := os.Stat(fullpath); !os.IsNotExist(err) {
			fmt.Println("File name is incorrect. Please try again: ")
			netctlFile, _ = reader.ReadString('\n')
			fullpath = filepath.Join("/etc/netctl/", netctlFile)
		} else {
			fileValid = true
		}
	}
	return strings.TrimSuffix(fullpath, "\n")
}

func getESSIDandKey(path string) (ESSID, Key string) {
	file_raw, err := os.Open(path)
	check(err)
	defer file_raw.Close()
	s := bufio.NewScanner(file_raw)
	for s.Scan() {
		substrs := strings.Split(s.Text(), "=")
		if substrs[0] == "ESSID" {
			ESSID = strings.Trim(substrs[1], "'")
		} else if substrs[0] == "Key" {
			Key_raw := substrs[1]
			regex, err := regexp.Compile("['\"\\\\]\\w+['\"\\\\]")
			check(err)
			Key = regex.FindString(Key_raw)
		}
	}
	return
}

func getPSK(ESSID, Key string) (PSK string) {
	cmdstr := exec.Command("wpa_passphrase", ESSID, Key)
	output, err := cmdstr.Output()
	check(err)
	regex, err2 := regexp.Compile("\\w{64}")
	check(err2)
	PSK_b := regex.Find(output)
	return string(PSK_b)
}

func keyToPSK(file *os.File) {
	s := bufio.NewScanner(file)
	for s.Scan() {
		if strings.Contains(s.Text(), "]") {

		}
	}
}

func main() {
	/* checkUser()
	ESSID := "Test"                       // Test data
	Key := "Thisisatest"
	fmt.Println(getPSK(ESSID, Key)) */
	fmt.Println(getPSK(getESSIDandKey(getPath())))
}

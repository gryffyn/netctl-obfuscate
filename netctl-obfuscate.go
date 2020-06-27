package main

import (
	"bufio"
	"fmt"
	"io"
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
		fmt.Println("[!] netctl-obfuscate requires running as root or using sudo")
		os.Exit(126)
	}
}

func getPath(netctlFile string) (fullpath string) {
	fullpath = filepath.Join("/etc/netctl/", netctlFile)
	fileValid := false
	for fileValid == false {
		if _, err := os.Stat(fullpath); os.IsNotExist(err) {
			check(err)
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

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
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
		if strings.Contains(s.Text(), "Key=") {

		} else {

		}
	}
}

func main() {
	/* checkUser()
	ESSID := "Test"                       // Test data
	Key := "Thisisatest"
	fmt.Println(getPSK(ESSID, Key)) */
	netctlName := os.Args[1]
	path := getPath(netctlName)
	fmt.Println("Saving backup copy of file as " + netctlName + ".orig ...")
	_, err := copyFile(path, string(path+".bak"))
	check(err)
}

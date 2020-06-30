package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
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

func getPath(inputstr string) (fullpath string) {
	tempfile := strings.Split(inputstr, "/")
	netctlFile := tempfile[len(tempfile)-1]
	fullpath = filepath.Join("/etc/netctl/", netctlFile)
	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		check(err)
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
			KeyRaw := substrs[1]
			regex, err2 := regexp.Compile("\\w{64}")
			check(err2)
			if regex.Match([]byte(KeyRaw)) {
				println("Profile already contains key in PSK format, exiting")
				os.Exit(1)
			}
			if KeyRaw[0] == '\\' {
				Key = strings.TrimPrefix(KeyRaw, "\\\"")
			} else if strings.Contains(KeyRaw, "'\"\"") {
				regex, err := regexp.Compile(`['"\\]\w+['"\\]`)
				check(err)
				Key = "\"" + regex.FindString(KeyRaw)
			} else {
				regex, err := regexp.Compile(`['"\\]\w+['"\\]`)
				check(err)
				tmpkey := regex.FindString(KeyRaw)
				Key = strings.Trim(tmpkey, "'")
			}
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

func replaceKey(path, psk string) {
	input, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "Key=") {
			lines[i] = ("Key='" + psk + "'")
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(path, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	checkUser()
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./netctl-obfuscate [netctl profile]")
		return
	}
	netctlName := os.Args[1]
	path := getPath(netctlName)
	essid, key := getESSIDandKey(path)
	psk := getPSK(essid, key)
	fmt.Println("Saving backup copy of file as " + netctlName + ".orig ...")
	_, err := copyFile(path, string(path+".orig"))
	check(err)
	fmt.Println("ESSID: " + essid + "\nKey: " + key + "\nPSK: " + psk)
	replaceKey(path, psk)
	println("done")
}

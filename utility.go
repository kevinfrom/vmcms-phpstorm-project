package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

/*
** Print an error message and then exit with status code 1 (error)
 */
func exitWithErrorMessage(errorMessage string) {
	fmt.Println("ERROR! " + errorMessage)
	os.Exit(1)
}

/*
** Check if given path exists
 */
func exitIfFileDoesNotExist(path string) {
	_, err := os.Stat(path)
	if err != nil {
		exitWithErrorMessage(fmt.Sprintf("%s does not exist!", path))
	}
}

/*
** Create directory if it doesnt exist already
 */
func createDirectoryIfNotExists(path string) {
	_, err := os.Stat(path)

	if err != nil {
		os.Mkdir(path, 0777)
	}
}

/*
** Set working directory relative to binary
 */
func setWorkingDirectory() {
	tries := 0
	found, _ := os.Stat("config.example.json")

	for found == nil {
		os.Chdir("..")
		found, _ = os.Stat("config.example.json")
		tries++

		if tries >= 3 {
			exitWithErrorMessage("Failed to set working directory!")
		}
	}
}

/*
** Open and parse a given json file
 */
func readJsonFile(filePath string) map[string]string {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		exitWithErrorMessage(fmt.Sprintf("Failed to open %s!", filePath))
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result = map[string]string{}
	json.Unmarshal([]byte(byteValue), &result)

	return result
}

/*
** Get path seperator based on operating system
 */
func getPathSeparator() string {
	var pathSeparator string
	osName := runtime.GOOS

	switch osName {
	case "windows":
		pathSeparator = "\\"
	case "darwin":
		pathSeparator = "/"
	default:
		exitWithErrorMessage(fmt.Sprintf("%s operating system is not supported!", osName))
	}

	return pathSeparator
}

/*
** Find files recursively from given path
 */
func findPhpstormExecutable(path string) string {
	// Get immediate files in given path
	files, err := ioutil.ReadDir(path)
	errorMessage := fmt.Sprintf("Failed to find files in %s!", path)
	if err != nil {
		exitWithErrorMessage(errorMessage)
	}

	// Get directories matching regex
	var directories []string
	regex := regexp.MustCompile("^[0-9\\.]+$")
	for _, file := range files {
		if file.IsDir() && regex.MatchString(file.Name()) {
			directories = append(directories, file.Name())
		}
	}

	// Sort by version descending
	sort.Slice(directories, func(i, j int) bool {
		return directories[i] > directories[j]
	})

	// Find phpstorm executable
	var expectedExecutableFileName string
	osName := runtime.GOOS

	switch osName {
	case "windows":
		expectedExecutableFileName = "phpstorm64.exe"
	case "darwin":
		expectedExecutableFileName = "phpstorm"
	default:
		exitWithErrorMessage(fmt.Sprintf("%s operating system is not supported!", osName))
	}

	path = path + getPathSeparator() + directories[0]
	path = path + getPathSeparator() + "bin"
	path = path + getPathSeparator() + expectedExecutableFileName
	exitIfFileDoesNotExist(path)

	return path
}

/*
** Prompts the user for input
 */
func getUserInput(message string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	result, _ := reader.ReadString('\n')

	return strings.TrimSpace(result)
}

/*
** Get domain from stdin
 */
func getDomain() string {
	var domain string

	if len(os.Args) >= 3 {
		domain = os.Args[2]
	}

	if len(domain) <= 0 {
		domain = getUserInput("Please enter a domain: ")

		if len(domain) <= 0 {
			exitWithErrorMessage("Please enter a valid domain!")
		}
	}

	return domain
}

/*
** Write array of strings to file
 */
func writeStringSliceToFile(fileName string, slice []string) {
	content := strings.Join(slice, "\n")
	content = strings.ReplaceAll(content, "\\t", "	")

	data := []byte(content)

	err := os.WriteFile(fileName, data, 0644)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to write %s", fileName)
		exitWithErrorMessage(errorMessage)
	}
}

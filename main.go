package main

func main() {
	setWorkingDirectory()
	conf := getConfig()
	domain := getDomain()
	projectPath := conf.parsedConfig["drive_path"] + getPathSeparator() + domain
	exitIfFileDoesNotExist(projectPath)

	project := getProject(domain, conf)
	project.writeConfigs()
	project.openInPhpstorm()
}

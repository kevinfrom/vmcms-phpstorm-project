package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type project struct {
	domain      string
	config      config
	projectPath string
}

/*
** Get project
 */
func getProject(domain string, config config) project {
	projectPath := config.parsedConfig["projects_path"] + getPathSeparator() + domain
	project := project{
		domain:      domain,
		config:      config,
		projectPath: projectPath,
	}

	createDirectoryIfNotExists(project.projectPath)
	createDirectoryIfNotExists(project.projectPath + getPathSeparator() + ".idea")

	return project
}

/*
** Get config file path
 */
func (project project) getConfigFilePath(path string) string {
	return project.projectPath + getPathSeparator() + ".idea" + getPathSeparator() + path
}

/*
** Replace placeholders in string splice
 */
func (project project) replacePlaceholdersInStringSplice(content []string) []string {
	for key, value := range content {
		value = strings.ReplaceAll(value, "%drive_path%", project.config.parsedConfig["drive_path"])
		value = strings.ReplaceAll(value, "%domain%", project.domain)

		content[key] = value
	}

	return content
}

/*
** Write configs
 */
func (project project) writeConfigs() {
	project.writeDomainConfig()
	project.writeMiscConfig()
	project.writeEncodingsConfig()
	project.writeModulesConfig()
	project.writePhpConfig()
	project.writeVcsConfig()
}

/*
** Write [domain.tld].iml config
 */
func (project project) writeDomainConfig() {
	// Skip if file already exists
	filePath := fmt.Sprintf("%s.iml", project.domain)
	filePath = project.getConfigFilePath(filePath)
	fileInfo, _ := os.Stat(filePath)
	if fileInfo != nil {
		return
	}

	var content []string

	content = append(content, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	content = append(content, "<module type=\"WEB_MODULE\" version=\"4\">")
	content = append(content, "\\t<component name=\"NewModuleRootManager\">")

	// Exclude local folder
	content = append(content, "\\t\\t<content url=\"file://$MODULE_DIR$\">")
	content = append(content, "\\t\\t\\t<excludeFolder url=\"file://$MODULE_DIR$\" />")
	content = append(content, "\\t\\t</content>")

	// Exclude _resized, uploads and backup folders
	content = append(content, "\\t\\t<content url=\"file://%drive_path%/%domain%\">")
	content = append(content, "\\t\\t\\t<excludeFolder url=\"file://%drive_path%/%domain%/www/app/webroot/_resized\" />")
	content = append(content, "\\t\\t\\t<excludeFolder url=\"file://%drive_path%/%domain%/www/app/webroot/upload\" />")
	content = append(content, "\\t\\t\\t<excludeFolder url=\"file://%drive_path%/%domain%/www/backup\" />")
	content = append(content, "\\t\\t</content>")

	content = append(content, "\\t\\t<orderEntry type=\"inheritedJdk\" />")
	content = append(content, "\\t\\t<orderEntry type=\"sourceFolder\" forTests=\"false\" />")

	content = append(content, "\\t</component>")
	content = append(content, "</module>")

	content = project.replacePlaceholdersInStringSplice(content)
	writeStringSliceToFile(filePath, content)
}

/*
** Write misc.xml config
 */
func (project project) writeMiscConfig() {
	// Skip if file already exists
	filePath := project.getConfigFilePath("misc.xml")
	fileInfo, _ := os.Stat(filePath)
	if fileInfo != nil {
		return
	}

	var content []string

	content = append(content, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	content = append(content, "<project version=\"4\">")
	content = append(content, "\\t<component name=\"JavaScriptSettings\">")
	content = append(content, "\\t\\t<option name=\"languageLevel\" value=\"ES6\" />")
	content = append(content, "\\t</component>")
	content = append(content, "</project>")

	writeStringSliceToFile(filePath, content)
}

/*
** Write encodings.xml config
 */
func (project project) writeEncodingsConfig() {
	// Skip if file already exists
	filePath := project.getConfigFilePath("encodings.xml")
	fileInfo, _ := os.Stat(filePath)
	if fileInfo != nil {
		return
	}

	var content []string

	content = append(content, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	content = append(content, "<project version=\"4\">")
	content = append(content, "\\t<component name=\"Encoding\" addBOMForNewFiles=\"with NO BOM\" />")
	content = append(content, "</project>")

	writeStringSliceToFile(filePath, content)
}

/*
** Write modules.xml config
 */
func (project project) writeModulesConfig() {
	// Skip if file already exists
	filePath := project.getConfigFilePath("modules.xml")
	fileInfo, _ := os.Stat(filePath)
	if fileInfo != nil {
		return
	}

	var content []string

	content = append(content, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	content = append(content, "<project version=\"4\">")
	content = append(content, "\\t<component name=\"ProjectModuleManager\">")
	content = append(content, "\\t\\t<modules>")
	content = append(content, "\\t\\t\\t<module fileurl=\"file://$PROJECT_DIR$/.idea/%domain%.iml\" filepath=\"$PROJECT_DIR$/.idea/%domain%.iml\" />")
	content = append(content, "\\t\\t</modules>")
	content = append(content, "\\t</component>")
	content = append(content, "</project>")

	content = project.replacePlaceholdersInStringSplice(content)
	writeStringSliceToFile(filePath, content)
}

/*
** Write php.xml config
 */
func (project project) writePhpConfig() {
	// Skip if file already exists
	filePath := project.getConfigFilePath("php.xml")
	fileInfo, _ := os.Stat(filePath)
	if fileInfo != nil {
		return
	}

	var content []string

	content = append(content, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	content = append(content, "<project version=\"4\">")
	content = append(content, "\\t<component name=\"PhpProjectSharedConfiguration\" php_language_level=\"7.4\">")
	content = append(content, "\\t\\t<option name=\"suggestChangeDefaultLanguageLevel\" value=\"false\" />")
	content = append(content, "\\t</component>")
	content = append(content, "</project>")

	writeStringSliceToFile(filePath, content)
}

/*
** Write vcs.xml
 */
func (project project) writeVcsConfig() {
	// Skip if file already exists
	filePath := project.getConfigFilePath("vcs.xml")
	fileInfo, _ := os.Stat(filePath)
	if fileInfo != nil {
		return
	}

	// Find first theme with is not AdminTheme
	pluginsDir := project.projectPath + getPathSeparator() + "www"
	pluginsDir = pluginsDir + getPathSeparator() + "app"
	pluginsDir = pluginsDir + getPathSeparator() + "plugins"
	pluginsDirInfo, _ := os.Stat(pluginsDir)

	var theme string
	if pluginsDirInfo != nil {
		// Get directories
		files, err := ioutil.ReadDir(pluginsDir)
		errorMessage := fmt.Sprintf("Failed to find theme for in %s!", project.domain)
		if err != nil {
			exitWithErrorMessage(errorMessage)
		}

		for _, file := range files {
			if file.IsDir() && strings.Contains(file.Name(), "Theme") && file.Name() != "AdminTheme" {
				theme = file.Name()
				break
			}
		}
		if len(theme) <= 0 {
			exitWithErrorMessage(errorMessage)
		}
	}

	// Write config
	var content []string

	content = append(content, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	content = append(content, "<project version=\"4\">")
	content = append(content, "\\t<component name=\"VcsDirectoryMappings\">")

	if len(theme) > 0 {
		content = append(content, "\\t\\t<mapping directory=\"%drive_path%/%domain%/www/app/plugins/"+theme+"\" vcs=\"Git\">")
	}

	content = append(content, "\\t</component>")
	content = append(content, "</project>")

	content = project.replacePlaceholdersInStringSplice(content)
	writeStringSliceToFile(filePath, content)
}

/*
** Opens project in phpstorm
 */
func (project project) openInPhpstorm() {
	cmd := exec.Command(project.config.parsedConfig["phpstorm_path"], project.projectPath)
	err := cmd.Start()

	if err != nil {
		errorMessage := fmt.Sprintf("An error occured while trying to open project %s in phpstorm", project.domain)
		exitWithErrorMessage(errorMessage)
	}
}

# Vmcms phpstorm project

A small golang program to easily setup a vmcms6 project for phpstorm and open it quickly.

## Installation

1. Clone this repository

````bash
git clone git@github.com:kevinfrom/vmcms-phpstorm-project.git
````

2. Copy the `config.example.json` file as `config.json` and edit your values.

- The `drive_path` config is where the vmcms projects live (often mounted at `L:` or `Y:`).
- The `projects_path` config is where the phpstorm projects live on your local machine.
- The `phpstorm_path` config is where you have phpstorm installed. The value in the example is used to automatically find the greatest version installed via [JetBrains Toolbox](https://www.jetbrains.com/toolbox-app/).

### Building the source

*DISCLAIMER: requires golang to be installed*

#### Windows

Run the following command in a bash terminal.

````bash
GOOS=windows GOARCH=amd64 go build -o bin/vmcms-phpstorm-project.exe main.go utility.go config.go project.go
```` 

#### Mac

Run the following command in a bash terminal.

````bash
GOOS=darwin GOARCH=amd64 go build -o bin/vmcms-phpstorm-project-mac main.go utility.go config.go project.go
````

## Usage

Simply run the provided executable for your chosen operating system.

### Windows

Double click the `vmcms-phpstorm-project.exe` executable or execute via a terminal (e.g. PowerShell).

### Mac

*DISCLAIMER: Steps 1 and 2 are only necessary for initial setup.*

1. Open a Terminal
    - Consider adding the `vmcms-phpstorm-project` binary to your `$PATH`
2. Set executable permission
    - `chmod +x /path/to/binary/vmcms-phpstorm-project`
3. Executable the binary
    - `/path/to/binary/vmcms-phpstorm-project`

### Arguments

Enter the domain as an argument:

````bash
/path/to/binary/vmcms-phpstorm-project -- vestjyskmarketing.dk
````

## Author

Kevin From <kf@vestjyskmarketing.dk>


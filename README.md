# geng Documentation [WIP]

Fullform: `Gen`erate `G`olang Project

## Introduction
- Inspired by Nest CLI `geng` is command-line interface tool that helps you to initialize, develop, and maintain your Go gin applications
- By default it will generate the clean code architecture of [wesionaryTEAM](https://github.com/wesionaryTEAM/go_clean_architecture)

## Installation
- Pre-requisite: You need to install the golang in your system: [golang installation](https://go.dev/doc/install)
```zsh
go install github.com/mukezhz/geng@latest
```

**Alternative Install**
- Download and execute binary by downloading from [assets](https://github.com/mukezhz/geng/releases)

**Add binary dir to your PATH variable [If geng command didn't work after installation]**
```zsh
// For zsh: [Open.zshrc] & For bash: [Open .bashrc]
// Add the following:
export GO_HOME="$HOME/go"
export PATH="$PATH:$GO_HOME/bin"

// For fish: [Open config.fish]
// Add the following:
fish_add_path -aP $HOME/go/bin
```

## Basic workflow
- Get help
```zsh
geng help
```
- Create a new project
```zsh
geng new <project-name> -m <project-module-name> [-d <directory>]
```
- Use Interactive to create new project
```zsh
geng new
```
- Start project
```zsh
cd <directory>
cp .env.example .env
go mod tidy

go run main.go
```
- Generate module in already existing project
```zsh
// Project needs to initialize git repository to work module generation
geng gen module <module-name>
```

**NOTE:** default supported version is of `golang 1.20`

### TODO LIST:
- [x] generate a new project
- [x] get project module, project name and directory as command line argument
- [x] generate a module
- [x] refactor the code -> make code clean
- [x] modify the parent features `module.go` when new module is added
- [x] implement CI for assets build
- [x] add option to generate constructor of aws, gcp while project generation
- [ ] allow different template options when building project
- [ ] generate test case template
- [ ] generate hexagonal architecture
- [ ] generate route, controller and dtos from bru file

### Diagrams

<img width="457" alt="image" src="https://github.com/mukezhz/geng/assets/43813670/a0b11d39-e077-4038-852f-7b5b0adb27c8">

### Logo

<img width="45" alt="image" src="https://github.com/mukezhz/geng/assets/43813670/da07d8cc-8896-4a13-9b31-099958e65cb4">



**Thank You!!!**

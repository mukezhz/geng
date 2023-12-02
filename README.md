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
- [ ] refactor the code -> make code clean
- [x] modify the parent features `module.go` when new module is added
- [x] implement CI for assets build
- [ ] allow different template options when building project
- [ ] generate test case template

### Diagrams
![image](https://github.com/mukezhz/geng/assets/43813670/98f13c33-320f-4a4d-aa80-92e76e66c2a3)




**Thank You!!!**

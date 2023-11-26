# geng Documentation [WIP]

Fullform: `Gen`erate `G`olang Project

## Introduction
- Similar to Nest CLI `geng` is command-line interface tool that helps you to initialize, develop, and maintain your Nest applications
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
- Generate module in already existing project
```zsh
geng gen module <module-name>
```
### TODO LIST:
- [x] generate a new project
- [x] get project module, project name and directory as command line argument
- [x] generate a module
- [ ] refactor the code -> make code clean
- [ ] modify the parent features `module.go` when new module is added
- [ ] implement CI for assets build
- [ ] allow different template options when building project

**Thank You!!!**
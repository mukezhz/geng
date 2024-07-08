---
outline: deep
title: getting started of geng
description: learn how to install geng in detail.
---

# Getting Started

## Installation

Pre-requisite: You need to install the golang in your system: [golang installation](https://go.dev/doc/install)

```bash
go install github.com/mukezhz/geng@latest
```

### Alternative Install

Download and execute binary by downloading from [assets](https://github.com/mukezhz/geng/releases)

Add binary dir to your PATH variable [If geng command didn't work after installation]

```bash
// For zsh: [Open.zshrc] & For bash: [Open .bashrc]
// Add the following:
export GO_HOME="$HOME/go"
export PATH="$PATH:$GO_HOME/bin"

// For fish: [Open config.fish]
// Add the following:
fish_add_path -aP $HOME/go/bin
```

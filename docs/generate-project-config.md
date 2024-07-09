---
outline: deep
title: Scaffold a project
description: Let's learn how to generate a project using Gen-G
---

# Generate a project using config

## Let's start

Instead of using the following command.
```bash
geng new
```

- Create a `geng.json` file which can contain following value:
```jsonc
{
  "projectName": "todo app",
  "projectModuleName": "github.com/mukezhz/todo",
  "author": "Mukesh Kumar Chaudhary <mukezhz@duck.com>",
  "projectDescription": "A simple todo project",
  "goVersion": "1.21"
  // "directory": ".",
  // "infrastructureName": [],
  // "serviceName": []
}
```

Now enter the project generation command.
```bash
geng new
```

You will get the similar output:
```bash

            GENG: GENERATE GOLANG PROJECT

         ██████╗ ███████╗███╗   ██╗       ██████╗ 
        ██╔════╝ ██╔════╝████╗  ██║      ██╔════╝ 
        ██║  ███╗█████╗  ██╔██╗ ██║█████╗██║  ███╗
        ██║   ██║██╔══╝  ██║╚██╗██║╚════╝██║   ██║
        ╚██████╔╝███████╗██║ ╚████║      ╚██████╔╝
         ╚═════╝ ╚══════╝╚═╝  ╚═══╝       ╚═════╝ 
                                                                                          


        The information you have provided:

        Project Name        💻: Todo App       
        Project Module      📂: github.com/mukezhz/todo
        Project Description 📚: A simple todo project            
        Go Version          🆚: 1.21          
        Author Detail       🤓: Mukesh Kumar Chaudhary <mukezhz@duck.com>       

        💻 Change directory to project:
            cd todo_app

        💾 Initalize git repository:
            git init

        📚 Sync dependencies:
            go mod tidy

        🕵 Copy .env.example to .env:
            cp .env.example .env

        🏃 Start Project 🏃:
            go run main.go app:serve

        Thank You For using 🙏🇳🇵🙏:
```

**Note:** I have filled the project description, author name and go version.

:tada: Now follow the instruction and run the project
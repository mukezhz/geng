---
outline: deep
title: Start a project
description: Let's learn how to execute a project using Gen-G
---

# Start the project

## Let's start

In order to start a project just follow the given instruction which has been printed on project generation step.

1. change the directory if you haven't specified project directory to . (ie. current directory)
```bash
cd [your project name]
```
2. initialize the git
```bash
git init
```
3. install the dependency using
```bash
go mod tidy
```
4. copy the .env.example to .env
```bash
cp .env.example .env
```
5. run the project
```
geng run app:serve
//or
go run . app:serve
```

### Output
```bash
➜  geng git:(main) ✗ cd todo_app                                               
➜  todo_app git:(main) ✗ go mod tidy
go: finding module for package github.com/aws/smithy-go
go: found github.com/aws/smithy-go in github.com/aws/smithy-go v1.20.3
...
➜  todo_app git:(main) ✗ git init   
Initialized empty Git repository in /Users/mukesh.chaudhary/GolandProjects/geng/todo_app/.git/
➜  todo_app git:(main) ✗ cp .env.example .env               
➜  todo_app git:(main) ✗ geng run app:serve
2024-07-08T20:27:08.001+0545    INFO    [GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.


2024-07-08T20:27:08.002+0545    INFO    [GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)


2024-07-08T20:27:08.002+0545    INFO    [GIN-debug] GET    /health-check             --> github.com/mukezhz/todo/pkg/infrastructure.NewRouter.func1 (5 handlers)

2024-07-08T20:27:08.002+0545    INFO    [GIN-debug] GET    /api/hello                --> github.com/mukezhz/todo/domain/hello.(*Controller).HandleRoot-fm (5 handlers)

2024-07-08T20:27:08.002+0545    INFO    commands/serve.go:29    +-----------------------+
2024-07-08T20:27:08.002+0545    INFO    commands/serve.go:30    | GO CLEAN ARCHITECTURE |
2024-07-08T20:27:08.002+0545    INFO    commands/serve.go:31    +-----------------------+
2024-07-08T14:42:08.002Z        INFO    commands/serve.go:49    Running server
2024-07-08T14:42:08.002Z        INFO    [GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.

2024-07-08T14:42:08.002Z        INFO    [GIN-debug] Listening and serving HTTP on :5000
```

### Run

Once you follow all the step correctly you can see your app running at port **:5000**

you can verify it using curl:
```bash
curl http://localhost:5000/health-check
# or simply enter the uri in browser
```
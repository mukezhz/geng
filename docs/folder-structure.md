---
outline: deep
title: Start a project
description: Let's learn how to execute a project using Gen-G
---

# Folder Structure of the project

## Tree of folder structure

```bash
.
├── README.md
├── bootstrap
│   ├── app.go
│   └── module.go
├── console
│   ├── commands
│   │   ├── migration.go
│   │   ├── random.go
│   │   ├── seed.go
│   │   └── serve.go
│   └── console.go
├── docker
│   ├── custom.cnf
│   ├── db.Dockerfile
│   ├── run.sh
│   └── web.Dockerfile
├── docker-compose.yml
├── domain
│   ├── dtos
│   │   └── hello.dto.go
│   ├── hello
│   │   ├── hello.controller.go
│   │   ├── hello.module.go
│   │   ├── hello.repository.go
│   │   ├── hello.route.go
│   │   └── hello.service.go
│   ├── middlewares
│   │   └── module.go
│   ├── models
│   │   └── hello.model.go
│   └── module.go
├── go.mod
├── go.sum
├── main.go
├── migrations
│   ├── hello.go
│   ├── migrator.go
│   └── module.go
├── pkg
│   ├── framework
│   │   ├── command.go
│   │   ├── constants.go
│   │   ├── env.go
│   │   ├── logger.go
│   │   ├── migration.go
│   │   ├── module.go
│   │   └── seed.go
│   ├── infrastructure
│   │   ├── module.go
│   │   └── router.go
│   ├── middlewares
│   │   ├── command.go
│   │   ├── module.go
│   │   └── rate_limitter.go
│   ├── module.go
│   ├── responses
│   │   └── response.go
│   ├── services
│   │   └── module.go
│   └── utils
│       ├── aws_error.go
│       └── functional_programming.go
└── seeds
    ├── hello.go
    ├── module.go
    └── seeder.go

19 directories, 48 files
```

## Folder Structure :file_folder:

| Folder Path                      | Description                                                                                                  |
| -------------------------------- | ------------------------------------------------------------------------------------------------------------ |
| `/bootstrap`                     | contains modules required to start the application                                                           |
| `/console`                       | server commands, run `go run main.go -help` for all the available server commands                            |
| `/docker`                        | `docker` files required for `docker compose`                                                                 |
| `/domain`                        | contains dtos, models, constants and folder for each domain with controller, repository, routes and services |
| `/domain/constants`              | global application constants                                                                                 |
| `/domain/models`                 | ORM models                                                                                                   |
| `/domain/<name>`                 | controller, repository, routes and service for a `domain`. In this template `user` is a domain               |
| `/pkg`                           | contains setup for api_errors, infrastructure, middlewares, external services, utils                         |
| `/pkg/api-errors`                | server error handlers                                                                                        |
| `/pkg/framework`                 | contains env parser, logger...                                                                               |
| `/pkg/infrastructure`            | third-party services connections like `gmail`, `firebase`, `s3-bucket`, ...                                  |
| `/pkg/middlewares`               | all middlewares used in the app                                                                              |
| `/pkg/responses`                 | different types of http responses are defined here                                                           |
| `/pkg/services`                  | service layers, contains the functionality that compounds the core of the application                        |
| `/pkg/types`                     | data types used throught the application                                                                     |
| `/pkg/utils`                     | global utility/helper functions                                                                              |
| `/seeds`                         | seeds for already migrated tables                                                                            |
| `/tests`                         | includes application tests                                                                                   |
| `.env.example`                   | sample environment variables                                                                                 |
| `dbconfig.yml`                   | database configuration file for `sql-migrate` command                                                        |
| `docker-compose.yml`             | `docker compose` file for service application via `Docker`                                                   |
| `main.go`                        | entry-point of the server                                                                                    |
| `Makefile`                       | stores frequently used commands; can be invoked using `make` command                                         |
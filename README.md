# 📖 Tutorial: Build a RESTful API on Go using Fiber and Postgres
Fiber, Postgres, JWT and Swagger docs in isolated Docker containers.

## Getting Started
These instructions will get you a copy of the project up and running on docker container and on your local machine.

### Prerequisites
Prequisites package:
* [Docker](https://www.docker.com/get-started)
* [Go](https://golang.org/) Go Programming Language
* [Swag](https://github.com/swaggo/swag). converts Go annotations to [Swagger Documentation 2.0](https://swagger.io/docs/specification/2-0/basic-structure/)
* [migrate](https://github.com/golang-migrate/migrate) tool for applying migrations.
* [Make](https://golang.org/) Automated Execution using Makefile
* Email for SMTP. If you use gmail, make sure [Less Secure App Setting](https://support.google.com/a/answer/6260879) is enabled.

Optional package:
* [gosec](https://github.com/securego/gosec) Golang Security Checker

### Running In Docker Container
1. Rename `Makefile.example` to `Makefile` and fill it with your make setting.
2. Run project by this command:
```bash
$ make run

# Process:
#   - Generate API docs by Swagger
#   - Create a new Docker network for containers
#   - Build and run Docker containers (Fiber, PostgreSQL)
#   - Apply database migrations (using github.com/golang-migrate/migrate)
```
Stop application by this command:
```bash
$ make stop

# Process:
#   - Stop and remove postgres container
#   - Stop and remove Fiber container
```

### Running On Local Machine
1. Create local database by running command `createdb some_database_name`
1. Rename `run.sh.example` to `run.sh` and fill it with your environment values.
2. Set `run.sh` file permission
```bash
$ chmod +x ./run.sh
```    
3. Run application from terminal by execute this command:
```bash
$ ./run.sh
```
4. Go to your API Docs page: 127.0.0.1:3000/swagger/index.html

To reset database, run following command.
```bash
$ migrate -path ./platform/migrations -database "postgres://username:password@localhost:5432/db_name?sslmode=disable" down
```

### API Access
Go to your API Docs page: [127.0.0.1:3000/swagger/index.html](http://127.0.0.1:3000/swagger/index.html)
<br>
API Docs page will be look like:
<br><img src="https://raw.githubusercontent.com/cooljar/go-fiber-postgres-jwt/main/sc.png" width="500">

## Testing
- Inspects source code for security problems using [gosec](https://github.com/securego/gosec). You need to install it first.
- Execute unit test
```bash
$ make test
```
## Testing
- Inspects source code for security problems using [gosec](https://github.com/securego/gosec). You need to install it first.
- Execute unit test by using following command:
```bash
$ make test
```

## Built With
* [Go](https://golang.org/) - Go Programming Languange
* [Go Modules](https://github.com/golang/go/wiki/Modules) - Go Dependency Management System
* [Make](https://www.gnu.org/software/make/) - GNU Make Automated Execution
* [Docker](https://www.docker.com/) - Application Containerization

## Authors
* **Fajar Rizky** - *Initial Work* - [cooljar](https://github.com/cooljar)

## More
-------
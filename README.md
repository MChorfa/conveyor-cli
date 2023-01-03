# conveyor-cli

conveyor-cli for artifacts

## Developer Setup

```shell
# install Mage a Universal Makefile 
https://magefile.org/
# Git upstream: Keep up-to-date and contribute
https://www.atlassian.com/git/tutorials/git-forks-and-upstreams
```

## Init Project

```shell
# initialize the project
go mod init github.com/MChorfa/conveyor-cli
go mod tidy
# open project in VSCode
code .
```

## CI

```shell
go get dagger.io/dagger@latest
go mod tidy
go run ./ci/main.go
```

## Mage

```shell
mage build
mage test
```

## Build

```shell
# print the version of conveyor
go run main.go --version
conveyor version v0.0.1-alpha

## 
# build the conveyor CLI in version v0.0.1-alpha
go build -o ./dist/conveyor -ldflags="-X 'github.com/MChorfa/conveyor-cli/cmd/conveyor.version=v0.0.1-alpha'" main.go

# verify version is being set correctly
./dist/conveyor --version
> conveyor version v0.0.1-alpha
```

## Run

```shell
# conveyor
go run main.go \
--commit-hash "" \
--owner-name "" \
--pipeline-run-id 124 \
--project-id 123 \
--project-name "" \
--provider-api-url "https://gitlab.youcompany.com/api/v4" \
--provider-token "000" \
--provider-type "gitlab" \
--ref-name "main" \
--storage-token "000" \
--storage-type "azure" \
--storage-account-name "azure" \
--storage-container-name "azure"  
```


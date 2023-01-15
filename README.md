## ---- POC --- DO NOT USE ----
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
--pipeline-id 124 \
--project-id 123 \
--project-name "" \
--job-name generate-dsl \
--job-name generate-otm \
--job-name semgrep-sast \
--provider-api-url "https://gitlab.youcompany.com/api/v4" \
--provider-token "000" \
--provider-type "gitlab" \
--ref-name "main" \
--storage-token "000" \
--storage-type "azure" \
--storage-account-name "azure" \
--storage-container-name "azure"  
```

## Release

```sh
# https://goreleaser.com/quick-start/

brew install goreleaser/tap/goreleaser
goreleaser init
goreleaser build --single-target --snapshot --rm-dist
goreleaser release --snapshot --rm-dist

# The minimum permissions the GITHUB_TOKEN should have to run this are write:packages
export GITHUB_TOKEN="YOUR_GH_TOKEN"

git tag -d v0.0.1-alpha
git push --delete origin v0.0.1-alpha

git tag -a v0.0.1-alpha -m "Alpha pre-release"
git push origin --tags
goreleaser release --rm-dist
```


## CI Examples

### Gitlabci

```yaml

conveyor:
  variables:
    CONVEYOR_PROVIDER_TOKEN: "$CONVEYOR_PROVIDER_TOKEN"
    CONVEYOR_STORAGE_TOKEN: "$CONVEYOR_STORAGE_TOKEN"
    CONVEYOR_STORAGE_ACCOUNT_NAME: "$CONVEYOR_STORAGE_ACCOUNT_NAME"
    CONVEYOR_STORAGE_CONTAINER_NAME: "$CONVEYOR_STORAGE_CONTAINER_NAME"    
  stage: conveyor
  tags:
    - YOUR_TAG

  image: golang:1.19.5
  script:
    - echo "run conveyor"
    - go install github.com/MChorfa/conveyor-cli@latest
    - | 
      conveyor-cli conveyor \
      --commit-hash "$CI_COMMIT_SHA" \
      --owner-name "$CI_COMMIT_AUTHOR" \
      --pipeline-run-id $CI_PIPELINE_ID \
      --project-id $CI_PROJECT_ID \
      --project-name "$CI_PROJECT_NAME" \
      --stage-job-name semgrep-sast \
      --provider-api-url "$CI_API_V4_URL" \
      --provider-token "${CONVEYOR_PROVIDER_TOKEN}" \
      --provider-type "gitlab" \
      --ref-name "$CI_COMMIT_REF_NAME" \
      --storage-token "${CONVEYOR_STORAGE_TOKEN}" \
      --storage-type "azure" \
      --storage-account-name "${CONVEYOR_STORAGE_ACCOUNT_NAME}" \
      --storage-container-name "${CONVEYOR_STORAGE_CONTAINER_NAME}"

```


### GithubAction

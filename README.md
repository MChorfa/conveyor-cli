# ---- POC --- DO NOT USE ----

## conveyor-cli

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

stages:
- conveyor
...

conveyor:
  # Conveyor need all prior jobs in the workflow to be done before running.
  # For example here, we need the semgrep report to be generated in semgrep-sast
  needs: ["semgrep-sast"]
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
      --pipeline-id $CI_PIPELINE_ID \
      --project-id $CI_PROJECT_ID \
      --project-name "$CI_PROJECT_NAME" \
      --job-name semgrep-sast \
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

```yaml
env:
  # see: [secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets)
  # Personal Access Token  | personal use only outside the github actions, recommended to use GITHUB_TOKEN
  CONVEYOR_PROVIDER_TOKEN: ${{ secrets.CONVEYOR_PROVIDER_TOKEN }}
  # Internal Github Token valid only during the workflow lifecycle
  GITHUB_TOKEN: ${{ github.token }}
  # Azure Storage Account token info stored in Environment secrets
  CONVEYOR_STORAGE_TOKEN: ${{ secrets.CONVEYOR_STORAGE_TOKEN }} 
  # See: https://docs.github.com/en/actions/learn-github-actions/contexts#vars-context
  # Setting an environment variable with the value of a configuration variable
  env_var: ${{ vars.ENV_CONTEXT_VAR }}

...

jobs:
  conveyor:
      # Conveyor need all prior jobs in the workflow to be done before running.
      # For example here, we need the sbom to be generated in sbom-stage
      needs: ["sbom-stage"]
      runs-on: ubuntu-latest
      steps:
        - name: Checkout
          uses: actions/checkout@v3
          with:
            fetch-depth: 0
        
        - name: Set up Go
          uses: actions/setup-go@v3
          with:
            go-version: '>=1.19'
        
        - name: Get Conveyor
          run: go install github.com/MChorfa/conveyor-cli@latest

        - name: Run Conveyor
          run: |
            conveyor-cli conveyor \
            --commit-hash "$GITHUB_SHA" \
            --ref-name "$GITHUB_REF_NAME" \
            --pipeline-id $GITHUB_RUN_ID \
            --project-id $GITHUB_REPOSITORY_ID \
            --project-name "conveyor-cli" \
            --owner-name "$GITHUB_REPOSITORY_OWNER" \
            --job-name sbom-stage \ 
            --provider-api-url "https://api.github.com" \
            --provider-token "${{ github.token }}" \
            --provider-type "github" \
            --storage-type "azure" \
            --storage-token "${{ secrets.CONVEYOR_STORAGE_TOKEN }}" \
            --storage-account-name "${{ vars.CONVEYOR_STORAGE_ACCOUNT_NAME }}" \
            --storage-container-name "{{ vars.CONVEYOR_STORAGE_CONTAINER_NAME }}"
```

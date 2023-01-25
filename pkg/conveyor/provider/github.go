package provider

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/MChorfa/conveyor-cli/pkg/conveyor/types"
	"github.com/google/go-github/v49/github"
	"golang.org/x/oauth2"
)

type CGithub struct {
	Configuration *types.Configuration
	Artifacts     []*types.Artifact
}

func NewCGithub(configuration *types.Configuration) IProvider {
	return &CGithub{
		Configuration: configuration,
		Artifacts:     []*types.Artifact{},
	}
}

func (cGithub *CGithub) GetArtifacts() []*types.Artifact {

	if len(cGithub.Configuration.Spec.JobsNames) > 0 {

		repoName := cGithub.Configuration.Spec.ProjectName
		workflowID := cGithub.Configuration.Spec.PipelineID
		// workflowName := cGithub.Configuration.Spec.PipelineName
		ownerName := cGithub.Configuration.Spec.OwnerName
		// refName := cGithub.Configuration.Spec.RefName

		ctx := context.Background()
		// Create an OAuth2 client
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: cGithub.Configuration.Spec.Provider.ProviderToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		// Create a new GitHub client
		client := github.NewClient(tc)
		// GitHub API docs: https://docs.github.com/en/rest/actions/workflow-jobs#list-jobs-for-a-workflow-run
		workflowJobsOptions := &github.ListWorkflowJobsOptions{Filter: *github.String("all")}
		workflowJobs, _, err := client.Actions.ListWorkflowJobs(ctx, ownerName, repoName, int64(workflowID), workflowJobsOptions)
		handleError(err)
		// GitHub API docs: https://docs.github.com/en/rest/actions/workflow-jobs#download-job-logs-for-a-workflow-run
		workflowArtifacts, _, err := client.Actions.ListWorkflowRunArtifacts(ctx, ownerName, repoName, int64(workflowID), &github.ListOptions{})
		handleError(err)

		for _, jobName := range cGithub.Configuration.Spec.JobsNames {

			for _, job := range workflowJobs.Jobs {
				if strings.EqualFold(strings.ToLower(job.GetName()), strings.ToLower(jobName)) {
					for _, workflowArtifact := range workflowArtifacts.Artifacts {

						// GEt Download url
						url, _, err := client.Actions.DownloadArtifact(ctx, ownerName, repoName, workflowArtifact.GetID(), true)
						if err != nil {
							fmt.Printf("The requested job %s do not seems to have an url artifact attached to it. \nPlease make sure the artifact section is configured within your job \nERROR: %v", jobName, err.Error())
						} else {
							// Download the artifact
							response, err := http.Get(url.String())
							handleError(err)
							// Read
							artifactBuf, err := io.ReadAll(response.Body)
							handleError(err)

							defer response.Body.Close()

							cGithub.Artifacts = append(cGithub.Artifacts, &types.Artifact{
								Id:      int(workflowArtifact.GetID()),
								Name:    workflowArtifact.GetName(),
								Payload: bytes.NewReader(artifactBuf),
							})
						}

					}
				}

			}

		}
	}
	fmt.Printf("\nConveyor collected %d artifacts from the workflow", len(cGithub.Artifacts))
	return cGithub.Artifacts
}

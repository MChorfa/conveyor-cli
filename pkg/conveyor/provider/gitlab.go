package provider

import (
	"fmt"
	"log"

	"github.com/MChorfa/conveyor-cli/pkg/conveyor/types"
	"github.com/xanzy/go-gitlab"
)

type CGitlab struct {
	Configuration *types.Configuration
	Artifacts     []*types.Artifact
}

func NewCGitlab(configuration *types.Configuration) IProvider {
	return &CGitlab{
		Configuration: configuration,
		Artifacts:     []*types.Artifact{},
	}
}

func (cGitlab *CGitlab) GetArtifacts() []*types.Artifact {

	artifacts := make([]*types.Artifact, 10)

	if len(cGitlab.Configuration.Spec.StagesAndJobsNames) > 0 {

		projectID := cGitlab.Configuration.Spec.ProjectID
		refName := cGitlab.Configuration.Spec.RefName
		pipelineID := cGitlab.Configuration.Spec.PipelineRunID

		client, err := gitlab.NewClient(cGitlab.Configuration.Spec.Provider.ProviderToken, gitlab.WithBaseURL(cGitlab.Configuration.Spec.Provider.ProviderApiURL))
		handleError(err)

		// TODO: refactor into subfunction to reduce complexcity
		for _, stageOrJobName := range cGitlab.Configuration.Spec.StagesAndJobsNames {
			fmt.Printf("%v", stageOrJobName)
			// This pipelines iteration is temporary once the project is stable will simplify to target a specific pipeline
			pipelines, _, err := client.Pipelines.ListProjectPipelines(projectID, &gitlab.ListProjectPipelinesOptions{Ref: gitlab.String(refName)})
			handleError(err)
			for _, pipeline := range pipelines {

				if pipeline.ID == pipelineID {
					// Get the actual jobs from the pipeline
					jobs, _, err := client.Jobs.ListPipelineJobs(projectID, pipeline.ID, &gitlab.ListJobsOptions{})
					handleError(err)
					// TODO: refactor to use filter
					for _, job := range jobs {

						if stageOrJobName == job.Name {
							// https://docs.gitlab.com/ee/api/jobs.html#get-a-trace-file
							// artifactBuf, _, err := client.Jobs.GetTraceFile(projectID, job.ID)
							// GET /projects/:id/jobs/artifacts/:ref_name/download?job=name
							optDownloadArtifactsFileOptions := &gitlab.DownloadArtifactsFileOptions{Job: &job.Name}
							artifactBuf, _, err := client.Jobs.DownloadArtifactsFile(projectID, cGitlab.Configuration.Spec.RefName, optDownloadArtifactsFileOptions)
							handleError(err)
							artifacts = append(artifacts, &types.Artifact{
								Id:      job.ID,
								Name:    job.Name,
								Payload: artifactBuf,
							})
						}
					}
				}

			}
		}

	}

	return artifacts
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
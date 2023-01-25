package provider

import (
	"fmt"
	"strings"

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

	if len(cGitlab.Configuration.Spec.JobsNames) > 0 {

		projectID := cGitlab.Configuration.Spec.ProjectID
		refName := cGitlab.Configuration.Spec.RefName
		pipelineID := cGitlab.Configuration.Spec.PipelineID

		client, err := gitlab.NewClient(cGitlab.Configuration.Spec.Provider.ProviderToken, gitlab.WithBaseURL(cGitlab.Configuration.Spec.Provider.ProviderApiURL))
		handleError(err)

		// Get the actual jobs from the pipeline
		jobs, _, err := client.Jobs.ListPipelineJobs(projectID, pipelineID, &gitlab.ListJobsOptions{})
		handleError(err)

		for _, jobName := range cGitlab.Configuration.Spec.JobsNames {
			for _, job := range jobs {

				if strings.EqualFold(strings.ToLower(job.Name), strings.ToLower(jobName)) {
					// GET /projects/:id/jobs/artifacts/:ref_name/download?job=name
					optDownloadArtifactsFileOptions := &gitlab.DownloadArtifactsFileOptions{Job: &job.Name}
					artifactBuf, _, err := client.Jobs.DownloadArtifactsFile(projectID, refName, optDownloadArtifactsFileOptions)
					if err != nil {
						fmt.Printf("The requested job %s do not seems to have an artifact attached to it. \nPlease make sure the artifact section is configured within your pipeline \nERROR: %v", jobName, err.Error())
					} else {
						cGitlab.Artifacts = append(cGitlab.Artifacts, &types.Artifact{
							Id:      job.ID,
							Name:    job.Name,
							Payload: artifactBuf,
						})
					}
				}
			}
		}
	}
	fmt.Printf("\nConveyor collected %d artifacts from the pipeline", len(cGitlab.Artifacts))
	return cGitlab.Artifacts
}

package storage

import "github.com/MChorfa/conveyor-cli/pkg/conveyor/types"

type IStorage interface {
	HandleArtifacts(artifacts []*types.Artifact)
	GetStatus() string
}

package storage

import "github.com/MChorfa/conveyor/pkg/conveyor/types"

type IStorage interface {
	HandleArtifacts(artifacts []*types.Artifact)
	GetStatus() string
}

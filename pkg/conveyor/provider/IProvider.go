package provider

import "github.com/MChorfa/conveyor/pkg/conveyor/types"

type IProvider interface {
	GetArtifacts() []*types.Artifact
}

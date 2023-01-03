package provider

import "github.com/MChorfa/conveyor-cli/pkg/conveyor/types"

type IProvider interface {
	GetArtifacts() []*types.Artifact
}

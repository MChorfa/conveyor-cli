package conveyor

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/MChorfa/conveyor-cli/pkg/conveyor/provider"
	"github.com/MChorfa/conveyor-cli/pkg/conveyor/storage"
	"github.com/MChorfa/conveyor-cli/pkg/conveyor/types"
	"github.com/google/uuid"
)

// Config represents Conveyor object.
type Conveyor struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	Nonce         string
	Provider      provider.IProvider
	Storage       storage.IStorage
	Configuration *types.Configuration
}

func NewConveyor(cfg *types.Configuration) *Conveyor {
	return &Conveyor{
		ID:            uuid.New(),
		CreatedAt:     time.Now(),
		Nonce:         "nonce-123564756780tyrgfdcbvxg",
		Configuration: cfg,
	}
}

func (c *Conveyor) SetProvider(cfg *types.Configuration) {
	p, err := getProvider(cfg)
	handleError(err)
	c.Provider = p
}

func (c *Conveyor) SetStorage(cfg *types.Configuration) {
	s, err := getStorage(cfg)
	handleError(err)
	c.Storage = s
}

func (c *Conveyor) ProcessArtifact() error {

	c.Storage.HandleArtifacts(c.Provider.GetArtifacts())

	return nil
}

func getProvider(cfg *types.Configuration) (provider.IProvider, error) {
	if strings.ToLower(string(cfg.Spec.Provider.ProviderType)) == "gitlab" {
		gitlab := provider.NewCGitlab(cfg)
		return gitlab, nil
	} else {
		return nil, fmt.Errorf("invalid provider type")
	}
}

func getStorage(cfg *types.Configuration) (storage.IStorage, error) {
	if strings.ToLower(string(cfg.Spec.Storage.StorageType)) == "azure" {
		storage := storage.NewCAZStorage(cfg)
		return storage, nil
	} else {
		return nil, fmt.Errorf("invalid storage type")
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

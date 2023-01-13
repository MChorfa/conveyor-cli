package types

import (
	"bytes"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Artifact struct {
	Id        int
	Name      string
	Payload   *bytes.Reader
	Nonce     string
	CreatedAt time.Time
}

// TracingSpec defines Tracing configurations.
type TracingSpec struct {
	SamplingRate string     `json:"samplingRate" yaml:"samplingRate"`
	Stdout       bool       `json:"stdout" yaml:"stdout"`
	Zipkin       ZipkinSpec `json:"zipkin" yaml:"zipkin"`
	Otel         OtelSpec   `json:"otel" yaml:"otel"`
}

// ZipkinSpec defines Zipkin exporter configurations.
type ZipkinSpec struct {
	EndpointAddress string `json:"endpointAddress" yaml:"endpointAddress"`
}

// OtelSpec defines Otel exporter configurations.
type OtelSpec struct {
	Protocol        string `json:"protocol" yaml:"protocol"`
	EndpointAddress string `json:"endpointAddress" yaml:"endpointAddress"`
	IsSecure        bool   `json:"isSecure" yaml:"isSecure"`
}

// MetricSpec configuration for metrics.
type MetricSpec struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
}

// APILoggingSpec defines the configuration for API logging.
type APILoggingSpec struct {
	// Default value for enabling API logging. Sidecars can always override this by setting `--enable-api-logging` to true or false explicitly.
	// The default value is false.
	Enabled bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	// If true, health checks are not reported in API logs. Default: false.
	// This option has no effect if API logging is disabled.
	OmitHealthChecks bool `json:"omitHealthChecks,omitempty" yaml:"omitHealthChecks,omitempty"`
}

// LoggingSpec defines the configuration for logging.
type LoggingSpec struct {
	// Configure API logging.
	APILogging APILoggingSpec `json:"apiLogging,omitempty" yaml:"apiLogging,omitempty"`
}

// ProviderSpec defines the configuration for Git Provider.
type ProviderSpec struct {
	ProviderType   RemoteProviderType `yaml:"providerType"  mapstructure:"PROVIDER_TYPE"`
	ProviderApiURL string             `yaml:"providerApiURL"  mapstructure:"PROVIDER_API_URL"`
	ProviderToken  string             `yaml:"providerToken" mapstructure:"PROVIDER_TOKEN"`
}

// Storage types
type RemoteStorageType string

// Provider types
type RemoteProviderType string

// StorageSpec defines the configuration for Remote Storage.
type StorageSpec struct {
	StorageType          RemoteStorageType `yaml:"storageType" mapstructure:"STORAGE_TYPE"`
	StorageToken         string            `yaml:"storageToken" mapstructure:"STORAGE_TOKEN"` //SharedKeySignature
	StorageAccountName   string            `yaml:"storageAccountName" mapstructure:"STORAGE_ACCOUNT_NAME"`
	StorageContainerName string            `yaml:"storageContainerName" mapstructure:"STORAGE_CONTAINER_NAME"`
}

// Enum types
const (
	Gitlab  RemoteProviderType = "gitlab"
	Github  RemoteProviderType = "github"
	azure   RemoteStorageType  = "azure"
	AWSS3   RemoteStorageType  = "awss3"
	MINIOS3 RemoteStorageType  = "minios3"
)

// Configuration specs
type ConfigurationSpec struct {
	PipelineRunID      int           `json:"pipelineRunID" yaml:"pipelineRunID" mapstructure:"PIPELINE_RUN_ID"`
	ProjectID          int           `json:"projectID" yaml:"projectID" mapstructure:"PROJECT_ID"`
	ProjectName        string        `json:"projectName" yaml:"projectName" mapstructure:"PROJECT_NAME"`
	RefName            string        `json:"refName" yaml:"refName" mapstructure:"REF_NAME"`
	OwnerName          string        `json:"ownerName,omitempty" yaml:"ownerName,omitempty"  mapstructure:"OWNER_NAME"`
	CommitHash         string        `json:"commitHash" yaml:"commitHash" mapstructure:"COMMIT_HASH"`
	StagesAndJobsNames []string      `json:"stagesAndJobsNames" yaml:"stagesAndJobsNames" mapstructure:"STAGES_AND_JOBS_NAMES"`
	Storage            *StorageSpec  `json:"storage" yaml:"storage"`
	Provider           *ProviderSpec `json:"provider" yaml:"provider"`
	TracingSpec        TracingSpec   `json:"tracing,omitempty" yaml:"tracing,omitempty"`
	MetricSpec         MetricSpec    `json:"metric,omitempty" yaml:"metric,omitempty"`
	LoggingSpec        LoggingSpec   `json:"logging,omitempty" yaml:"logging,omitempty"`
}

// Configuration that represents Conveyor config object.
type Configuration struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	// See https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	// See https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Spec ConfigurationSpec `json:"spec" yaml:"spec"`
}

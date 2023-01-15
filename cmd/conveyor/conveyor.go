package conveyor

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/MChorfa/conveyor-cli/pkg/conveyor"
	"github.com/MChorfa/conveyor-cli/pkg/conveyor/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	version                    = "v0.0.1-alpha"
	name                       = "conveyor"
	defaultConfigFilename      = "conveyor.yaml"
	envPrefix                  = "CONVEYOR"
	replaceHyphenWithCamelCase = false
)

type ConfigFromFlags struct {
	JobsNames            []string
	PipelineID           int
	PipelineName         string
	ProjectID            int
	ProjectName          string
	RefName              string
	OwnerName            string
	CommitHash           string
	ProviderType         types.RemoteProviderType
	ProviderAPIUrl       string
	ProviderToken        string
	StorageType          types.RemoteStorageType
	StorageToken         string
	StorageAccountName   string
	StorageContainerName string
}

func Execute() {
	initLogging(name, version)
	cmd := NewRootCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Sorry. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

func NewRootCommand() *cobra.Command {
	var jobsNames []string = []string{}
	var pipelineID int = 0
	var pipelineName string = ""
	var projectID int = 0
	var projectName string = ""
	var refName string = ""
	var ownerName string = ""
	var commitHash string = ""
	var providerType string = ""
	var providerAPIUrl string = ""
	var providerToken string = ""
	var storageType string = ""
	var storageToken string = ""
	var storageAccountName string = ""
	var storageContainerName string = ""

	rootCmd := &cobra.Command{
		Use:     "conveyor",
		Version: version,
		Short:   "conveyor - a super easy CLI to upload artifacts from a pipeline or a run ",
		Long:    `conveyor is a super easy CLI to upload gitlab or github artifacts resulting from pipeline or a run to remote storage for more elaborate analysis`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initConfigWithViper(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			config := initConfigFromFlags(
				&ConfigFromFlags{
					PipelineID:           pipelineID,
					PipelineName:         pipelineName,
					JobsNames:            jobsNames,
					ProjectID:            projectID,
					ProjectName:          projectName,
					RefName:              refName,
					OwnerName:            ownerName,
					CommitHash:           commitHash,
					ProviderType:         types.RemoteProviderType(providerType),
					ProviderAPIUrl:       providerAPIUrl,
					ProviderToken:        providerToken,
					StorageType:          types.RemoteStorageType(storageType),
					StorageToken:         storageToken,
					StorageAccountName:   storageAccountName,
					StorageContainerName: storageAccountName,
				})
			err := handleConveyrCmd(config)
			if err != nil {
				fmt.Printf("%v", err)
			}
			out := cmd.OutOrStdout()
			fmt.Fprintln(out)
		},
	}

	rootCmd.Flags().IntVar(&pipelineID, "pipeline-id", 000, "What is the pipeline id or workflow id if available?")
	rootCmd.Flags().StringVar(&pipelineName, "pipeline-name", "", "What is the pipeline name  or workflow name if available?")
	rootCmd.Flags().IntVar(&projectID, "project-id", 000, "What is the project id or repo id ?")
	rootCmd.Flags().StringVar(&projectName, "project-name", "conveyor", "What is the project or repo name?")
	rootCmd.Flags().StringVar(&refName, "ref-name", "main", "What is the project ref name or repo ref name?")
	rootCmd.Flags().StringVar(&ownerName, "owner-name", "conveyor", "What is the project owner name or repo owner name?")
	rootCmd.Flags().StringVar(&commitHash, "commit-hash", "", "What is the latest commit hash?")
	rootCmd.Flags().StringArrayVar(&jobsNames, "job-name", []string{}, "What is the job name?")
	rootCmd.Flags().StringVar(&providerType, "provider-type", "gitlab", "What is the provider type [Gitlab | Github]?")
	rootCmd.Flags().StringVar(&providerAPIUrl, "provider-api-url", "https://gitlab.youcompany.com/api/v4", "What is the provider api url?")
	rootCmd.Flags().StringVar(&providerToken, "provider-token", "000", "What is the provider api token?")
	rootCmd.Flags().StringVar(&storageType, "storage-type", "azure", "What is the storage type?")
	rootCmd.Flags().StringVar(&storageToken, "storage-token", "000", "What is the storage token?")
	rootCmd.Flags().StringVar(&storageAccountName, "storage-account-name", "dev0relay0data", "What is the storage account name?")
	rootCmd.Flags().StringVar(&storageContainerName, "storage-container-name", "raw-data", "What is the storage container name?")
	return rootCmd
}

func handleConveyrCmd(cfg *types.Configuration) error {

	conveyor := &conveyor.Conveyor{
		ID:            uuid.New(),
		CreatedAt:     time.Now(),
		Configuration: cfg,
	}
	conveyor.SetProvider(cfg)
	conveyor.SetStorage(cfg)
	conveyor.ProcessArtifact()
	return nil
}

func initLogging(name, version string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.TimestampFieldName = "ts"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "msg"
	zerolog.ErrorFieldName = "err"
	zerolog.CallerFieldName = "caller"
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func initConfigFromFlags(cfgFromFlags *ConfigFromFlags) *types.Configuration {
	config := &types.Configuration{}
	config.APIVersion = "conveyor.io/v1alpha1"
	config.Kind = "Configuration"
	config.SetName("conveyorConfig")
	config.SetNamespace("conveyor")
	config.Spec.PipelineID = cfgFromFlags.PipelineID
	config.Spec.PipelineName = cfgFromFlags.PipelineName
	config.Spec.ProjectID = cfgFromFlags.ProjectID
	config.Spec.ProjectName = cfgFromFlags.ProjectName
	config.Spec.RefName = cfgFromFlags.RefName
	config.Spec.OwnerName = cfgFromFlags.OwnerName
	config.Spec.CommitHash = cfgFromFlags.CommitHash
	config.Spec.JobsNames = cfgFromFlags.JobsNames
	config.Spec.Provider = &types.ProviderSpec{
		ProviderType:   cfgFromFlags.ProviderType,
		ProviderApiURL: cfgFromFlags.ProviderAPIUrl,
		ProviderToken:  cfgFromFlags.ProviderToken,
	}
	config.Spec.Storage = &types.StorageSpec{
		StorageType:          cfgFromFlags.StorageType,
		StorageToken:         cfgFromFlags.StorageToken,
		StorageAccountName:   cfgFromFlags.StorageAccountName,
		StorageContainerName: cfgFromFlags.StorageContainerName,
	}

	return config
}

func initConfigWithViper(cmd *cobra.Command) error {
	v := viper.New()

	v.SetConfigName(defaultConfigFilename)
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()
	bindFlags(cmd, v)

	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := f.Name
		if replaceHyphenWithCamelCase {
			configName = strings.ReplaceAll(f.Name, "-", "")
		}
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

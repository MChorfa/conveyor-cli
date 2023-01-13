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

func Execute() {
	initLogging(name, version)
	cmd := NewRootCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Sorry. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

func NewRootCommand() *cobra.Command {
	var stagesAndJobsNames []string = []string{}
	var pipelineRunID int = 0
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
			config := initConfigFromFlags(pipelineRunID, projectID, projectName, refName, ownerName, commitHash, stagesAndJobsNames,
				types.RemoteProviderType(providerType), providerAPIUrl, providerToken,
				types.RemoteStorageType(storageType), storageToken, storageAccountName, storageContainerName)
			err := handleConveyrCmd(config)
			if err != nil {
				fmt.Printf("%v", err)
			}
			out := cmd.OutOrStdout()
			fmt.Fprintln(out)
		},
	}

	rootCmd.Flags().IntVarP(&pipelineRunID, "pipeline-run-id", "", 000, "What is the pipeline or the run id?")
	rootCmd.Flags().IntVarP(&projectID, "project-id", "", 000, "What is the project id?")
	rootCmd.Flags().StringVarP(&projectName, "project-name", "", "conveyor", "What is the project name?")
	rootCmd.Flags().StringVarP(&refName, "ref-name", "", "main", "What is the project ref name?")
	rootCmd.Flags().StringVarP(&ownerName, "owner-name", "", "conveyor", "What is the project owner name?")
	rootCmd.Flags().StringVarP(&commitHash, "commit-hash", "", "000", "What is the latest commit hash?")
	rootCmd.Flags().StringArrayVarP(&stagesAndJobsNames, "stage-job-name", "", []string{}, "What is the stages or jobs names?")
	rootCmd.Flags().StringVarP(&providerType, "provider-type", "", "gitlab", "What is the provider type [Gitlab | Github]?")
	rootCmd.Flags().StringVarP(&providerAPIUrl, "provider-api-url", "", "https://gitlab.youcompany.com/api/v4", "What is the provider api url?")
	rootCmd.Flags().StringVarP(&providerToken, "provider-token", "", "000", "What is the provider api token?")
	rootCmd.Flags().StringVarP(&storageType, "storage-type", "", "azure", "What is the storage type?")
	rootCmd.Flags().StringVarP(&storageToken, "storage-token", "", "000", "What is the storage token?")
	rootCmd.Flags().StringVarP(&storageAccountName, "storage-account-name", "", "dev0relay0data", "What is the storage account name?")
	rootCmd.Flags().StringVarP(&storageContainerName, "storage-container-name", "", "raw-data", "What is the storage container name?")
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

func initConfigFromFlags(pipelineRunID int, projectID int, projectName string, refName string, ownerName string, commitHash string,
	stagesAndJobsNames []string, providerType types.RemoteProviderType, providerAPIUrl string, providerToken string,
	storageType types.RemoteStorageType, storageToken string, storageAccountName string, storageContainerName string) *types.Configuration {
	config := &types.Configuration{}
	config.APIVersion = "conveyor.io/v1alpha1"
	config.Kind = "Configuration"
	config.SetName("conveyorConfig")
	config.SetNamespace("conveyor")
	config.Spec.PipelineRunID = pipelineRunID
	config.Spec.ProjectID = projectID
	config.Spec.RefName = refName
	config.Spec.OwnerName = ownerName
	config.Spec.CommitHash = commitHash
	config.Spec.StagesAndJobsNames = stagesAndJobsNames
	config.Spec.Provider = &types.ProviderSpec{
		ProviderType:   providerType,
		ProviderApiURL: providerAPIUrl,
		ProviderToken:  providerToken,
	}
	config.Spec.Storage = &types.StorageSpec{
		StorageType:          storageType,
		StorageToken:         storageToken,
		StorageAccountName:   storageAccountName,
		StorageContainerName: storageContainerName,
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

package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/MChorfa/conveyor-cli/pkg/conveyor/types"
)

type CAZStorage struct {
	Configuration *types.Configuration
	Status        string
}

func NewCAZStorage(configuration *types.Configuration) IStorage {
	return &CAZStorage{
		Configuration: configuration,
		Status:        "none",
	}
}

func (cAZStorage *CAZStorage) GetStatus() string {
	return "idle"
}

func (cAZStorage *CAZStorage) HandleArtifacts(artifacts []*types.Artifact) {

	// Create Transaction Folder
	transaction_dir, err := os.MkdirTemp("", "transaction")
	handleError(err)
	defer os.RemoveAll(transaction_dir)

	for _, artifact := range artifacts {

		if artifact != nil && artifact.Payload != nil {

			// Create The Artifact Metadata File
			artifactMetadata, err := newArtifactMetadata(artifact, cAZStorage.Configuration)
			handleError(err)
			// Define file name
			artifactMetadataFileName := path.Join(transaction_dir, fmt.Sprintf("%s-cfg.yaml", artifact.Name))
			artifactMetadataFile, err := os.Create(artifactMetadataFileName)
			handleError(err)
			// Write the file
			artifactMetadataFile.Write(artifactMetadata)
			defer artifactMetadataFile.Close()

			// Register the artifact payload into a new file
			// Read the payload
			artifactData, err := io.ReadAll(artifact.Payload)
			handleError(err)
			// Define file name
			artifacFileName := path.Join(transaction_dir, fmt.Sprintf("%s.zip", artifact.Name))
			artifactFile, err := os.Create(artifacFileName)
			handleError(err)
			// Write the file
			artifactFile.Write(artifactData)
			defer artifactMetadataFile.Close()

			// Hash filemane
			hashedFileName := getHashedFileName(artifact.Name, artifact.Id)

			// Create archive file object
			archiveFileName := fmt.Sprintf("%x.tar.gz", hashedFileName)
			archiveFilePath := path.Join(transaction_dir, archiveFileName)
			archiveFile, err := os.Create(archiveFilePath)
			handleError(err)

			//handle compression
			err = handleCompression([]string{artifactMetadataFileName, artifacFileName}, archiveFile)
			handleError(err)

			defer archiveFile.Close()

			// Upload File to the remote storage
			upload(cAZStorage.Configuration, archiveFilePath, archiveFileName)
		}
	}

	fmt.Printf("\nConveyor uploaded %d artifacts to the remote storage", len(artifacts))

}

func upload(configuration *types.Configuration, blobFilePath string, blobFileName string) {

	file, err := os.Open(blobFilePath) // Open the file we want to upload
	handleError(err)

	defer func(file *os.File) {
		err := file.Close()
		handleError(err)
	}(file)

	fileSize, err := file.Stat() // Get the size of the file (stream)
	handleError(err)

	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s?%s", configuration.Spec.Storage.StorageAccountName, configuration.Spec.Storage.StorageContainerName, configuration.Spec.Storage.StorageToken)
	containerClient, err := container.NewClientWithNoCredential(containerURL, nil)
	handleError(err)

	// Upload a simple blob.
	blockBlobClient := containerClient.NewBlockBlobClient(blobFileName)
	handleError(err)

	// Pass the Context, stream, stream size, block blob URL, and options to StreamToBlockBlob
	response, err := blockBlobClient.UploadFile(context.TODO(), file,
		&blockblob.UploadFileOptions{
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				fmt.Printf("\nUploaded %d of %d bytes.\n", bytesTransferred, fileSize.Size())
			},
		})
	handleError(err)
	_ = response // Avoid compiler's "declared and not used" error | will be in future iterations
}

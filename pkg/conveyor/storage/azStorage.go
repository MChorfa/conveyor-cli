package storage

import (
	"compress/gzip"
	"context"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/MChorfa/conveyor-cli/pkg/conveyor/types"
)

type CAZStorage struct {
	Configuration *types.Configuration
	Status        string
}

func (cAZStorage *CAZStorage) HandleArtifacts(artifacts []*types.Artifact) {

	// Create Trnasition Folder
	transition_dir, err := os.MkdirTemp("", "transition")
	handleError(err)
	defer os.RemoveAll(transition_dir)
	for _, artifact := range artifacts {

		if artifact != nil && artifact.Payload != nil {
			// read artifact data
			artifactData, err := ioutil.ReadAll(artifact.Payload)
			handleError(err)
			// Hash filemane
			buildFileName := fmt.Sprintf("%s-%d", artifact.Name, artifact.Id)
			hash := sha256.New()
			hash.Write([]byte(buildFileName))
			hashFileName := hash.Sum(nil)
			// Build File object
			finalFileName := fmt.Sprintf("%x.gz", hashFileName)
			transitionFilePath := filepath.Join(transition_dir, finalFileName)
			transitionFile, err := os.Create(transitionFilePath)
			handleError(err)
			// compress file
			w := gzip.NewWriter(transitionFile)
			w.Write(artifactData)
			w.Close()
			// Upload File
			upload(cAZStorage.Configuration, transitionFile.Name(), finalFileName)
		}

	}

}

func (cAZStorage *CAZStorage) GetStatus() string {
	return "idle"
}

func NewCAZStorage(configuration *types.Configuration) IStorage {
	return &CAZStorage{
		Configuration: configuration,
		Status:        "none",
	}
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
				fmt.Printf("Uploaded %d of %d bytes.\n", bytesTransferred, fileSize.Size())
			},
		})
	handleError(err)
	_ = response // Avoid compiler's "declared and not used" error | will be in future iterations
}

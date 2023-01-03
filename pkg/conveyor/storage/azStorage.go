package storage

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/MChorfa/conveyor/pkg/conveyor/types"
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
			// compress data
			data, err := ioutil.ReadAll(artifact.Payload)
			handleError(err)
			compressedData, compressedDataErr := gZipData(data)
			handleError(compressedDataErr)
			// Build File object
			finalFileName := fmt.Sprintf("%s-%d.gz", artifact.Name, artifact.Id)
			transitionFilePath := filepath.Join(transition_dir, finalFileName)
			transitionFile, err := os.Create(transitionFilePath)
			handleError(err)
			w := gzip.NewWriter(transitionFile)
			w.Write(compressedData)
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

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
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
	_ = response // Avoid compiler's "declared and not used" error
}

func gZipData(data []byte) (compressedData []byte, err error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	_, err = gz.Write(data)
	if err != nil {
		return
	}

	if err = gz.Flush(); err != nil {
		return
	}

	if err = gz.Close(); err != nil {
		return
	}

	compressedData = b.Bytes()

	return
}

package storage

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/MChorfa/conveyor-cli/pkg/conveyor/types"
	"gopkg.in/yaml.v2"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func newArtifactMetadata(artifact *types.Artifact, configuration *types.Configuration) ([]byte, error) {
	artifactMetadata, err := yaml.Marshal(&types.ArtifactMetadata{
		Id:           artifact.Id,
		Name:         artifact.Name,
		Nonce:        artifact.Nonce,
		CreatedAt:    artifact.CreatedAt,
		PipelineID:   configuration.Spec.PipelineID,
		PipelineName: configuration.Spec.PipelineName,
		ProjectID:    configuration.Spec.ProjectID,
		ProjectName:  configuration.Spec.ProjectName,
		RefName:      configuration.Spec.RefName,
		OwnerName:    configuration.Spec.OwnerName,
		CommitHash:   configuration.Spec.CommitHash,
	})
	return artifactMetadata, err
}

func handleCompression(fileNames []string, buffer io.Writer) error {

	gzipWriter := gzip.NewWriter(buffer)
	defer gzipWriter.Close()
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()
	for _, f := range fileNames {
		err := appnedFileToCompressor(tarWriter, f)
		if err != nil {
			return err
		}
	}

	return nil
}

func appnedFileToCompressor(tarWriter *tar.Writer, filename string) error {

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}
	header.Name = filename
	err = tarWriter.WriteHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(tarWriter, file)
	if err != nil {
		return err
	}

	return nil
}

func getHashedFileName(name string, id int) []byte {

	buildFileName := fmt.Sprintf("%s-%d", name, id)
	hash := sha256.New()
	hash.Write([]byte(buildFileName))
	hashFileName := hash.Sum(nil)

	return hashFileName
}

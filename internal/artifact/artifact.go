package artifact

import (
	"clyft/internal/utils"
	"context"
	"fmt"
	"os"
	"path/filepath"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/content/oci"
)

const (
	fileMediaType      = "application/octet-stream"
	directoryMediaType = "application/vnd.oci.image.layer.v1.tar+gzip"
	artifactType       = "application/vnd.clyft.artifact.v1"
)

func addPathToOCIStore(filestore *file.Store, ctx context.Context, path string) (ocispec.Descriptor, error) {
	p, err := os.Stat(path)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("error: %w", err)
	}

	kind := "file"
	mediaType := fileMediaType
	if p.IsDir() {
		kind = "directory"
		mediaType = directoryMediaType
	}

	descriptor, err := filestore.Add(ctx, path, mediaType, path)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("error adding %s %s to OCI filestore: %w", kind, path, err)
	}
	return descriptor, nil
}

func Artifact(ctx context.Context, tag string, paths []string) error {
	clyftStoragePath, err := utils.GetClyftStoragePath()
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	filestore, err := file.New("")
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	defer filestore.Close()

	var layers []ocispec.Descriptor

	for _, path := range paths {
		descriptor, err := addPathToOCIStore(filestore, ctx, path)
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}
		layers = append(layers, descriptor)
	}

	ociStore, err := oci.New(filepath.Join(clyftStoragePath, tag))
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	manifestDescriptor, err := oras.PackManifest(ctx, filestore, oras.PackManifestVersion1_1, artifactType, oras.PackManifestOptions{
		Layers: layers,
	})
	if err != nil {
		return fmt.Errorf("error packing manifest: %w", err)
	}

	if err := filestore.Tag(ctx, manifestDescriptor, tag); err != nil {
		return fmt.Errorf("error tagging manifest: %w", err)
	}

	if _, err := oras.Copy(ctx, filestore, tag, ociStore, tag, oras.DefaultCopyOptions); err != nil {
		return fmt.Errorf("error copying artifact to oci store: %w", err)
	}

	return nil
}

package artifact

import (
	"clyft/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func Artifact(tag string, paths []string) error {
	files, err := expandPaths(paths)
	if err != nil {
		return fmt.Errorf("expand paths: %w", err)
	}
	if len(files) == 0 {
		return fmt.Errorf("no files found under %v", paths)
	}

	index, err := buildLayout(tag, files)
	if err != nil {
		return fmt.Errorf("build OCI layout: %w", err)
	}

	clyftStoragePath, err := utils.GetClyftStoragePath()
	if err != nil {
		return fmt.Errorf("detect storage path: %w", err)
	}

	dest := filepath.Join(clyftStoragePath, tag)
	if err := writeLayout(index, dest); err != nil {
		return fmt.Errorf("write layout: %w", err)
	}
	return nil
}

func writeLayerBlob(layer *OCILayer, destPath string) error {
	layerSource, err := os.Open(layer.path)
	if err != nil {
		return fmt.Errorf("open layer source %q: %w", layer.path, err)
	}
	defer layerSource.Close()

	dest, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("create layer blob %q: %w", destPath, err)
	}
	defer dest.Close()

	if _, err := io.Copy(dest, layerSource); err != nil {
		return fmt.Errorf("copy layer %q: %w", layer.path, err)
	}

	if err := dest.Sync(); err != nil {
		return fmt.Errorf("sync layer blob %q: %w", destPath, err)
	}

	return nil
}

func blobName(digest string) (string, error) {
	const prefix = "sha256:"
	if !strings.HasPrefix(digest, prefix) {
		return "", fmt.Errorf("unexpected digest format %q (want %q)", digest, prefix)
	}
	return strings.TrimPrefix(digest, prefix), nil
}

func writeLayout(index *OCIIndex, dest string) (retErr error) {
	parent := filepath.Dir(dest)
	if err := os.MkdirAll(parent, 0755); err != nil {
		return fmt.Errorf("create parent %q: %w", parent, err)
	}

	tmp, err := os.MkdirTemp(parent, filepath.Base(dest)+".tmp.*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}

	defer func() {
		if retErr != nil {
			_ = os.RemoveAll(tmp)
		}
	}()

	blobsDir := filepath.Join(tmp, "blobs", "sha256")
	if err := os.MkdirAll(blobsDir, 0755); err != nil {
		return err
	}

	ociLayout := NewOCILayout()
	ociLayoutData, err := json.Marshal(ociLayout)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(tmp, "oci-layout"), ociLayoutData, 0644); err != nil {
		return err
	}

	for _, manifest := range index.Manifests {
		config := manifest.Config
		configDigest, err := blobName(config.Digest)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(blobsDir, configDigest), []byte(config.Content), 0644); err != nil {
			return err
		}
		for _, layer := range manifest.Layers {
			layerDigest, err := blobName(layer.Digest)
			if err != nil {
				return err
			}
			if err := writeLayerBlob(layer, filepath.Join(blobsDir, layerDigest)); err != nil {
				return err
			}
		}
		manifestDigest, err := blobName(manifest.Digest)
		if err != nil {
			return err
		}
		manifestData, err := json.Marshal(manifest)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(blobsDir, manifestDigest), manifestData, 0644); err != nil {
			return err
		}
	}

	indexData, err := json.Marshal(index)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(tmp, "index.json"), indexData, 0644); err != nil {
		return err
	}

	if err := os.Rename(tmp, dest); err != nil {
		return fmt.Errorf("rename %q -> %q: %w", tmp, dest, err)
	}
	return nil
}

func buildLayout(tag string, files []string) (*OCIIndex, error) {
	var layers []*OCILayer
	for _, path := range files {
		layer, err := NewOCILayer(path)
		if err != nil {
			return nil, err
		}
		layers = append(layers, layer)
	}
	config := NewOCIConfig([]byte(`{}`))
	manifest, err := NewOCIManifest(config, layers)
	if err != nil {
		return nil, err
	}
	manifests := []*OCIManifest{manifest}
	index := NewOCIIndex(manifests, tag)
	return index, nil
}

func expandPaths(paths []string) ([]string, error) {
	var files []string
	for _, root := range paths {
		info, err := os.Stat(root)
		if err != nil {
			return nil, fmt.Errorf("stat %q: %w", root, err)
		}
		if info.IsDir() {
			err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if !d.IsDir() {
					files = append(files, path)
				}
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("walk %q: %w", root, err)
			}
		} else {
			files = append(files, root)
		}
	}
	return files, nil
}

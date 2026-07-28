package artifact

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type OCILayout struct {
	ImageLayoutVersion string `json:"imageLayoutVersion"`
}

func NewOCILayout() *OCILayout {
	return &OCILayout{ImageLayoutVersion: "1.0.0"}
}

type OCILayer struct {
	path        string
	MediaType   string            `json:"mediaType"`
	Annotations map[string]string `json:"annotations"`
	Size        int64             `json:"size"`
	Digest      string            `json:"digest"`
}

func NewOCILayer(path string) (*OCILayer, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat layer file: %w", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read layer file: %w", err)
	}
	h := sha256.Sum256(data)
	return &OCILayer{
		path:      path,
		MediaType: "application/octet-stream",
		Annotations: map[string]string{
			"org.opencontainers.image.title": filepath.Base(path),
		},
		Size:   info.Size(),
		Digest: fmt.Sprintf("sha256:%x", h),
	}, nil
}

type OCIConfig struct {
	MediaType string `json:"mediaType"`
	Content   []byte `json:"-"`
	Digest    string `json:"digest"`
	Size      int64  `json:"size"`
}

func NewOCIConfig(content []byte) *OCIConfig {
	h := sha256.Sum256(content)
	return &OCIConfig{
		MediaType: "application/vnd.oci.empty.v1+json",
		Content:   content,
		Digest:    fmt.Sprintf("sha256:%x", h),
		Size:      int64(len(content)),
	}
}

type OCIManifest struct {
	SchemaVersion int               `json:"schemaVersion"`
	MediaType     string            `json:"mediaType"`
	Config        *OCIConfig        `json:"config"`
	Layers        []*OCILayer       `json:"layers"`
	Annotations   map[string]string `json:"annotations,omitempty"`
	Digest        string            `json:"-"`
	Size          int64             `json:"-"`
}

func NewOCIManifest(config *OCIConfig, layers []*OCILayer) (*OCIManifest, error) {
	m := &OCIManifest{
		SchemaVersion: 2,
		MediaType:     "application/vnd.oci.image.manifest.v1+json",
		Config:        config,
		Layers:        layers,
	}
	data, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("marshal manifest: %w", err)
	}
	h := sha256.Sum256(data)
	m.Digest = fmt.Sprintf("sha256:%x", h)
	m.Size = int64(len(data))
	return m, nil
}

type OCIIndex struct {
	SchemaVersion int               `json:"schemaVersion"`
	Manifests     []*OCIManifest    `json:"manifests"`
	RefName       string            `json:"-"`
	Annotations   map[string]string `json:"annotations,omitempty"`
}

func NewOCIIndex(manifests []*OCIManifest, refName string) *OCIIndex {
	annotations := map[string]string{
		"org.opencontainers.image.ref.name": refName,
		"org.opencontainers.image.created":  time.Now().UTC().Format(time.RFC3339),
	}
	return &OCIIndex{
		SchemaVersion: 2,
		Manifests:     manifests,
		RefName:       refName,
		Annotations:   annotations,
	}
}

type ociIndexManifestDescriptor struct {
	MediaType   string            `json:"mediaType"`
	Digest      string            `json:"digest"`
	Size        int64             `json:"size"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

func (idx *OCIIndex) MarshalJSON() ([]byte, error) {
	descs := make([]ociIndexManifestDescriptor, len(idx.Manifests))
	for i, m := range idx.Manifests {
		descs[i] = ociIndexManifestDescriptor{
			MediaType:   m.MediaType,
			Digest:      m.Digest,
			Size:        m.Size,
			Annotations: idx.Annotations,
		}
	}
	type alias OCIIndex
	return json.Marshal(&struct {
		*alias
		Manifests []ociIndexManifestDescriptor `json:"manifests"`
	}{
		alias:     (*alias)(idx),
		Manifests: descs,
	})
}

package push

import (
	"clyft/internal/utils"
	"context"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"strings"

	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/oci"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"
)

func Push(ctx context.Context, tag string) error {
	clyftStoragePath, err := utils.GetClyftStoragePath()
	if err != nil {
		return fmt.Errorf("resolve clyft storage path: %w", err)
	}

	localStore, err := oci.New(filepath.Join(clyftStoragePath, tag))
	if err != nil {
		return fmt.Errorf("failed to open local OCI layout: %w", err)
	}

	repoURI := utils.GetRepoURI(tag)

	repo, err := remote.NewRepository(repoURI)
	if err != nil {
		return fmt.Errorf("failed to create remote repository client: %w", err)
	}

	registryAuth, err := utils.GetRegistryAuth(repo.Reference.Registry)
	if err != nil {
		return fmt.Errorf("get registry credentials: %w", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(registryAuth.Auth)
	if err != nil {
		return fmt.Errorf("failed to decode registry credentials: %w", err)
	}
	username, password, found := strings.Cut(string(decoded), ":")
	if !found {
		return fmt.Errorf("malformed registry credentials: expected 'username:password'")
	}

	repo.Client = &auth.Client{
		Client: retry.DefaultClient,
		Cache:  auth.NewCache(),
		Credential: func(_ context.Context, _ string) (auth.Credential, error) {
			return auth.Credential{
				Username: username,
				Password: password,
			}, nil
		},
	}

	opts := oras.DefaultCopyOptions
	desc, err := oras.Copy(ctx, localStore, tag, repo, tag, opts)
	if err != nil {
		return fmt.Errorf("failed to push OCI layout to remote registry: %w", err)
	}
	_ = desc
	return nil
}

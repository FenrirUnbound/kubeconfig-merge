package source_test

import (
	"path/filepath"
	"testing"

	"github.com/fenrirunbound/kubeconfig-merge/internal/source"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/tools/clientcmd"
)

func loadFile(filename string) (source.Source, error) {
	absPath, err := filepath.Abs("./test_data/" + filename)
	if err != nil {
		return nil, err
	}

	return source.New(absPath)
}

func TestLocalSourceInterface(t *testing.T) {
	src, err := loadFile("fullconfig.yml")
	assert.NoError(t, err)

	t.Run("contexts", func(t *testing.T) {
		a := assert.New(t)
		contexts := src.Contexts()
		a.Contains(contexts, "exp-scratch")
		a.Contains(contexts, "dev-storage")
		a.Contains(contexts, "dev-frontend")
	})

	t.Run("auth info", func(t *testing.T) {
		a := assert.New(t)
		contexts := src.AuthInfos()
		a.Contains(contexts, "developer")
		a.Contains(contexts, "experimenter")
	})

	t.Run("clusters", func(t *testing.T) {
		a := assert.New(t)
		contexts := src.Clusters()
		a.Contains(contexts, "scratch")
		a.Contains(contexts, "development")
	})
}

func TestCombineWithoutOverwrites(t *testing.T) {
	r := require.New(t)
	first, err := loadFile("config00.yml")
	r.NoError(err)

	second, err := loadFile("config01.yml")
	r.NoError(err)

	fullConfigPath, err := filepath.Abs("./test_data/fullconfig.yml")
	r.NoError(err)

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: fullConfigPath},
		&clientcmd.ConfigOverrides{})
	cfg, err := clientConfig.RawConfig()
	r.NoError(err)

	first.Combine(second, false)

	t.Run("combined contexts", func(t *testing.T) {
		a := assert.New(t)
		ctx := first.Contexts()

		for name, context := range cfg.Contexts {
			a.Contains(ctx, name)

			target := ctx[name]
			a.NotNil(target)
			a.Equal(context.AuthInfo, target.AuthInfo)
			a.Equal(context.Cluster, target.Cluster)
			a.Equal(context.Extensions, target.Extensions)
			a.Equal(context.Namespace, target.Namespace)
		}
	})

	t.Run("combined auth info", func(t *testing.T) {
		a := assert.New(t)
		auth := first.AuthInfos()

		for name, expected := range cfg.AuthInfos {
			a.Contains(auth, name)

			target := auth[name]
			a.NotNil(target)

			// spot check
			a.Equal(expected.ClientCertificate, target.ClientCertificate)
			a.Equal(expected.ClientKey, target.ClientKey)
			a.Equal(expected.Password, target.Password)
			a.Equal(expected.Username, target.Username)
		}
	})

	t.Run("combined clusters", func(t *testing.T) {
		a := assert.New(t)
		clr := first.Clusters()

		for name, expected := range cfg.Clusters {
			a.Contains(clr, name)

			target := clr[name]
			a.NotNil(target)

			// spot check
			a.Equal(expected.CertificateAuthority, target.CertificateAuthority)
			a.Equal(expected.Server, target.Server)
			a.Equal(expected.InsecureSkipTLSVerify, target.InsecureSkipTLSVerify)
		}
	})
}

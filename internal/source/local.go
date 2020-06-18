package source

import (
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type localSource struct {
	config api.Config
}

func (l *localSource) AuthInfos() map[string]*api.AuthInfo {
	return l.config.AuthInfos
}

func (l *localSource) Clusters() map[string]*api.Cluster {
	return l.config.Clusters
}

func (l *localSource) Contexts() map[string]*api.Context {
	return l.config.Contexts
}

func (l *localSource) Combine(src Source, overwrite bool) error {
	for name, ctx := range src.Contexts() {
		l.mergeContexts(name, ctx, overwrite)
	}

	for name, auth := range src.AuthInfos() {
		l.mergeAuthInfos(name, auth, overwrite)
	}

	for name, cluster := range src.Clusters() {
		l.mergeClusters(name, cluster, overwrite)
	}

	return nil
}

func (l *localSource) RawConfig() api.Config {
	return l.config
}

func (l *localSource) mergeAuthInfos(name string, auth *api.AuthInfo, overwrite bool) {
	localAuth := l.AuthInfos()

	if _, exists := localAuth[name]; exists && !overwrite {
		return
	}

	localAuth[name] = auth
}

func (l *localSource) mergeClusters(name string, cluster *api.Cluster, overwrite bool) {
	localCluster := l.Clusters()

	if _, exists := localCluster[name]; exists && !overwrite {
		return
	}

	localCluster[name] = cluster
}

func (l *localSource) mergeContexts(name string, ctx *api.Context, overwrite bool) {
	localCtx := l.config.Contexts

	if _, exists := localCtx[name]; exists && !overwrite {
		return
	}

	localCtx[name] = ctx
}

func newLocalSource(uri string) (Source, error) {
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: uri},
		&clientcmd.ConfigOverrides{})
	cfg, err := clientConfig.RawConfig()
	if err != nil {
		return nil, err
	}

	return &localSource{
		config: cfg,
	}, nil
}

package source

import "k8s.io/client-go/tools/clientcmd/api"

type Source interface {
	AuthInfos() map[string]*api.AuthInfo
	Clusters() map[string]*api.Cluster
	Contexts() map[string]*api.Context
	Combine(Source, bool) error
	RawConfig() api.Config
}

func New(uri string) (Source, error) {
	return newLocalSource(uri)
}

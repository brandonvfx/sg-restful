package main

type clientConfig struct {
	shotgunHost string
	version     string
}

func newClientConfig(version, shotgunHost string) clientConfig {
	return clientConfig{
		shotgunHost: shotgunHost,
		version:     version,
	}
}

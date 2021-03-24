module jolicode.com/docker-starter

go 1.16

require (
	github.com/docker/cli v20.10.5+incompatible
	github.com/docker/compose-cli v1.0.10
	github.com/docker/docker v20.10.5+incompatible // indirect
	github.com/spf13/cobra v1.1.3
	github.com/yuin/gopher-lua v0.0.0-20200816102855-ee81675732da
)

replace github.com/jaguilar/vt100 => github.com/tonistiigi/vt100 v0.0.0-20190402012908-ad4c4a574305

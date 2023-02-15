module github.com/MaibornWolff/elcep

require (
	github.com/MaibornWolff/elcep/main/config v1.2.0
	github.com/MaibornWolff/elcep/main/plugin v1.2.0
	github.com/olivere/elastic v6.2.37+incompatible
	github.com/prometheus/client_golang v1.11.1
)

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/mailru/easyjson v0.0.0-20190614124828-94de47d64c63 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.26.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	golang.org/x/sys v0.0.0-20210603081109-ebe580a85c40 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/protobuf v1.26.0-rc.1 // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

replace (
	github.com/MaibornWolff/elcep/main/config v1.2.0 => ./config
	github.com/MaibornWolff/elcep/main/plugin v1.2.0 => ./plugin
)

go 1.17

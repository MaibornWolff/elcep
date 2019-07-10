module github.com/MaibornWolff/elcep/main/plugin

go 1.12

require (
	github.com/MaibornWolff/elcep/main/config v1.2.0
	github.com/mailru/easyjson v0.0.0-20190614124828-94de47d64c63 // indirect
	github.com/olivere/elastic v6.2.19+incompatible
	github.com/prometheus/client_golang v1.0.0
)

replace github.com/MaibornWolff/elcep/main/plugin => ./

replace github.com/MaibornWolff/elcep/main/config => ../config

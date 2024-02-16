module github.com/gopherguides/hype

go 1.22

toolchain go1.22.0

require (
	github.com/gobuffalo/flect v1.0.2
	github.com/gofrs/uuid/v5 v5.0.0
	github.com/markbates/clam v0.0.0-20220808175708-ef60f46826fb
	github.com/markbates/cleo v0.0.0-20230821202903-72220ef5f7f0
	github.com/markbates/fsx v1.3.0
	github.com/markbates/garlic v1.0.0
	github.com/markbates/hepa v0.0.0-20211129002629-856d16f89b9d
	github.com/markbates/iox v0.0.0-20230829013604-e0813da73cc6
	github.com/markbates/plugins v0.0.0-20230821202759-9443baa9b3df
	github.com/markbates/sweets v0.0.0-20210926032915-062eb9bcc0e5
	github.com/markbates/syncx v1.5.1
	github.com/markbates/table v0.0.0-20230314205021-441ed58296d1
	github.com/mattn/go-shellwords v1.0.12
	github.com/russross/blackfriday/v2 v2.1.0
	github.com/stretchr/testify v1.8.4
	golang.org/x/net v0.21.0
	golang.org/x/sync v0.6.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/exp v0.0.0-20240213143201-ec583247a57a // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/markbates/clam => ../clam

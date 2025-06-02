module github.com/mongodb/mongodb-cli/mongocli/v2

go 1.24.1

require (
	github.com/AlecAivazis/survey/v2 v2.3.7
	github.com/Delta456/box-cli-maker/v2 v2.3.0
	github.com/Masterminds/semver/v3 v3.3.1
	github.com/Netflix/go-expect v0.0.0-20220104043353-73e0943537d2
	github.com/PaesslerAG/jsonpath v0.1.1
	github.com/briandowns/spinner v1.23.2
	github.com/creack/pty v1.1.24
	github.com/evergreen-ci/shrub v0.0.0-20240215220116-3f233ddeff2a
	github.com/fatih/color v1.18.0
	github.com/gemalto/kmip-go v0.0.10
	github.com/go-test/deep v1.1.1
	github.com/golang-jwt/jwt/v4 v4.5.2
	github.com/golang/mock v1.6.0
	github.com/google/go-github/v61 v61.0.0
	github.com/hinshun/vt10x v0.0.0-20220301184237-5011da428d02
	github.com/klauspost/compress v1.18.0
	github.com/mattn/go-isatty v0.0.20
	github.com/mongodb-forks/digest v1.1.0
	github.com/mongodb-labs/cobra2snooty v1.19.1
	github.com/olekukonko/tablewriter v1.0.4
	github.com/pelletier/go-toml v1.9.5
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c
	github.com/spf13/afero v1.14.0
	github.com/spf13/cobra v1.9.1
	github.com/spf13/viper v1.20.1
	github.com/stretchr/testify v1.10.0
	github.com/tangzero/inflector v1.0.0
	go.mongodb.org/atlas v0.38.0
	go.mongodb.org/mongo-driver v1.17.3
	go.mongodb.org/ops-manager v0.60.0
	golang.org/x/crypto v0.38.0
	golang.org/x/mod v0.24.0
	golang.org/x/tools v0.33.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/PaesslerAG/gval v1.0.0 // indirect
	github.com/ansel1/merry v1.6.2 // indirect
	github.com/ansel1/merry/v2 v2.0.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/gemalto/flume v0.13.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gookit/color v1.5.2 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/olekukonko/errors v0.0.0-20250405072817-4e6d85265da6 // indirect
	github.com/olekukonko/ll v0.0.6-0.20250511102614-9564773e9d27 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	go.uber.org/goleak v1.3.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/term v0.32.0 // indirect
	golang.org/x/text v0.25.0 // indirect
)

replace go.mongodb.org/ops-manager => ../go-client-mongodb-ops-manager

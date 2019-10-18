module cdpctl

go 1.12

require (
	github.com/glory-cd/server v0.0.0-20190814085509-565c322cb218
	github.com/mitchellh/go-homedir v1.1.0
	github.com/modood/table v0.0.0-20181112072225-499dc7fba710
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/glory-cd/server => ../server

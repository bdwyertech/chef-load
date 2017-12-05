package main

import (
	"fmt"
	"os"

	"github.com/naoina/toml"
)

type chefLoadConfig struct {
	RunChefClient              bool
	LogFile                    string
	ChefServerURL              string `toml:"chef_server_url"`
	ClientName                 string
	ClientKey                  string
	DataCollectorURL           string `toml:"data_collector_url"`
	DataCollectorToken         string
	OhaiJSONFile               string `toml:"ohai_json_file"`
	ConvergeStatusJSONFile     string `toml:"converge_status_json_file"`
	ComplianceStatusJSONFile   string `toml:"compliance_status_json_file"`
	NumNodes                   int
	Interval                   int
	NodeNamePrefix             string
	ChefEnvironment            string
	RunList                    []string
	SleepDuration              int
	DownloadCookbooks          string
	APIGetRequests             []string `toml:"api_get_requests"`
	ChefVersion                string
	ChefServerCreatesClientKey bool `toml:chef_server_creates_client_key`
	EnableReporting            bool
	RandomData                 bool
}

func defaultConfig() chefLoadConfig {
	return chefLoadConfig{
		RunChefClient:              false,
		LogFile:                    "/var/log/chef-load/chef-load.log",
		ChefServerURL:              "",
		DataCollectorURL:           "",
		DataCollectorToken:         "93a49a4f2482c64126f7b6015e6b0f30284287ee4054ff8807fb63d9cbd1c506",
		OhaiJSONFile:               "",
		ConvergeStatusJSONFile:     "",
		ComplianceStatusJSONFile:   "",
		NumNodes:                   30,
		Interval:                   30,
		NodeNamePrefix:             "chef-load",
		ChefEnvironment:            "_default",
		RunList:                    make([]string, 0),
		SleepDuration:              0,
		DownloadCookbooks:          "never",
		ChefVersion:                "13.2.20",
		ChefServerCreatesClientKey: false,
		EnableReporting:            false,
		RandomData:                 false,
	}
}

func printSampleConfig() {
	sampleConfig := `# log_file specifies the location to log API requests
# log_file = "/var/log/chef-load/chef-load.log"

# The chef_server_url, client_name and client_key parameters must be set if you want
# to make API requests to a Chef Server.
#
# chef-load will also automatically attempt to connect to the Chef Server authenticated data collector proxy.
# If you enabled this feature on the Chef Server, Chef Client run data will automatically be forwarded to Automate.
# If you do not have Automate or the feature is disabled on the Chef Server, chef-load will detect this and
# disable data collection.
#
# Be sure to include the organization name
# For example: chef_server_url = "https://chef.example.com/organizations/demo/"
# chef_server_url = ""
#
# The client defined by client_name needs to be an admin user of the Chef Server org.
# client_name = "CLIENT_NAME"
# client_key = "/path/to/CLIENT_NAME.pem"

# The data_collector_url must be set if you want to make API requests directly to an Automate server.
# For example: data_collector_url = "https://automate.example.org/data-collector/v0/"
# data_collector_url = ""

# The Authorization token for the Automate server.
# The following default value is sufficient unless you set your own token in your Automate server.
# data_collector_token = "93a49a4f2482c64126f7b6015e6b0f30284287ee4054ff8807fb63d9cbd1c506"

# Ohai data will be loaded from this file and used for the nodes' automatic attributes.
# See the chef-load README for instructions for creating an ohai JSON file.
# ohai_json_file = "/path/to/example-ohai.json"

# Data from a converge status report will be loaded from this file and used
# for each node's converge status report that is sent to the Automate server.
# See the chef-load README for instructions for creating a converge status JSON file.
# converge_status_json_file = "/path/to/example-converge-status.json"

# Data from a compliance status report will be loaded from this file and used
# for each node's compliance status report that is sent to the Automate server.
# See the chef-load README for instructions for creating a compliance status JSON file.
# compliance_status_json_file = "/path/to/example-compliance-status.json"

# chef-load will evenly distribute the number of nodes across the desired interval (minutes)
# Examples:
#   30 nodes / 30 minute interval =  1 chef-client run per minute
# 1800 nodes / 30 minute interval = 60 chef-client runs per minute
# num_nodes = 30
# interval = 30

# This prefix will go at the beginning of each node name.
# This enables running multiple instances of chef-load without affecting each others' nodes
# For example, a value of "chef-load" will result in nodes named "chef-load-1", "chef-load-2", ...
# node_name_prefix = "chef-load"

# Chef environment used for each node
# chef_environment = "_default"

# run_list is the run list used for each node. It should be a list of strings.
# For example: run_list = [ "role[role_name]", "recipe_name", "recipe[different_recipe_name@1.0.0]" ]
# The default value is an empty run_list.
# run_list = [ ]

# sleep_duration is an optional setting that is available to provide a delay to simulate
# the amount of time a Chef Client takes actually converging all of the run list's resources.
# sleep_duration is measured in seconds
# sleep_duration = 0

# download_cookbooks controls which chef-client run downloads cookbook files.
# Options are: "never", "first" (first chef-client run only), "always"
#
# Downloading cookbooks can significantly increase the number of API requests that chef-load
# makes depending on the run_list. If you aren't concerned with simulating the download of cookbook files
# then the recommendation is to use "never" or "first".
#
# download_cookbooks = "never"

# api_get_requests is an optional list of API GET requests that are made during the chef-client run.
# This is used to simulate the API requests that the cookbooks would make.
# For example, it can make Chef Search or data bag item requests.
# The values can be either full URLs that include the chef_server_url portion or just the portion of
# the URL that comes after the chef_server_url.
# For example, to make a Chef Search API request that searches for all nodes you can use either of the
# following values:
#
# "https://chef.example.com/organizations/orgname/search/node?q=*%253A*&sort=X_CHEF_id_CHEF_X%20asc&start=0"
# or
# "search/node?q=*%253A*&sort=X_CHEF_id_CHEF_X%20asc&start=0"
#
# api_get_requests = [ ]

# chef_version sets the value of the X-Chef-Version HTTP header in API requests sent to the Chef Server.
# This value represents the version of the Chef Client making the API requests. The default is "13.2.20"
# chef_version = "13.2.20"

# Ever since Chef Client 12.x was released the default behavior has been for the Chef Client to create its
# own client key locally and then upload the public side to the Chef Server when it creates the client object.
# chef-load simulates this behavior. However, if you want chef-load to ask the Chef Server to create a client key
# when the client object is created then set chef_server_creates_client_key to true.
# chef_server_creates_client_key = false

# Send data to the Chef server's Reporting service
# enable_reporting = false

# Generate Random Data
# random_data = true
`
	fmt.Print(sampleConfig)
}

func loadConfig(file string) (*chefLoadConfig, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Initialize default configuration values
	config := defaultConfig()

	if err = toml.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

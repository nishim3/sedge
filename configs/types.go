/*
Copyright 2022 Nethermind

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package configs

// LogConfig : Struct Represent configuration data for logging
type LogConfig struct {
	Level string `mapstructure:"logLevel"`
}

type NetworkConfig struct {
	Name                      string
	NoJWT                     bool
	NetworkService            string
	GenesisForkVersion        string
	DefaultECBootnodes        []string
	DefaultCCBootnodes        []string
	DefaultCustomChainSpecSrc string
	DefaultCustomConfigSrc    string
	DefaultCustomGenesisSrc   string
	DefaultCustomDeployBlock  string
	SupportsMEVBoost          bool
	CheckpointSyncURL         string
	RelayURLs                 []string
	Weight                    int // Weight of the network for sorting purposes
}

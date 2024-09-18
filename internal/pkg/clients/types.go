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
package clients

import "github.com/NethermindEth/sedge/configs"

// Client : Struct Represent a client like geth, prysm, etc
type Client struct {
	Name      string
	Type      string
	Image     string
	Endpoint  string
	Supported bool
}

func (c *Client) SetImageOrDefault(image string) {
	switch c.Type {
	case "validator":
		c.setValidatorImage(image)
	case "consensus":
		c.setConsensusImage(image)
	case "execution":
		c.setExecutionImage(image)
	case "optimism":
		c.setOptimismImage(image)
	}
}

func (c *Client) setExecutionImage(image string) {
	switch c.Name {
	case "geth":
		c.Image = valueOrDefault(image, configs.ClientImages.Execution.Geth.String())
	case "besu":
		c.Image = valueOrDefault(image, configs.ClientImages.Execution.Besu.String())
	case "nethermind":
		c.Image = valueOrDefault(image, configs.ClientImages.Execution.Nethermind.String())
	case "erigon":
		c.Image = valueOrDefault(image, configs.ClientImages.Execution.Erigon.String())
	case "reth":
		c.Image = valueOrDefault(image, configs.ClientImages.Execution.Reth.String())
	}
}

func (c *Client) setConsensusImage(image string) {
	switch c.Name {
	case "lighthouse":
		c.Image = valueOrDefault(image, configs.ClientImages.Consensus.Lighthouse.String())
	case "prysm":
		c.Image = valueOrDefault(image, configs.ClientImages.Consensus.Prysm.String())
	case "teku":
		c.Image = valueOrDefault(image, configs.ClientImages.Consensus.Teku.String())
	case "lodestar":
		c.Image = valueOrDefault(image, configs.ClientImages.Consensus.Lodestar.String())
	}
}

func (c *Client) setValidatorImage(image string) {
	switch c.Name {
	case "lighthouse":
		c.Image = valueOrDefault(image, configs.ClientImages.Validator.Lighthouse.String())
	case "prysm":
		c.Image = valueOrDefault(image, configs.ClientImages.Validator.Prysm.String())
	case "teku":
		c.Image = valueOrDefault(image, configs.ClientImages.Validator.Teku.String())
	case "lodestar":
		c.Image = valueOrDefault(image, configs.ClientImages.Validator.Lodestar.String())
	}
}

func (c *Client) setOptimismImage(image string) {
	switch c.Name {
	case "optimism":
		c.Image = valueOrDefault(image, configs.ClientImages.Optimism.OpNode.String())
	}
}

func valueOrDefault(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

// Clients : Struct Represent a combination of execution, consensus and validator clients
type Clients struct {
	Execution   *Client
	Consensus   *Client
	Validator   *Client
	ExecutionOP *Client
	Optimism    *Client
}

type ClientMap map[string]*Client

type OrderedClients map[string]ClientMap

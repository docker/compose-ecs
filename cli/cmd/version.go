/*
   Copyright 2020 Docker Compose CLI authors

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

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/docker/compose-ecs/internal"
)

const formatOpt = "format"

// VersionCommand command to display version
func VersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show the Docker version information",
		Args:  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			runVersion()
			return nil
		},
	}
	// define flags for backward compatibility with com.docker.cli
	flags := cmd.Flags()
	flags.StringP(formatOpt, "f", "", "Format the output. Values: [pretty | json]. (Default: pretty)")

	return cmd
}

func runVersion() {
	if formatOpt == "json" {
		fmt.Printf("{\"version\":%q}\n", internal.Version)
	} else {
		fmt.Printf("Compose ECS %s\n", internal.Version)
	}
}

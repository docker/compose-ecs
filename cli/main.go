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

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/docker/compose-ecs/ecs"
	compose2 "github.com/docker/compose/v2/cmd/compose"
	"github.com/spf13/cobra"

	"github.com/docker/compose-ecs/api/backend"
	"github.com/docker/compose-ecs/cli/cmd"
	"github.com/docker/compose-ecs/cli/cmd/volume"
	// Backend registrations
)

func main() {
	root := &cobra.Command{
		Use:              "compose-ecs",
		SilenceErrors:    true,
		SilenceUsage:     true,
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			return fmt.Errorf("unknown command: %q", args[0])
		},
	}

	root.AddCommand(
		cmd.VersionCommand(),
		cmd.SecretCommand(),
		volume.Command(),
	)

	ctx, cancel := newSigContext()
	defer cancel()

	service, err := ecs.NewComposeECS()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	backend.WithBackend(service)

	command := compose2.RootCommand(service.ComposeService())

	for _, c := range command.Commands() {
		switch c.Name() {
		case "convert", "down", "logs", "ps", "up": // compose-ecs only implement a subset of compose commands
			root.AddCommand(c)
		}
	}

	err = root.ExecuteContext(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func newSigContext() (context.Context, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-s
		cancel()
	}()
	return ctx, cancel
}

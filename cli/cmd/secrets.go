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
	"io"
	"os"

	"github.com/docker/compose-ecs/api/backend"
	"github.com/docker/compose-ecs/api/secrets"

	"github.com/docker/compose/v2/cmd/formatter"
	"github.com/spf13/cobra"
)

// SecretCommand manage secrets
func SecretCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secret",
		Short: "Manages secrets",
	}

	cmd.AddCommand(
		createSecret(),
		inspectSecret(),
		listSecrets(),
		deleteSecret(),
	)
	return cmd
}

func createSecret() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [OPTIONS] SECRET [file|-]",
		Short: "Creates a secret.",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := backend.Current()
			file := "-"
			if len(args) == 2 {
				file = args[1]
			}
			fmt.Println("create secret from " + file)
			if len(file) == 0 {
				return fmt.Errorf("secret data source empty: %q", file)
			}
			var in io.ReadCloser
			switch file {
			case "-":
				in = os.Stdin
			default:
				f, err := os.Open(file)
				if err != nil {
					return err
				}
				in = f
				defer in.Close() //nolint:errcheck
			}
			content, err := io.ReadAll(in)
			if err != nil {
				return fmt.Errorf("failed to read content from %q: %v", file, err)
			}
			name := args[0]
			secret := secrets.NewSecret(name, content)
			id, err := c.SecretsService().CreateSecret(cmd.Context(), secret)
			if err != nil {
				return err
			}
			fmt.Println(id)
			return nil
		},
	}
	return cmd
}

func inspectSecret() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect ID",
		Short: "Displays secret details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := backend.Current()
			secret, err := c.SecretsService().InspectSecret(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			out, err := secret.ToJSON()
			if err != nil {
				return err
			}
			fmt.Println(out)
			return nil
		},
	}
	return cmd
}

type listSecretsOpts struct {
	format string
	quiet  bool
}

func listSecrets() *cobra.Command {
	var opts listSecretsOpts
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List secrets stored for the existing account.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := backend.Current()
			secretsList, err := c.SecretsService().ListSecrets(cmd.Context())
			if err != nil {
				return err
			}
			if opts.quiet {
				for _, s := range secretsList {
					fmt.Println(s.ID)
				}
				return nil
			}
			view := viewFromSecretList(secretsList)
			return formatter.Print(view, opts.format, os.Stdout, func(w io.Writer) {
				for _, secret := range view {
					_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", secret.ID, secret.Name, secret.Description)
				}
			}, "ID", "NAME")
		},
	}
	cmd.Flags().StringVar(&opts.format, "format", "", "Format the output. Values: [pretty | json]. (Default: pretty)")
	cmd.Flags().BoolVarP(&opts.quiet, "quiet", "q", false, "Only display IDs")
	return cmd
}

type secretView struct {
	ID          string
	Name        string
	Description string
}

func viewFromSecretList(secretList []secrets.Secret) []secretView {
	retList := make([]secretView, len(secretList))
	for i, s := range secretList {
		retList[i] = secretView{
			ID:   s.ID,
			Name: s.Name,
		}
	}
	return retList
}

type deleteSecretOptions struct {
	recover bool
}

func deleteSecret() *cobra.Command {
	opts := deleteSecretOptions{}
	cmd := &cobra.Command{
		Use:     "delete NAME",
		Aliases: []string{"rm", "remove"},
		Short:   "Removes a secret.",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := backend.Current()
			return c.SecretsService().DeleteSecret(cmd.Context(), args[0], opts.recover)
		},
	}
	cmd.Flags().BoolVar(&opts.recover, "recover", false, "Enable recovery.")
	return cmd
}

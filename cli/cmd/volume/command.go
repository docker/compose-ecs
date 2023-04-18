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

package volume

import (
	"context"
	"fmt"

	"github.com/docker/compose-ecs/api/backend"
	"github.com/docker/compose-ecs/ecs"

	format "github.com/docker/compose/v2/cmd/formatter"
	"github.com/docker/compose/v2/pkg/progress"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
)

// Command manage volumes
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volume",
		Short: "Manages volumes",
	}

	cmd.AddCommand(
		createVolume(),
		listVolume(),
		rmVolume(),
		inspectVolume(),
	)
	return cmd
}

func createVolume() *cobra.Command {
	var opts interface{}
	cmd := &cobra.Command{
		Use:   "create [OPTIONS] VOLUME",
		Short: "Creates an EFS filesystem to use as AWS volume.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			c := backend.Current()
			result, err := progress.RunWithStatus(ctx, func(ctx context.Context) (string, error) {
				volume, err := c.VolumeService().Create(ctx, args[0], opts)
				if err != nil {
					return "", err
				}
				return volume.ID, nil
			})
			if err != nil {
				return err
			}
			fmt.Println(result)
			return nil
		},
	}

	ecsOpts := ecs.VolumeCreateOptions{}
	cmd.Flags().StringVar(&ecsOpts.KmsKeyID, "kms-key", "", "ID of the AWS KMS CMK to be used to protect the encrypted file system")
	cmd.Flags().StringVar(&ecsOpts.PerformanceMode, "performance-mode", "", "performance mode of the file system. (generalPurpose|maxIO)")
	cmd.Flags().Float64Var(&ecsOpts.ProvisionedThroughputInMibps, "provisioned-throughput", 0, "throughput in MiB/s (1-1024)")
	cmd.Flags().StringVar(&ecsOpts.ThroughputMode, "throughput-mode", "", "throughput mode (bursting|provisioned)")
	opts = &ecsOpts
	return cmd
}

func rmVolume() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm [OPTIONS] VOLUME [VOLUME...]",
		Short: "Remove one or more volumes.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := backend.Current()
			var errs *multierror.Error
			for _, id := range args {
				err := c.VolumeService().Delete(cmd.Context(), id, nil)
				if err != nil {
					errs = multierror.Append(errs, err)
					continue
				}
				fmt.Println(id)
			}
			format.SetMultiErrorFormat(errs)
			return errs.ErrorOrNil()
		},
	}
	return cmd
}

func inspectVolume() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect VOLUME [VOLUME...]",
		Short: "Inspect one or more volumes.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := backend.Current()
			v, err := c.VolumeService().Inspect(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			outJSON, err := format.ToStandardJSON(v)
			if err != nil {
				return err
			}
			fmt.Print(outJSON)
			return nil
		},
	}
	return cmd
}

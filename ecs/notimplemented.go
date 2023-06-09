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

package ecs

import (
	"context"

	"github.com/compose-spec/compose-go/types"
	"github.com/docker/compose/v2/pkg/api"
)

func (b *ComposeECS) Build(ctx context.Context, project *types.Project, options api.BuildOptions) error {
	return api.ErrNotImplemented
}

func (b *ComposeECS) Push(ctx context.Context, project *types.Project, options api.PushOptions) error {
	return api.ErrNotImplemented
}

func (b *ComposeECS) Pull(ctx context.Context, project *types.Project, options api.PullOptions) error {
	return api.ErrNotImplemented
}

func (b *ComposeECS) Create(ctx context.Context, project *types.Project, opts api.CreateOptions) error {
	return api.ErrNotImplemented
}

func (b *ComposeECS) Start(ctx context.Context, project *types.Project, options api.StartOptions) error {
	return api.ErrNotImplemented
}

func (b *ComposeECS) Restart(ctx context.Context, project *types.Project, options api.RestartOptions) error {
	return api.ErrNotImplemented
}

func (b *ComposeECS) Stop(ctx context.Context, project *types.Project, options api.StopOptions) error {
	return api.ErrNotImplemented
}

func (b *ComposeECS) Pause(ctx context.Context, project string, options api.PauseOptions) error {
	return api.ErrNotImplemented
}

func (b *ComposeECS) UnPause(ctx context.Context, project string, options api.PauseOptions) error {
	return api.ErrNotImplemented
}

func (b *ComposeECS) Events(ctx context.Context, project string, options api.EventsOptions) error {
	return api.ErrNotImplemented
}

func (b *ComposeECS) Port(ctx context.Context, project string, service string, port int, options api.PortOptions) (string, int, error) {
	return "", 0, api.ErrNotImplemented
}

func (b *ComposeECS) Copy(ctx context.Context, project *types.Project, options api.CopyOptions) error {
	return api.ErrNotImplemented
}

func (b *ComposeECS) RunOneOffContainer(ctx context.Context, project *types.Project, opts api.RunOptions) (int, error) {
	return 0, api.ErrNotImplemented
}

func (b *ComposeECS) Remove(ctx context.Context, project *types.Project, options api.RemoveOptions) error {
	return api.ErrNotImplemented
}

func (b *ComposeECS) Images(ctx context.Context, projectName string, options api.ImagesOptions) ([]api.ImageSummary, error) {
	return nil, api.ErrNotImplemented
}

func (b *ComposeECS) Top(ctx context.Context, projectName string, services []string) ([]api.ContainerProcSummary, error) {
	return nil, api.ErrNotImplemented
}

func (b *ComposeECS) Exec(ctx context.Context, project string, opts api.RunOptions) (int, error) {
	return 0, api.ErrNotImplemented
}

func (b *ComposeECS) Kill(ctx context.Context, project *types.Project, options api.KillOptions) error {
	return api.ErrNotImplemented
}

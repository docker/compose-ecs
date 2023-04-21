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

package backend

import (
	"github.com/docker/compose/v2/pkg/api"

	"github.com/docker/compose-ecs/api/secrets"
	"github.com/docker/compose-ecs/api/volumes"
)

var instance Service

// Current return the active backend instance
func Current() Service {
	return instance
}

// WithBackend set the active backend instance
func WithBackend(s Service) {
	instance = s
}

// Service aggregates the service interfaces
type Service interface {
	ComposeService() api.Service
	SecretsService() secrets.Service
	VolumeService() volumes.Service
}

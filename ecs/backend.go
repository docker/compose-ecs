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
	"os"

	"github.com/docker/compose-ecs/api/secrets"
	"github.com/docker/compose-ecs/api/volumes"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/docker/compose/v2/pkg/api"
)

func NewComposeECS() (*ComposeECS, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           os.Getenv("AWS_PROFILE"),
	})
	if err != nil {
		return nil, err
	}

	sdk := newSDK(sess)
	return &ComposeECS{
		Region: *sess.Config.Region,
		aws:    sdk,
	}, nil
}

type ComposeECS struct {
	Region string
	aws    API
}

func (b *ComposeECS) ComposeService() api.Service {
	return b
}

func (b *ComposeECS) SecretsService() secrets.Service {
	return b
}

func (b *ComposeECS) VolumeService() volumes.Service {
	return ecsVolumeService{backend: b}
}

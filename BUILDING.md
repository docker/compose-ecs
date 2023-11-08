
### Prerequisites

* Windows:
  * [Docker Desktop](https://hub.docker.com/editions/community/docker-ce-desktop-windows)
  * make
* macOS:
  * [Docker Desktop](https://hub.docker.com/editions/community/docker-ce-desktop-mac)
  * make
* Linux:
  * [Docker 19.03 or later](https://docs.docker.com/engine/install/)
  * make

### Building the CLI

Once you have the prerequisites installed, you can build the CLI using:

```console
make
```

This will output a CLI for your host machine in `./bin`.

You will then need to make sure that you have the existing Docker CLI in your
`PATH` with the name `com.docker.cli`. A make target is provided to help with
this:

```console
make moby-cli-link
```

This will create a symbolic link from the existing Docker CLI to
`/usr/local/bin` with the name `com.docker.cli`.

You can statically cross compile the CLI for Windows, macOS, and Linux using the
`cross` target.

### Updating the API code

The API provided by the CLI is defined using protobuf. If you make changes to
the `.proto` files in [`protos/`](./protos), you will need to regenerate the API
code:

```console
make protos
```

### Unit tests

To run all of the unit tests, run:

```console
make test
```

If you need to update a golden file simply do `go test ./... -test.update-golden`.

### End to end tests

#### Local tests

To run the local end to end tests, run:

```console
make e2e-local
```

Note that this requires the CLI to be built and a local Docker Engine to be
running.

#### ECS tests

To run the end to end ECS tests, you will need to have an AWS account and have
credentials for it in the `~/.aws/credentials` file.

You can then use the `e2e-ecs` target:

```console
TEST_AWS_PROFILE=myProfile TEST_AWS_REGION=eu-west-3 make e2e-ecs
```

## Releases

To create a new release:
* Check that the CI is green on the main branch for commit you want to release
* Create a new tag of the form vx.y.z, following existing tags, and push the tag

Pushing the tag will automatically create a new release and make binaries for
Windows, macOS, and Linux available for download on the
[releases page](https://github.com/docker/compose-ecs/releases).


# Pogo a GitLab Executor using Podman built with Go

Run GitLab CI jobs on a user-level with containers using podman. For ease of
use the custom executor is managed and processed using a Go package named
`pogo`.

![](./pogo-24fps.gif)

No docker required: pogo provides you with the flexibility of using container
images for testing but on a user level as it is based on
[podman](https://podman.io/). Therefor you can easily set it up on server or
even just register a CI job runner for private use on your local machine that
will run when you are online. Service images work as well, they are attached to
the job container so any service ports is available on `localhost`.

## Usage

We added a [minimal but complete example](./example) of a GitLab project
utilizing a pogo runner on a Fedora machine.

Use the configuration file to provide a default fallback container image for
jobs that have none defined, providing an auth-file that allows to authenticate
podman for private registries, mount directories for specific jobs identified
by tags or add extra arguments to podman identified by tags.

### Install & Register

In order to use pogo you have to install and register a custom gitlab-runner.
Doing so, configure the pogo binary with tasks config, prepare, run and cleanup
steps as arguments.

Inspect the [Makefile](./Makefile) to get details on this process.

```sh
# Pack and register a runner on a server or development machine. The scripts
# provided expect a user named runner and install in the HOME directory of this
# user.
$ make zip # Pack the local files (binary, config).
$ scp pogo-gitlab-runner.zip:<remote-location> # Move zip to the remote server.
$ make setup # Run setup on the remote server.
```

Pogo can be adapted either using a configuration file in YAML format or
environment variables. The file has to be located in the home directory
`~/.pogo.yaml` or in a custom path when using the `--config` flag. Use the
configuration to mount custom volumes (by tags) or apply extra arguments to
podman (by tags) as well as provide an auth-file for podman when dealing with
private registries.

```yaml
# .pogo.yaml
default_image: fedora:latest
working_dir: /home/runner
cache_dir: /home/runner/cache
auth_file: /home/runner/.docker-auth.json
mounts:
  -
    volume: /home/test/dir:/mnt/dir
    tags:
      - build
      - deploy
  -
    volume: /home/test/fixtures:/mnt/fixtures
    tags:
      - test
    readonly: true
extra_arguments:
  help: "--help"
  buildah: '--security-opt label=disable --security-opt seccomp=unconfined --device /dev/fuse:rw'
```

We suggest running the gitlab-runner using systemd, see the `Makefile` and
`gitlab-runner.service` file on how to setup a proper task on user-level.

In case you run `pogo` on a server as user, the systemd task only runs
when the user is logged in. You can adapt this behaviour to run `loginctl
enable-linger <username>`.

If you want to build containers inside a container we recommend to use
[buildah](https://buildah.io/) and checking out the [best practice of Red
Hat](https://developers.redhat.com/blog/2019/08/14/best-practices-for-running-buildah-in-a-container)
doing so.

If your project CI job requires a secret you can either use [GitLab to add env
variables](https://docs.gitlab.com/ee/ci/variables/#add-a-cicd-variable-to-a-project)
or store it on the exectur host and define a mount for your job containers.
Ensure you never store the secret inside of your container image.

## Development

```sh
$ make test # run tests
$ make run ARGS="help" # run pogo
```

## References

- [GitLab Predefined CI Variables](https://docs.gitlab.com/ee/ci/variables/predefined_variables.html)

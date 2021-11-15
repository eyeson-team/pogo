
# Pogo a GitLab Executor using Podman built with Go

Run GitLab CI Jobs on a user-level using podman. For ease of use the custom
executor is managed and processed using a Go package named `pogo`.

![](./pogo-24fps.gif)

## Usage

In order to use pogo you have to install and register a custom gitlab-runner.
Doing so configure the pogo binary with tasks config, prepare, run and cleanup
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
```

We suggest running the gitlab-runner using systemd, see the `Makefile` and
`gitlab-runner.service` file on how to setup a proper task on user-level.

Note: In case you run `pogo` on a server as user, the systemd task only runs
when the user is logged in. You can adapt this behaviour to run `loginctl
enable-linger <username>`.

## Development

```sh
$ make test # run tests
$ make run ARGS="help" # run pogo
```

## References

- [GitLab Predefined CI Variables](https://docs.gitlab.com/ee/ci/variables/predefined_variables.html)

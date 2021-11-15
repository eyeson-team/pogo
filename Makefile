ARGS=help
WORKDIR=`pwd`
GITLAB_RUNNER_SOURCE="https://gitlab-runner-downloads.s3.amazonaws.com/latest/binaries/gitlab-runner-linux-amd64"
SYSTEMD_TARGET_DIR=~/.config/systemd/user
ZIPPED_FILE=pogo-gitlab-runner.tar.gz
GITLAB_URL=https://gitlab.com

.PHONY: run
run:
	go run main.go $(ARGS)

.PHONY: test
test:
	go test pogo/gitlab pogo/podman pogo/cmd

bin/pogo:
	go build -o bin/pogo

zip: bin/pogo
	tar czf $(ZIPPED_FILE) bin/ Makefile gitlab-runner.service .pogo.yaml
	@echo ""
	@echo "USAGE:"
	@echo ">  Move $(ZIPPED_FILE) to server, extract with 'tar xzf $(ZIPPED_FILE)'"
	@echo ">  and run 'make setup'. Ensure you have the GitLab registry token at"
	@echo ">  hand."

.PHONY: install
install:
	mkdir -p $(SYSTEMD_TARGET_DIR) builds cache && \
		mv gitlab-runner.service $(SYSTEMD_TARGET_DIR) && \
		systemctl --user enable gitlab-runner && \
		systemctl --user start gitlab-runner

.PHONY: setup
setup: bin/gitlab-runner register install

bin/gitlab-runner:
	wget -O $@ $(GITLAB_RUNNER_SOURCE) && chmod +x $@

.PHONY: clean
clean:
	rm -f bin/pogo bin/gitlab-runner $(ZIPPED_FILE) && rm -rf build/ cache/

.PHONY: register
register:
	gitlab-runner register \
		--url $(GITLAB_URL) \
		--name "$(HOSTNAME) runner" \
		--executor custom \
		--shell bash \
		--builds-dir /builds \
		--cache-dir /cache \
		--custom-config-exec "/home/runner/bin/pogo" \
		--custom-config-args "config" \
		--custom-config-args "--config" \
		--custom-config-args "/home/runner/.pogo.yaml" \
		--custom-prepare-exec "/home/runner/bin/pogo" \
		--custom-prepare-args "prepare" \
		--custom-prepare-args "--config" \
		--custom-prepare-args "/home/runner/.pogo.yaml" \
		--custom-run-exec "/home/runner/bin/pogo" \
		--custom-run-args "run" \
		--custom-run-args "--config" \
		--custom-run-args "/home/runner/.pogo.yaml" \
		--custom-cleanup-exec "/home/runner/bin/pogo" \
		--custom-cleanup-args "clean" \
		--custom-cleanup-args "--config" \
		--custom-cleanup-args "/home/runner/.pogo.yaml"

.PHONY: unregister
unregister:
	systemctl --user stop gitlab-runner && \
		systemctl --user disable gitlab-runner && \
		rm -rf .gitlab-runner/

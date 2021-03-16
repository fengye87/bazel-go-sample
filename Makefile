.DEFAULT_GOAL := test

gazelle:
	bazel run //:gazelle -- update-repos \
		-from_file=go.mod -to_macro=repos.bzl%go_repositories \
		-build_file_generation=on -build_file_proto_mode=disable
	bazel run //:gazelle

test:
	bazel test --test_output=errors //...

dev:
	bazel run //deploy:deploy.apply

push:
	bazel run //build:images.push

manifest:
	bazel run //deploy:manifest | tee greeter-operator.yaml

clean:
	bazel clean --expunge

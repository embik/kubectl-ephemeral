clean:
	rm -rf _build

build:
	mkdir -p _build
	go build -o ./_build/kubectl-ephemeral ./cmd/kubectl-ephemeral/

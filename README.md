# kubectl-ephemeral

A simple `kubectl` plugin that allows launching ephemeral containers from a YAML file (as `kubectl debug` does not expose all options).

## Installation

`kubectl-ephemeral` is currently available from [my krew index](https://github.com/embik/krew-index) or my [Homebrew tap](https://github.com/embik/homebrew-tap).

### Krew

This installation method requires [krew](https://krew.sigs.k8s.io) installed.

Add the krew index and synchronize the local copy of it:

```sh
$ kubectl krew index add embik https://github.com/embik/krew-index.git
$ kubectl krew update
```

`kubectl-ephemeral` can now be installed with the following command:

```sh
$ kubectl krew install embik/ephemeral
```

### Homebrew

Installation via [Homebrew](https://brew.sh) is possible with the following command:

```sh
$ brew install embik/tap/kubectl-ephemeral
```

### Manual

Run `make build` to get the binary `_build/kubectl-ephemeral` built. Then run:

```sh
$ chmod +x _build/kubectl-ephemeral
$ cp _build/kubectl-ephemeral /usr/local/bin/kubectl-ephemeral # or somewhere else where it will be available in your PATH
```

## Usage

`kubectl-ephemeral` requires you to provide a YAML file that describes the ephemeral container you want to launch. An example on how to launch [Delve](https://github.com/go-delve/delve) to debug a running Go application is provided in [examples/delve.yaml](./examples/delve.yaml).

The ephemeral container specification mostly aligns with normal containers, but some differences exist. In doubt, check the [corev1.EphemeralContainer](https://pkg.go.dev/k8s.io/api/core/v1#EphemeralContainer) type.

```sh
$ kubectl ephemeral <target pod name> -f <path to ephemeral container>.yaml -c <target container name>
```



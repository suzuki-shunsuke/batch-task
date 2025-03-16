# Remove .github/workflows/wc-update-aqua-checksums.yaml

Requirements

- [multi-guitter](https://github.com/lindell/multi-gitter)

You can install them by [aqua](https://aquaproj.github.io/)

```sh
aqua i -l -a
```

## Create pull requests

```sh
bash run.sh <org or user name>
```

e.g.

```sh
bash run.sh aquaproj
```

## Merge pull requests

```sh
bash merge.sh <org or user name>
```

e.g.

```sh
bash merge.sh aquaproj
```

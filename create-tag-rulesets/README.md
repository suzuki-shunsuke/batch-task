# Create Tag Rulesets across all repositories

1. List owners
1. List repositories
1. Create Tag Rulesets

## 1. List owners

```sh
bash list_owners.sh
```

Then `orgs.json` and `owners.txt` are created.

- Exclude some organizations from `owners.txt` if you want
- Add the authenticated user to `owners.txt` manually

## 2. List repositories

```sh
bash list_repos.sh
```

```sh
bash list_repos_txt.sh
```

Then `owners/*.json` and `repos.txt` are created.

- If you want to exclude some repositories, please add `#` to lines in `repos.txt`

e.g.

```
# suzuki-shunsuke/tfcmt # Exclude suzuki-shunsuke/tfcmt
```

## 3. Create Tag Rulesets

```sh
bash main.sh
```

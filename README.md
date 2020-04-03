# GitHub CLI Test

## Links

* https://dhdersch.github.io/golang/2016/01/23/golang-when-to-use-string-pointers.html
* https://httpie.org/doc#explicit-json
* https://github.com/toml-lang/toml
* https://godoc.org/github.com/google/go-github/github
* https://github.com/pelletier/go-toml



## To Build

```sh
make build
```

## To Run

```sh
./ghcs
```

## To Test

Create an API token, and this to the file `apitoken`, without extension.
This is in the `.gitignore` file, so should not be commited ever, but double check to be sure.

Update the contents of `config.toml` to your repository info.

```sh
CONFIG_LOCATION=$(pwd)/config.toml
TOKEN_LOCATION=$(pwd)/apitoken
```

```sh
http :1323/ser number:=2 email=joostvdg@gmail.com course=mine 
```

## Create A File In An Existing Repo Without Cloning

* https://stackoverflow.com/questions/22312545/github-api-to-create-a-file
* https://developer.github.com/v3/repos/contents/#create-a-file

```sh
curl -X PUT -H 'Authorization: token yadayada' \
    -d '{"message": "Initial Commit","content": "bXkgbmV3IGZpbGUgY29udGVudHM="}' \
    https://api.github.com/repos/user/test/contents/so-test.txt
```

### Manual Process

```sh
FILE_NAME=request.json
```

```sh
USER=joostvdg
REPO=github-cli-test-env
```

```sh
CONTENT=$(base64 --input request.json)
echo ${CONTENT}
```

#### Create New File

```json
{
    "message": "my commit message",
    "committer": {
        "name": "Joost van der Griendt",
        "email": "joostvdg@gmail.com"
    },
    "content": "bXkgbmV3IGZpbGUgY29udGVudHM="
}
```

```sh
curl -X PUT -H 'Authorization: token ${TOKEN}' \
    -d '{"message": "Initial Commit","content": \"${CONTENT}"}' \
    https://api.github.com/repos/${USER}/${REPO}/contents/${FILE_NAME}
```

#### Update File

```json
{
    "message": "my commit message",
    "committer": {
        "name": "Joost van der Griendt",
        "email": "joostvdg@gmail.com"
    },
    "content": "bXkgbmV3IGZpbGUgY29udGVudHM=",
    "sha": "95b966ae1c166bd92f8ae7d1c313e738c731dfc3"
}
```

First retrieve current `sha`.

```sh
SHA=$(curl -H 'Authorization: token ${TOKEN}' \
    https://api.github.com/repos/joostvdg/github-cli-test-env/contents/test.json \
    | jq .sha)
```

```sh
curl -X PUT -H 'Authorization: token ${TOKEN}' \
    -d '{"message": "Seconf","content": "ewogICAgImZvbyI6ICJiYXIyIgp9", "committer": {"name": "Joost van der Griendt","email": "joostvdg@gmail.com"}, "sha": "8a79687628fe86b467ec0e87e7e155c4661caa4f"}' \
    https://api.github.com/repos/joostvdg/github-cli-test-env/contents/test.json
```

## Kubernetes

### Create Secret

```sh
kubectl create secret generic ghcs-api-token --from-literal=apitoken=mysecrettoken
```

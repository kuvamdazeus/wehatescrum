# We Hate Scrum v0.0.1 [WIP 🚧]
---
## Usage Steps
1. Clone this repo
1. Get an access token & save it in a file `token.txt` in project root with repo read/write access from GitHub under Settings > Developer Settings > Personal access tokens > Tokens (classic)
1. Create an `entrypoint.sh` file in project root (take reference from below `entrypoint.sh` template). This file is run while building Docker image & pre-loads your git repositories
1. Build Docker image by running: `docker build --build-arg gh_token="$(cat token.txt)" -t wehatescrum .`
1. Upload the built image to your preferred registry & run the Docker image `docker run -p 8080:8080 wehatescrum`

### `entrypoint.sh` template
```bash
#!/bin/bash

# ONLY PARTS with "<...>" present to be edited & modified, it works, please don't touch it...

# List of repositories to clone
repositories=(
    "--branch dev_stable https://$TOKEN@github.com/<USERNAME>/<REPO_NAME_1>"
    "--branch dev_stable https://$TOKEN@github.com/<USERNAME>/<REPO_NAME_2>"
    ...
)

# Use printf to convert the array into a newline-separated string and pipe it into xargs
printf "%s\n" "${repositories[@]}" | xargs -P 3 -I {} sh -c 'git clone --single-branch {}'
```
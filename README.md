## Github create pull resource resource

Custom concourse resource to create a pull request. It expects 

## Source Configuration

| Parameter                   | Required | Example                          | Description                                                                                                                                                                                                                                                                                |
|-----------------------------|----------|----------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `remote_repository`         | Yes      | `paketo-buildpacks/nginx`       | The repository to target.                                                                                                                                                                                                                                                                  |
| `github_token`              | Yes      |                                  | A Github Access Token with repository access (required for setting status on commits). N.B. If you want github-pr-resource to work with a private repository. Set `repo:full` permissions on the access token you create on GitHub. If it is a public repository, `repo:status` is enough. |

#### `check`

The resource does not implement check script. 

#### `get`

The get script is only implemented to work alongside of put script. Thus no params are explicitly needed. One wouldn't need to add get in the task configuration for this resource. 

#### `put`

| Parameter                  | Required | Example                                    | Description                                                                                                                                                    |
|----------------------------|----------|--------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `repo_location`            | Yes      | `nginx`                                    | The name given to the output of previous step. The output MUST be a git repo with local commit(s) already made on existing branch e.g. master.                 |
| `title`                    | No       | `Updating version of nginx from 1.2 to 1.4`| Set title of PR. default "Pull request by bot"                                                                                                                 |
| `description`              | No       | `nginx 1.4 version is available to use`    | Set description of PR. default "This is default description of the PR"                                                                                         |                                                                               |
| `branch_prefix`            | No       | `pr-by-ci`                                 | To create a PR a new branch will be created, this will be prefix of the PR. epoch.time.now will be appended at the end to avoid conflicts. default "pr-by-bot" |
| `base`                     | No       | `develop`                                  | To create a PR, use this as base branch. Default "master"                                                                                                      |
| `auto_merge`               | No       | `true`                                     | Set true to auto-merge the PR. Default is false                                                                                                                |


## Example

```yaml
resources:
- name: repo-name
  type: git
  source:
    uri: git@github.com:org-name/repo-name.git
    private_key: ((private_key))
    branch: master

- name: pull-request
  type: create-pull-request-resource
  source:
    remote_repository: org-name/repo-name
    github_token: ((access_token))

jobs:
- name: create-a-pull-request
  serial: true
  plan:
    - get: repo-name
    - task: make-commit
      config:
        image_resource:
          source:
            repository: cfbuildpacks/bre-ci
          type: registry-image
        inputs:
          - name: devopsdays
        outputs:
          - name: devopsdays
        platform: linux
        run:
          path: sh
          args:
          - -exc
          - |
            cd devopsdays
            git checkout master
            touch onemoretest1
            git add onemoretest1
            git commit -m 'One final commit'
    - put: pull-request
      params:
        repo_location: devopsdays
        title: "Title of PR is this"
        description: "Descritption of PR"
        branch_prefix: "pr-by-bot"
        auto_merge: true
```

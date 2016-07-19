# kranen

This tool allows you to easily create callback endpoints for Docker Hub.

## Usage

Simply execute `kranen -config <path/to/config/yaml> [-httpAddress <http address string>]`.
The httpAddress flag is optional and defaults to `:8080`.

## Configuration

A sample configuration looks like this:

```
- api_key: foobar      # Secret key, needs to be url compatible
  tag: latest          # The tag to react to, other tags are ignored
  name: connctd/test   # The name of the repository, if names don't match calls are ignored
  script: "/foo/bar.sh {{.ENV.HOME}}" # The command to execute, templating and env vars are supported

- <Next hook definition...>
```

## Script templating

The script string can be templated. Environment variables are available as `.ENV.<var>` and the data from
the callback payload is available as `.Hub.<path to data>` (for example `.Hub.Repo.RepoName` for the repository name).
Additionally the specified command is called with all available environment variables.

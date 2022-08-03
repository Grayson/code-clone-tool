# code-clone-tool

This is a simple Go-based command line utility that aids in keeping code
repositories up to date.  Do you have an organization or user account that has
a *lot* of repositories?  Have you ever had to migrate to a new computer or
potentially work with teams that update multiple repos per day?  This tool
might be useful for you.

## How to install

If you have `go` installed, you can install the latest version using
`go install`:

    go install github.com/grayson/code-clone-tool@latest

If not, you can open the [list of releases][releases] in your web browser,
locate the latest release, select the `.tar.gz` archive appropriate to your
computer, and download the latest release.  From there, you should be able to
readily copy or move the `code-clone-tool` binary to a `$PATH` directory.

[releases]: https://github.com/Grayson/code-clone-tool/releases

## What does it do?

Very simply, `code-clone-tool` will receive repo information from the Github API
and then attempt to either clone or pull each of those repos.  The algorithm is
simple.  `code-clone-tool` will attempt to `git clone` repos to the current
working directory.  If the repos already exist, then it will attempt to do a
`git pull` instead.

At present, there's no branch management.  It just attempts to clone or pull.

## Configuration

`code-clone-tool` needs a bit of information in order to work.  At default,
you'll need to specify a Github API URL and a personal access token.  Let's
start with the API URL.

If you want to clone your personal repos, you can simply use
`https://api.github.com/user/repos`.  Your personal access token will identify
you to the service and `code-clone-tool` will receive information about your
repos.

If you want to clone from an organization, you'll want to use a URL similar to
`https://api.github.com/orgs/<ORG>/repos`.  You might notice the `<ORG>` tag in
the example.  You'll need to replace that with your Github Organization name.

In order to identify yourself to the Github API service, you'll need to generate
a personal access token.  Those can be generated in your [Tokens Settings][ts].
You'll be prompted for a scope of access for these tokens.  The following
scopes should be sufficient: `repo` and `read:org`.

[ts]: https://github.com/settings/tokens

### Defining the configuration

There are three ways to get information into `code-clone-tool` in order of
precedence.  First, you can inject data via environment variables.  Those will
be picked up by the application if specified in the environment.  Second, you
can use an `.env` file to specify configuration or *override* environment
variables.  .env file is simply a YAML document with several top-level keys and 
string values.  Finally, there are command line options that will override both
.env file settings and environment variables.

| Env Var               | Config file key       | CLI flag                    |
|-----------------------|-----------------------|-----------------------------|
|`PERSONAL_ACCESS_TOKEN`|`personal_access_token`|`personalaccesstoken`, `t`   |
|`API_URL`              |`api_url`              |`url`, `u`                   |
|`WORKING_DIRECTORY`    |`working_directory`    |`workingdirectory`, `wd`, `d`|
|`CONFIG_PATH`          |n/a                    |`config`, `c`                |

The access token and url were discussed above.  The working directory allows you
to specify a relative or absolute path to use as the working directory.  The
tool will set it as the root from which clones and pulls are subsequently
executed.

There is one flag that cannot be set by the config file (default `.env`).
That's the `--config` CLI flag (`CONFIG_PATH` environment variable).  This
allows you to specify a config file location.  One expected use case is for
users that want to have a single root directory (e.g. `/code`) that contains
cloned repositories across multiple organizations or accounts without extraneous
intermediary directories that contain `.env` files.  This should also allow for
simpler shell aliases by moving configuration into files rather than repeating
command line arguments or environment flags.
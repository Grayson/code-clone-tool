package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"

	"grayson/cct/lib"
	git "grayson/cct/lib/GitApi"
	githubapi "grayson/cct/lib/GithubApi"
	"grayson/cct/lib/fs"
	"grayson/cct/lib/stage"
)

func main() {
	flagsEnv := lib.Env{}
	cliFlagConfigPath := ""

	app := &cli.App{
		Name:  "code-clone-tool",
		Usage: "easily clone repos",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "personalaccesstoken",
				Usage:       "Github Personal Access Token generated at https://github.com/settings/tokens",
				Aliases:     []string{"pat", "token", "t"},
				Destination: &flagsEnv.PersonalAccessToken,
			},
			&cli.StringFlag{
				Name:        "url",
				Usage:       "URL to Github API for an org or a user similar to: https://api.github.com/orgs/<ORG>/repos or https://api.github.com/user/repos",
				Aliases:     []string{"u"},
				Destination: &flagsEnv.ApiUrl,
			},
			&cli.StringFlag{
				Name:        "workingdirectory",
				Usage:       "Change internal working directory",
				Aliases:     []string{"dir", "wd"},
				Destination: &flagsEnv.WorkingDirectory,
			},
			&cli.StringFlag{
				Name:        "config",
				Usage:       "Select config file to load (default `.env`)",
				Aliases:     []string{"c"},
				Destination: &cliFlagConfigPath,
			},
		},
		Action: func(*cli.Context) error {
			config := determineConfigPath(cliFlagConfigPath, func() (string, bool) { return os.LookupEnv("CONFIG_PATH") })
			return run(mergeEnvs(&flagsEnv, loadEnv(config)))
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(env *lib.Env) error {
	gc := git.CreateGitClient(log.Default())
	fs := fs.OsFs{}

	start := stage.Start(func() (bool, error) { return cwd(fs, env.WorkingDirectory) })
	repos := stage.Then(
		start,
		func(bool) (*githubapi.GithubOrgReposResponse, error) {
			client := githubapi.NewClient(http.DefaultClient, env.PersonalAccessToken)
			return fetchRepoInformation(client, env.ApiUrl)
		},
	)
	actions := stage.Then(
		repos,
		func(repos *githubapi.GithubOrgReposResponse) (actions []lib.Action, err error) {
			return mapActions(fs, repos)
		},
	)
	performedTasks := stage.Iterate(
		actions,
		func(a lib.Action) (lib.Task, error) {
			return performGitActions(a, gc)
		},
	)
	counts := stage.Then(
		performedTasks,
		countTasks,
	)
	log.Println()
	_, err := stage.Finally(
		counts,
		func(m map[lib.Task]int) (bool, error) {
			for k, v := range m {
				log.Printf("%v: %v", k, v)
				log.Println()
			}
			return true, nil
		},
	)
	return err
}

func determineConfigPath(initial string, fallback func() (string, bool)) string {
	if initial != "" {
		return initial
	}

	if fb, ok := fallback(); ok {
		return fb
	}

	return ".env"
}

func countTasks(tasks []lib.Task) (map[lib.Task]int, error) {
	m := make(map[lib.Task]int)
	for _, t := range tasks {
		v := m[t]
		m[t] = v + 1
	}
	return m, nil
}

func performGitActions(action lib.Action, gc git.GitApi) (lib.Task, error) {
	var output string
	var err error
	task := action.Task
	switch action.Task {
	case lib.Clone:
		output, err = gc.Clone(action.GitUrl, action.Path)
	case lib.Pull:
		output, err = gc.Pull(action.Path)
	default:
		err = fmt.Errorf("unexpected task: %v", action.Task.String())
	}
	if err != nil {
		return lib.Invalid, err
	}
	log.Print(output)
	return task, nil
}

func fetchRepoInformation(client githubapi.GithubApi, url string) (*githubapi.GithubOrgReposResponse, error) {
	resp, err := client.FetchOrgInformation(url)
	if err != nil {
		return nil, err
	}

	if errResp, ok := resp.GetRight(); ok {
		return nil, fmt.Errorf("service error with the following message:\n%v\n\n%v", errResp.Message, errResp.DocumentationURL)
	}
	repos, _ := resp.GetLeft()
	return repos, nil
}

func cwd(f fs.Fs, p string) (bool, error) {
	err := f.ChangeWorkingDirectory(p)
	return err == nil, err
}

func loadEnv(path string) *lib.Env {
	readers := []lib.ReadYamlFile{
		func() ([]byte, error) {
			return os.ReadFile(path)
		},
	}
	return lib.NewEnv(os.LookupEnv, readers)
}

func mergeEnvs(change *lib.Env, into *lib.Env) *lib.Env {
	if change == nil {
		return into
	}

	if into == nil {
		return change
	}

	if change.ApiUrl != "" {
		into.ApiUrl = change.ApiUrl
	}
	if change.PersonalAccessToken != "" {
		into.PersonalAccessToken = change.PersonalAccessToken
	}
	if change.WorkingDirectory != "" {
		into.WorkingDirectory = change.WorkingDirectory
	}
	return into
}

func mapActions(fs fs.Fs, repos *githubapi.GithubOrgReposResponse) (actions []lib.Action, err error) {
	for _, repo := range *repos {
		action := lib.Action{
			Task:   lib.DiscernTask(repo.FullName, fs),
			Path:   repo.FullName,
			GitUrl: repo.SshUrl,
		}
		actions = append(actions, action)
	}
	return
}

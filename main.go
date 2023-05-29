package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"

	"github.com/grayson/code-clone-tool/app"
	"github.com/grayson/code-clone-tool/lib"
	git "github.com/grayson/code-clone-tool/lib/GitApi"
	githubapi "github.com/grayson/code-clone-tool/lib/GithubApi"
	"github.com/grayson/code-clone-tool/lib/fs"
)

var (
	version = "dev"
	date    = "unknown"
)

func main() {
	flagsEnv := lib.Env{}
	cliFlagConfigPath := ""
	shouldPrintVersion := false
	shouldMirror := false

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
			&cli.BoolFlag{
				Name:        "mirror",
				Usage:       "Mirror a repo rather than clone",
				Aliases:     []string{"m"},
				Destination: &shouldMirror,
			},
			&cli.BoolFlag{
				Name:        "version",
				Usage:       "Print version information and quit",
				Aliases:     []string{"v"},
				Destination: &shouldPrintVersion,
			},
		},
		Action: func(*cli.Context) error {
			if shouldPrintVersion {
				log.Printf("Version %v built on %v", version, date)
				log.Println()
				return nil
			}

			flagsEnv.IsMirror = lib.NewBoolString(shouldMirror)

			fileconfigPath := determineConfigPath(cliFlagConfigPath, func() (string, bool) {
				return os.LookupEnv("CONFIG_PATH")
			})
			fileconfig := lib.LoadEnvironmentYamlFile(func() ([]byte, error) {
				return os.ReadFile(fileconfigPath)
			})
			envconfig := lib.LoadEnvironmentVariables(os.LookupEnv)
			config := flagsEnv.Merge(fileconfig).Merge(envconfig)
			return run(config)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(env *lib.Env) error {
	model := app.InitAppModel(env, version, fs.OsFs{})
	_, err := tea.NewProgram(model).Run()
	if err == nil && model.Error != nil {
		err = model.Error
	}
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

func cwd(f fs.Fs, p string) (bool, error) {
	if p == "" {
		return true, nil
	}
	err := f.ChangeWorkingDirectory(p)
	return err == nil, err
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

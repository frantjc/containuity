package actions

import "path/filepath"

type CtxOpt func(*GlobalContext) error

func WithToken(token string) CtxOpt {
	return func(gc *GlobalContext) error {
		if gc.GitHubContext == nil {
			gc.GitHubContext = &GitHubContext{}
		}
		gc.GitHubContext.Token = token

		if gc.SecretsContext == nil {
			gc.SecretsContext = map[string]string{}
		}

		gc.SecretsContext[EnvVarToken] = token

		return nil
	}
}

func WithSecrets(secrets map[string]string) CtxOpt {
	return func(gc *GlobalContext) error {
		if gc.SecretsContext == nil {
			gc.SecretsContext = secrets
		} else {
			for k, v := range secrets {
				gc.SecretsContext[k] = v
			}
		}
		return nil
	}
}

func WithEnv(env map[string]string) CtxOpt {
	return func(gc *GlobalContext) error {
		if gc.EnvContext == nil {
			gc.EnvContext = env
		} else {
			for k, v := range env {
				gc.EnvContext[k] = v
			}
		}
		return nil
	}
}

func WithWorkdir(workdir string) CtxOpt {
	return func(gc *GlobalContext) error {
		if gc.GitHubContext == nil {
			gc.GitHubContext = &GitHubContext{}
		}
		gc.GitHubContext.ActionPath = filepath.Join(workdir, "action")
		gc.GitHubContext.Workspace = filepath.Join(workdir, "workspace")
		gc.RunnerContext.Temp = filepath.Join(workdir, "runner", "temp")
		gc.RunnerContext.ToolCache = filepath.Join(workdir, "runner", "toolcache")
		return nil
	}
}

func WithJobName(job string) CtxOpt {
	return func(gc *GlobalContext) error {
		if gc.GitHubContext == nil {
			gc.GitHubContext = &GitHubContext{}
		}
		gc.GitHubContext.Job = job
		return nil
	}
}

func WithInputs(inputs map[string]string) CtxOpt {
	return func(gc *GlobalContext) error {
		if gc.InputsContext == nil {
			gc.InputsContext = inputs
		} else {
			for k, v := range inputs {
				gc.InputsContext[k] = v
			}
		}
		return nil
	}
}

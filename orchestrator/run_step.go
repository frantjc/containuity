package orchestrator

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/frantjc/sequence"
	"github.com/frantjc/sequence/env"
	"github.com/frantjc/sequence/meta"
	"github.com/frantjc/sequence/plan"
	"github.com/frantjc/sequence/runtime"
	"github.com/rs/zerolog/log"
)

// TODO should we be using runtime.SpecOpts here?
func RunStep(ctx context.Context, r runtime.Runtime, s *sequence.Step, opts ...runtime.SpecOpt) error {
	spec, err := plan.PlanStep(ctx, s)
	if err != nil {
		return err
	}

	// TODO always pull; this conditional is for dev purposes ONLY
	if spec.Image != meta.Image() {
		_, err = r.Pull(ctx, spec.Image)
		if err != nil {
			return err
		}
	}

	for _, mount := range spec.Mounts {
		if mount.Type == runtime.MountTypeBind {
			err = os.MkdirAll(mount.Source, 0777)
			if err != nil {
				return err
			}
		}
	}

	copts := append([]runtime.SpecOpt{runtime.WithSpec(spec)}, opts...)
	container, err := r.Create(ctx, copts...)
	if err != nil {
		return err
	}

	var (
		eopts = []runtime.ExecOpt{runtime.WithStdio}
		buf   = new(bytes.Buffer)
	)
	if s.IsStdoutResponse() {
		eopts = []runtime.ExecOpt{runtime.WithStreams(os.Stdin, buf, os.Stderr)}
	}
	err = container.Exec(ctx, eopts...)
	if err != nil {
		return err
	}

	if s.IsStdoutResponse() {
		resp := &sequence.StepResponse{}
		if err = json.NewDecoder(buf).Decode(resp); err != nil {
			return err
		}

		if resp.Step != nil {
			mergedStep := s.Merge(resp.Step)
			// encoderBuf   := new(bytes.Buffer)
			b, err := json.Marshal(mergedStep)
			if err != nil {
				return err
			}
			stepString := string(b)
			// todo replace
			replacedStep := Expand(stepString, Mapping)
			log.Info().Msgf("%s", replacedStep)
			step := &sequence.Step{}
			err = json.Unmarshal([]byte(replacedStep), step)
			if err != nil {
				return err
			}
			return RunStep(ctx, r, step, append(opts, runtime.WithMounts(spec.Mounts...), runtime.WithEnv(env.ArrToMap(spec.Env)))...)
		}
	}

	return nil
}

func Mapping(name string) string {
	trimmed := strings.Trim(name, "{ }")
	if trimmed == "github.repository" {
		return "asdf"
	} else {
		return "qwert"
	}
}

func Expand(s string, mapping func(string) string) string {
	var buf []byte
	// ${} is all ASCII, so bytes are fine for this operation.
	i := 0
	for j := 0; j < len(s); j++ {
		if s[j] == '$' && j+1 < len(s) {
			if buf == nil {
				buf = make([]byte, 0, 2*len(s))
			}
			buf = append(buf, s[i:j]...)
			name, w := getShellName(s[j+1:])
			if name == "" && w > 0 {
				// Encountered invalid syntax; eat the
				// characters.
			} else if name == "" {
				// Valid syntax, but $ was not followed by a
				// name. Leave the dollar character untouched.
				buf = append(buf, s[j])
			} else {
				buf = append(buf, mapping(name)...)
			}
			j += w
			i = j + 1
		}
	}
	if buf == nil {
		return s
	}
	return string(buf) + s[i:]
}

func getShellName(s string) (string, int) {
	switch {
	case s[0] == '{':
		// if len(s) > 2 && isShellSpecialVar(s[1]) && s[2] == '}' {
		// 	return s[1:2], 3
		// }
		// Scan to closing brace
		for i := 1; i < len(s); i++ {
			if s[i] == '}' {
				if i == 1 {
					return "", 2 // Bad syntax; eat "${}"
				}
				return s[1:i], i + 2
			}
		}
		return "", 1 // Bad syntax; eat "${"
		// case isShellSpecialVar(s[0]):
		// 	return s[0:1], 1
	}
	// Scan alphanumerics.
	// var i int
	// for i = 0; i < len(s) && isAlphaNum(s[i]); i++ {
	// }
	// return s[:i], i
	return "", 1 // Bad syntax; eat "${"

}

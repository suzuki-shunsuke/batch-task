// replace-go-test-full-workflow replaces suzuki-shunsuke/go-test-full-workflow
// with suzuki-shunsuke/go-test-workflow and suzuki-shunsuke/test-workflow.
// It must be run at the root directory of a repository.
// https://github.com/suzuki-shunsuke/batch-task/issues/22
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/suzuki-shunsuke/go-yamledit/yamledit"
)

const (
	workflowDir = ".github/workflows"
	// A reusable workflow created when go-test-full-workflow is called from a pull_request workflow.
	workflowCallFile = "workflow_call_test.yaml"
	jobName     = "test-workflow"
	// Update the version when go-test-workflow is released.
	goTestWorkflow = "suzuki-shunsuke/go-test-workflow/.github/workflows/test.yaml@da03362365c429417c3af65c7b224b06cb25ece4 # v3.2.1"
	// Update the version when test-workflow is released.
	testWorkflowJob = `uses: suzuki-shunsuke/test-workflow/.github/workflows/test.yaml@51ece8dcaabff59820bfb2053c56e3cb281cf7b2 # v0.0.1
with:
  aqua_version: v2.64.0
permissions:
  contents: read
  pull-requests: read`
)

var usesPattern = regexp.MustCompile(`suzuki-shunsuke/go-test-full-workflow/\.github/workflows/test\.yaml@\S+(\s+#.*)?`)

type Workflow struct {
	Jobs map[string]*Job `yaml:"jobs"`
}

type Job struct {
	Uses    string         `yaml:"uses"`
	Needs   any            `yaml:"needs"`
	With    map[string]any `yaml:"with"`
	Secrets any            `yaml:"secrets"`
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	files, err := listWorkflows()
	if err != nil {
		return err
	}
	found := false
	for _, file := range files {
		ok, err := edit(file)
		if err != nil {
			return fmt.Errorf("edit %s: %w", file, err)
		}
		if ok {
			found = true
		}
	}
	if !found {
		return errors.New("no workflow calls go-test-full-workflow")
	}
	return nil
}

func listWorkflows() ([]string, error) {
	files := []string{}
	for _, pattern := range []string{"*.yaml", "*.yml"} {
		matches, err := filepath.Glob(filepath.Join(workflowDir, pattern))
		if err != nil {
			return nil, fmt.Errorf("find workflows: %w", err)
		}
		files = append(files, matches...)
	}
	return files, nil
}

// edit edits a workflow calling go-test-full-workflow.
// It returns true if the workflow calls go-test-full-workflow.
func edit(file string) (bool, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return false, fmt.Errorf("read a file: %w", err)
	}
	wf := &Workflow{}
	if err := yaml.Unmarshal(b, wf); err != nil {
		return false, fmt.Errorf("unmarshal a workflow: %w", err)
	}
	targets := []string{}
	for name, job := range wf.Jobs {
		if strings.HasPrefix(job.Uses, "suzuki-shunsuke/go-test-full-workflow/") {
			targets = append(targets, name)
		}
	}
	if len(targets) == 0 {
		return false, nil
	}

	if hasDependent(wf, targets) {
		// If the job calling go-test-full-workflow is required by other jobs such as status-check,
		// move it to a reusable workflow so that status-check depends only on one job.
		// https://github.com/suzuki-shunsuke/batch-task/issues/22
		if err := split(file, b, wf, targets); err != nil {
			return false, err
		}
		return true, nil
	}

	actions := []yamledit.Action{}
	for _, name := range targets {
		with := wf.Jobs[name].With
		if _, ok := with["go-version"]; ok {
			continue
		}
		if _, ok := with["go-version-file"]; ok {
			continue
		}
		// go-test-full-workflow sets go-version-file to go.mod by default, but go-test-workflow doesn't.
		actions = append(actions, yamledit.MapAction(fmt.Sprintf("$.jobs.%s.with", name), yamledit.SetKey("go-version-file", "go.mod", nil)))
	}

	// The job already exists if the script is re-run on the pull request branch.
	if _, ok := wf.Jobs[jobName]; !ok {
		actions = append(actions, yamledit.MapAction("$.jobs", yamledit.SetKey(jobName, yamledit.NewBytes([]byte(testWorkflowJob)), nil)))
	}

	s := usesPattern.ReplaceAllString(string(b), goTestWorkflow)
	s = strings.ReplaceAll(s, "# go-test-full-workflow requires go-licenses", "# go-test-workflow requires go-licenses")
	if len(actions) > 0 {
		s, err = yamledit.EditBytes(file, []byte(s), actions...)
		if err != nil {
			return false, fmt.Errorf("edit a workflow: %w", err)
		}
	}
	if err := os.WriteFile(file, []byte(fixBlankLines(s, b)), 0o644); err != nil { //nolint:gosec
		return false, fmt.Errorf("write a workflow: %w", err)
	}
	return true, nil
}

// hasDependent returns true if any job depends on the given jobs.
func hasDependent(wf *Workflow, targets []string) bool {
	for _, job := range wf.Jobs {
		needs := []string{}
		switch v := job.Needs.(type) {
		case string:
			needs = append(needs, v)
		case []any:
			for _, need := range v {
				if s, ok := need.(string); ok {
					needs = append(needs, s)
				}
			}
		}
		if slices.ContainsFunc(needs, func(s string) bool { return slices.Contains(targets, s) }) {
			return true
		}
	}
	return false
}

// split moves the job calling go-test-full-workflow to a new reusable workflow
// and makes the job call the reusable workflow.
// https://github.com/aquaproj/aqua/blob/main/.github/workflows/test.yaml
func split(file string, b []byte, wf *Workflow, targets []string) error {
	if len(targets) != 1 {
		return errors.New("multiple jobs call go-test-full-workflow")
	}
	target := targets[0]
	newFile := filepath.Join(workflowDir, workflowCallFile)
	if _, err := os.Stat(newFile); err == nil {
		return fmt.Errorf("%s already exists", newFile)
	}
	secrets := []string{}
	switch v := wf.Jobs[target].Secrets.(type) {
	case nil:
	case map[string]any:
		for k := range v {
			secrets = append(secrets, k)
		}
		slices.Sort(secrets)
	default:
		return errors.New("secrets must be a map")
	}

	// Copy the job to the reusable workflow as is.
	lines := strings.Split(string(b), "\n")
	start := slices.Index(lines, "  "+target+":")
	if start == -1 {
		return fmt.Errorf("job %s isn't found", target)
	}
	end := start + 1
	for ; end < len(lines); end++ {
		l := lines[end]
		if strings.TrimSpace(l) != "" && !strings.HasPrefix(l, "   ") {
			break
		}
	}
	for end > start && strings.TrimSpace(lines[end-1]) == "" {
		end--
	}
	header := []string{
		"---",
		"name: test (workflow_call)",
		"on:",
	}
	if len(secrets) == 0 {
		header = append(header, "  workflow_call:")
	} else {
		header = append(header, "  workflow_call:", "    secrets:")
		for _, k := range secrets {
			header = append(header, "      "+k+":", "        required: true")
		}
	}
	header = append(header, "permissions: {}", "jobs:")
	content := strings.Join(append(header, lines[start:end]...), "\n") + "\n"
	if err := os.WriteFile(newFile, []byte(content), 0o644); err != nil { //nolint:gosec
		return fmt.Errorf("write %s: %w", newFile, err)
	}
	if _, err := edit(newFile); err != nil {
		return fmt.Errorf("edit %s: %w", newFile, err)
	}

	// Call the reusable workflow.
	s := usesPattern.ReplaceAllString(string(b), "./"+newFile)
	if wf.Jobs[target].With != nil {
		var err error
		s, err = yamledit.EditBytes(file, []byte(s), yamledit.MapAction(fmt.Sprintf("$.jobs.%s", target), yamledit.RemoveKeys("with")))
		if err != nil {
			return fmt.Errorf("edit a workflow: %w", err)
		}
	}
	if err := os.WriteFile(file, []byte(s), 0o644); err != nil { //nolint:gosec
		return fmt.Errorf("write a workflow: %w", err)
	}
	return nil
}

// fixBlankLines separates test-workflow from the previous job with a blank line
// if the original file separates jobs with blank lines.
func fixBlankLines(s string, orig []byte) string {
	if !strings.Contains(string(orig), "\n\n  ") {
		return s
	}
	lines := strings.Split(s, "\n")
	ret := make([]string, 0, len(lines)+1)
	for _, line := range lines {
		if line == "  "+jobName+":" && len(ret) > 0 && ret[len(ret)-1] != "" {
			ret = append(ret, "")
		}
		ret = append(ret, line)
	}
	return strings.Join(ret, "\n")
}

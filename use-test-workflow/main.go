// use-test-workflow replaces wc-renovate-config-validator.yaml, wc-ghalint.yaml, and typos.yaml
// with suzuki-shunsuke/test-workflow.
// It must be run at the root directory of a repository.
// https://github.com/suzuki-shunsuke/batch-task/issues/19
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
	jobName     = "test-workflow"
	// Update the version when test-workflow is released.
	testWorkflowJob = `uses: suzuki-shunsuke/test-workflow/.github/workflows/test.yaml@51ece8dcaabff59820bfb2053c56e3cb281cf7b2 # v0.0.1
with:
  aqua_version: v2.64.0
permissions:
  contents: read
  pull-requests: read`
)

// Reusable workflows replaced by test-workflow and paths-filter keys used for them.
var targets = map[string]string{
	"wc-renovate-config-validator.yaml": "renovate-config-validator",
	"renovate-config-validator.yaml":    "renovate-config-validator",
	"wc-ghalint.yaml":                   "ghalint",
	"typos.yaml":                        "typos",
}

// aqua packages that are no longer needed if the reusable workflow is removed.
var aquaPackages = map[string]string{
	"wc-ghalint.yaml": "ghalint",
	"typos.yaml":      "typos",
}

type Workflow struct {
	Jobs map[string]*Job `yaml:"jobs"`
}

type Job struct {
	Uses        string            `yaml:"uses"`
	Needs       any               `yaml:"needs"`
	Outputs     map[string]string `yaml:"outputs"`
	Permissions any               `yaml:"permissions"`
	Steps       []*Step           `yaml:"steps"`
}

type Step struct {
	Uses string         `yaml:"uses"`
	With map[string]any `yaml:"with"`
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
	// go-test-full-workflow v6.1.0 or later runs test-workflow, so test-workflow shouldn't be added.
	addJob := true
	for _, file := range files {
		b, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read %s: %w", file, err)
		}
		if strings.Contains(string(b), "suzuki-shunsuke/go-test-full-workflow/") {
			addJob = false
			break
		}
	}
	removedFiles := map[string]struct{}{}
	callers := []string{}
	for _, file := range files {
		found, err := editCaller(file, addJob)
		if err != nil {
			return fmt.Errorf("edit %s: %w", file, err)
		}
		if len(found) == 0 {
			continue
		}
		callers = append(callers, file)
		for _, f := range found {
			removedFiles[f] = struct{}{}
		}
	}
	if len(callers) == 0 {
		return errors.New("no workflow calls the target reusable workflows")
	}
	for _, file := range files {
		if !addJob {
			break
		}
		if err := addPermission(file, callers); err != nil {
			return fmt.Errorf("add permissions to %s: %w", file, err)
		}
	}
	for f := range removedFiles {
		if err := os.Remove(filepath.Join(workflowDir, f)); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("remove a workflow: %w", err)
		}
		if pkg, ok := aquaPackages[f]; ok {
			if err := removeAquaPackage(pkg); err != nil {
				return fmt.Errorf("remove an aqua package %s: %w", pkg, err)
			}
		}
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

func readWorkflow(file string) (*Workflow, []byte, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, nil, fmt.Errorf("read a file: %w", err)
	}
	wf := &Workflow{}
	if err := yaml.Unmarshal(b, wf); err != nil {
		return nil, nil, fmt.Errorf("unmarshal a workflow: %w", err)
	}
	return wf, b, nil
}

// editCaller edits a workflow calling target reusable workflows.
// It returns the file names of called target reusable workflows.
// If addJob is false, the test-workflow job isn't added.
func editCaller(file string, addJob bool) ([]string, error) {
	wf, b, err := readWorkflow(file)
	if err != nil {
		return nil, err
	}
	removedJobs := []string{}
	removedFilterKeys := []string{}
	found := []string{}
	for name, job := range wf.Jobs {
		f, ok := strings.CutPrefix(localPath(job.Uses), workflowDir+"/")
		if !ok {
			continue
		}
		key, ok := targets[f]
		if !ok {
			continue
		}
		removedJobs = append(removedJobs, name)
		removedFilterKeys = append(removedFilterKeys, key)
		found = append(found, f)
	}
	if len(found) == 0 {
		return nil, nil
	}
	if _, ok := wf.Jobs[jobName]; ok {
		// The job was already added. e.g. The script is re-run on the pull request branch.
		addJob = false
	}

	actions := []yamledit.Action{
		yamledit.MapAction("$.jobs", yamledit.RemoveKeys(toAny(removedJobs)...)),
	}

	// Edit the paths-filter job.
	removedNeeds := slices.Clone(removedJobs)
	for name, job := range wf.Jobs {
		idx := slices.IndexFunc(job.Steps, func(s *Step) bool {
			return strings.HasPrefix(s.Uses, "dorny/paths-filter@")
		})
		if idx == -1 {
			continue
		}
		filters, _ := job.Steps[idx].With["filters"].(string)
		if strings.TrimSpace(removeFilters(filters, removedFilterKeys)) == "" {
			// No filter is left, so the paths-filter job is no longer needed.
			actions = append(actions, yamledit.MapAction("$.jobs", yamledit.RemoveKeys(name)))
			removedNeeds = append(removedNeeds, name)
			continue
		}
		actions = append(actions, yamledit.MapAction(fmt.Sprintf("$.jobs.%s.outputs", name), yamledit.RemoveKeys(toAny(removedFilterKeys)...)))
		// The filters block is edited as text later to keep the literal style.
	}

	if addJob {
		actions = append(actions, yamledit.MapAction("$.jobs", yamledit.SetKey(jobName, yamledit.NewBytes([]byte(testWorkflowJob)), nil)))
	}

	// Remove removed jobs from needs.
	// Jobs depending on removed reusable workflows (e.g. status-check) depend on test-workflow instead.
	for name, job := range wf.Jobs {
		if slices.Contains(removedNeeds, name) {
			continue
		}
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
		if !slices.ContainsFunc(needs, func(s string) bool { return slices.Contains(removedNeeds, s) }) {
			continue
		}
		newNeeds := []any{}
		for _, need := range needs {
			if !slices.Contains(removedNeeds, need) {
				newNeeds = append(newNeeds, need)
			}
		}
		if addJob && slices.ContainsFunc(needs, func(s string) bool { return slices.Contains(removedJobs, s) }) {
			newNeeds = append(newNeeds, jobName)
		}
		if len(newNeeds) == 0 {
			actions = append(actions, yamledit.MapAction(fmt.Sprintf("$.jobs.%s", name), yamledit.RemoveKeys("needs")))
			continue
		}
		if _, ok := job.Needs.(string); ok {
			actions = append(actions, yamledit.MapAction(fmt.Sprintf("$.jobs.%s", name), yamledit.SetKey("needs", newNeeds[0], nil)))
			continue
		}
		listActions := []yamledit.SequenceNodeAction{
			yamledit.RemoveValuesFromList(func(node *yamledit.Node[string]) (bool, error) {
				return slices.Contains(removedNeeds, node.Value), nil
			}),
		}
		if slices.Contains(newNeeds, any(jobName)) {
			listActions = append(listActions, yamledit.AddValuesToList(-1, jobName))
		}
		actions = append(actions, yamledit.ListAction(fmt.Sprintf("$.jobs.%s.needs", name), listActions...))
	}

	s, err := yamledit.EditBytes(file, []byte(removeFiltersInFile(string(b), removedFilterKeys)), actions...)
	if err != nil {
		return nil, fmt.Errorf("edit a workflow: %w", err)
	}
	if err := os.WriteFile(file, []byte(fixBlankLines(s, b)), 0o644); err != nil { //nolint:gosec
		return nil, fmt.Errorf("write a workflow: %w", err)
	}
	return found, nil
}

// fixBlankLines removes blank lines just after "jobs:" and
// separates test-workflow from the previous job with a blank line if the original file separates jobs with blank lines.
func fixBlankLines(s string, orig []byte) string {
	lines := strings.Split(s, "\n")
	ret := make([]string, 0, len(lines)+1)
	separated := strings.Contains(string(orig), "\n\n  ")
	afterJobs := false
	for _, line := range lines {
		if afterJobs && strings.TrimSpace(line) == "" {
			continue
		}
		afterJobs = line == "jobs:"
		if separated && line == "  "+jobName+":" && len(ret) > 0 && ret[len(ret)-1] != "" {
			ret = append(ret, "")
		}
		ret = append(ret, line)
	}
	// Collapse consecutive blank lines between jobs left by removed jobs.
	out := make([]string, 0, len(ret))
	for i, line := range ret {
		if line == "" && i > 0 && ret[i-1] == "" {
			continue
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

var filtersPattern = regexp.MustCompile(`^(\s+)filters:\s*\|[-+]?\s*$`)

// removeFiltersInFile removes given keys from dorny/paths-filter's filters written as a literal block in a workflow.
func removeFiltersInFile(s string, keys []string) string {
	lines := strings.Split(s, "\n")
	ret := make([]string, 0, len(lines))
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		ret = append(ret, line)
		m := filtersPattern.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		// Collect the literal block.
		block := []string{}
		j := i + 1
		for ; j < len(lines); j++ {
			l := lines[j]
			if strings.TrimSpace(l) != "" && len(l)-len(strings.TrimLeft(l, " ")) <= len(m[1]) {
				break
			}
			block = append(block, l)
		}
		// Keep trailing blank lines outside of the block.
		k := len(block)
		for k > 0 && strings.TrimSpace(block[k-1]) == "" {
			k--
		}
		indent := ""
		if k > 0 {
			indent = block[0][:len(block[0])-len(strings.TrimLeft(block[0], " "))]
		}
		body := make([]string, k)
		for n, l := range block[:k] {
			body[n] = strings.TrimPrefix(l, indent)
		}
		for _, l := range strings.Split(removeFilters(strings.Join(body, "\n"), keys), "\n") {
			if l == "" {
				ret = append(ret, "")
				continue
			}
			ret = append(ret, indent+l)
		}
		ret = append(ret, block[k:]...)
		i = j - 1
	}
	return strings.Join(ret, "\n")
}

var filterKeyPattern = regexp.MustCompile(`^([^\s#][^:]*):\s*$`)

// removeFilters removes given keys and their file patterns from dorny/paths-filter's filters.
func removeFilters(filters string, keys []string) string {
	lines := strings.Split(filters, "\n")
	ret := make([]string, 0, len(lines))
	skip := false
	for _, line := range lines {
		if m := filterKeyPattern.FindStringSubmatch(line); m != nil {
			skip = slices.Contains(keys, m[1])
		}
		if skip {
			continue
		}
		ret = append(ret, line)
	}
	return strings.Join(ret, "\n")
}

// addPermission grants pull-requests:read to jobs calling the edited workflows
// because test-workflow requires it.
func addPermission(file string, callers []string) error {
	wf, _, err := readWorkflow(file)
	if err != nil {
		return err
	}
	actions := []yamledit.Action{}
	for name, job := range wf.Jobs {
		if f := localPath(job.Uses); f == "" || !slices.Contains(callers, f) {
			continue
		}
		perms, ok := job.Permissions.(map[string]any)
		if !ok {
			// permissions isn't set or is set to read-all or write-all.
			continue
		}
		if _, ok := perms["pull-requests"]; ok {
			continue
		}
		if len(perms) == 0 {
			actions = append(actions, yamledit.MapAction(fmt.Sprintf("$.jobs.%s", name), yamledit.SetKey("permissions", yamledit.NewBytes([]byte("pull-requests: read")), nil)))
			continue
		}
		actions = append(actions, yamledit.MapAction(fmt.Sprintf("$.jobs.%s.permissions", name), yamledit.SetKey("pull-requests", "read", nil)))
	}
	if len(actions) == 0 {
		return nil
	}
	if _, err := yamledit.EditFile(file, actions...); err != nil {
		return fmt.Errorf("edit a workflow: %w", err)
	}
	return nil
}

// removeAquaPackage removes aqua/imports/<pkg>.yaml and the package from aqua.yaml.
func removeAquaPackage(pkg string) error {
	for _, f := range []string{"aqua/imports/" + pkg + ".yaml", "aqua/imports/" + pkg + ".yml"} {
		if err := os.Remove(f); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("remove %s: %w", f, err)
		}
	}
	names := map[string]string{
		"ghalint": "suzuki-shunsuke/ghalint",
		"typos":   "crate-ci/typos",
	}
	for _, f := range []string{"aqua.yaml", ".aqua.yaml", "aqua/aqua.yaml", ".aqua/aqua.yaml"} {
		b, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		cfg := struct {
			Packages []any `yaml:"packages"`
		}{}
		if err := yaml.Unmarshal(b, &cfg); err != nil {
			return fmt.Errorf("unmarshal %s: %w", f, err)
		}
		if len(cfg.Packages) == 0 {
			// yamledit.ListAction panics if the path isn't found.
			continue
		}
		if _, err := yamledit.EditFile(f, yamledit.ListAction("$.packages", yamledit.RemoveValuesFromList(func(node *yamledit.Node[map[string]any]) (bool, error) {
			name, _ := node.Value["name"].(string)
			return strings.HasPrefix(name, names[pkg]+"@"), nil
		}))); err != nil {
			return fmt.Errorf("edit %s: %w", f, err)
		}
	}
	return nil
}

// localPath returns the path of a reusable workflow in the same repository.
// Both "./" and "$/" are supported.
// If uses isn't a local reusable workflow, it returns an empty string.
func localPath(uses string) string {
	for _, prefix := range []string{"./", "$/"} {
		if f, ok := strings.CutPrefix(uses, prefix); ok {
			return f
		}
	}
	return ""
}

func toAny(s []string) []any {
	ret := make([]any, len(s))
	for i, v := range s {
		ret[i] = v
	}
	return ret
}

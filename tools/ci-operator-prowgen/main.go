package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ghodss/yaml"

	cioperatorapi "github.com/openshift/ci-operator/pkg/api"
	kubeapi "k8s.io/api/core/v1"
	prowconfig "k8s.io/test-infra/prow/config"
	prowkube "k8s.io/test-infra/prow/kube"
)

type options struct {
	ciOperatorConfigPath string

	fullRepoMode bool

	ciOperatorConfigDir string
	prowJobConfigDir    string

	help    bool
	verbose bool
}

func bindOptions(flag *flag.FlagSet) *options {
	opt := &options{}

	flag.StringVar(&opt.ciOperatorConfigPath, "source-config", "", "Path to ci-operator configuration file in openshift/release repository.")

	flag.BoolVar(&opt.fullRepoMode, "full-repo", false, "If set to true, the generator will walk over all ci-operator config files in openshift/release repository and regenerate all component prow job config files")

	flag.StringVar(&opt.ciOperatorConfigDir, "config-dir", "", "Path to a root of directory structure with ci-operator config files (ci-operator/config in openshift/release)")
	flag.StringVar(&opt.prowJobConfigDir, "prow-jobs-dir", "", "Path to a root of directory structure with Prow job config files (ci-operator/jobs in openshift/release)")

	flag.BoolVar(&opt.help, "h", false, "Show help for ci-operator-prowgen")
	flag.BoolVar(&opt.verbose, "v", false, "Show verbose output")

	return opt
}

func generatePodSpec(org, repo, branch, target string, additionalArgs []string) *kubeapi.PodSpec {
	configMapKeyRef := kubeapi.EnvVarSource{
		ConfigMapKeyRef: &kubeapi.ConfigMapKeySelector{
			LocalObjectReference: kubeapi.LocalObjectReference{
				Name: fmt.Sprintf("ci-operator-%s-%s", org, repo),
			},
			Key: fmt.Sprintf("%s.json", branch),
		},
	}

	return &kubeapi.PodSpec{
		ServiceAccountName: "ci-operator",
		Containers: []kubeapi.Container{
			kubeapi.Container{
				Image:   "ci-operator:latest",
				Command: []string{"ci-operator"},
				Args:    append([]string{"--artifact-dir=$(ARTIFACTS)", fmt.Sprintf("--target=%s", target)}, additionalArgs...),
				Env: []kubeapi.EnvVar{
					kubeapi.EnvVar{
						Name:      "CONFIG_SPEC",
						ValueFrom: &configMapKeyRef,
					},
				},
			},
		},
	}
}

type testDescription struct {
	Name   string
	Target string
}

func generatePresubmitForTest(test testDescription, org, repo, branch string) *prowconfig.Presubmit {
	presubmit := prowconfig.Presubmit{
		Agent:        "kubernetes",
		AlwaysRun:    true,
		Brancher:     prowconfig.Brancher{Branches: []string{branch}},
		Context:      fmt.Sprintf("ci/prow/%s", test.Name),
		Name:         fmt.Sprintf("pull-ci-%s-%s-%s", org, repo, test.Name),
		RerunCommand: fmt.Sprintf("/test %s", test.Name),
		// Spec:         generatePodSpec(org, repo, branch, test.Target, []string{}),
		Trigger: fmt.Sprintf("((?m)^/test( all| %s),?(\\\\s+|$))", test.Name),
		UtilityConfig: prowconfig.UtilityConfig{
			DecorationConfig: &prowkube.DecorationConfig{SkipCloning: true},
			Decorate:         true,
		},
	}
	return &presubmit
}

func generatePostsubmitForTest(test testDescription, org, repo, branch string, additionalArgs []string) *prowconfig.Postsubmit {
	postsubmit := prowconfig.Postsubmit{
		Agent: "kubernetes",
		Name:  fmt.Sprintf("branch-ci-%s-%s-%s", org, repo, test.Name),
		// Spec:  generatePodSpec(org, repo, branch, test.Target, additionalArgs),
		UtilityConfig: prowconfig.UtilityConfig{
			DecorationConfig: &prowkube.DecorationConfig{SkipCloning: true},
			Decorate:         true,
		},
	}
	return &postsubmit
}

func generateJobs(
	configSpec *cioperatorapi.ReleaseBuildConfiguration,
	org, repo, branch string,
) (*map[string][]prowconfig.Presubmit, *map[string][]prowconfig.Postsubmit) {

	orgrepo := fmt.Sprintf("%s/%s", org, repo)
	presubmits := map[string][]prowconfig.Presubmit{}
	postsubmits := map[string][]prowconfig.Postsubmit{}

	imagesTest := false

	for _, element := range configSpec.Tests {
		// Check if config file has "images" test defined to avoid name clash
		if element.As == "images" {
			imagesTest = true
		}
		test := testDescription{Name: element.As, Target: element.As}
		presubmits[orgrepo] = append(presubmits[orgrepo], *generatePresubmitForTest(test, org, repo, branch))
		postsubmits[orgrepo] = append(postsubmits[orgrepo], *generatePostsubmitForTest(test, org, repo, branch, []string{}))
	}

	if len(configSpec.Images) > 0 && !imagesTest {
		// TODO: somehow handle the images case better than just not creating this job when there is name conflict
		test := testDescription{Name: "images", Target: "[images]"}
		presubmits[orgrepo] = append(presubmits[orgrepo], *generatePresubmitForTest(test, org, repo, branch))
		postsubmits[orgrepo] = append(postsubmits[orgrepo], *generatePostsubmitForTest(test, org, repo, branch, []string{"--promote"}))
	}

	return &presubmits, &postsubmits
}

// these are unnecessary, and make the config larger so we strip them out
// TODO: Remove once PR TODO is merged
func yamlBytesStripNulls(yamlBytes []byte) []byte {
	nullRE := regexp.MustCompile("(?m)[\n]+^[^\n]+: null$")
	return nullRE.ReplaceAll(yamlBytes, []byte{})
}

func generateProwJobsFromConfigFile(configFilePath string) []byte {
	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read ci-operator config (%v)\n", err)
		os.Exit(1)
	}

	var configSpec *cioperatorapi.ReleaseBuildConfiguration
	if err := json.Unmarshal(data, &configSpec); err != nil {
		fmt.Printf("failed to load ci-operator config (%v)\n", err)
	}

	configSpecDir := filepath.Dir(configFilePath)
	repo := filepath.Base(configSpecDir)
	org := filepath.Base(filepath.Dir(configSpecDir))
	branch := strings.TrimSuffix(filepath.Base(configFilePath), filepath.Ext(configFilePath))

	presubmits, postsubmits := generateJobs(configSpec, org, repo, branch)

	jobConfig := prowconfig.JobConfig{
		Presubmits:  *presubmits,
		Postsubmits: *postsubmits,
	}

	jobConfigAsYaml, err := yaml.Marshal(jobConfig)
	if err != nil {
		fmt.Printf("failed to marshal the job config (%v)", err)
		os.Exit(1)
	}

	jobConfigAsYaml = yamlBytesStripNulls(jobConfigAsYaml)

	return jobConfigAsYaml
}

func generateAllProwJobs(configDir, jobDir string) {
	err := filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error encontered while generating Prow job config: %v\n", err)
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			jobConfigAsYaml := generateProwJobsFromConfigFile(path)
			suffixPath := filepath.Dir(strings.TrimPrefix(path, configDir))
			jobDirForComponent := filepath.Join(jobDir, suffixPath)
			os.MkdirAll(jobDirForComponent, os.ModePerm)
			branch := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
			target := filepath.Join(jobDirForComponent, fmt.Sprintf("%s.yaml", branch))
			err := ioutil.WriteFile(target, jobConfigAsYaml, 0664)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write job config to '%s' (%v)\n", target, err)
				return err
			}
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate all Prow jobs\n")
		os.Exit(1)
	}
}

func inferConfigDirectories() (string, string) {
	repoRootRaw, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to determine repository root with 'git rev-parse --show-toplevel")
		os.Exit(1)
	}
	repoRoot := strings.TrimSpace(string(repoRootRaw))
	configDir := filepath.Join(repoRoot, "ci-operator", "config")
	jobDir := filepath.Join(repoRoot, "ci-operator", "jobs")

	return configDir, jobDir
}

func main() {
	flagSet := flag.NewFlagSet("", flag.ExitOnError)
	opt := bindOptions(flagSet)
	flagSet.Parse(os.Args[1:])

	if opt.help {
		flagSet.Usage()
		os.Exit(0)
	}

	if len(opt.ciOperatorConfigPath) > 0 {
		jobConfigAsYaml := generateProwJobsFromConfigFile(opt.ciOperatorConfigPath)
		fmt.Printf(string(jobConfigAsYaml))
	} else if opt.fullRepoMode {
		configDir, jobDir := inferConfigDirectories()
		generateAllProwJobs(configDir, jobDir)
	} else if len(opt.ciOperatorConfigDir) > 0 && len(opt.prowJobConfigDir) > 0 {
		generateAllProwJobs(opt.ciOperatorConfigDir, opt.prowJobConfigDir)
	} else {
		fmt.Fprintf(os.Stderr, "ci-operator-prowgen needs --source-config, --full-repo or --{config,prow-jobs}-dir option\n")
		os.Exit(1)
	}
}

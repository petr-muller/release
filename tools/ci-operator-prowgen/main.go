package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/ghodss/yaml"

	cioperatorapi "github.com/openshift/ci-operator/pkg/api"
	kubeapi "k8s.io/api/core/v1"
	prowconfig "k8s.io/test-infra/prow/config"
	prowkube "k8s.io/test-infra/prow/kube"
)

type options struct {
	ciOperatorConfigPath string

	help    bool
	verbose bool
}

func bindOptions(flag *flag.FlagSet) *options {
	opt := &options{}

	flag.StringVar(&opt.ciOperatorConfigPath, "source-config", "", "Path to ci-operator configuration file in openshift/release repository.")
	flag.BoolVar(&opt.help, "h", false, "Show help for ci-operator-prowgen")
	flag.BoolVar(&opt.verbose, "v", false, "Show verbose output")

	return opt
}

func generatePodSpec(org, repo, branch string) *kubeapi.PodSpec {
	podspec := kubeapi.PodSpec{ServiceAccountName: "ci-operator"}
	command := []string{"ci-operator"}
	args := []string{"--artifact-dir=$(ARTIFACTS)"}

	configMapKeyRef := kubeapi.EnvVarSource{
		ConfigMapKeyRef: &kubeapi.ConfigMapKeySelector{
			LocalObjectReference: kubeapi.LocalObjectReference{Name: fmt.Sprintf("ci-operator-%s-%s", org, repo)},
			Key:                  fmt.Sprintf("%s.json", branch),
		},
	}
	configSpec := kubeapi.EnvVar{Name: "CONFIG_SPEC", ValueFrom: &configMapKeyRef}
	env := []kubeapi.EnvVar{configSpec}

	container := kubeapi.Container{
		Name:    "test",
		Image:   "ci-operator:latest",
		Command: command,
		Args:    args,
		Env:     env,
	}
	podspec.Containers = []kubeapi.Container{container}

	return &podspec
}

func generatePresubmitForTest(test, org, repo, branch string) prowconfig.Presubmit {
	presubmit := prowconfig.Presubmit{
		Agent:        "kubernetes",
		AlwaysRun:    true,
		Brancher:     prowconfig.Brancher{Branches: []string{branch}},
		Context:      fmt.Sprintf("ci/prow/%s", test),
		Name:         fmt.Sprintf("ci-operator-%s-%s-%s", org, repo, test),
		RerunCommand: fmt.Sprintf("/test %s", test),
		Spec:         generatePodSpec(org, repo, branch),
		Trigger:      fmt.Sprintf("((?m)^/test( all| %s),?(\\\\s+|$))", test),
		UtilityConfig: prowconfig.UtilityConfig{
			DecorationConfig: &prowkube.DecorationConfig{SkipCloning: true},
			Decorate:         true,
		},
	}

	return presubmit
}

func generatePresubmits(tests []cioperatorapi.TestStepConfiguration, org, repo, branch string) map[string][]prowconfig.Presubmit {
	orgrepo := fmt.Sprintf("%s/%s", org, repo)
	presubmits := map[string][]prowconfig.Presubmit{}
	for _, element := range tests {
		presubmits[orgrepo] = append(presubmits[orgrepo], generatePresubmitForTest(element.As, org, repo, branch))
	}
	return presubmits
}

func generateJobConfig(configSpec *cioperatorapi.ReleaseBuildConfiguration, org, repo, branch string) *prowconfig.JobConfig {
	jobConfig := prowconfig.JobConfig{}
	jobConfig.Presubmits = generatePresubmits(configSpec.Tests, org, repo, branch)
	// TODO: postsubmits
	// jobConfig.Postsubmits = make(map[string][]prowconfig.Postsubmit)
	// jobConfig.Postsubmits[repo] = make([]prowconfig.Postsubmit, 0)
	// ...etc...

	return &jobConfig
}

func main() {
	flagSet := flag.NewFlagSet("", flag.ExitOnError)
	opt := bindOptions(flagSet)
	flagSet.Parse(os.Args[1:])

	if opt.help {
		flagSet.Usage()
		os.Exit(0)
	}

	if len(opt.ciOperatorConfigPath) == 0 {
		fmt.Fprintf(os.Stderr, "ci-operator-prowgen needs --source-config option to read ci-operator configuration\n")
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(opt.ciOperatorConfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read ci-operator config (%v)\n", err)
		os.Exit(1)
	}

	var configSpec *cioperatorapi.ReleaseBuildConfiguration
	if err := json.Unmarshal(data, &configSpec); err != nil {
		fmt.Printf("failed to load ci-operator config (%v)\n", err)
	}

	configSpecDir := path.Dir(opt.ciOperatorConfigPath)
	repo := path.Base(configSpecDir)
	org := path.Base(path.Dir(configSpecDir))
	branch := strings.TrimSuffix(path.Base(opt.ciOperatorConfigPath), path.Ext(opt.ciOperatorConfigPath))

	jobConfig := generateJobConfig(configSpec, org, repo, branch)

	jobConfigAsYaml, err := yaml.Marshal(jobConfig)
	if err != nil {
		fmt.Printf("failed to marshal the job config (%v)", err)
		os.Exit(1)
	}

	fmt.Printf(string(jobConfigAsYaml))
}

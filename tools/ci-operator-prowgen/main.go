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
				Name:    "test",
				Image:   "ci-operator:latest",
				Command: []string{"ci-operator"},
				Args:    []string{"--artifact-dir=$(ARTIFACTS)"},
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
		Spec:         generatePodSpec(org, repo, branch),
		Trigger:      fmt.Sprintf("((?m)^/test( all| %s),?(\\\\s+|$))", test.Name),
		UtilityConfig: prowconfig.UtilityConfig{
			DecorationConfig: &prowkube.DecorationConfig{SkipCloning: true},
			Decorate:         true,
		},
	}
	presubmit.Spec.Containers[0].Args = append(
		presubmit.Spec.Containers[0].Args,
		fmt.Sprintf("--target=%s", test.Target),
	)

	return &presubmit
}

func generatePostsubmitForTest(test testDescription, org, repo, branch string) *prowconfig.Postsubmit {
	postsubmit := prowconfig.Postsubmit{
		Agent: "kubernetes",
		Name:  fmt.Sprintf("branch-ci-%s-%s-%s", org, repo, test.Name),
		Spec:  generatePodSpec(org, repo, branch),
		UtilityConfig: prowconfig.UtilityConfig{
			DecorationConfig: &prowkube.DecorationConfig{SkipCloning: true},
			Decorate:         true,
		},
	}

	postsubmit.Spec.Containers[0].Args = append(
		postsubmit.Spec.Containers[0].Args,
		fmt.Sprintf("--target=%s", test.Target),
	)
	postsubmit.Spec.Containers[0].Args = append(postsubmit.Spec.Containers[0].Args, "--promote")

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
		postsubmits[orgrepo] = append(postsubmits[orgrepo], *generatePostsubmitForTest(test, org, repo, branch))
	}

	if len(configSpec.Images) > 0 && !imagesTest {
		// TODO: somehow handle the images case better than just not creating this job when there is name conflict
		test := testDescription{Name: "images", Target: "[images]"}
		presubmits[orgrepo] = append(presubmits[orgrepo], *generatePresubmitForTest(test, org, repo, branch))
		postsubmits[orgrepo] = append(postsubmits[orgrepo], *generatePostsubmitForTest(test, org, repo, branch))
	}

	return &presubmits, &postsubmits
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

	fmt.Printf(string(jobConfigAsYaml))
}

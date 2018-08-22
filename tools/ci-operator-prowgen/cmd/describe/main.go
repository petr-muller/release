package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	jc "github.com/openshift/release/tools/ci-operator-prowgen/pkg/jobconfig"

	prowconfig "k8s.io/test-infra/prow/config"
)

type options struct {
	prowJobConfig string

	help bool
}

func bindOptions(flag *flag.FlagSet) *options {
	opt := &options{}

	flag.StringVar(&opt.prowJobConfig, "prow-jobs", "", "Path to a file containing Prow jobs")
	flag.BoolVar(&opt.help, "h", false, "Show help for ci-operator-prowgen")

	return opt
}

func describeJobs(prowJobConfigPath string) error {
	var jobConfig *prowconfig.JobConfig
	var err error
	if jobConfig, err = jc.ReadFromFile(prowJobConfigPath); err != nil {
		return fmt.Errorf("Failed to read Prow jobs from '%s' (%v)", prowJobConfigPath, err)
	}

	if jobConfig.Presubmits != nil && len(jobConfig.Presubmits) > 0 {
		repos := []string{}
		for repo, jobs := range jobConfig.Presubmits {
			repos = append(repos, repo)
		}
		fmt.Printf("Presubmits for %d repos (%s)\n", len(jobConfig.Presubmits), strings.Join(repos, ","))
	}

	if jobConfig.Postsubmits != nil {
	}

	return nil
}

func main() {
	flagSet := flag.NewFlagSet("", flag.ExitOnError)
	opt := bindOptions(flagSet)
	flagSet.Parse(os.Args[1:])

	if opt.help {
		flagSet.Usage()
		os.Exit(0)
	}

	if len(opt.prowJobConfig) > 0 {
		if err := describeJobs(opt.prowJobConfig); err != nil {
			fmt.Fprintf(os.Stderr, "describe failed (%v)\n", err)

		}
	} else {
		fmt.Fprintf(os.Stderr, "describe tool needs the --prow-jobs\n")
		os.Exit(1)
	}
}

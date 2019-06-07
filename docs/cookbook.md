# OKD CI Cookbook

#### Rules

Each cookbook item is at most 5 lines (100-character wide), can be reasonably prolonged adjust for
the markdown markup, if needed. The items should be goal-oriented (how to do X), briefly and
formally written and ideally should provide many links to find out more about the topic. Each item
can provide a single snippet, not longer than 10 lines.

## Component configuration for ci-operator

### Enroll a repository to follow OpenShift branching rules

TBD

## Repositories on GitHub

### Create "official" OpenShift release branch

Run the [repo-brancher](https://github.com/openshift/ci-operator-prowgen/tree/master/cmd/repo-brancher)
tool from the [ci-operator-prowgen](https://github.com/openshift/ci-operator-prowgen) repository
over the ci-operator config files in the [openshift/release](https://github.com/openshift/ci-operator-prowgen)
repository. The tool only works on [official](#enroll-a-repository-to-follow-openshift-branching-rules)
repositories. The tool finds a branch matching the current release and creates new branches from it.
The tool can be limited to operate only on certain orgs or repos. The branches will be pushed to GH
using the provided credentials. If a branch exists already, it is not touched unless the tool
was run in the fast-forward mode. 

#### Example:
```
$ repo-brancher --config-dir ci-operator/config/ \
    --current-release=4.0 \
    --future-release=4.1 --future-release=4.2 \
    --username openshift-merge-robot --token-path oauth-token \
    --org openshift --repo client-go \
    --fast-forward --confirm
```

#### Documentation

1. [Centralized Release Branching and Config Management](https://docs.google.com/document/d/1USkRjWPVxsRZNLG5BRJnm5Q1LSk-NtBgrxl2spFRRU8/edit#heading=h.3myk8y4544sk)

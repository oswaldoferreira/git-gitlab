package main

import (
	"fmt"
	"github.com/plouc/go-gitlab-client"
	"os/exec"
	"strconv"
	"strings"
)

type GitLabClient struct {
	GitLab *gogitlab.Gitlab
}

func NewGitLabClient(config GitConfig) (*GitLabClient, error) {
	host, e := config.Host()
	if e != nil {
		return nil, e
	}
	path, e := config.ApiPath()
	if e != nil {
		return nil, e
	}
	token, e := config.Token()
	if e != nil {
		return nil, e
	}
	client := GitLabClient{
		GitLab: gogitlab.NewGitlab(host, path, token),
	}
	return &client, e
}

func (this *GitLabClient) clone(remote string, path string) (string, error) {
	if this == nil {
		return "", nil
	}
	fmt.Println("search project " + this.GitLab.BaseUrl + "/" + remote)
	// search remote repository
	projectID := strings.Replace(remote, "/", "%2F", -1)
	project, e := this.GitLab.Project(projectID)
	if e != nil {
		return "", e
	}
	// get remote URL
	remoteURL := project.SshRepoUrl
	// clone
	var dest string
	if len(path) < 0 {
		dest = project.Path
	} else {
		dest = path + "/" + project.Path
	}

	fmt.Println("clone repository from " + remoteURL + " to " + dest)
	repo, e := exec.Command("git", "clone", remoteURL, dest).Output()
	if e != nil {
		return "", e
	}
	res := string(repo)
	return res, e
}

func (this *GitLabClient) show(config GitConfig, dashboardFlag string, issuablePath string) {
	projectPath, err := config.Project()
	if err != nil {
		fmt.Println(err.Error())
	}

	hostPath, err := config.Host()
	if err != nil {
		fmt.Println(err.Error())
	}

	if issuablePath != "" {
		exec.Command("open", hostPath+"/"+projectPath+"/"+issuablePath).Output()
	} else if dashboardFlag != "" {
		fmt.Println("Fetching user information..")

		user, err := this.GitLab.CurrentUser()
		if err != nil {
			fmt.Println(err.Error())
		}

		exec.Command("open", hostPath+"/dashboard/"+dashboardFlag+"?assignee_id="+strconv.Itoa(user.Id)).Output()
	}
}

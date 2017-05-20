package util

import (
	"github.com/libgit2/git2go"
	"time"
)

type Repository struct {
	repo *git.Repository
}

func GitClone(url, branch, dest string, credentials git.CredentialsCallback) (*Repository, error) {
	cloneOptions := &git.CloneOptions{
		CheckoutBranch: branch,
		FetchOptions: &git.FetchOptions{
			RemoteCallbacks: git.RemoteCallbacks{
				CredentialsCallback: credentials,
			},
		},
	}
	repo, err := git.Clone(url, dest, cloneOptions)
	if err != nil {
		return nil, err
	}
	return &Repository{repo}, nil
}

func MakeGitCredentials(username, password string) git.CredentialsCallback {
	called := false
	return func(url string, username_from_url string, allowed_types git.CredType) (git.ErrorCode, *git.Cred) {
		if called {
			return git.ErrAuth, nil
		}
		called = true
		errCode, cred := git.NewCredUserpassPlaintext(username, password)
		return git.ErrorCode(errCode), &cred
	}
}

func (r *Repository) GitRemoteAdd(remote, url string) (*git.Remote, error) {
	return r.repo.Remotes.Create(remote, url)
}

func (r *Repository) GitRemoteRm(remote string) error {
	return r.repo.Remotes.Delete(remote)
}

func (r *Repository) GitAdd(path ...string) error {
	index, err := r.repo.Index()
	if err != nil {
		return err
	}
	for _, p := range path {
		err = index.AddByPath(p)
		if err != nil {
			return err
		}
	}
	return nil
}

func MakeGitSignature(name, email string) *git.Signature {
	return &git.Signature{
		Name:  name,
		Email: email,
		When:  time.Now(),
	}
}

func (r *Repository) GitCommit(branch, message string, author, committer *git.Signature) error {
	var err error

	index, err := r.repo.Index()
	if err != nil {
		return err
	}

	treeId, err := index.WriteTree()
	if err != nil {
		return err
	}

	err = index.Write()
	if err != nil {
		return err
	}

	tree, err := r.repo.LookupTree(treeId)
	if err != nil {
		return err
	}

	head, err := r.repo.Head()
	if err != nil {
		return err
	}

	commitTarget, err := r.repo.LookupCommit(head.Target())
	if err != nil {
		return err
	}

	_, err = r.repo.CreateCommit("refs/heads/"+branch, author, committer, message, tree, commitTarget)
	return err
}

func (r *Repository) GitPush(remote, branch string, credentials git.CredentialsCallback) error {
	options := git.PushOptions{RemoteCallbacks: git.RemoteCallbacks{
		CredentialsCallback: credentials,
	}}
	rem, err := r.repo.Remotes.Lookup(remote)
	if err != nil {
		return err
	}
	return rem.Push([]string{"refs/heads/" + branch}, &options)
}

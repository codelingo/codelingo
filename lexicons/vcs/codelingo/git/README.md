# codelingo/git vcs lexicon



##  git facts
<details><summary>git.commit</summary><p>

#### Example of finding every commit and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_commit
    doc:  Example query to find all instances of commit
    flows:
      codelingo/review
	       comments: This is a commit.
	   query: |
	     import codelingo/vcs/git

	     @ review.comment
	     git.commit
```
</p></details>

<details><summary>git.patch</summary><p>

#### Example of finding every patch and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_patch
    doc:  Example query to find all instances of patch
    flows:
      codelingo/review
	       comments: This is a patch.
	   query: |
	     import codelingo/vcs/git

	     @ review.comment
	     git.patch
```
</p></details>

<details><summary>git.repo</summary><p>

#### Example of finding every repo and having a review flow comment on it:

```yaml
tenets:
  - name: find_all_repo
    doc:  Example query to find all instances of repo
    flows:
      codelingo/review
	       comments: This is a repo.
	   query: |
	     import codelingo/vcs/git

	     @ review.comment
	     git.repo
```
</p></details>


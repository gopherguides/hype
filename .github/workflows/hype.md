# Using Github Actions to update your README

This repo uses the action to keep the README up to date.

## Requirements

For this action to work, you need to either configure your repo with specific permissions, or use a `personal access token`.

### Repo Permissions

You need to give permission to your GitHub Actions to create a pull request in your GitHub repo settings *(Settings -> Actions -> General)*.

Under `Workflow Permissions`

- Check `Allow GitHub Actions to create and approve pull requests`.
- Check `Read and write permissions`

### Personal Access Token

Alternately, you can use tokens to give permission to your action.

It is recommend to use a GitHub [Personnal Acces Token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token#creating-a-fine-grained-personal-access-token) like: `${{`{{secrets.PAT}}`}}` instead of using `{{`${{secrets.GITHUB_TOKEN}}`}}` in GitHub Actions.

## The Action

The current action is set to only generate the readme on a pull request and commit it back to that same pull request.  You can modify this to your own needs.

<code src="hype.yml"></code>

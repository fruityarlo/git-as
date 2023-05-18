# git-as

git-as is a custom git command that greatly simplifies the usage of
multiple GitHub accounts. All entries are stored in the git-as sections in your
per-user .gitconfig file.

Keep in mind that I wrote this tool mostly for myself and I use it on Linux (WSL2). I have no intetion to test it on Windows or MacOS. But thanks to `go` it _should_ be easy to use on these platforms as well.


## Installation

Run `go install github.com/fruityarlo/git-as` and you're good to go as long as your PATH is set up correctly. Maybe I'll use Releases on Github in the future so it's easier to just download a file.


## Usage

First, be sure to have a global default git user set up correctly. You can check with

```sh
git --global config user.name
# and
git --global config user.email
```

Then, be sure to always use the SSH stuff when cloning or pushing to GitHub (or any other provider). It won't work correctly otherwise.

And, please read the _Add a new user_ section below.

With that out of the way, let's see how to use git-as.

```sh
# CHANGE YOUR GIT USER for the current directory with <id> being
# the identifier used when adding the entry. 'default' (without
# the quotes) can also be used to get rid of the local git-as
# entries and use your global defaults.
#
# Examples: `git as johnd` or `git as default`
git as <id>

# CLONE A REPO AS a specific user with <id> being the identifier
# used when adding the entry and <repo> being the git repo to
# clone. All [options] are directly passed on to git-clone.
#
# Example: `git as johnd clone git@github.com:fruityarlo/git-as.git`
git as <id> clone <repo> [options]

# SHOW VERSION
git as version

# SHOW HELP
git as help
```

### Add a new user

First, let's generate a new SSH key for you new user (Let's call him _John Doe_ for now). There's this [GitHub page](https://docs.github.com/en/authentication/connecting-to-github-with-ssh/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent?platform=linux) with lots of info, but we can just

```sh
ssh-keygen -t ed25519 -C "john.doe@example.com"
```

and answer all questions. Oh, and please remember where you saved the key file, probably in `~/.ssh`. Also, don't overwrite your existing keys by rushing through the generation step! Then we can add _johnd_ to our per-user .gitconfig, which is probably `~/.gitconfig`. You can either edit `.gitconfig` directly or run some git commands:

In both cases, you have to generate the SSH certificate yourself. An you can mix
and match using commands or editing the config file directly, of course.

#### Edit `.gitconfig` directly

Add a section per user with all required info, e.g.

```ini
[git-as "johnd"]
        name: John Doe
        email: john.doe@example.com
        cert: ~/.ssh/id_ed25519_johnd
```

You can of course add as many section as you like. Keep in mind that the name in quotes (here: johnd) corresponds to the <id> argument. If you copy the example from above, make sure to fix the indentation.

#### Run some `git` commands

You can use git itself to add new users with 3 commands, e.g.

```sh
git config --global git-as.johnd.name "John Doe"
git config --global git-as.johnd.email john.doe@example.com
git config --global git-as.johnd.cert ~/.ssh/id_ed25519_johnd
```

You have to run all 3 commands for every user that you want to add. Remember: the string after git-as. (here: johnd) corresponds to the <id> argument.

And that's basically it. You can now use _johnd_ in a git repo with `git as johnd`


## Tips

While you can always run `git config user.name` to find out which persona you are using currently, you can display this directly in your shell (you're using oh-my-zsh or something, right?). Here's a custom `.zsh-theme` idea (you have to have the git plugin activated):

```sh
git_custom_status() {
  local cb=$(git_current_branch)
  local gituser=$(git config user.name)

  if [ -n "$cb" ]; then
    echo -n "$ZSH_THEME_GIT_PROMPT_PREFIX$cb$(parse_git_dirty)"
    if [ -n "$gituser" ]; then
      echo -n " $gituser"
    fi
    echo -n "$ZSH_THEME_GIT_PROMPT_SUFFIX"
  fi
}

ZSH_THEME_GIT_PROMPT_PREFIX=" ("
ZSH_THEME_GIT_PROMPT_SUFFIX=")"
ZSH_THEME_GIT_PROMPT_DIRTY="*"

PROMPT='%3~$(git_custom_status) %% '
```

Feel free to add fancy colors!

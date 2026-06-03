# Git App

Installs and configures Git with devgita integration.

## Recovering Lost Commits

Git rarely deletes commits immediately. A `reset --hard`, bad rebase, or force-push makes commits **unreferenced** — they survive ~90 days before GC.

### 1. Check reflog

```bash
git reflog --date=iso              # current branch
git reflog --date=iso --all        # all refs including remotes
```

Look for the line **before** the reset/rebase.

### 2. Scan dangling commits (if reflog isn't enough)

```bash
git fsck --no-reflogs --lost-found
```

List them readable, newest first:

```bash
for c in $(git fsck --no-reflogs 2>/dev/null | awk '/dangling commit/{print $3}'); do
  echo "$(git show -s --format='%ci %h %an | %s' $c)"
done | sort -r | head -30
```

### 3. Inspect before trusting

```bash
git show <hash>            # full diff
git show --stat <hash>     # files changed
git log --oneline <hash>   # chain behind it
```

### 4. Pin it (prevents GC)

```bash
git branch recovered <hash>
```

### 5. Restore

```bash
git reset --hard <hash>                      # move branch to commit
git push --force-with-lease origin <branch>  # restore remote
```

## Special Cases

**Staged but never committed (`git add` only):**

```bash
git fsck --lost-found                # writes blobs to .git/lost-found/other/
git show <blob-hash>                 # inspect contents
```

You get file contents but not names.

**Never staged:** Unrecoverable from git. Check editor undo/local history.

## Prevention

```bash
git push --force-with-lease    # refuses if remote changed
git branch backup              # before risky rebases
```

## Common Workflows

### Create a clean branch

```bash
git fetch origin
git checkout -b <branch> origin/main
git add . && git commit -m "feat: description"
git push -u origin <branch>
```

### Re-sync with main (preserving uncommitted work)

```bash
git reset --soft <commit-before-your-work>
git stash
git merge main
# resolve conflicts if any, then:
git stash pop
git restore --staged .   # optional: unstage
```

### Squash merge into clean branch

```bash
git fetch origin
git checkout -b <branch> origin/main
git merge --squash origin/<source-branch>
git commit -m "feat: combined description"
git push -u origin <branch>
```

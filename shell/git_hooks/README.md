# Git Hooks 

Check update file size before commit:

```bash
wget -qO .git/hooks/pre-commit  https://raw.githubusercontent.com/Dokiys/example/main/shell/git_hooks/pre_commit_filesize_check && git config core.hooksPath .git/hooks && chmod +x .git/hooks/pre-commit
```

Do golangci-lint check before push：

```bash
wget -qO .git/hooks/pre-push  https://raw.githubusercontent.com/Dokiys/example/main/shell/git_hooks/pre_push_golint_check && git config core.hooksPath .git/hooks && chmod +x .git/hooks/pre-push
```




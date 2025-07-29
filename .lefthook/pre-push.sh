#!/bin/sh

# 定义一个函数来处理错误退出
handle_error() {
    echo 1>&2 "  (use \"git push --no-verify\" to discard)";
    kill 0  # 发送 SIGTERM 信号给脚本中启动的所有进程
    exit 2
}

# Check WIP commits -- compatible with Lefthook (no stdin)
# ============================================
red_color='\033[0;31m'
reset_color='\033[0m'

# 获取当前本地分支名
local_branch=$(git symbolic-ref --quiet --short HEAD)

# 获取当前分支的上游远程分支，如 origin/main
remote_branch=$(git for-each-ref --format='%(upstream:short)' "refs/heads/${local_branch}")

# 推送的提交范围
if [ -z "$remote_branch" ]; then
  # 若无 upstream，推测为新建或第一次 push，使用最近 N 次提交
  range="HEAD~10..HEAD"
else
  # 正常使用远程分支为比较基线
  range="${remote_branch}..${local_branch}"
fi

# 查找所有包含 "--wip--" 的提交（表示工作中、未完成，不应推送）
commits=$(git log --pretty=format:'%H %s' $range | grep -- '--wip--')

if [ -n "$commits" ]; then
  echo >&2 "🚫 检测到包含 \"--wip--\" 的提交，阻止 push"
  echo >&2 "以下为匹配的提交记录（可使用 \"git push --no-verify\" 忽略检查）："
  echo "$commits" | while read -r line; do
    echo 1>&2 "$line"
  done
  echo >&2 ""
  exit 2
fi


# Check wire
(
  if command -v wire >/dev/null && find . -type f -name "wire.go" -print | grep -q .; then
    # 对所有 wire.go 所在目录执行 wire diff
    find . -type f -name "wire.go" | while read -r file; do
      dir=$(dirname "$file")
      wire diff "$dir" 1>&2
      if [ $? -ne 0 ]; then
        handle_error
      fi
    done
  fi
) &

# 等待所有后台任务完成
wait

exit 0
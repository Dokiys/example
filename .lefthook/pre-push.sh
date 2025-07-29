#!/bin/sh

# 定义一个函数来处理错误退出
handle_error() {
    echo 1>&2 "  (use \"git push --no-verify\" to discard)";
    kill 0  # 发送 SIGTERM 信号给脚本中启动的所有进程
    exit 2
}

# Check WIP commits && DELETE ME
(
  red_color='\033[0;31m' # red_color
  zero=$(git hash-object --stdin </dev/null | tr '[0-9a-f]' '0')
  while read local_ref local_oid remote_ref remote_oid; do
    if test "$local_oid" = "$zero"; then
      # Handle delete
      :
    else
      if test "$remote_oid" = "$zero"; then
        # New branch, examine parent branch
        split_id=$(git --no-pager reflog $local_ref | tail -n 1 | awk '{print $1}')
        range="$split_id..$local_ref"
        fixme_range="$split_id..$local_ref"
      else
        # Update to existing branch, examine new commits
        range="$remote_oid..$local_oid"
        fixme_range="$remote_oid..$local_oid"
      fi

  		# Check for WIP commits
  		remote_commit=$(git --no-pager log --pretty=format:'%s' $remote_ref -1 | grep -q -- '--wip--') # remote WIP commit
  		commits=$(git rev-list -n 1 --grep '--wip--' "$range") # local WIP commit
  		if [ -n "$commits" ] || [ -n "$remote_commit" ]; then
  			echo >&2 "Found [WIP] commits in $local_ref, not pushing"
        echo >&2 "  (use \"git push --no-verify\" to discard)"
        for commit in $commits; do
          echo 1>&2 "${red_color}commit ${commit}"
        done
        echo >&2 ""
        exit 2
  		fi

      # Check for FIXME ME
#      fixme_list=$(git --no-pager diff $fixme_range | grep '+' | grep 'FIXME\[Dokiy\]')
#      if test -n "$fixme_list" ; then
#        echo >&2 "Found [FIXME] commits in $local_ref, not pushing"
#        echo >&2 "  (use \"git push --no-verify\" to discard)"
#        echo "${red_color}${fixme_list}"
#        echo >&2 ""
#        exit 2
#      fi
    fi
  done
) &

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
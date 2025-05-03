# quick bash function to read .env file
# use it via:
# source readenv
# readenv
#
# or
#
# readenv <filename>
#
# modified from https://gist.github.com/mihow/9c7f559807069a03e302605691f85572
# fixed for whitespace issues, posix compliance (e.g. \t on mac means t)
#
# NOT a standalone script as when used as a standalone script, it'll read in the ENV variables into a sub-process, not the
# calling process

readenv() {
  local filePath="${1:-.env}"

  if [ ! -f "$filePath" ]; then
    return 0
  fi

  while read -r line; do
    if [[ "$line" =~ ^\s*#.*$ || -z "$line" ]]; then
      continue
    fi

    # 分割key和value
    key=$(echo "$line" | cut -d '=' -f 1 | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//')
    value=$(echo "$line" | cut -d '=' -f 2- | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//')

    # 处理引号
    # 1. 如果值以单引号或双引号开始和结束，则移除它们
    # 2. 保留值中间的引号
    if [[ $value =~ ^[\'\"].*[\'\"]$ ]]; then
        # 只移除首尾的引号
        value="${value:1:${#value}-2}"
    fi

    export "$key=$value"
  done < "$filePath"
}

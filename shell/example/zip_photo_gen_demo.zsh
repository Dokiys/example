#!/bin/zsh

# 检查参数数量
if [[ $# -ne 1 ]]; then
    echo "Usage: $0 <number_of_files>"
    exit 1
fi

# 参数赋值
num_files=$1

# 检查是否为有效的数值
if ! [[ "$num_files" =~ ^[0-9]+$ ]]; then
    echo "Error: The parameter must be a positive integer."
    exit 1
fi

# 检查 demo 文件夹是否存在，如果不存在则创建
folder="demo"
if [[ ! -d $folder ]]; then
    mkdir $folder
fi

# 文件生成循环
for i in {1..$num_files}; do
    # 使用 dd 命令创建 1 MB 文件
    dd if=/dev/zero of="${folder}/demo_file_${i}.bin" bs=1M count=1 2>/dev/null

    echo "Created ${folder}/demo_file_${i}.bin (1 MB)"
done

echo "Successfully created $num_files files in the '$folder' folder."

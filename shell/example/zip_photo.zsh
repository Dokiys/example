#!/bin/zsh

# 参数检查
if [[ $# -ne 3 ]]; then
    echo "Usage: $0 <password> <source_directory> <output_directory>"
    exit 1
fi

password=$1
directory=$2
output_directory="${$directory}_zip"  # 新增的输出目录参数
#max_size_bytes=$((3*1024*1024))  # 3MB
max_size_bytes=$((4*1024*1024*1024))  # 4GB

# 确认源目录存在
if [[ ! -d "$directory" ]]; then
    echo "Error: Source directory $directory does not exist."
    exit 1
fi

# 确认输出目录存在，如果不存在，则创建它
if [[ ! -d "$output_directory" ]]; then
    mkdir -p "$output_directory"
    if [[ $? -ne 0 ]]; then
        echo "Error: Unable to create output directory $output_directory."
        exit 1
    fi
fi

current_zip_index=1
current_zip_file="${output_directory}/${directory:t}_${current_zip_index}.zip"  # 修改文件的输出路径

current_zip_size=0

# 使用 find, stat (adapted for macOS) 和 sort 提取并排序文件
files_sorted_by_creation=($(find "$directory" -type f -print0 | xargs -0 stat -f "%m:%N" | sort -n | cut -d: -f2-))

for file in $files_sorted_by_creation; do
    # 估计添加当前文件后的 zip 大小
    estimated_size=$(stat -f%z "$file")
    new_size=$((current_zip_size + estimated_size))

    # 检查是否需要创建新的 zip 文件
    if [[ $new_size -gt $max_size_bytes ]]; then
        current_zip_index=$((current_zip_index + 1))
        current_zip_file="${output_directory}/${directory:t}_${current_zip_index}.zip"
        current_zip_size=0
    fi

    # 向当前 zip 文件添加文件
    zip -P "$password" --grow "$current_zip_file" "$file"
    current_zip_size=$((current_zip_size + estimated_size))
done

echo "打包完成。"
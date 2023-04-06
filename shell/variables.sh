#!/usr/bin/env bash

## basic
PRICE_PER_APPLE=5
MyFirstLetters="ABC"
MySecondLetters=${MySecondLetters:-"DEF"}
single_quote='Hello        world!'		# Won't interpolate anything.
double_quote="Hello   $MyFirstLetters     world!"		# Will interpolate.
escape_special_character="Hello   \$MyFirstLetters     world!"		# Use backslash.

one=`echo 1`
FileWithTimeStamp=/tmp/my-dir/file_$(/bin/date +%Y-%m-%d).txt

echo $PRICE_PER_APPLE                 # => 5
echo $MyFirstLetters                  # => ABC
echo $MySecondLetters                 # => DEF
echo $single_quote                    # => Hello world!
echo $double_quote                    # => Hello ABC world!
echo $escape_special_character        # => Hello $MyFirstLetters world!
echo $one                             # => 1
echo $FileWithTimeStamp               # => /tmp/my-dir/file_2022-09-06.txt

## array
my_array=(apple banana "Fruit Basket" orange)
new_array[2]=apricot

echo ${new_array[2]}                    # => apricot
echo ${my_array[3]}                     # => orange
# adding another array element
my_array[4]="carrot"                    # value assignment without a $ and curly brackets
echo "${my_array[@]}"                   # => apple banana Fruit Basket carrot
echo ${#my_array[@]}                    # => 5
echo ${my_array[${#my_array[@]}-1]}     # => carrot

## special variables
# ./variables.sh 1 "2 2" 3
echo "Script Name: $0"              # => Script Name: ./variables.sh
echo "Second arg: $2"               # => Second arg: 2 2
echo "Arg Num: $#"                  # => Arg Num: 3
echo "Exit status: $?"              # => Exit status: 0
echo "The process ID: $$"           # => The process ID: 4589
echo "As a separate argument: $@"   # => As a separate argument: 1 2 2 3
echo "As a one argument: $*"        # => As a one argument: 1 2 2 3

# --- $*
# 1 2 2 3
echo '--- $*'
for ARG in "$*"; do
  echo $ARG
done

# --- $@
# 1
# 2 2
# 3
echo '--- $@'
for ARG in "$@"; do
  echo $ARG
done

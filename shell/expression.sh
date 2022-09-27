#!/usr/bin/env bash

# decision making
NAME="Zhangsan"
if [ "$NAME" = "Zhangsan" ]; then
  echo "Hello Zhangsan!"
elif [ "$NAME" = "Lisi" ]; then
  echo "Hi Lisi!"
else
  echo "Wow"
fi

CASE_CONDITION="one"
case $CASE_CONDITION in
"one") echo "You selected bash" ;;
"two") echo "You selected perl" ;;
"three") echo "You selected python" ;;
"four") echo "You selected c++" ;;
"five") exit ;;
esac

VAR_A=(1 2 3)
VAR_B="be"
VAR_C="cat"
if [[ ${VAR_A[0]} -eq 1 && ($VAR_B = "bee" || $VAR_C = "cat") ]]; then
  echo "True"
fi

# loops
# for
echo
NAMES=('Joe Ham' Jenny Sara Tony)
for N in "${NAMES[@]}"; do
  echo "My name is $N"
done
for N in $(echo -e 'Joe Ham' Jenny Sara Tony); do
  echo "My name is $N"
done


# while
echo
WHILE_COUNT=3
while [ $WHILE_COUNT -gt 0 ]; do
  echo "Value of count is: $WHILE_COUNT"
  WHILE_COUNT=$(($WHILE_COUNT - 1))
done

# until
echo
UNTIL_COUNT=1
until [ $UNTIL_COUNT -gt 3 ]; do
  echo "Value of count is: $UNTIL_COUNT"
  UNTIL_COUNT=$(($UNTIL_COUNT + 1))
done

# break and continue
echo
CTRL_COUNT=3
while [ $CTRL_COUNT -gt 0 ]; do
  if [ $CTRL_COUNT -eq 1 ]; then
    CTRL_COUNT=$(($CTRL_COUNT - 1))
    continue
  fi
  echo "Value of count is: $CTRL_COUNT"
  if [ $CTRL_COUNT -eq 3 ]; then
      break
  fi
done

#!/bin/bash

# "----------------------------------------------"
#current_path=$(cd `dirname $0`; pwd)
current_path=`pwd`
#default temp file path
tmp_file_name="/dev/shm/go_test_file_list"
#default test file path
test_result_file_path="/dev/shm"
# "----------------------------------------------"

# "----------------------------------------------"
#list go test file
find $current_path -name "*_test.go" > $tmp_file_name
# "----------------------------------------------"

# "----------------------------------------------"
#filter folder
declare -A map=()
while read LINE
do
    file_name=${LINE##*/}
    test_path=`echo $LINE |awk -F $file_name '{print $1}'`
    map[$test_path]=$test_path
done < $tmp_file_name
# "----------------------------------------------"

# "----------------------------------------------"
#go test
count=1
path="/dev/shm"
for key in ${!map[@]}
do
    cd ${map[$key]}
    go test -v | tee $test_result_file_path/$count.gotest.log
    ((count+=1))
done
# "----------------------------------------------"

# "----------------------------------------------"
#check result
echo "=============================================="
run_num=`grep "=== RUN" $test_result_file_path/*.gotest.log -n | wc -l`
pass_num=`grep "\--- PASS" $test_result_file_path/*.gotest.log -n |wc -l`
fail_num=`grep "\--- FAIL" $test_result_file_path/*.gotest.log -n |wc -l`

echo RUN : $run_num
echo PASS : $pass_num
echo FAIL : $fail_num
echo "failed case:"
grep "\--- FAIL" $test_result_file_path/*.gotest.log | awk -F : '{print $2, $3}'
echo "=============================================="
# "----------------------------------------------"

# "----------------------------------------------"
rm $tmp_file_name
rm $test_result_file_path/*.gotest.log
# "----------------------------------------------"

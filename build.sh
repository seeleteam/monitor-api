#!/bin/sh

build_log="buildLog"
current_path=$(cd `dirname $0`; pwd)
shell_path=$current_path/tools/prebuild/shell

PROCESS_NAME=monitor-api

# --------------------------------------------------
function echo_red(){
    local str_info=$@
    echo -e "\033[31m $str_info \033[0m"
    return 0
}

function echo_green
{
    local content=$@
    echo -e "\033[32m $content \033[0m"
    return 0
}

function exec_func(){
    eval $@
    [[ $? -ne 0 ]] && {
        echo_red "cmd[$@] execute fail!"
        exit 1
    }
    return 0
}

# --------------------------------------------------
function usage
{
    echo_green "
Usage:
    $0 [\$1]
Options:
    h|-h|help|-help   usage help
    buildd            compile project(debug)
    buildr            compile project(release)
    test              execute junit test
    ############################################################################# 
    package           output
                      directory：
                      xxxx.tar.gz ─┬conf
                                   ├bin
                                   └log
    "
    return 0
}


# --------------------------------------------------
mkdir -p $build_log
chmod +x $shell_path/*.sh 2>/dev/null

case $1 in
    h|help|-h|-help)
        usage
    ;;
    buildd)
        $shell_path/2_compile.sh debug | tee $build_log/2_compile.log
    ;;
    buildr)
        $shell_path/2_compile.sh release | tee $build_log/2_compile.log
    ;;
    test)
        $shell_path/3_unit_test.sh | tee $build_log/3_unit_test.log
    ;;
    package)
        $shell_path/4_package.sh | tee $build_log/4_package.log
    ;;
    *)
        usage
    ;;
esac

exit 0
# --------------------------------------------------

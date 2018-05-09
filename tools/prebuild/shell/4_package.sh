#!/bin/bash

# "----------------------------------------------"
current_path=`pwd`
#default package file name
package_name="monitor-api"
package_file="$package_name.tar.gz"
# "----------------------------------------------"

# "----------------------------------------------"
cd $current_path
if [ -d "$package_file" ]
then
    rm -rf "$package_file"
fi
mkdir "$package_name"
mkdir "$package_name/log"
cp -R bin/* "$package_name"
cp -R "conf" "$package_name"

cd $current_path
if [ -f $package_file ]
then
    rm -f $package_file
fi
tar -zcvf $package_file "$package_name"
rm -rf "$package_name"
# "----------------------------------------------"

# "----------------------------------------------"
echo "ok"
# "----------------------------------------------"

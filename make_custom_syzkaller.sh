#!/usr/bin/env bash
set -x

# Pre-requisites: linux git repo cloned and buil + syzkaller image

# example usage
# ./make_custom_syzkaller.sh [syzkaller-folder-name] [name-of-the-branch]

# default values
#REPOSITORY_BRANCH=${2:-'master'}
REPOSITORY_NAME=${1:-'syz'}
REPOSITORY_BRANCH=${2:-'profiling-all-logs-dump'}


# pull and build the fuzzer
git clone https://github.com/RaulinN/syzkaller.git $REPOSITORY_NAME
cd $REPOSITORY_NAME
git checkout $REPOSITORY_BRANCH

# cp /root/run_syzkaller.sh .
cp /root/syzkaller_docker.cfg .
chmod +x run_syzkaller.sh

# cp /root/make_manager.sh .
chmod +x make_manager.sh

# Edit the configuration file
sed -i "s/{{SYZKALLER_FOLDER_NAME}}/$REPOSITORY_NAME/g" ./syzkaller_docker.cfg

echo ">> Run 'cd $REPOSITORY_NAME' and ./make_manager to build the project"

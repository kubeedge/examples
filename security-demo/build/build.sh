#!/bin/bash

check_error_and_exit()
{
  if [ $1 -ne 0 ]; then
    echo $2
    exit 1
  else
    echo $3
  fi
}

build_components()
{
  components=`cat ./component.list`
  while read component
  do
    echo $component
    name=`echo $component | awk -F'::' '{print $1}'`
    gitpath=`echo $component | awk -F'::' '{print $2}'`
    buildcmd=`echo $component | awk -F'::' '{$1=$2=""; print $0}'`
    build_component $name "$gitpath" "$buildcmd"
  done <<< $components
}

get_project_dir()
{
  removehttps=`echo $2 | awk -F"https://" {'print $2'}`
  ret=`echo $removehttps | awk -F"$1.git" {'print $1'}`
  echo $ret
}

build_component()
{
  echo "building component $1...\n"

  dir=$(get_project_dir $1 $2)

  projectdir=$GOPATH/src/$dir

  succ_msg="build is successful for components in $projectdir/$1"
  err_msg="failed to build components in $projectdir/$1"
  if [ -d $projectdir/$1 ]; then
    cd $projectdir/$1
  else
    ret=$(mkdir -p $projectdir)
    err_code=$?
    if [ ! -d $projectdir ]; then
      check_error_and_exit "$err_code" "$err_msg" "$succ_msg"
    fi
    cd $projectdir
    git clone $2
    cd $projectdir/$1
  fi
  echo "build command is => $3"
  ret=$(eval $3)
  err_code=$?
  check_error_and_exit "$err_code" "$err_msg" "$succ_msg"
}

build_components

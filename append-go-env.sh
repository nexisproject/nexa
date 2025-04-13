#!/bin/bash

function join() {
    local IFS="$1"
    shift
    echo "$*"
}

function append() {
  #arr=$(awk -F',' '{ for( i=1; i<=NF; i++ ) print $i }' <<<"$(go env "$1")")
  #echo "export $1=${arr[*]}"
  arr=()
  IFS=',' read -r -a arr <<< "$(go env "$1")"
  arr+=("$2")
  #arr=($(echo "${arr[@]}" | tr ' ' '\n' | sort -u | tr '\n' ' '))
  read -r -a arr <<< "$(echo "${arr[@]}" | tr ' ' '\n' | sort -u | tr '\n' ' ')"
  #echo "${arr[*]}"
  #echo "${arr[@]}"
  str=$(join , "${arr[@]}")
  #echo "$1=$str"
  go env -w "$1=$str"
  echo "$1=$(go env "$1")"
}

append GOPRIVATE "$@"
append GONOPROXY "$@"
append GONOSUMDB "$@"


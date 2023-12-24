#!/bin/bash

mkdir -p "crds/$1"

kubectl get crd --no-headers | grep "$1" | awk '{print $1}' | xargs -I {} bash -c "kubectl get crd {} -o yaml | kubectl neat > crds/$1/{}.yaml"

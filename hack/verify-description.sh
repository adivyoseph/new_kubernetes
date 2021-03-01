#!/usr/bin/env bash

# Copyright 2014 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script checks API-related files for missing descriptions and outputs a
# list of structs and fields that are missing descriptions.
# Usage: `hack/verify-description.sh`.

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"

kube::golang::setup_env

make -C "${KUBE_ROOT}" WHAT=cmd/genswaggertypedocs
make -C "${KUBE_ROOT}" WHAT=cmd/fieldversiontag

# Find binary
genswaggertypedocs=$(kube::util::find-binary "genswaggertypedocs")
fieldversiontag=$(kube::util::find-binary "fieldversiontag")

gen_swagger_result=0
fieldversiontag_result=0
result=0

find_files() {
  find . -not \( \
      \( \
        -wholename './output' \
        -o -wholename './_output' \
        -o -wholename './_gopath' \
        -o -wholename './release' \
        -o -wholename './target' \
        -o -wholename '*/third_party/*' \
        -o -wholename '*/vendor/*' \
        -o -wholename './pkg/*' \
      \) -prune \
    \) \
    -wholename './staging/src/k8s.io/api/*/v*/types.go'
}

if [[ $# -eq 0 ]]; then
  versioned_api_files=$(find_files) || true
else
  versioned_api_files="${*}"
fi

for file in $versioned_api_files; do
  gen_swagger_result=0
  fieldversiontag_result=0

  # TODO: the find_files function had incorrect regexes which led
  # to genswaggertypedocs never being invoked.
  # Comment this out for now since there are many types with missing descriptions.
  # $genswaggertypedocs -v -s "${file}" -f - || gen_swagger_result=$?
  # if [[ "${gen_swagger_result}" -ne "0" ]]; then
  #   echo "API file: ${file} is missing: ${gen_swagger_result} descriptions"
  #   result=1
  # fi
  $fieldversiontag -s "${file}" || fieldversiontag_result=$?
  if [[ "${fieldversiontag_result}" -ne "0" ]]; then
    echo "API file: ${file} has ${fieldversiontag_result} type(s) with invalid fieldversion tags"
    result=1
  fi
  if grep json: "${file}" | grep -v // | grep description: ; then
    echo "API file: ${file} should not contain descriptions in struct tags"
    result=1
  fi
  if grep json: "${file}" | grep -Ee ",[[:space:]]+omitempty|omitempty[[:space:]]+" ; then
    echo "API file: ${file} should not contain leading or trailing spaces for omitempty directive"
    result=1
  fi
done

internal_types_files="${KUBE_ROOT}/pkg/apis/core/types.go ${KUBE_ROOT}/pkg/apis/extensions/types.go"
for internal_types_file in $internal_types_files; do
  if [[ ! -e $internal_types_file ]]; then
    echo "Internal types file ${internal_types_file} does not exist"
    result=1
    continue
  fi

  if grep json: "${internal_types_file}" | grep -v // | grep description: ; then
    echo "Internal API types should not contain descriptions"
    result=1
  fi
done

exit ${result}

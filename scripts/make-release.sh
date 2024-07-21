#!/bin/bash

# Exit script if a statement returns a non-true return value.
set -o errexit
# Use the error status of the first failure, rather than that of the last item in a pipeline.
set -o pipefail

GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Display the options to the user
echo -e  "${BLUE}Which version are you releasing? (Omit the v)${NC}"
echo "Version:  "

read -r version

echo -e  "${BLUE}Supported builds:${NC}\nWindows\tLinuxARM64\tLinuxARM\tLinux"
echo "Your choice: "

# Read user input
read -r build_os

# Check if the input is one of the supported options
if [[ "$build_os" == "Windows" || "$build_os" == "LinuxARM64" || "$build_os" == "LinuxARM" || "$build_os" == "Linux" ]]; then
    echo "You chose: $build_os"
else
    echo -e "${RED}Error: Invalid choice.${NC}"
    exit 1
fi

# Setting the folder name
releaseFolder="Grafana GT7 - v${version}${build_os}Release"

# Making the folder
echo "${GREEN}Making the ${build_os} release folder${NC}"
cd ..
mkdir -p "${releaseFolder}"

npm run build
mage "${build_os}"

# Including files
echo "${GREEN}Copying dist...r${NC}"
cp -r ./dist "${releaseFolder}"

echo "${GREEN}Copying provisioning...r${NC}"
cp -r ./provisioning "${releaseFolder}"

echo "${GREEN}Copying docker-compose.yaml...r${NC}"
cp docker-compose.yaml "${releaseFolder}"

zip -r "grafana-gt7-v${version}-${build_os}-release.zip" "${releaseFolder}"

rm -rf ./dist
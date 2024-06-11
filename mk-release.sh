#!/usr/bin/env bash

set -o errexit
set -o nounset

#!/usr/bin/env bash

# Fetch latest from origin and checkout main
git fetch origin && git checkout main

# Pull latest from main
git pull origin main

# Run the make targets
make generate test testacc

# Get the latest version tag and propose a new one
base=$(git describe --abbrev=0 --tags)
minor=${base##*.}
major=${base%.*}
new_version="$major.$(($minor + 1))"

# Ask for a new release tag
read -p "Enter a new release tag [default: $new_version]: " version
version="${version:-$new_version}"

# Ask for an annotation for the tag
default_annotation="Release version $version"
read -p "Enter an annotation for the tag [default: '$default_annotation']: " annotation
annotation="${annotation:-$default_annotation}"

# Create the new tag
git tag -a "$version" -m "$annotation"
echo "Version $version tagged!"

# Ask if push the new tag
while true; do
    read -p "Do you want to push the new tag to origin? [y/n]: " yn
    case $yn in
        [Yy]* ) git push origin "$version"; break;;
        [Nn]* ) exit;;
        * ) echo "Please answer yes or no.";;
    esac
done

#!/usr/bin/env bash

plugin_version="${PLUGIN_VERSION}"
ci_tag="${TAG_NAME}"

if [ ! "${plugin_version}" ] && [ ! "${ci_tag}" ]; then
    echo "No version or tag specified."
    exit 1
fi

if [ "${plugin_version}" != "${ci_tag}" ]; then
    echo "Versions do not match: plugin@${plugin_version} tag@${ci_tag}"
    exit 1
fi

echo "Versions match: plugin@${plugin_version} tag@${ci_tag}"
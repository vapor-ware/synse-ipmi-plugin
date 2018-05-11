#!/usr/bin/env bash

plugin_version="${PLUGIN_VERSION}"
circle_tag="${CIRCLE_TAG}"

if [ ! "${plugin_version}" ] && [ ! "${circle_tag}" ]; then
    echo "No version or tag specified."
    exit 1
fi

if [ "${plugin_version}" != "${circle_tag}" ]; then
    echo "Versions do not match: plugin@${plugin_version} tag@${circle_tag}"
    exit 1
fi

echo "Versions match: plugin@${plugin_version} tag@${circle_tag}"
#!/usr/bin/env python3

import os
import sys

builds = (
    ('darwin', 'amd64'),
    ('linux', '386'),
    ('linux', 'amd64'),
    ('linux', 'arm'),
    ('windows', '386'),
    ('windows', 'amd64'),
)

arches = ('arm', 'amd64', '386')
platforms = ('darwin', 'linux', 'windows')

cmd = "GOOS={platform} GOARCH={arch} go build -o builds/{platform}-{arch}.{ext} {package}"  # noQA

for platform, arch in builds:
    os.system(cmd.format(
        platform=platform,
        arch=arch,
        ext='bin' if platform != 'windows' else 'exe',
        package=sys.argv[1],
    ))

# yumreposync
A tool for managing s3 yum repositories.

## Overview

RPM -> yumreposync server -> s3 yum repo

## Using

### Server
`yumreposync-server -addr :8080 -repodir ./myrepo -bucket s3://yumreposync-repo/public`

### Client
`yumreposync -server http://localhost:8080 package1.rpm [package2.rpm]...`

# Minio Client (mc) [![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/minio/minio?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

``mc`` provides minimal tools to work with Amazon S3 compatible cloud storage and filesystems. It has features to resume partial downloads, progress bar, parallel copy and network reconnect. ``mc`` is written in golang and released under [Apache license v2](./LICENSE).

## Commands

``mc`` implements the following commands
```
  ls		List files and folders.
  mb		Make a bucket or folder.
  cat		Display contents of a file.
  cp		Copy files and folders from many sources to a single destination.
  mirror	Mirror folders recursively from a single source to many destinations.
  session	Manage sessions for cp and mirror.
  share		Share documents via URL.
  diff		Compute differences between two files or folders.
  access	Set access permissions.
  config	Modify, add, remove alias from default configuration file [~/.mc/config.json].
  update	Check for new software updates.
  version	Print version.
```

## Install [![Build Status](https://api.travis-ci.org/minio/mc.svg?branch=master)](https://travis-ci.org/minio/mc) [![Build status](https://ci.appveyor.com/api/projects/status/3ng8bef7b3e1v763?svg=true)](https://ci.appveyor.com/project/harshavardhana/mc)

#### GNU/Linux

Download ``mc`` from https://dl.minio.io:9000/updates/2015/Sept/linux-amd64/mc

~~~
$ wget https://dl.minio.io:9000/updates/2015/Sept/linux-amd64/mc
$ chmod +x mc
$ ./mc
~~~

#### OS X

Download ``mc`` from https://dl.minio.io:9000/updates/2015/Sept/darwin-amd64/mc

~~~
$ curl https://dl.minio.io:9000/updates/2015/Sept/darwin-amd64/mc > mc
$ chmod +x mc
$ ./mc
~~~

#### Windows

Download ``mc`` from https://dl.minio.io:9000/updates/2015/Sept/windows-amd64/mc.exe

~~~
C:\Users\Username\Downloads> mc.exe
~~~

#### Source

If you do not have a working Golang environment, please follow [Install Golang](./INSTALLGO.md).

```sh
$ go get github.com/minio/mc
```

## Public Minio Server

Minio server is hosted at ``https://play.minio.io:9000`` for public use. This service is primarily intended for developers and users to familiarize themselves with Amazon S3 compatible cloud storage. Minio runs in memory mode with auto expiry of objects in about an hour.  No account signup is required, which means S3 compatible tools and applications can access this service without access and secret keys.

## How to use mc?

[![asciicast](https://asciinema.org/a/21576.png)](https://asciinema.org/a/21576?async)

## Configuring mc for Amazon S3

Update AccessKeyID and SecretAccessKey fields in your ``~/.mc/config.json`` configuration file by following [AWS Credentials Guide](http://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSGettingStartedGuide/AWSCredentials.html).

## Contribute to Minio Client
Please follow Minio [Contributor's Guide](./CONTRIBUTING.md)

### Jobs
If you think in Lisp or Haskell and hack in go, you would blend right in. Send your github link to callhome@minio.io.

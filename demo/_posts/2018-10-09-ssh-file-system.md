---
title: SSH File System
subtitle: Share, copy, deploy stuff to a remote folder over ssh.
layout: post
date: 2018-10-09T12:00:00.00Z
publish: true
---
> Dig deep!

> <a href="https://xkcd.com/1360">![XKCD Documents](https://imgs.xkcd.com/comics/old_files.png)</a>

Get this article read to you.

# Contents

* [Goal](#goal)
* [The SSH File System](#sshfs)
* [User Access Workflow](#user-access)
* [Setup](#setup)
  * [Debian Linux](#setup-debian)
  * [macOS](#setup-macos)


# <a name="goal"></a>Goal

A couple of weeks ago I found myself setting up an e-books server that was supposed to store my ever growing *"Books I'd like to but I'll rarely find the time to read"* collection. That imposed an interesting network problem I had to solve. I had a large external drive connected to another machine that I wanted to utilize, however, setting up a Samba server wasn't something I wanted to spend several cups of coffee on and NFS wasn't an option for me either. It occurred to me that I should be able to achieve my goal using SSH.

# <a name="sshfs"></a>The SSH File System

First, I must say, I just love this idea. Using SSH as means to provide seamless file I/O operations and access control via asymmetric crypto is something I'll always consider from now on. The tool behind this is [sshfs](https://github.com/libfuse/sshfs). It's being frequently [updated](https://github.com/libfuse/sshfs/blob/master/ChangeLog.rst) and one may be sure to find it in their package manager.

However, there are a couple of things to consider.

  1. It's expected that SSH transfers will be slower due to crypto (and of course network) operations, although **sshfs** does use [caching](https://github.com/libfuse/sshfs/blob/master/cache.c) and multi-threading to address this. In any case, this does not seem to be an issue for me in a local network over a wired ethernet connection.
  2. **sshfs** leads to some very interesting user management cases, but this may quickly get too complex to handle. I'd say that using this within a complex access permissions scenario could be painful.

# <a name="user-access"></a>User Access Workflow

Here's a simplified workflow on what happens behind the scenes.

<a href="http://i.imgur.com/7YPQPBA.png">![SSHFS User Access](http://i.imgur.com/7YPQPBA.png)</a>

  * **alice** on **Server A** needs to access **Server B**'s storage.
  * **bob** has full access to **Server B**'s storage.
  * **alice** generates an ssh key pair and uses that to establish an SSH connection and impersonate *Bob* on **Server B**.
  * **alice** creates a mount point on **Server A** and can induce file read/write operations on Storage via the established SSH channel.

# <a name="setup"></a>Setup

I have prepared two setup use cases - one based on Linux that I'm currently using and another one for macOS, which I mostly did for research purposes. It is most certainly possible to use **sshfs** on **Windows**, however, I have not done any setup research in that direction.

## <a name="setup-debian"></a>Debian Linux

Install **sshfs** on **Server A** via aptitude:

    sudo apt-get install sshfs

Install FUSE (Filesystem in Userspace) via aptitude:

    sudo apt-get install fuse

You would need to load the `fuse` module via:

    modprobe fuse

Check if the module is loaded via:

    lsmod | grep fuse

Run `ssh-keygen` on **Server A** to generate a key pair.

Copy the public key in `$HOME/.ssh/id_rsa.pub` from **Server A** and paste it to `$HOME/.ssh/authorized_keys` on **Server B**. Or just use `scp` to copy the file over.

Create a target mount point and mount it to **Server A**, e.g.,

    mkdir $HOME/mountpoint
    sshfs bob@server-b:/var/storage $HOME/mountpoint

You now have access to **Server B**'s storage space. To unmount run:

    fusermount -u $HOME/mountpoint

### <a name="setup-debian"></a>Mount via systemd.service

It'd be great to have this run as a service, so if **Server A** reboots, the mount point gets setup automatically. For systems using **systemd** we can do that the following way:

Create a `mount` file in your `$HOME` (or any other accessible) directory.

```bash
#!/bin/sh

CMD="$1"
MNT="/home/mountpoint"
TARGET="/var/storage/"

start() {
  exec sshfs bob@server-b:$TARGET $MNT
}

stop() {
  fusermount -u $MNT
}

case $1 in
  start|stop) "$1" ;;
esac
```

As a `root` user, create a `bob-mount.service` file in `/lib/systemd/system` and add the following inside:

```
[Unit]
Description=Server B Mount Service
After=network.target

[Service]
User=alice
Group=alice
Type=forking
ExecStart=/home/alice/mount start
ExecStop=/home/alice/mount stop

[Install]
WantedBy=multi-user.target
```

Start and stop the service via:

    systemctl start bob-mount
    systemctl stop bob-mount

Enable the service to run at boot time:

    systemctl enable bob-mount

## <a name="setup-macos"></a>macOS

Installation on macOS is pretty similar. Use **brew** to install the following:

    brew install Caskroom/cask/osxfuse
    brew install Caskroom/cask/sshfs

You should see `FUSE for macOS` in your **System Preferences** afterwards. To mount a share run:

    mkdir $HOME/mountpoint
    sshfs bob@server-b:/var/storage $HOME/mountpoint

To unmount just run:

    unmount $HOME/mountpoint

An easy way to mount after reboot is to create a mount script and put it in the `crontab`.
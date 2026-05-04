# Shipping Your Machine: Building a Container from Scratch

This repository contains the live-code demonstration for the presentation 
**"Shipping Your Machine: Building a Container in 40 Lines of Code."** It demonstrates how Docker and other container runtimes work under the hood by recreating process and filesystem isolation from first principles using standard Linux features.

> ## ⚠️ Important Note: Linux Required
>
> Containers are not a universal concept—they are a clever illusion built entirely out of **Linux kernel features**.
>
> **This code will NOT run natively on macOS (Darwin) or Windows.** It requires direct access to Linux namespaces. If you are on a Mac, you must run this inside a Linux Virtual Machine (like Multipass, Vagrant, or a cloud EC2 instance).

## Prerequisites

* A Linux environment (Ubuntu recommended)
* [Go](https://go.dev/doc/install) installed
* Docker (ironically, we only use Docker to steal a pristine filesystem!)

## Setup: Getting the Filesystem

Before running the Go code, we need an isolated filesystem to "jail" our process inside. We can use Docker to quickly export a bare-bones Ubuntu filesystem into a local directory.

Run these commands in your terminal:

```bash
# 1. Create a directory for the filesystem
mkdir -p ubuntu-rootfs

# 2. Spin up a dummy container
docker run -d --name fs-extractor ubuntu:latest sleep 60

# 3. Export the filesystem to a tar file
docker export fs-extractor -o ubuntu-fs.tar

# 4. Extract the tarball into our new directory and clean up
tar -xf ubuntu-fs.tar -C ubuntu-rootfs
rm ubuntu-fs.tar

# 5. Kill the dummy container
docker rm -f fs-extractor
```

You should now have an unpacked Linux directory structure inside ubuntu-rootfs.

## Running the Code

To execute the container wrapper, build and run the Go program, passing it the command you want to run inside your isolated environment (e.g., /bin/bash):

```bash
go run main.go run /bin/bash
```

## What's Happening Under the Hood?

This code strips away the magic of Docker and relies on two primary Linux primitives:

* `chroot` (The Filesystem Jail): We change the root directory of our Go process to point to the ubuntu-rootfs directory we unpacked earlier. If the process looks for /bin/bash, it finds it inside our isolated folder, completely unaware of the host machine's real hard drive.
* Namespaces (The Invisibility Cloak): We use syscall.CLONE_NEWPID and syscall.CLONE_NEWUTS to isolate the process tree and hostname. This ensures our containerized bash shell cannot see the host's background processes (like Chrome or Spotify).

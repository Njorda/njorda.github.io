---
layout: post
title: "Overlay filesystem"
subtitle: "Using rust to overlay file systems"
date: 2024-01-11
author: "Niklas Hansson"
URL: "/2024/01/11/overlay-file-system-rust"
---

In this blog post we will explore overlaying file systems with the goal of adding a binary to a file systems. First step is to have a filesystem to work with:

```bash
curl -o root-drive-with-ssh.img https://s3.amazonaws.com/spec.ccfc.min/ci-artifacts/disks/${ARCH}/ubuntu-18.04.ext4
```

# Copy the filesystem

The simples approach is to just copy the filesystem and add the binary. To do this the first step is to make a directory where we temporary will open the filesystem and then copy the file over.  


```bash
sudo mkdir /mnt/my-image
```

The next step is to mount the filesystem and then copy the binary over.


```bash
sudo mount -o loop /path/to/your/image.img /mnt/my-image
```


Now we are ready to copy over the binary and make it executable. 


```bash
# Copy the binary
sudo cp /path/to/your/binary /mnt/my-image/path/within/image

# (Optional) Set the binary to be executable
sudo chmod +x /mnt/my-image/path/within/image/binary
```


The last step is to unmount the image. 

```sydo 
sudo umount /mnt/my-image
```

The downside with this approach is that we get a new filesystem for each time we want a filesystem with a new binary. In our use case running a firecracker vmm, this is not that nice of an approach. Instead we would like to overlay it. 
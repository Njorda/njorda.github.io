---
layout: post
title: "Modal Tunnel"
subtitle: "Interactive shell"
date: 2025-08-16
author: "Niklas Hansson"
URL: "/2025/08/16/modal_tunnel"
---

Long time no sea ... 

[Modal](https://modal.com) offers great cloud services and one of the best things is the fast start up(and avilability) of GPU:s. I hace historically use some of the hyperscalers out there but it is just so slow and so many steps. Also modal offers 30$ a month for free. 

However I have been looking for an option to get a interactive terminal up on a machine with a GPU but not had a good way to do it. Today I found [Tunnel](https://modal.com/docs/guide/tunnels) which allows you to expose live TCP ports on a Modal container.  So today we will try to hack it out in order to get a shell with a GPU that can start up and shut down quickly!


The first version we will try out is(we will wait with the GPU until we have something working): 


```python
# terminal.py
import modal, subprocess, os

image = (
    modal.Image.debian_slim()
    .apt_install("bash", "curl", "ca-certificates")
    .run_commands(
        "set -eux; arch=$(dpkg --print-architecture); case $arch in amd64) asset=ttyd.x86_64 ;; arm64) asset=ttyd.aarch64 ;; *) echo Unsupported arch: $arch; exit 1 ;; esac; curl -fsSL -o /usr/local/bin/ttyd https://github.com/tsl0922/ttyd/releases/latest/download/$asset; chmod +x /usr/local/bin/ttyd"
    )
)

app = modal.App("modal-terminal", image=image)

@app.function()
def terminal():
    with modal.forward(7681) as tunnel:
        url = tunnel.url
        print(f"Open your terminal at: {url}")
        # ttyd serves your shell over the web; it runs until you stop it
        env = {**os.environ, "SHELL": "/bin/bash"}
        subprocess.run(
            ["ttyd", "-p", "7681", "-t", "title=Modal Terminal", "/bin/bash"],
            env=env,
            check=True,
        )

```

In order to run it(however you will have to [login](https://modal.com/login) to modal and set up an account if you dont alreay have one),[modal run](https://modal.com/docs/reference/cli/run): 

```bash
modal run terminal.py::terminal
```

In mycase I use UV and actually endup with `uv run modal run terminal.py::terminal`


After the build succeds you should get something like(Yes I have shut it down): 

```bash
✓ Created objects.
├── 🔨 Created mount /Users/niklashansson/Documents/github/njorda.github.io/python/modal-gpu/terminal.py
└── 🔨 Created function terminal.
Open your terminal at: https://pahywszy0fbshx.r448.modal.host
```

However this will open it in readonly and thus we need to update the `ttyd` commands. 

```python
    [
        "ttyd",
        "--writable",
        "-p",
        "7681",
        "-t",
        "title=Modal Terminal",
        "/bin/bash",
    ],
```


And that is it! Seems to work amazing. The defaule timeout is 300 seconds, however it can be [extended to 24 hours](https://modal.com/docs/guide/timeouts).

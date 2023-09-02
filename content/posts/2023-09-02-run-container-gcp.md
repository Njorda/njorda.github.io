---
layout: post
title: "How to Quickly Run a Docker Image on GCP"
subtitle: "When you need to test a docker container fast"
date: 2023-09-02
author: "Johan Hansson"
URL: "/2023/09/02/docker-iterations"
---

# Run docker container in GCP cloud shell

In this blog post, we will briefly go over how to build and run a Docker container as quickly as possible on Google Cloud Platform (GCP).

First, you need to create a Docker image. Something very simple should suffice. Here is a sample Dockerfile:

```Dockerfile
FROM python:3
RUN apt-get update && apt-get install -y vim
RUN python -m pip install pandas
```

After creating the Dockerfile, the next step is to build and push the image to Artifact Registry. One convenient way to do this is by using the `gcloud` command-line tool. This approach abstracts away the need to worry about the type of processor your local machine has and whether it is compatible with GCP.

Replace `<gcp_project>` with your GCP project ID and `<docker_image_tag>` with the desired Artifact Registry location and name for the container:

```bash
gcloud builds submit --project <gcp_project> --tag <docker_image_tag> .
```

Once the build is complete, you can check to ensure that the container has been pushed to Artifact Registry.

Now it's time to run the container. The quickest way to do this is using Cloud Shell. ![Cloud Shell](/img/add_code_shell.png)

Before running the container, authenticate the shell by executing the command below:

```bash
gcloud auth configure-docker europe-west4-docker.pkg.dev
```

Finally, you can run the Docker container with an interactive terminal. Replace `<image>` with the name of the image you pushed to Artifact Registry. You can specify a tag or use `:latest` if you didn't specify a tag:

```bash
docker run -it <image>:latest bash
```

Buy using this flow you can itterate and test you docker containers in a very fast manor. 
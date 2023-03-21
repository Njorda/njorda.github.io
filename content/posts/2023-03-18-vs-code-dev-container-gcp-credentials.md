
---
layout: post 
title: "Configuring Your Local Dev Container with GCP Default Credentials"
subtitle: "A Step-by-Step Guide to Setting Up Your Development Environment for Seamless Integration with Google Cloud Platform"
date: 2023-03-20
author: "Johan Hansson"
URL: "/2023/03/20/vs-code-dev-container-gcp-credentials/"
image: "/img/background_2022_07_17.png"
---

This blog post provides a step-by-step guide for setting up your VS Code dev container to work with Google Cloud Platform (GCP) services and APIs by configuring default GCP credentials. By authenticating your application with your GCP credentials, you can access the necessary resources without requiring additional authentication steps, saving time and streamlining your development workflow. 

## Configuring Default GCP Credentials

To use default credentials with GCP, you can follow the steps below:

1. Run the following command to add default credentials:

   ```bash
   gcloud auth application-default login
   ```

2. To check the contents of the default credentials, run the following command:
    ```bash
    cat ~/.config/gcloud/application_default_credentials.json
    ```
This will display the contents of the JSON file containing the credentials.


3. [Clone the base environment repository from the previous blog post.](https://github.com/Njorda/test_dev_containers).



4. To copy the default credentials to the local folder, run the following command:
```bash 
cp ~/.config/gcloud/application_default_credentials.json .devcontainer/application_default_credentials.json
```

5. To prevent accidental commits of sensitive information to Git, add a .gitignore file that excludes JSON files.

6. In the Dockerfile, add the following code to mount the default credentials and set the GOOGLE_APPLICATION_CREDENTIALS environment variable:

    ```bash 
        FROM python:3.9-slim

        COPY requirements.txt requirements.txt 
        RUN pip install --upgrade pip
        RUN pip install -r requirements.txt 
        COPY application_default_credentials.json application_default_credentials.json
        ENV GOOGLE_APPLICATION_CREDENTIALS=/application_default_credentials.json
    ```

7. To test the setup, rebuild the Docker container and add the following Python script:
    ```python 
    from google.cloud.resourcemanager import ProjectsClient

    for project in ProjectsClient().search_projects():
        print(project.display_name)

    ```

This should display a list of your projects. If the API is not enabled, you will receive a warning, indicating that your keys are working correctly.

While there are ways to improve this setup further, such as avoiding the need to copy the default credentials, this guide provides a simple way to configure your VS Code dev container to work with GCP services and APIs.



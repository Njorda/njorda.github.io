
---
layout:     post 
title: "Supercharge Your Development Workflow with VS Code Dev Containers"
subtitle: "Streamline your development process and improve consistency with container-based development environments"
date:       2023-03-18
author:     "Johan Hansson"
URL: "/2023/03/18/dev-containers/"
image:      "/img/background_2022_07_17.png"
---



## Setting Up default gcp credentials

[Read more gcp](https://cloud.google.com/docs/authentication/application-default-credentials#GAC)


To add default crendtials run the following command
```bash
gcloud auth application-default login
```

To check what is in the default credentials run 

```bash
cat ~/.config/gcloud/application_default_credentials.json
```

In this example we will build on the previous blog post, where we created an base enviroment `2023-03-17-vs-code-dev-container`. 

To copy the default crentials the local folder run the following command. 
```bash 
cp ~/.config/gcloud/application_default_credentials.json .devcontainer/application_default_credentials.json
```
But before things get out of hand and we buy misstake pushes our credentials to git lets add a git ignore file that excludes .json files! 


It is now possible to mount the default crednetials to the docker file and add the path to the env GOOGLE_APPLICATION_CREDENTIALS. Below is the update docker file. 


```bash 
    FROM python:3.9-slim

    COPY requirements.txt requirements.txt 
    RUN pip install --upgrade pip
    RUN pip install -r requirements.txt 
    COPY application_default_credentials.json application_default_credentials.json
    ENV GOOGLE_APPLICATION_CREDENTIALS=/application_default_credentials.json
```

To test run it rebuild the docker container and add the following small python script 


```python 
from google.cloud.resourcemanager import ProjectsClient

for project in ProjectsClient().search_projects():
    print(project.display_name)

```

You should now see a list of the projects you have, it might be that you haven't enabled this api but the you will get a warning and you keys are working as they should! 


## Can we do it nicer without copying default credentials into our folder? 


YES of course so let's lookinto how we update this! 



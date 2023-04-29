---
layout: post
title: "Setting up a Basic dbt Development Container for BigQuery in GCP"
subtitle: "Learn how to set up a basic dbt project in GCP and share a development container to kickstart your project."
date: 2023-04-29
author: "Johan Hansson"
URL: "/2023/04/29/setting-up-basic-dbt-dev-container"
image: "/img/background_2022_07_17.png"
---

In this post, you will learn how to set up a basic dbt project in Google Cloud Platform (GCP) and share a development container to kickstart your project. While there are numerous blog posts out there about dbt and BigQuery, none of them share how to set it up in a development container without using any of the dbt-cloud services (at least to my knowledge).


## Setting upp your enviorment

1. Set up a gcp project/ or take one you allready have

2. Make sure enabled billing and bq, we will not do anything expensive, it should actually cost close 0 euros to do this short demo. But you need to have it enabled.

3. You can check if bq works by running these queries. 

```sql
    select * from `dbt-tutorial.jaffle_shop.customers`;
    select * from `dbt-tutorial.jaffle_shop.orders`;
    select * from `dbt-tutorial.stripe.payment`;
```


4. Create a bq data_set that is multi regional in the EU(I am based in eu, but you can use US if you want)
![Create dataset](/img/artifact_registry_kubeflow.png)

## Running dbt

5. Ok time to create a dev container! There are two blogg post going over this in depth but the one I used can be found [here](https://github.com/Njorda/basic-dbt-setup). Rember to run the comand below to push your gcp credentials into the container folder. 
```bash
    cp ~/.config/gcloud/application_default_credentials.json .devcontainer/application_default_credentials.json
```

6. Time to run initiate the dbt project. 
```bash
    dbt init 
```

You will now get a number of questions, below you can see the path I suggest taking. However you need to set your project, your dataset and another dbt_project(at least if cloned my container) the mine. The gcp project need to match the key you added in previous step. 

```
Which database would you like to use?
[1] bigquery

(Don't see the one you want? https://docs.getdbt.com/docs/available-adapters)

Enter a number: 1
[1] oauth
[2] service_account
Desired authentication method option (enter a number): 1
project (GCP project id): johan-kubeflow
dataset (the name of your dbt dataset): dbt_test
threads (1 or more): 1
job_execution_timeout_seconds [300]: 30
[1] US
[2] EU
Desired location option (enter a number): 2
13:51:38  Profile test_3 written to profiles.yml using target's profile_template.yml and your supplied values. Run 'dbt debug' to validate the connection.
13:51:38  
Your new dbt project "test_3" was created!

For more information on how to configure the profiles.yml file,
please consult the dbt documentation here:

  https://docs.getdbt.com/docs/configure-your-profile

One more thing:

Need help? Don't hesitate to reach out to us via GitHub issues or on Slack:

  https://community.getdbt.com/

Happy modeling!
```

Everything should work now so lets test it! 

```Bash
    cd test_3
```

```bash
    dbt run
```

![This should run dbt and populate your bq dataset.](/img/dbt_test_bq.png)

## Conclusion
In this post, we learned how to set up a basic dbt project in GCP and share a development container to kickstart our project. With this setup, you can easily
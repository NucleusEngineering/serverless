<walkthrough-metadata>
  <meta name="title" content="Serverless Journey>: Building and deploying container images with Cloud Build" />
  <meta name="description" content="Learn how to build container images, store them in the registry and deploy them to Cloud Run" />
  <meta name="keywords" content="deploy, containers, console, run" />
</walkthrough-metadata>

# Serverless Journey: Building and deploying container images with Cloud Build

![Cloud Build Logo](https://storage.googleapis.com/gweb-cloudblog-publish/images/finserve_2022_OvLe6x5.max-700x700.jpg)

In this tutorial we'll learn how to use Cloud Build, the serverless CI system on Google Cloud. Instead of using Build Packs, we'll be using Dockerfiles and Cloud Builds own declarative configuration to leverage higher flexibility and control over how we build our images. Finally, we'll use Cloud Build to also deploy to Cloud Run, so we can continuously deliver updates to our Cloud Run services.

<walkthrough-tutorial-difficulty difficulty="2"></walkthrough-tutorial-difficulty>

Estimated time:
<walkthrough-tutorial-duration duration="15"></walkthrough-tutorial-duration>

To get started, click **Start**.

## Project setup

First, let's make sure we got the correct project selected. Go ahead and select the provided project ID.

<walkthrough-project-setup billing="true"></walkthrough-project-setup>

<walkthrough-enable-apis apis="cloudbuild.googleapis.com,
run.googleapis.com,
artifactregistry.googleapis.com">
</walkthrough-enable-apis>

## Cloud Build

Cloud Build is a serverless CI system in Google Cloud. It is a general purpose compute platform that allows you to execute all sorts of batch workflows. However, it really excels at building code and creating container images.

You don't need to provision anything to get started with using Cloud Build, it's serverless: simply enable the API and submit your jobs. Cloud Build manages all the required infrastructure for you. Per default, Cloud Build schedules build jobs on shared infrastructure, but it can be configured to run on a dedicated worker pool that is not shared with other users.

In the previous section of the Serverless journey we saw how to use Cloud Build to automatically build container images from source code using the magic of Build Packs. These are great to get you started quickly. Almost as quickly, you will realize that you need more control over your build process, so you will start writing your own Dockerfiles. Let's see how that works with Cloud Build.

## Building with Dockerfiles

## Building with a cloudbuild.yaml file

## Deploying to Cloud Run

## Setting up an automatic source triggers from Github

## Summary

You now know how Cloud Build can help you automate integrating your artifacts.

When you are ready to proceed to the next chapter to learn more about how to implement Google Cloud APIs from your code, execute the following in the shell:

```bash
cloudshell launch-tutorial .journey/03-extend.neos.md
```

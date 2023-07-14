
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

In order to get started let's set up some configuration:

```bash
# Set default locations for Cloud Run and Artifact Registry to europe-north1, Finland.

gcloud config set run/region europe-north1 
gcloud config set artifacts/location europe-north1 
```

We are building some container images in this tutorial, so let's set up a dedicated Docker iamge registry using Google Cloud Artifact Registry:

```bash
gcloud artifacts repositories create my-repo --repository-format docker

# Also update config to use this repo
gcloud config set artifacts/repository my-repo
```

Great! Let's get started.

## Building with Dockerfiles

Previously, we build our container images by spceifying the `--source` flag, which will cause Cloud Build to use [Google Cloud Buildpacks](https://cloud.google.com/docs/buildpacks/overview) to automatically determine the toolchain of our application. It then uses templated container build instructions which are based on best practices to build and containerize the application.

Sometimes it is required to exert more control over the Docker build process. Google Cloud Build also understands Dockerfiles. With Dockerfiles it is possible to precisely specify which instructions to execute in order to create the container image. Cloud Build typically inspects the build context and looks for build instructions in the following order:

1. cloudbuild.yaml: If this file is present, Cloud Build will execute the instructions define in it. We'll learn about cloudbuild.yaml files later in this tutorial.
2. Dockerfile: If this file is present, Cloud Build will use it to execute a container build process that is equivalent to running `docker build`. Finally, it will tag and store the resulting image.
3. Buildpacks: If no explicit build instructions are available, Cloud Build can still be used to integrate by using Google Cloud Buildpacks.

Create a new file called `Dockerfile` and place it in the same directory as the Go application. Let's start with a simple, single-stage build for Go:

```Dockerfile
FROM golang:1.20-bullseye
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY main.go ./
RUN CGO_ENABLED=0 go build -v -o server
CMD ["/app/server"]
```
This Dockerfile:
1. uses Debian Bullseye with Golang 1.19 as base image
2. copies the definition of the Go module and installs all dependencies
3. copies the sources and builds the Go application
4. specifies the application binary to run

We can now submit the build context to Cloud Build and use the instructions in the Dockerfile to create a new image by running:

```bash
gcloud builds submit -t $(gcloud config get-value artifacts/location)-docker.pkg.dev/$(gcloud config get-value project)/my-repo/dockerbuild .
```

Next, let's navigate to [Artifact Registry in the Google cloud Console](https://console.cloud.google.com/artifacts/docker/), find the repo and inspect the image.

Huh, it seems like our image is quite big! We can fix this by running a [multi-stage Docker build](https://docs.docker.com/build/building/multi-stage/). Let's extend the Dockerfile and replace its contents with the following:

```Dockerfile
FROM golang:1.20-bullseye as builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY main.go ./
RUN CGO_ENABLED=0 go build -v -o server

FROM gcr.io/distroless/static-debian11
WORKDIR /app
COPY --from=builder /app/server /app/server
CMD ["/app/server"]
```

Great! We'll keep the first stage as build stage. Once the statically-linked Go binary is built, it is copied to a fresh stage that is based on `gcr.io/distroless/static-debian11`. The _distroless_ images are Google-provided base images that are very small. That means there is less to store and images typically have a much smaller attack surfaces. Let's build again:

```bash
gcloud builds submit -t $(gcloud config get-value artifacts/location)-docker.pkg.dev/$(gcloud config get-value project)/my-repo/dockermultistage .
```

Navigate back to Artifact Registry in the Google Cloud Console and inspect the new image `dockermultistage`. The resulting image is much smaller.

## Building with a cloudbuild.yaml file

Cloud Build will build, tag and store your container image when you submit a Dockerfile, but it can do a lot more. You can instruct Cloud Build to execute any arbitrary sequence of instructions by specifying and submitting a `cloudbuild.yaml` file. When detected in the build context, this file will take precedence over any other configuration and will ultimately define what your Cloud Build task does.

In a nutshell, you are able to specify a sequence of steps with each step being executed in a container. You can specify which image and process to execute for every individual step. The steps can collaborate and share data via a shared volume that is automatically mounted into the filesystem at `/workspace`.

Go ahead and create a `cloudbuild.yaml` and place it in the same directory as the Dockerfile and the Go application. Add the following contents:

```yaml
steps:
- name: 'ubuntu'
  args: ['echo', 'hi there']
```

This configuration asks Cloud Build to run a single step: Use the container image provided by [ubuntu](https://hub.docker.com/_/ubuntu) and run `echo`.

While this might not be a terribly useful or sophisticated CI pipeline, it is a useful illustration to understand how versatile Cloud Build is. Let's run it by submitting a new Cloud Build task:

```bash
gcloud builds submit
```

OK. Let's modify the `cloudbuild.yaml` to do something that's actually useful. Replace the file contents with the following:

```yaml
steps:
- name: 'gcr.io/cloud-builders/docker'
  args: ['build', '-t', '$_REGION-docker.pkg.dev/$PROJECT_ID/my-repo/dockercloudbuildyaml', '.']
- name: 'gcr.io/cloud-builders/docker'
  args: ['push', '$_REGION-docker.pkg.dev/$PROJECT_ID/my-repo/dockercloudbuild']
```

Pay attention to the `$`-prefixed `$PROJECT_ID` and `$_REGION`. `$PROJECT_ID` is a pseudo-variable which is automatically set by the Cloud Build environment. `$_REGION` is what Cloud Build calls a _substitution_ and it allows users to override otherwise static content of the build definition by injecting configuration during invokation of the build. Let's take a look:

```bash
gcloud builds submit --substitutions _REGION=$(gcloud config get-value artifacts/location)
```

This will build, tag and store a container image.

To fully understand what else can go into your `cloudbuild.yaml`, please check out the [schema documentation](https://cloud.google.com/build/docs/build-config-file-schema).

## Deploying to Cloud Run

Cloud Build is a versatile tool and is suited to run a wide variety of batch-like jobs. Until now, we only used Cloud Build to accomplish Continuous Integration (CI) tasks, but we don need to stop there.

We can extend the `cloudbuild.yaml` definition and automatically deploy the newly created image to Cloud Run like this:

```yaml

```

<!-- TODO extending cloudbuild.yaml, gcloud command, service agent perms in console-->

## Setting up an automatic source triggers from Github

<!-- TODO fork on GH, connect repo, create trigger, commit and push code and see it execute -->

## Summary

You now know how Cloud Build can help you automate integrating your artifacts.

When you are ready to proceed to the next chapter to learn more about how to implement Google Cloud APIs from your code, execute the following in the shell:

```bash
cloudshell launch-tutorial .journey/03-extend.neos.md
```

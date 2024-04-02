<walkthrough-metadata>
  <meta name="title" content="Serverless Journey>: Building and deploying container images with Cloud Build" />
  <meta name="description" content="Learn how to build container images using Cloud Build" />
  <meta name="keywords" content="deploy, containers, console, run" />
</walkthrough-metadata>

# Serverless Journey Module 2: Building and deploying container images with Cloud Build

![Tutorial header image](https://storage.googleapis.com/gweb-cloudblog-publish/images/finserve_2022_OvLe6x5.max-700x700.jpg)

Check out this [Overview on Cloud Build](https://images.squarespace-cdn.com/content/v1/65a6226068668c33fe4a4676/b21eb739-3638-416c-9356-7ace972f11f4/Chapter05Spread05Figure01.png?format=2500w)

In this tutorial we'll learn how to use Cloud Build, the serverless CI system on Google Cloud. Instead of using Build Packs, we'll be using Dockerfiles and Cloud Builds own declarative configuration to leverage higher flexibility and control over how we build our images. Finally, we'll use Cloud Build to also deploy to Cloud Run, so we can continuously deliver updates to our Cloud Run services.

<walkthrough-tutorial-difficulty difficulty="3"></walkthrough-tutorial-difficulty>

Estimated time:
<walkthrough-tutorial-duration duration="35"></walkthrough-tutorial-duration>

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

In the previous section of the Serverless journey we saw how to use build and deploy directly to Cloud Run from source code using the magic of Build Packs. Actually, this deployment via `gcloud run deploy` already leverages Cloud Build in the background, as you can see here in the [Cloud Build Dashboard](https://console.cloud.google.com/cloud-build/dashboard). These are great to get you started quickly. Almost as quickly, you will realize that you need more control over your build process, so you will start writing your own Dockerfiles. Let's see how that works with Cloud Build.

In order to get started let's set up some configuration:

```bash
gcloud config set project <walkthrough-project-id/>
gcloud config set run/region europe-north1 
gcloud config set artifacts/location europe-north1 
```

We are building some container images in this tutorial, so let's set up a dedicated Docker image registry using Google Cloud Artifact Registry and remember it in `gcloud config`:

```bash
gcloud artifacts repositories create my-repo --repository-format docker
gcloud config set artifacts/repository my-repo
```

Great! Let's get started.

## Building with Dockerfiles

Previously, we build our container images by specifying the `--source` flag, which will cause Cloud Build to use [Google Cloud Buildpacks](https://cloud.google.com/docs/buildpacks/overview) to automatically determine the tool chain of our application. It then uses templated container build instructions which are based on best practices to build and containerize the application.

Sometimes it is required to exert more control over the Docker build process. Google Cloud Build also understands Dockerfiles. With Dockerfiles it is possible to precisely specify which instructions to execute in order to create the container image. Cloud Build typically inspects the build context and looks for build instructions in the following order:

1. cloudbuild.yaml: If this file is present, Cloud Build will execute the instructions define in it. We'll learn about cloudbuild.yaml files later in this tutorial.
2. Dockerfile: If this file is present, Cloud Build will use it to execute a container build process that is equivalent to running `docker build`. Finally, it will tag and store the resulting image.
3. Buildpacks: If no explicit build instructions are available, Cloud Build can still be used to integrate by using Google Cloud Buildpacks.

In the 
<walkthrough-editor-spotlight spotlightId="activity-bar-explorer">
File Explorer
</walkthrough-editor-spotlight>
create a new file called `Dockerfile` and place it in the same directory as the Go application. Let's start with a simple, single-stage build for Go:

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
1. uses Debian Bullseye with Golang 1.20 as base image
2. copies the definition of the Go module and installs all dependencies
3. copies the sources and builds the Go application
4. specifies the application binary to run

We can now submit the build context to Cloud Build and use the instructions in the Dockerfile to create a new image by running:

```bash
gcloud builds submit -t $(gcloud config get-value artifacts/location)-docker.pkg.dev/$(gcloud config get-value project)/my-repo/dockerbuild .
```

Next, let's navigate first to the [Cloud Build Dashboard](https://console.cloud.google.com/cloud-build/dashboard) to see the build we just started. As soon as that is finished we go to [Artifact Registry in the Google cloud Console](https://console.cloud.google.com/artifacts/docker?project=<walkthrough-project-id/>) to find the repository and inspect the image.

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

In a nutshell, you are able to specify a sequence of steps with each step being executed in a container. You can specify which image and process to execute for every individual step. The steps can collaborate and share data via a shared volume that is automatically mounted into the file system at `/workspace`.

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
  args: ['build', '-t', '$_ARTIFACT_REGION-docker.pkg.dev/$PROJECT_ID/my-repo/dockercloudbuildyaml', '.']
- name: 'gcr.io/cloud-builders/docker'
  args: ['push', '$_ARTIFACT_REGION-docker.pkg.dev/$PROJECT_ID/my-repo/dockercloudbuildyaml']
```

Pay attention to the `$`-prefixed `$PROJECT_ID` and `$_ARTIFACT_REGION`. `$PROJECT_ID` is a pseudo-variable which is automatically set by the Cloud Build environment. `$_ARTIFACT_REGION` is what Cloud Build calls a _substitution_ and it allows users to override otherwise static content of the build definition by injecting configuration during invocation of the build. Let's take a look:

```bash
gcloud builds submit --substitutions _ARTIFACT_REGION=$(gcloud config get-value artifacts/location)
```

This will build, tag and store a container image.

To fully understand what else can go into your `cloudbuild.yaml`, please check out the [schema documentation](https://cloud.google.com/build/docs/build-config-file-schema).

## Deploying to Cloud Run

Cloud Build is a versatile tool and is suited to run a wide variety of batch-like jobs. Until now, we only used Cloud Build to accomplish Continuous Integration (CI) tasks, but we don't need to stop there.

We can extend the `cloudbuild.yaml` definition and automatically deploy the newly created image to Cloud Run like this:

```yaml
steps:
- name: 'gcr.io/cloud-builders/docker'
  args: ['build', '-t', '$_ARTIFACT_REGION-docker.pkg.dev/$PROJECT_ID/my-repo/dockercloudbuildyaml', '.']
- name: 'gcr.io/cloud-builders/docker'
  args: ['push', '$_ARTIFACT_REGION-docker.pkg.dev/$PROJECT_ID/my-repo/dockercloudbuildyaml']
- name: 'gcr.io/cloud-builders/gcloud'
  args:
    - 'run'
    - 'deploy'
    - '--region=$_RUN_REGION'
    - '--image=$_ARTIFACT_REGION-docker.pkg.dev/$PROJECT_ID/my-repo/dockercloudbuildyaml'
    - 'jokes'
```

In order for this to work, we need to grant some permissions to the service account used by the Cloud Build service agent, so it can interact with Cloud Run resources on our behalf.

Navigate to the [settings section of Cloud Build in the Google Cloud Console](https://console.cloud.google.com/cloud-build/settings/service-account) and assign both the **Cloud Run Admin** and the **Service Account User** roles to Cloud Build's service account.

Finally, let's supply all the substitutions and invoke Cloud Build once again to execute a fresh build and deploy to Cloud Run automatically.

```bash
gcloud builds submit --substitutions \
    _ARTIFACT_REGION=$(gcloud config get-value artifacts/location),_RUN_REGION=$(gcloud config get-value run/region)
```

## Setting up an automatic source triggers from GitHub

Now that we have successfully created a full build definition on Cloud Build for a simple yet very useful CI/CD pipeline, let's have a look at creating automatic build triggers that execute the pipeline whenever code gets updated on a remote Git repository. In order to put everything in place, we'll now look at how to authenticate with GitHub, fork the code into a new repository, authorize GitHub to connect with Cloud Build and create an automatic build trigger on Cloud Build.

<walkthrough-info-message>This part of the tutorial is **optional**, as it requires you to have a GitHub account. If you don't have one you can [register a new one](https://github.com/signup).</walkthrough-info-message>

### Authenticating with GitHub

First, let's authenticate using `gh`, GitHub's CLI, which comes pre-installed on Cloud Shell. In your terminal begin the authentication handshake and complete the CLI wizard.

We'd like to connect in the following matter:
 * Connect to `github.com`
 * Use HTTPS for authentication
 * Configure `git` to use GitHub credentials

Begin the authentication exchange by running the following:

```bash
gh auth login
```

Copy the authentication code, navigate to [GitHub's device authentication page](https://github.com/login/device), paste the token and authorize the application.

Once authentication is complete, you should be able to test access to GitHub by listing all of you repositories, like so:

```bash
gh repo list
```

### Creating a fork of the upstream repo

Next, we'll create a fork of the repository we are currently working with, so that we have a separate, remote repository that we can push changes to. Create the fork by running the following:

```bash
gh repo fork
```

This will create a fork of the upstream repository from `github.com/nucleusengineering/serverless` and place it at `github.com/YOUR_GITHUB_HANDLE/serverless`. Additionally, the command automatically reconfigures the remotes of the local repository clone. Take a look the the configured remotes and see how it configured both `upstream` and the new `origin` by running the following:

```bash
git remote -v
```

<walkthrough-info-message>Make sure that your new `origin` points to the work belonging to your own user and push commits to the new `origin` from now on.</walkthrough-info-message>

### Setting up the Cloud Build triggers

Next, Cloud Build needs to be configured with a trigger resource. The trigger contains everything Cloud Build needs to automatically run new builds whenever the remote repository gets updated.

First, navigate to the [Cloud Build triggers section of the Google Cloud Console](https://console.cloud.google.com/cloud-build/triggers) and click on 'Connect Repository'. Follow the wizard on the right to connect to github.com, authenticate Google's GitHub app, filter repositories based on your user name, find the forked repository called 'serverless', tick the box to accept the policy and connect your repository. Once completed you should see a connection confirmation message displayed. Now it's time to create the trigger.

Now, hit 'Create Trigger' and create a new trigger. In the wizard, specify that the trigger should read configuration from the provided `./cloudbuild.yaml` and **add all the substitutions** you used previously to trigger your build.

<walkthrough-info-message>When using the "Autodetect" configuration option, there is no possibility to add substituiton variables through the UI. So make sure to specify the "Cloud Build configuration file (yaml or json)" option explicitly, and then continue to fill in the substitution variables.</walkthrough-info-message>


### Pushing some changes

We should now have everything in place to automatically trigger a new build whenever changes are pushed to `main` on the remote repository fork.

If you haven't done so already, you will need to configure git on your Cloud Shell. To do so, run the following and configure your email and name, so git know who you are.

```bash
git config --global user.email "you@example.com"
git config --global user.name "Your Name"
```

Use git to add all the changes to the index and create a commit like this:

```bash
git add Dockerfile
git add cloudbuild.yaml
git commit -m 'add build config'
```

Finally, push the commit to the remote repossitory:

```bash
git push origin main
```

Changes should automatically be detected and trigger a new Cloud Build task. Navigate to the [Cloud Build dashboard](https://console.cloud.google.com/cloud-build/dashboard) and explore the running build.

## Summary

You now know how Cloud Build can help you automate integrating your artifacts.

<walkthrough-conclusion-trophy></walkthrough-conclusion-trophy>

This completes Module 2. You can now wait for the live session to resume or continue by yourself and on-demand.

When you are ready to proceed to the next chapter to learn more about how to implement Google Cloud APIs from your code, execute the following in the shell:

```bash
cloudshell launch-tutorial .journey/02-extend.neos.md
```

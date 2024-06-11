<!--markdownlint-disable MD024 MD033 MD036 MD041 -->
<walkthrough-metadata>
  <meta name="title" content="Giggle Cloud Platform" />
  <meta name="description" content="Learn how to build and run containerized
  services, use automated build and deployment tools, implement Google Cloud
  APIs to leverage GenAI capabilities and operate your service in production
  using SRE practices." />
  <meta name="keywords" content="cloud, go, containers, console, run, AI,
  GenAI, APIs, SRE" />
</walkthrough-metadata>

# From Prototype to Production

_"Your ML models are looking real fit, they must do a lot of weighted
training."_

Welcome to this hands-on tutorial in which you will learn how to build hilarious
services on Google Cloud Platform!

<walkthrough-info-message>Google is phasing out 3rd party cookies in Chrome to
better protect the privacy of its users. If you are attempting this tutorial on
Chrome, please make sure that 3rd party cookies are enabled in Chrome's address
bar (the icon that looks like an eye), since the tutorial still requires these.
The same is required for any other browser.</walkthrough-info-message>

## Overview

![Tutorial header image](https://raw.githubusercontent.com/NucleusEngineering/serverless/main/.images/containers.jpg)

In this interactive tutorial, you will learn how to build and run containerized
services, use automated build and deployment tools, implement Google Cloud APIs
to leverage GenAI capabilities and operate your service in production using SRE
practices.

The tutorial is split into **five modules**:

**Module 1** will explore Cloud Run, a compute platform for container workloads,
which you can use to easily run your services in a scalable, secure,
cost-effective and production-grade environment in a matter of minutes.

**Module 2** will get you familiar with Cloud Build, a fully-managed
container-based build system. In this module you will add more and more
automation to your service delivery. Code will be stored on a remote version
control system which automatically executes multi-stage Docker builds followed
by safe deployments to Cloud Run.

**Module 3** shows how to interact with Google Cloud APIs. You'll then use those
skills to implement interactions with Gemini, a large language model hosted on
Google cloud Vertex AI, to bring Generative AI capabilities to your
applications. The module also explores the usage of IAM and service accounts to
secure usage of APIs.

In **module 4**, we'll look at the practices of Site Reliability Engineering and
how you can practically apply these to reliably operate your services and manage
how you introduce changes to its functionality.

In **module 5**, explores how to strengthen the security posture of our software
delivery process by applying best practices for dependency management,
leveraging container analysis tools and signing produce images. Additionally,
we'll learn about the SLSA framework to better understand how to improve supply
chain security of our software delivery processes and subsequently gain more
confidence in the integrity and correctness of the artifacts we build, deploy
and operate.

Hope you're not afraid to get your hands dirty. Make sure get a beverage and
some snacks of your choice, strap in and hit 'Start'.

## Module 1: Run a freshly built container image on Cloud Run

![Tutorial header image](https://raw.githubusercontent.com/NucleusEngineering/serverless/main/.images/run.jpg)

In this tutorial we'll learn how to remotely build a container image from source
code, store the image in Artifact Registry and deploy it to Cloud Run. After
that, we'll familiarize ourselves with the Cloud Run UI, its API and explore
some of the available options to tweak our service.

[Overview on Cloud Run](https://raw.githubusercontent.com/NucleusEngineering/serverless/main/.images/overview-run.png)

<walkthrough-tutorial-difficulty difficulty="2"></walkthrough-tutorial-difficulty>

Estimated time:
<walkthrough-tutorial-duration duration="30"></walkthrough-tutorial-duration>

To get started, click **Start**.

## Project setup

First, let's make sure we got the correct project selected. Go ahead and select
the provided project ID.

<walkthrough-project-setup billing="true"></walkthrough-project-setup>

<walkthrough-enable-apis apis="cloudbuild.googleapis.com,
run.googleapis.com,
artifactregistry.googleapis.com"> </walkthrough-enable-apis>

You can use the builtin
<walkthrough-editor-spotlight spotlightId="menu-terminal-new-terminal"> Cloud
Shell Terminal </walkthrough-editor-spotlight> to run the gcloud command.

## Exploring the code

Take some time and
<walkthrough-editor-open-file filePath="cloudshell_open/serverless/main.go">
familiarize yourself with the code. </walkthrough-editor-open-file>

Also have a look at the dependencies reference in the Go module. You can
navigate to its [upstream repository](https://github.com/helloworlddan/tortune)
to figure out what it does.

Once you've understood what's going on, you can try to run that app directly in
Cloud Shell, running the following in the terminal:

```bash
go run main.go
```

This downloads dependencies, compiles and starts the web server. You can now use
Cloud Shell Web Preview on port 8080 to check out your application. You can find
the Web Preview <walkthrough-web-preview-icon></walkthrough-web-preview-icon> at
the top right in Cloud Shell. If you don't see a response from the web server,
try waiting a little longer for the Go compiler to build and start your server.

Finally, focus the terminal again and terminate the web server with `Ctrl-C`.

## Building and deploying the code

We established that our code runs locally, great! Let's put it on the cloud.

In the following steps, we'll deploy our app to Cloud Run, which is a compute
platform for running modern, cloud-first, containerized applications. Cloud Run
schedules and runs container images that expose services over HTTP. You can use
any programming language or framework, as long as you can bind an HTTP server on
a port provided by the `$PORT` environment variable. Please note that container
images must be in Docker or OCI format and will be run on Linux x86_64.

Cloud Run is available in all Google Cloud Platform regions globally. If you are
unsure where to deploy to, you can use the
[Region Picker](https://cloud.withgoogle.com/region-picker/) tool to find the
most suitable one.

You can configure your project and preferred regions and zones in `gcloud` so
its invocation becomes more convenient.

```bash
gcloud config set project <walkthrough-project-id/>
gcloud config set run/region europe-north1
gcloud config set artifacts/location europe-north1
```

Note that our code is not yet built or containerized, but Cloud Run requires
that. The `gcloud` CLI has a convenient shortcut for deploying Cloud Run which
quickly allows us to do so.

We can use a single command to easily:

- tarball a directory (build context)
- upload it to Cloud Build and use it to untar the context, build the app,
  containerize and store it on Artifact Registry
- create a new Cloud Run service with a fresh revision
- bind an IAM policy to the service
- route traffic to the new endpoint

```bash
gcloud run deploy jokes --source .
```

The command requires additional information and hence switches to an interactive
mode. You can have a look at the `gcloud` CLI
[reference](https://cloud.google.com/sdk/gcloud/reference/run/deploy) in case
some of the options may be unclear to you.

We want our application to be publicly available on the internet, so we should
allow unauthenticated invocations from any client.

Wait for the deployment to finish and then navigate to the `*.a.run.app`
endpoint it created. You should be able to call the endpoint from your browser
or by using any other HTTP client like cURL.

Next, let's use `gcloud` to retrieve the auto-generated service URL and then
`curl` it:

```bash
curl $(gcloud run services describe jokes --format 'value(status.url)')
```

Cloud Run services consist of one or more revisions. Whenever you update your
service or its configuration, you are creating a new revision. Revisions are
immutable.

Navigate to the
[Cloud Run section in the Google Cloud Console](https://console.cloud.google.com/run)
to explore the service and its active revision.

Building from `--source` uses Google Cloud Buildpacks. Buildpacks are predefined
archetypes that take your application source code and transform it into
production-ready container images without any additional boilerplate (e.g. a
Dockerfile).
[Learn more about supported languages for Buildpacks](https://cloud.google.com/run/docs/deploying-source-code).

## Using Cloud Code

Alternatively, you can deploy and otherwise manage the life cycle of your Cloud
Run services and other resources using Cloud Code. Cloud Code is an IDE plugin
for both VS Code and Intellj-based IDEs and is designed to make interaction with
Google Cloud more convenient.

Cloud Code is pre-installed in Cloud Shell.

Click the <walkthrough-editor-spotlight spotlightId="activity-bar-cloud-code">
Cloud Code icon in the activity bar </walkthrough-editor-spotlight> to expand
it.

Within the activity bar
<walkthrough-editor-spotlight spotlightId="cloud-code-cloud-run-deploy"> you can
directly deploy to Cloud Run </walkthrough-editor-spotlight> too, without the
need to use the CLI.

Take a moment and familiarize yourself with the wizard. You can also use this
wizard to build and deploy directly to Cloud Run.

## Scaling your app

Cloud Run automatically scales your application based on how many web requests
are coming in via the HTTPS endpoint. Cloud Run's horizontal auto-scaler is
extremely fast and can launch 100s of new instances in seconds.

Let's put some load on our newly created service and learn about scaling while
we wait. We'll start by pushing 50.000 requests using `hey`:

```bash
hey -n 50000 $(gcloud run services describe jokes --format 'value(status.url)')
```

In order to build modern, cloud-first applications that scale well horizontally,
we need to watch out for some design considerations.

**Applications should be engineered to boot quickly.** Cloud Run can start your
containers very quickly, but it is your responsibility to bring up a web server
and you should engineer your code to do so quickly. The earlier this happens,
the earlier Cloud Run is able to route HTTP requests to the new instance, reduce
stress on the older instances and hence scale out effectively. Cloud Run
considers the life cycle stage from starting your application to the moment it
can serve HTTP requests as 'Startup'. During this time Cloud Run will
periodically check if your application has bound the port provided by `$PORT`
and if so, Cloud Run considers startup complete and routes live traffic to the
new instance. Depending on your application code, you can further cut down
startup time by enabling _Startup CPU boost_. When enabled, Cloud Run will
temporarily allocate additional CPU resources during startup of your
application.

**Applications should ideally be stateless.** Cloud Run will also automatically
scale in again, should application traffic decrease. Container instances will be
terminated, if Cloud Run determines that too many are active to deal with the
current request load. When an instance is scheduled for termination, Cloud Run
will change its life cycle stage to 'Shutdown'. You can trap the `SIGTERM`
signal in your code and begin a graceful shutdown of your instance. Requests
will no longer be routed to your container and your application has 10 seconds
to persist data over the network, flush caches or complete some other remaining
write operations.

**Applications are generally request-driven**. During the 'Startup' and
'Shutdown' stages of each container life cycle, your application can expect to
be able to fully use the allocated CPU. During the 'Serving' life cycle, the CPU
is only available when there is at least one active request being processed on a
container instance. If there is no active request on the instance, Cloud Run
will throttle the CPU and use it elsewhere. You will not be charged for CPU time
if it's throttled. Occasionally, you might create applications that require a
CPU to be always available, for instance when running background services. In
this scenario, you want to switch from Cloud Run's default _CPU allocated during
requests_ to the alternative mode _CPU always allocated_. Note that this will
also switch Cloud Run to a
[different pricing model](https://cloud.google.com/run/pricing#tables). The
following diagram shows the two pricing models and their effect on how CPUs are
throttled throughout the life cycle of a container instance.

![Cloud Run container life cycle and CPU throttling](https://raw.githubusercontent.com/NucleusEngineering/serverless/main/.images/lifecycle.png)

Take a look at the deployment wizard for your service in the
[Cloud Run section in the Google Cloud Console](https://console.cloud.google.com/run/deploy/europe-north1/jokes)
and make sure **Startup CPU boost** is enabled.

While you are in the Google Cloud Console, have a look at
[the metrics section of your service](https://console.cloud.google.com/run/detail/europe-north1/jokes/metrics)
and explore how the previously executed load test affected scaling.

## Summary

You now have a good feel for what it means to build and deploy your code to
Cloud Run.

<walkthrough-conclusion-trophy></walkthrough-conclusion-trophy>

This completes Module 1. You can now wait for the live session to resume or
continue by yourself and on-demand.

## Module 2: Building and deploying container images with Cloud Build

![Tutorial header image](https://raw.githubusercontent.com/NucleusEngineering/serverless/main/.images/build.jpg)

In this tutorial we'll learn how to use Cloud Build, the fully-managed CI system
on Google Cloud. Instead of using Build Packs, we'll be using Dockerfiles and
Cloud Builds own declarative configuration to leverage higher flexibility and
control over how we build our images. Finally, we'll use Cloud Build to also
deploy to Cloud Run, so we can continuously deliver updates to our Cloud Run
services.

[Overview on Cloud Build](https://raw.githubusercontent.com/NucleusEngineering/serverless/main/.images/overview-build.jpg)

<walkthrough-tutorial-difficulty difficulty="3"></walkthrough-tutorial-difficulty>

Estimated time:
<walkthrough-tutorial-duration duration="35"></walkthrough-tutorial-duration>

To get started, click **Start**.

## Project setup

First, let's make sure we got the correct project selected. Go ahead and select
the provided project ID.

<walkthrough-project-setup billing="true"></walkthrough-project-setup>

<walkthrough-enable-apis apis="cloudbuild.googleapis.com,
run.googleapis.com,
artifactregistry.googleapis.com"> </walkthrough-enable-apis>

## Cloud Build

Cloud Build is a fully-managed CI system in Google Cloud. It is a general
purpose compute platform that allows you to execute all sorts of batch
workflows. However, it really excels at building code and creating container
images.

You don't need to provision anything to get started with using Cloud Build, it's
fully-managed: simply enable the API and submit your jobs. Cloud Build manages
all the required infrastructure for you. Per default, Cloud Build schedules
build jobs on shared infrastructure, but it can be configured to run on a
[dedicated worker pool](https://cloud.google.com/build/docs/private-pools/private-pools-overview)
that is not shared with other users.

In the previous section of our journey we saw how to use build and deploy
directly to Cloud Run from source code using the magic of Build Packs. Actually,
this deployment via `gcloud run deploy` already leverages Cloud Build in the
background, as you can see here in the
[Cloud Build History](https://console.cloud.google.com/cloud-build/builds;region=global).
These are great to get you started quickly. Almost as quickly, you will realize
that you need more control over your build process, so you will start writing
your own Dockerfiles. Let's see how that works with Cloud Build.

In order to get started let's set up some configuration:

```bash
gcloud config set project <walkthrough-project-id/>
gcloud config set run/region europe-north1
gcloud config set artifacts/location europe-north1
```

We are building some container images in this tutorial, so let's set up a
dedicated Docker image registry using Google Cloud Artifact Registry and
remember it in `gcloud config`:

```bash
gcloud artifacts repositories create my-repo --repository-format docker
gcloud config set artifacts/repository my-repo
```

Great! Let's get started.

## Building with Dockerfiles

Previously, we built our container images by specifying the `--source` flag,
which will cause Cloud Build to use
[Google Cloud Buildpacks](https://cloud.google.com/docs/buildpacks/overview) to
automatically determine the tool chain of our application. It then uses
templated container build instructions which are based on best practices to
build and containerize the application.

Sometimes it is required to exert more control over the Docker build process.
Google Cloud Build also understands Dockerfiles. With Dockerfiles it is possible
to precisely specify which instructions to execute in order to create the
container image. Cloud Build typically inspects the build context and looks for
build instructions in the following order:

1. cloudbuild.yaml: If this file is present, Cloud Build will execute the
   instructions defined in it. We'll learn about cloudbuild.yaml files later in
   this tutorial.
2. Dockerfile: If this file is present, Cloud Build will use it to execute a
   container build process that is equivalent to running `docker build`.
   Finally, it will tag and store the resulting image.
3. Buildpacks: If no explicit build instructions are available, Cloud Build can
   still be used to integrate by using Google Cloud Buildpacks.

In the <walkthrough-editor-spotlight spotlightId="activity-bar-explorer"> File
Explorer </walkthrough-editor-spotlight> create a new file called `Dockerfile`
and place it in the same directory as the Go application. Let's start with a
simple, single-stage build for Go:

```Dockerfile
FROM golang:1.22-bookworm
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY main.go ./
RUN CGO_ENABLED=0 go build -v -o server
CMD ["/app/server"]
```

This Dockerfile:

1. uses Debian bookworm with Golang 1.22 as base image
2. copies the definition of the Go module and installs all dependencies
3. copies the sources and builds the Go application
4. specifies the application binary to run

We can now submit the build context to Cloud Build and use the instructions in
the Dockerfile to create a new image by running:

```bash
LOCATION="$(gcloud config get-value artifacts/location)"
PROJECT="$(gcloud config get-value project)"
gcloud builds submit . \
  --tag ${LOCATION}-docker.pkg.dev/${PROJECT}/my-repo/dockerbuild
```

Next, let's navigate first to the
[Cloud Build History](https://console.cloud.google.com/cloud-build/builds;region=global)
to see the build we just started. As soon as that is finished we go to
[Artifact Registry in the Google cloud Console](https://console.cloud.google.com/artifacts/docker?project=<walkthrough-project-id/>)
to find the repository and inspect the image.

Huh, it seems like our image is quite big! We can fix this by running a
[multi-stage Docker build](https://docs.docker.com/build/building/multi-stage/).
Let's extend the Dockerfile and replace its contents with the following:

```Dockerfile
FROM golang:1.22-bookworm as builder
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

Great! We'll keep the first stage as the build stage. Once the statically-linked
Go binary is built, it is copied to a fresh stage that is based on
`gcr.io/distroless/static-debian11`. The _distroless_ images are Google-provided
base images that are very small. That means there is less to store and images
typically have much smaller attack surfaces and improved start-up time. Let's
build again:

```bash
LOCATION="$(gcloud config get-value artifacts/location)"
PROJECT="$(gcloud config get-value project)"
gcloud builds submit . \
  --tag ${LOCATION}-docker.pkg.dev/${PROJECT}/my-repo/dockermultistage
```

Navigate back to Artifact Registry in the Google Cloud Console and inspect the
new image `dockermultistage`. The resulting image is much (70x) smaller.

## Building with a cloudbuild.yaml file

Cloud Build will build, tag and store your container image when you submit a
Dockerfile, but it can do a lot more. You can instruct Cloud Build to execute
any arbitrary sequence of instructions by specifying and submitting a
`cloudbuild.yaml` file. When detected in the build context, this file will take
precedence over any other configuration and will ultimately define what your
Cloud Build task does.

In a nutshell, you are able to specify a sequence of steps with each step being
executed in a container. You can specify which image and process to execute for
every individual step. The steps can collaborate and share data via a shared
volume that is automatically mounted into the file system at `/workspace`.

Go ahead and create a `cloudbuild.yaml` and place it in the same directory as
the Dockerfile and the Go application. Add the following contents:

```yaml
steps:
  - name: "ubuntu"
    args:
      - "echo"
      - "hi there"
```

This configuration asks Cloud Build to run a single step: Use the container
image provided by [ubuntu](https://hub.docker.com/_/ubuntu) and run `echo`.

While this might not be a terribly useful or sophisticated CI pipeline, it is a
useful illustration to understand how versatile Cloud Build is. Let's run it by
submitting a new Cloud Build task:

```bash
gcloud builds submit
```

OK. Let's modify the `cloudbuild.yaml` to do something that's actually useful.
Replace the file contents with the following:

```yaml
steps:
  - name: "gcr.io/cloud-builders/docker"
    args:
      - "build"
      - "-t"
      - "$_ARTIFACT_REGION-docker.pkg.dev/$PROJECT_ID/my-repo/dockercloudbuildyaml"
      - "."
  - name: "gcr.io/cloud-builders/docker"
    args:
      - "push"
      - "$_ARTIFACT_REGION-docker.pkg.dev/$PROJECT_ID/my-repo/dockercloudbuildyaml"
```

Pay attention to the `$`-prefixed `$PROJECT_ID` and `$_ARTIFACT_REGION`.
`$PROJECT_ID` is a pseudo-variable which is automatically set by the Cloud Build
environment. `$_ARTIFACT_REGION` is what Cloud Build calls a _substitution_ and
it allows users to override otherwise static content of the build definition by
injecting configuration during invocation of the build. Let's take a look:

```bash
LOCATION="$(gcloud config get-value artifacts/location)"
gcloud builds submit \
  --substitutions _ARTIFACT_REGION=${LOCATION}
```

This will build, tag and store a container image.

To fully understand what else can go into your `cloudbuild.yaml`, please check
out the
[schema documentation](https://cloud.google.com/build/docs/build-config-file-schema).

## Deploying to Cloud Run

Cloud Build is a versatile tool and is suited to run a wide variety of
batch-like jobs. Until now, we only used Cloud Build to accomplish Continuous
Integration (CI) tasks, but we don't need to stop there.

We can extend the `cloudbuild.yaml` definition and automatically deploy the
newly created image to Cloud Run like this:

```yaml
steps:
  - name: "gcr.io/cloud-builders/docker"
    args:
      - "build"
      - "-t"
      - "$_ARTIFACT_REGION-docker.pkg.dev/$PROJECT_ID/my-repo/dockercloudbuildyaml"
      - "."
  - name: "gcr.io/cloud-builders/docker"
    args:
      - "push"
      - "$_ARTIFACT_REGION-docker.pkg.dev/$PROJECT_ID/my-repo/dockercloudbuildyaml"
  - name: "gcr.io/cloud-builders/gcloud"
    args:
      - "run"
      - "deploy"
      - "--region"
      - "$_RUN_REGION"
      - "--image"
      - "$_ARTIFACT_REGION-docker.pkg.dev/$PROJECT_ID/my-repo/dockercloudbuildyaml"
      - "jokes"
```

In order for this to work, we need to grant some permissions to the service
account used by the Cloud Build service agent, so it can interact with Cloud Run
resources on our behalf.

Navigate to the
[settings section of Cloud Build in the Google Cloud Console](https://console.cloud.google.com/cloud-build/settings/service-account)
and assign both the **Cloud Run Admin** and the **Service Account User** roles
to Cloud Build's service account.

Finally, let's supply all the substitutions and invoke Cloud Build once again to
execute a fresh build and deploy to Cloud Run automatically.

```bash
LOCATION="$(gcloud config get-value artifacts/location)"
REGION="$(gcloud config get-value run/region)"
gcloud builds submit \
  --substitutions _ARTIFACT_REGION=${LOCATION},_RUN_REGION=${REGION}
```

## Setting up an automatic source triggers from GitHub

Now that we have successfully created a full build definition on Cloud Build for
a simple yet very useful CI/CD pipeline, let's have a look at creating automatic
build triggers that execute the pipeline whenever code gets updated on a remote
Git repository. In order to put everything in place, we'll now look at how to
authenticate with GitHub, fork the code into a new repository, authorize GitHub
to connect with Cloud Build and create an automatic build trigger on Cloud
Build.

<walkthrough-info-message>This part of the tutorial is **optional**, as it
requires you to have a GitHub account. If you don't have one you can
[register a new one](https://github.com/signup).</walkthrough-info-message>

### Authenticating with GitHub

First, let's authenticate using `gh`, GitHub's CLI, which comes pre-installed on
Cloud Shell. In your terminal begin the authentication handshake and complete
the CLI wizard.

We'd like to connect in the following matter:

- Connect to `github.com`
- Use HTTPS for authentication
- Configure `git` to use GitHub credentials

Begin the authentication exchange by running the following:

```bash
gh auth login
```

Copy the authentication code, navigate to
[GitHub's device authentication page](https://github.com/login/device), paste
the token and authorize the application.

Once authentication is complete, you should be able to test access to GitHub by
listing all of your repositories, like so:

```bash
gh repo list
```

### Creating a fork of the upstream repo

Next, we'll create a fork of the repository we are currently working with, so
that we have a separate, remote repository that we can push changes to. Create
the fork by running the following:

```bash
gh repo fork
```

This will create a fork of the upstream repository from
`github.com/nucleusengineering/serverless` and place it at
`github.com/YOUR_GITHUB_HANDLE/serverless`. Additionally, the command
automatically reconfigures the remotes of the local repository clone. Take a
look the the configured remotes and see how it configured both `upstream` and
the new `origin` by running the following:

```bash
git remote -v
```

<walkthrough-info-message>Make sure that your new `origin` points to the work
belonging to your own user and push commits to the new `origin` from now
on.</walkthrough-info-message>

### Setting up the Cloud Build triggers

Next, Cloud Build needs to be configured with a trigger resource. The trigger
contains everything Cloud Build needs to automatically run new builds whenever
the remote repository gets updated.

- Navigate to the
  [Cloud Build triggers section of the Google Cloud Console](https://console.cloud.google.com/cloud-build/triggers)
  and click on 'Connect Repository'.
- Follow the wizard on the right to connect to github.com, authenticate Google's
  GitHub app.
- Filter repositories based on your user name.
- Find the forked repository called 'serverless'.
- Tick the box to accept the policy and connect your repository.

Once completed you should see a connection confirmation message displayed. Now
it's time to create the trigger.

- Hit 'Create Trigger' and create a new trigger.
- In the wizard, specify that the trigger should read configuration from the
  provided `./cloudbuild.yaml` and **add all the substitutions** you used
  previously to trigger your build.

<walkthrough-info-message>When using the "Autodetect" configuration option,
there is no possibility to add substitution variables through the UI. So make
sure to specify the "Cloud Build configuration file (yaml or json)" option
explicitly, and then continue to fill in the substitution
variables.</walkthrough-info-message>

### Pushing some changes

We should now have everything in place to automatically trigger a new build
whenever changes are pushed to `main` on the remote repository fork.

If you haven't done so already, you will need to configure git on your Cloud
Shell. To do so, run the following and configure your email and name, so git
know who you are.

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

Finally, push the commit to the remote repository:

```bash
git push origin main
```

Changes should automatically be detected and trigger a new Cloud Build task.
Navigate to the
[Cloud Build dashboard](https://console.cloud.google.com/cloud-build/dashboard)
and explore the running build.

## Summary

You now know how Cloud Build can help you automate integrating your artifacts.

<walkthrough-conclusion-trophy></walkthrough-conclusion-trophy>

This completes Module 2. You can now wait for the live session to resume or
continue by yourself and on-demand.

## Module 3: Extend your code to call Cloud APIs

![Tutorial header image](https://raw.githubusercontent.com/NucleusEngineering/serverless/main/.images/apis.jpg)

In this tutorial we'll learn how to extend the existing code to call Cloud APIs
directly. Currently, the deployed application uses a library which contains a
static set of jokes. Whenever the library is used it randomly selects a joke and
returns it. After a while we will surely start to see the same jokes again and
the only way to see new jokes is when a human would actually implement them in
the library.

Luckily, there is a thing called _generative AI_ now. Google Cloud Vertex AI
contains the Google-built, pre-trained, Gemini 1.5 Flash model which is a
general-purpose, multimodal, speed-optimized, generative large language model
(LLM) that can be queried with text prompts in natural language to generate all
sorts of outputs, including text. In this tutorial we'll implement the
`model:predict` endpoint of Vertex AI to execute this model in order to add new
dad jokes in a generative manner.

<!-- TODO: include on-demand clip for module 3 -->

Additionally, we'll learn a little bit about custom service accounts, IAM
permissions and how to use the principle of least privilege to secure our
services on Cloud Run.

<walkthrough-tutorial-difficulty difficulty="3"></walkthrough-tutorial-difficulty>

Estimated time:
<walkthrough-tutorial-duration duration="45"></walkthrough-tutorial-duration>

To get started, click **Start**.

## Project setup

First, let's make sure we got the correct project selected. Go ahead and select
the provided project ID.

<walkthrough-project-setup billing="true"></walkthrough-project-setup>

Run the following to make sure all required APIs are enabled. Note that
`aiplatform.googleapis.com` is added.

<walkthrough-enable-apis apis="cloudbuild.googleapis.com,
run.googleapis.com,aiplatform.googleapis.com,
artifactregistry.googleapis.com"> </walkthrough-enable-apis>

## About the different ways to call Google APIs

In the cloud we often need to implement API calls to interact with other
services. Platforms like Google Cloud are a large collection of APIs; if we
learn how to use them we can use the entire wealth of Google Cloud's
functionality in our applications and build almost anything!

There are typically three different ways to interact with Google APIs
programmatically and we should choose them in the following order:

1. **Cloud Client Libraries**: These are the recommended options. Cloud Client
   Libraries are SDK that you can use in a language-native, idiomatic style.
   They give you a high-level interface to the most important operation and
   allow you to quickly and comfortably get the job done. An example in the Go
   ecosystem would be the `cloud.google.com/go/storage` package, which
   implements the most commonly used operations of Google Cloud Storage. Have a
   look at the
   [documentation of the package](https://pkg.go.dev/cloud.google.com/go/storage)
   and see how it respects and implements native language concepts like the
   `io.Writer` interface of the Go programming language.

2. **Google API Client Libraries**: Should no Cloud Client Library be available
   for what you are trying to accomplish you can fall back to using a Google API
   Client Library. These libraries are auto-generated and should be available
   for almost all Google APIs.

3. **Direct Implementation**: You can always choose to implement exposed APIs
   directly with your own client code. While most APIs are available as
   traditional RESTful services communicating JSON-serialized data, some APIs
   also implement [gRPC](https://grpc.io).

Have a look at
[this page of the Google Cloud documentation](https://cloud.google.com/apis/docs/client-libraries-explained)
to learn more about the differences between the available libraries.

Let's configure `gcloud` for the project and default regions:

```bash
gcloud config set project <walkthrough-project-id/>
gcloud config set run/region europe-north1
gcloud config set artifacts/location europe-north1
```

Okay, all set!

## Extending the code

In order to be able to replace the statically created jokes with jokes generated
by Vertex AI, the code needs to be extended in the following ways:

1. **Create a client to execute the remote model**: The first step is to safely
   instantiate the correct client type, that we are going to use later on to
   interact with the API. This type holds all the configuration for endpoints
   and handles authentication, too. Have a look at the
   [documentation of the package](https://pkg.go.dev/cloud.google.com/go/vertexai/genai#NewClient)
   for `genai.Client`. The client will automatically look for credentials
   according to the rules of
   [Google's Application Default Credentials scheme](https://cloud.google.com/docs/authentication/application-default-credentials).

2. **Execute the call against the API to run the model prediction**: Next, we'll
   use the previously instantiated client to actually invoke the remote model
   with `[]genai.Parts` containing the prompt and execute the API by calling
   `genai.GenerateContent`
   [as described in the docs](https://pkg.go.dev/cloud.google.com/go/vertexai/genai#GenerativeModel.GenerateContent).
   Take a look at the
   [documentation of the package](https://pkg.go.dev/cloud.google.com/go/vertexai/genai#GenerativeModel)
   to see what else can be done with `genai.GenerativeModel`.

Now, it's time to make some changes to the code.

<walkthrough-info-message>You may now attempt to **implement the above code
changes yourself** for which you should have a good understanding of the Go
programming language. If you choose to do so, you should stop reading now and
give it your best shot.</walkthrough-info-message>

Alternatively, you can **use a prebuilt library** to accomplish the same. If
that's more your cup of tea, go hit 'Next' and proceed to the next section.

Have a look at the
[Vertex AI Model Garden](https://console.cloud.google.com/vertex-ai/publishers/google/model-garden/gemini-1.5-flash-preview-0514)
to learn more about Gemini Flash and other available models.

## Using a library to call the Gemini Flash model

To get started, first have a look at the Go module
`github.com/helloworlddan/tortuneai`. The module provides the package
`github.com/helloworlddan/tortuneai/tortuneai` which implements the Cloud Client
Library for Vertex AI to call the _Gemini Flash_ model.
[Read through the code](https://github.com/helloworlddan/tortuneai/blob/main/tortuneai/tortuneai.go)
of `tortuneai.HitMe` and see how it implements the aforementioned steps to
interact with the API.

In order to use the package we need to first get the module like this:

```bash
go get github.com/helloworlddan/tortuneai@v0.0.4
```

Once that is completed, we can update
<walkthrough-editor-open-file filePath="cloudshell_open/serverless/main.go"> the
main application source file main.go </walkthrough-editor-open-file> and change
all references to the previously used package `tortune` with `tortuneai`.

Notice that the signature for `tortuneai.HitMe()` is different from the previous
`tortune.HitMe()`. While the original function did not require any parameters,
you are required to pass two `string` values into the new one: One with an
actual text prompt for the Gemini Flash model and one with your Google Cloud
project ID. Additionally, the function now returns multiple return values: a
`string` containing the response from the API and an `error`. If everything goes
well, the error will be `nil`, if not it will contain information about what
went wrong.

Here is a possible implementation:

```golang
joke, err := tortuneai.HitMe("", "<walkthrough-project-id/>")
if err != nil {
    fmt.Fprintf(w, "error: %v\n", err)
    return
}
fmt.Fprint(w, joke)
```

Update the implementation of `http.HandleFunc()` in
<walkthrough-editor-open-file filePath="cloudshell_open/serverless/main.go"> the
main application source file main.go </walkthrough-editor-open-file> with the
code snippet.

Let's check if the modified code compiles by running it:

```bash
go run main.go
```

This recompiles and starts the web server. Let's check the application with the
Web Preview <walkthrough-web-preview-icon></walkthrough-web-preview-icon> at the
top right in Cloud Shell and see if we can successfully interact with the Gemini
Flash model.

If you are satisfied you can focus the terminal again and terminate the web
server with `Ctrl-C`.

It's good practice to clean up old dependencies from the `go.mod` file. You can
do this automatically my running:

```bash
go mod tidy
```

If you like you can stay at this point for a moment, change the prompt (the
first argument to `tortuneai.HitMe()`), re-run with `go run main.go` and use the
Web Preview <walkthrough-web-preview-icon></walkthrough-web-preview-icon> to
have a look at how the change in prompt affected the model's output.

## Creating a custom service account for the Cloud Run service

Good! The code changes we made seem to work, now it's time to deploy the changes
to the cloud.

When running the code from Cloud Shell, the underlying implementation used
Google
[Application Default Credentials (ADC)](https://cloud.google.com/docs/authentication/provide-credentials-adc)
to find credentials. In this case it was using the credentials of the Cloud
Shell user identity (yours).ADC to find credentials. In this case it was using
the c redentials of the Cloud Shell user identity (yours).

Cloud Run can be configured to use a service account, which exposes credentials
to the code running in your container. Your application can then make
authenticated requests against Google APIs.

Per default, Cloud Run uses the
[Compute Engine default service account](https://cloud.google.com/compute/docs/access/service-accounts#default_service_account).
This service account has wide permissions and should generally be replaced by a
service account with the least amount of permissions required to do whatever
your service needs to do. The Compute Engine default service has a lot of
permissions, but it does not have the required permissions to execute
[the API call](https://cloud.google.com/vertex-ai/docs/reference/rest/v1beta1/projects.locations.publishers.models/predict).

When it comes to identifying the correct IAM roles to attach to identities
[this reference page on the IAM documentation](https://cloud.google.com/iam/docs/understanding-roles)
is an extremely useful resource. On the page you can check the section for
[Vertex AI roles](https://cloud.google.com/iam/docs/understanding-roles#vertex-ai-roles)
and learn that
[_Vertex AI User_](https://cloud.google.com/iam/docs/understanding-roles#aiplatform.user)
(`roles/aiplatform.user`) is a suitable role for our Cloud Run service, because
this role contains the permission `aiplatform.endpoints.predict`. If you are
unsure which permission you require, you can always check the
[API reference for the required operation](https://cloud.google.com/vertex-ai/docs/reference/rest/v1beta1/projects.locations.publishers.models/predict).

<walkthrough-info-message>**Note**: You could argue that the role _Vertex AI
User_ has too many permissions for the Cloud Run service and a
security-conscious person would probably agree with you. If you really wanted to
make sure that the Cloud Run service only had least amount of privilege to
execute the absolutely required permissions, you would have to create a
[IAM custom role](https://cloud.google.com/iam/docs/creating-custom-roles) to
achieve this.</walkthrough-info-message>

For now, we'll stick with _Vertex AI User_.

Next, let's create a brand new customer IAM service account like this:

```bash
gcloud iam service-accounts create tortune
```

We can bind the identified role to the various resource levels. For the sake of
simplicity, let's attach it at the project level by executing:

```bash
PROJECT=$(gcloud config get-value project)
ACCOUNT=tortune@${PROJECT}.iam.gserviceaccount.com
gcloud projects add-iam-policy-binding ${PROJECT} \
    --member serviceAccount:${ACCOUNT} \
    --role "roles/aiplatform.user"
```

The service account will now be able to use all the permissions in _Vertex AI
User_ on all resources in our current project. Finally, we need to deploy a new
Cloud Run revision by updating the service configuration so that our Cloud Run
service will use the newly-created service account:

```bash
PROJECT=$(gcloud config get-value project)
ACCOUNT=tortune@${PROJECT}.iam.gserviceaccount.com
gcloud run services update jokes \
    --service-account ${ACCOUNT}
```

Now, that all IAM resources and configurations are in place, we can trigger a
new Cloud Build execution to deploy changes in a CI/CD fashion, like this:

```bash
LOCATION="$(gcloud config get-value artifacts/location)"
REGION="$(gcloud config get-value run/region)"
gcloud builds submit \
  --substitutions _ARTIFACT_REGION=${LOCATION},_RUN_REGION=${REGION}
```

If you have previously configured Github, you can achieve the same by committing
and pushing your changes, like this:

```bash
git add .
git commit -m 'upgrade to Gemini Flash for generative jokes'
git push origin main
```

Navigate to
[Cloud Build's dashboard](https://console.cloud.google.com/cloud-build/dashboard)
and click into the active build to monitor it's progress.

Once completed, you should be able to get fresh generated content by cURLing the
endpoint of the Cloud Run Service:

```bash
curl $(gcloud run services describe jokes --format 'value(status.url)')
```

Amazing!

## Summary

You now know how to call Google APIs directly from your code and understand how
to secure your services with least-privilege service accounts.

<walkthrough-conclusion-trophy></walkthrough-conclusion-trophy>

This completes Module 3. You can now wait for the live session to resume or
continue by yourself and on-demand.

## Module 4: Operate your services in production

![Tutorial header image](https://raw.githubusercontent.com/NucleusEngineering/serverless/main/.images/sre.jpg)

In this part of our journey, we are going to look at some basic principles of
how to operate your services in production. We'll learn about
[Site Reliability Engineering (SRE)](https://sre.google/sre-book/introduction/)
and define some custom Service Level Objectives (SLOs) to keep track of the
health of our services.

Check out this
[Overview on SLOs, SLIs and SLAs](https://www.youtube.com/watch?v=tEylFyxbDLE)

Before you continue let's load your CloudRun service with traffic to ensure that
historic metrics data is available for later exercises. Use
[`hey`](https://github.com/rakyll/hey) again to do this, which will send 180
requests per minute for the next 15 minutes to generate a base load without
triggering
[quota limits](https://cloud.google.com/vertex-ai/generative-ai/docs/quotas) of
the Gemini API:

```bash
hey -z 15m -c 3 -q 1 $(gcloud run services describe jokes --format 'value(status.url)')
```

At first we'll define what reliability means for the service we operate. Using
these definitions, we can then proceed and configure service monitoring to keep
track of the health of our service over time.

After that we'll have a look at Cloud Run's traffic splitting capabilities. This
allows us to deploy new revisions of a Cloud Run service and then gradually move
portion of production traffic over. Cloud Run's programmable network control
plane allows you to split traffic between revisions; you can use this built-in
feature to implement strategies like blue/green deployments, canary releases, or
rollback to previous revisions in seconds should you ever push a bad release.

To bring everything to life, we'll be deploying a faulty version as a canary
release with some amount of live traffic, get an alert, observe how this will
burn through the
[error budget](https://cloud.google.com/blog/products/management-tools/sre-error-budgets-and-maintenance-windows)
and quickly rollback to a safe state.

Check out this
[Overview on Risks and Error Budgets](https://www.youtube.com/watch?v=y2ILKr8kCJU)

<walkthrough-tutorial-difficulty difficulty="3"></walkthrough-tutorial-difficulty>

Estimated time:
<walkthrough-tutorial-duration duration="45"></walkthrough-tutorial-duration>

To get started, click **Start**.

## Project setup

First, let's make sure we got the correct project selected. Go ahead and select
the provided project ID.

<walkthrough-project-setup billing="true"></walkthrough-project-setup>

Run the following to make sure all required APIs are enabled.

<walkthrough-enable-apis apis="cloudbuild.googleapis.com,
run.googleapis.com,aiplatform.googleapis.com,
artifactregistry.googleapis.com"> </walkthrough-enable-apis>

## SRE: Hope is not a strategy

Google's Site Reliability Engineering is what you get when you treat operations
as if it’s a software problem. The mission is to protect, provide for, and
progress the software and systems behind all of Google’s public services —
Google Search, Ads, Gmail, Android, YouTube and many Cloud Services, to name
just a few — with an ever-watchful eye on their availability, latency,
performance, and capacity.

Many of the practices and tools of SRE have been integrated into Google Cloud
Platform so that everyone running services in production can benefit from what
Google learned over decades.

SRE begins with the idea that a prerequisite to success is reliability. A system
that is reliable cannot perform its function and will fail by default.
Reliability can be measured in many dimensions like latency, freshness,
correctness etc. Availability, in SRE terms, is one of the most important
measurements, it defines whether a system is able to fulfill its intended
function at a point in time. In addition to being used as a reporting tool, the
historical availability measurement can also describe the probability that your
system will perform as expected in the future.

When we set out to define the terms of SRE, we wanted to set a precise numerical
target for system availability. We term this target the availability SLO of our
system. Any discussion we have in the future about whether the system is running
sufficiently reliably and what design or architectural changes we should make to
it must be framed in terms of our system continuing to meet this SLO.

e also have a direct measurement of a service’s behavior: the proportion of
successful responses of our system. This is a Service-Level Indicator (SLI).
When we evaluate whether our system has been running within SLO for the past
week, we look at the SLI to get the service availability percentage. If it goes
below the specified SLO, we have a problem and may need to make the system more
available in some way, such as rolling back to a previous (good) revision or
running a second instance of the service in a different city and load-balancing
between the two. If you want to know how reliable your service is, you must be
able to measure the rates of successful and unsuccessful queries as your SLIs.

### The four golden signals

SRE defines _four golden signals_ as a good starting point for measuring
relevant metrics of your systems.

**Latency**: The time it takes to service a request. It’s important to
distinguish between the latency of successful requests and the latency of failed
requests. For example, an HTTP 500 error triggered due to loss of connection to
a database or other critical backend might be served very quickly; however, as
an HTTP 500 error indicates a failed request, factoring 500s into your overall
latency might result in misleading calculations. On the other hand, a slow error
is even worse than a fast error! Therefore, it’s important to track error
latency, as opposed to just filtering out errors.

**Traffic**: A measure of how much demand is being placed on your system,
measured in a high-level system-specific metric. For a web service, this
measurement is usually HTTP requests per second, perhaps broken out by the
nature of the requests (e.g., static versus dynamic content). For an audio
streaming system, this measurement might focus on network I/O rate or concurrent
sessions. For a key-value storage system, this measurement might be transactions
and retrievals per second.

**Errors**: The rate of requests that fail, either explicitly (e.g., HTTP 500s),
implicitly (for example, an HTTP 200 success response, but coupled with the
wrong content), or by policy (for example, "If you committed to one-second
response times, any request over one second is an error"). Where protocol
response codes are insufficient to express all failure conditions, secondary
(internal) protocols may be necessary to track partial failure modes. Monitoring
these cases can be drastically different: catching HTTP 500s at your load
balancer can do a decent job of catching all completely failed requests, while
only end-to-end system tests can detect that you’re serving the wrong content.

**Saturation**: How "full" your service is. A measure of your system fraction,
emphasizing the resources that are most constrained (e.g., in a
memory-constrained system, show memory; in an I/O-constrained system, show I/O).
Note that many systems degrade in performance before they achieve 100%
utilization, so having a utilization target is essential.

In complex systems, saturation can be supplemented with higher-level load
measurement: can your service properly handle double the traffic, handle only
10% more traffic, or handle even less traffic than it currently receives? For
very simple services that have no parameters that alter the complexity of the
request (e.g., "Give me a nonce" or "I need a globally unique monotonic
integer") that rarely change configuration, a static value from a load test
might be adequate. As discussed in the previous paragraph, however, most
services need to use indirect signals like CPU utilization or network bandwidth
that have a known upper bound. Latency increases are often a leading indicator
of saturation. Measuring your 99th percentile response time over some small
window (e.g., one minute) can give a very early signal of saturation.

Finally, saturation is also concerned with predictions of impending saturation,
such as "It looks like your database will fill its hard drive in 4 hours."

You can
[read more about the four golden signals](https://sre.google/sre-book/monitoring-distributed-systems/#xref_monitoring_golden-signals).

A more precise and nuanced approach for measuring reliability however are SLIs
and their SLOs. Let's define one!

## Defining SLOs and alerts

Let's begin by defining an SLO for our service that keeps track of its
availability, specifically of its error rate. A service that is unable to serve
requests because requests are errors, are unavailable. We are about to define an
objective that says: We aim at serving 99% of all requests without any error.
Cloud Run already keeps track of application errors. You can see the log-based
metric _Request count_ in the metrics section of your Cloud Run service. You can
find it by clicking
[into your service in the Cloud Run section](https://console.cloud.google.com/run)
of the Google Cloud Console. The metrics counts HTTP status codes and `5xx`
status codes are considered to be server-side errors.

Right next to the _Metrics_ tab, you'll find the tab that says _SLOs_. Get in
there and start the new SLO creation wizard. Define an SLO, that:

- uses a **Request-based, availability metric** as SLI
- this should track the metric `run.googleapis.com/request_count`
- Observe the definition for **goodness** being response codes 200-299 while bad
  requests are constituted from the 500-599 range.
- the SLO definition should look at a **rolling 1 day period** with and
  **objective of 99%**
- This means that your error budget is 1% of all requests or 14 minutes per day
  if your traffic is equally distributed throughout the day.

Finally, also create an alert that will notify you in case a
[fast-burn](https://cloud.google.com/stackdriver/docs/solutions/slo-monitoring/alerting-on-budget-burn-rate#burn-rates)
of yourerror budget gets detected. Create an new alert, that:

- uses a **look back duration of 5 minutes**
- breaches with a **fast-burn threshold 10(x baseline)**
- creates a **new alert notification channel using your email** in Cloud
  Monitoring

Done. You should now be immediately notified as soon as your service is
responding with too many errors in a given moment and you're about to burn
through your error budget quickly.

## Deploying a canary with traffic splitting

Before continuing, make sure `gcloud` configuration is set to our preferences:

```bash
gcloud config set project <walkthrough-project-id/>
gcloud config set run/region europe-north1
gcloud config set artifacts/location europe-north1
```

Okay! Good to go.

Cloud Run comes with a built-in traffic control plane, which lets operators
programmatically assign traffic to individual revisions of the same service.
This allows you to deploy new changes to your service with none or only
single-digit percentage traffic routed to them. Once you've slowly gained
confidence in your new revision you can gradually increase traffic until you
reach 100% of traffic on the new revision and your rollout is complete. This
strategy is commonly referred to as a
[canary release](https://cloud.google.com/deploy/docs/deployment-strategies/canary).

<walkthrough-info-message>Cloud Run's default deployment strategy is to
automatically route all traffic to the new revision if it should pass minimum
health checks.</walkthrough-info-message>

Let's begin by deploying an image we know works well, in this case the previous
version that uses the static tortune library.

```bash
gcloud run deploy jokes \
    --image europe-north1-docker.pkg.dev/<walkthrough-project-id/>/my-repo/dockermultistage
```

In Cloud Run, you can tag your individual revisions. Each tagged revision can be
easily referenced to by its tag and also gets an individual endpoint, just for
itself. That means we can preview the revision even without assigning production
traffic to it.

Let's deploy a new revision of the jokes service using a different container
image. We've been told that requirements have changed dramatically for this new
version. We don't trust the image and decide to implement a canary release
strategy. We begin by marking the current revision as 'good'.

```bash
gcloud run services update jokes \
    --tag good
```

Good! The revision is tagged and gets a separate URL that looks something like
`https://good---jokes-xxxxxxxxxx-xx.a.run.app`. As long as we keep the tag
around we can always directly reach this revision with this endpoint. You can
try it by directly cURLing it:

```bash
curl $(gcloud run services describe jokes --format 'value(status.traffic[0].url)')
```

Now, let's deploy the new image, tag it and **route no traffic to it**, like so:

```bash
gcloud run deploy jokes \
    --image gcr.io/serverless-journey/failing \
    --tag next \
    --no-traffic
```

Remember: we don't trust this new container image. So, let's begin by routing
only 10% of traffic to the new revision tagged as `next`, by running the
following:

```bash
gcloud run services update-traffic jokes \
    --to-tags next=10
```

Cloud Run will quickly reconfigure the traffic pattern and print out the current
pattern. It's time to simulate some traffic to see how the canary behaves under
traffic. Use `hey` again to do this:

```bash
hey -n 500000 $(gcloud run services describe jokes --format 'value(status.url)')
```

Eventually, the trajectory with which we will be burning though the error budget
will cause an alert to fire.

<walkthrough-info-message>Now is probably a good time to catch a little break,
get a snack, get hydrated, maybe take a stroll through the park. Don't forget:
you are still on-call!</walkthrough-info-message>

Proceed to the next section, once you've received the alert email.

## Responding to the incident

Oh, no! We have an active incident!

Follow the link in the alert mail or navigate to the
[incidents section of Cloud Monitoring](https://console.cloud.google.com/monitoring/alerting/incidents).
You should see the active incident marked in red: First "Acknowledge Incident"
then go ahead and explore it. Notice how the rate at which we are burning error
budget is way past the threshold.

Also look at the SLOs section of your Cloud Run service in the
[Cloud Run section of the Google Cloud Console](https://console.cloud.google.com/run).
The page tells you that the percentage of healthy, non-error responses from your
service has dropped to 90%, which is below your SLO target of 99%!

Let's investigate a bit further. Head to the _Logs_ tab of the service. Next to
your application logs, Cloud Run itself logs system events here. You can see how
every invocation of your service produces an event. Most lines are marked with a
blue `i` (Loglevel `info`). You see that these are the requests to which the
service responded with `HTTP 200 OK`. About 10% of the log lines are marked with
a red `!!` (Loglevel `error`), though!

Have a closer look at one of each log line and expand into the fields `resource`
and further down into `labels`. The structured log message tells you exactly
which revision of the service belongs to the message. Turns out that all the
errors are caused by the revision we previously marked as `next`.

Now, before we tell the engineering team about the problem, we first keep calm
and remember Cloud Run traffic management and mitigate the incident. To bring
the system back to a safe state, let's rollback to the previous revision, the
one we tagged with `good`. To do so, simply reconfigured the traffic pattern on
the service like this:

```bash
gcloud run services update-traffic jokes \
    --to-tags good=100
```

Good! Run a bit more simulated traffic over the service to check if this
remediated the issue.

```bash
hey -n 50000 $(gcloud run services describe jokes --format 'value(status.url)')
```

Go back to the _SLOs_ section of your service and observe how the SLI for your
error rate climbs back over 99%. Crisis averted!

## Summary

Great! You've learned the basic of SRE, defined meaningful SLOs, mastered Cloud
Run traffic management and responded to a failing canary release by rolling back
to a known good state of the system.

<walkthrough-conclusion-trophy></walkthrough-conclusion-trophy>

This completes Module 4. You can now wait for the live session to resume or
continue by yourself and on-demand.

## Module 5: Securing the delivery process of our software artifacts

![Tutorial header image](https://raw.githubusercontent.com/NucleusEngineering/serverless/main/.images/secure.jpg)

Our application is already running in a scalable fashion adapting to user
demand. It is automatically built and we can dynamically deploy new versions
after safely testing them in the same environment with full parity. Cloud Run's
nimble programmable network plane allows us to gradually phase in the new
version or roll back in case our monitoring setup reports spikes in errors
caused by the introduced changes.

Although we have limited the permissions on the service account in use by the
Cloud Run service itself, we have almost completely ignored an emerging security
threat: attacks on our software supply chain.

In this **module 5**, we'll explore how we can use several tools to strengthen
the overall security posture of our software delivery process.

Let's start by applying some general best practices around pinning specific
versions of all the dependencies in use by our application. Additionally, we
will also leverage both the Go and the Docker tool chains to enforce integrity
checks of every software package loaded as dependency by explicitly tagging
their check sums. This exercise will bring us closer to achieving almost
completely hermetic builds and will prevent loading dependencies that have been
tampered with.

Next, we'll have a look at the open framework _Security Levels for Software
Artifacts_ (SLSA), which can be used to rate how secure our software delivery
process is at any given time. Based on a rating we can then leverage it to
improve the security posture of our process and increase our own confidence in
the container images we run.

After that, we will setup a policy for Binary Authorization on Cloud Run. This
will cause Cloud Run to check for a cryptographic signature to be present on all
deployed container images or reject them. We can use this technique to make sure
that all deployed software artifacts have been integrated and packaged by our
trusted CI system, in this case Cloud Build.

Finally, we'll have some fun by building a completely different container image
based on a vulnerable Java application running on an outdated Debian base. We
can then use the Google Cloud Console to explore the various vulnerabilities
present in the image.

<walkthrough-tutorial-difficulty difficulty="2"></walkthrough-tutorial-difficulty>

Estimated time:
<walkthrough-tutorial-duration duration="40"></walkthrough-tutorial-duration>

To get started, click **Start**.

## Project setup

First, let's make sure we got the correct project selected. Go ahead and select
the provided project ID.

<walkthrough-project-setup billing="true"></walkthrough-project-setup>

Run the following to make sure all required APIs are enabled.

<walkthrough-enable-apis apis="cloudbuild.googleapis.com,
run.googleapis.com,binaryauthorization.googleapis.com,
containerscanning.googleapis.com,
containeranalysis.googleapis.com,
artifactregistry.googleapis.com"> </walkthrough-enable-apis>

Let's also make sure our gcloud configuration is present by running:

```bash
gcloud config set project <walkthrough-project-id/>
gcloud config set run/region europe-north1
gcloud config set artifacts/location europe-north1
```

## Precisely loading software dependencies

As a developer, you likely understand the importance of stability, security, and
predictability in software development. Version pinning is a powerful tool to
achieve these goals, especially in large, complex projects. By explicitly
defining the software versions you rely on, you create a more robust and
maintainable codebase.

In the previously created Dockerfile, we have already been slightly specific
about which base images we would like to use during both the build as well as
the runtime stage. We note only want the `golang` image, but we want the one
that is using Go 1.22 based on Debian bookworm `golang:1.22-bookworm`. Let's
inspect the image we mean on
[Dockerhub](https://hub.docker.com/layers/library/golang/1.22-bookworm/images/sha256-4304d944c421bb279ad1faada14d03ac7e7edca61793d2f6a7ad94681c457887).
The hosting registry Dockerhub provides us with the precise SHA256 checksum of
the image. We can use the check sum in our Dockerfile to instruct the building
container daemon to verify all the contents of the pulled container image file
system layers to ensure that absolutely every single bit of the image in use has
not changed in the meantime.

Let's modify our existing
<walkthrough-editor-open-file filePath="cloudshell_open/serverless/Dockerfile">
`Dockerfile`</walkthrough-editor-open-file> to include the SHA256 checksum like
this:

<!-- markdownlint-disable MD013 -->

```Dockerfile
FROM golang:1.22-bookworm@sha256:5c56bd47228dd572d8a82971cf1f946cd8bb1862a8ec6dc9f3d387cc94136976 as builder
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

<!-- markdownlint-enable MD013 -->

Let's do the same for the runtime images. The image
`gcr.io/distroless/static-debian11` is hosted
[here](https://console.cloud.google.com/gcr/images/distroless/global/static-debian11@sha256:d27f8df264c425579aaf8d996ddbf65db6c6dc1ef8deeea8c1b7b69b2ebddf02/details).
Have a look and copy the SHA256 checksum to then also explicitly pin it in the
`FROM` directive for the final runtime stage in the Dockerfile.

<!-- markdownlint-disable MD013 -->

```Dockerfile
FROM golang:1.22-bookworm@sha256:5c56bd47228dd572d8a82971cf1f946cd8bb1862a8ec6dc9f3d387cc94136976 as builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY main.go ./
RUN CGO_ENABLED=0 go build -v -o server

FROM gcr.io/distroless/static-debian11@sha256:d27f8df264c425579aaf8d996ddbf65db6c6dc1ef8deeea8c1b7b69b2ebddf02
WORKDIR /app
COPY --from=builder /app/server /app/server
CMD ["/app/server"]
```

<!-- markdownlint-enable MD013 -->

Perfect! We have now ensured that the building container daemon checks the
integrity of the images it downloads when the docker build is run.

The base images for the docker stages are not the only software dependencies
that are being pulled into scope during the build of our application: we are
also using various Go dependencies.

Luckily, Go already has us covered. Since our application is using Go Modules,
Go automatically keeps track of it's dependencies. Take a look at the
<walkthrough-editor-open-file filePath="cloudshell_open/serverless/go.mod">
`go.mod`</walkthrough-editor-open-file> file. This file lists all the
dependencies our application requires and also pins their specific versions.

Furthermore, Go Modules uses a second automatically generated file which keeps
note of versions and checksums for all transitive dependencies used. Have a look
at the
<walkthrough-editor-open-file filePath="cloudshell_open/serverless/go.sum">
`go.sum`</walkthrough-editor-open-file> file.

It is highly recommended, that you add both files to your version control system
so that they are available to Cloud Build or any other remote build after
cloning the repository. Some of the first instructions in the Dockerfile only
copy these two files into the build file system and instruct Go Modules to
download all required dependencies using their specific versions and verify
their integrity using the supplied checksums. This is where it happens:

```Dockerfile
COPY go.* ./
RUN go mod download
```

Good. We have finally made sure that every dependency required is explicitly
pinned with it's version number and also checked against remote tampering by
verifying checksums.

Another best practice is to not only be specific about the dependencies we need
to produce an artifact, but also be specific about what we produce. We can
explicitly tag images at the end of our <walkthrough-editor-open-file
  filePath="cloudshell_open/serverless/cloudbuild.yaml">
`cloudbuild.yaml`</walkthrough-editor-open-file> file. Let's update the contents
like this:

```yaml
steps:
  - id: "Building image"
    name: "gcr.io/cloud-builders/docker"
    args:
      - "build"
      - "-t"
      - "$_ARTIFACT_REGION-docker.pkg.dev/$PROJECT_ID/my-repo/dockercloudbuildyaml:build-$BUILD_ID"
      - "."
  - id: "Tagging image"
    name: "gcr.io/cloud-builders/docker"
    args:
      - "tag"
      - "$_ARTIFACT_REGION-docker.pkg.dev/$PROJECT_ID/my-repo/dockercloudbuildyaml:build-$BUILD_ID"
      - "$_ARTIFACT_REGION-docker.pkg.dev/$PROJECT_ID/my-repo/dockercloudbuildyaml:latest"
  - id: "Pushing image:build-id"
    name: "gcr.io/cloud-builders/docker"
    args:
      - "push"
      - "$_ARTIFACT_REGION-docker.pkg.dev/$PROJECT_ID/my-repo/dockercloudbuildyaml:build-$BUILD_ID"
  - id: "Pushing image:latest"
    name: "gcr.io/cloud-builders/docker"
    args:
      - "push"
      - "$_ARTIFACT_REGION-docker.pkg.dev/$PROJECT_ID/my-repo/dockercloudbuildyaml:latest"
  - id: "Deploying to Cloud Run"
    name: "gcr.io/cloud-builders/gcloud"
    args:
      - "run"
      - "deploy"
      - "--region"
      - "$_RUN_REGION"
      - "--image"
      - "$_ARTIFACT_REGION-docker.pkg.dev/$PROJECT_ID/my-repo/dockercloudbuildyaml:build-$BUILD_ID"
      - "jokes"

images:
  - "$_ARTIFACT_REGION-docker.pkg.dev/$PROJECT_ID/my-repo/dockercloudbuildyaml:build-$BUILD_ID"
  - "$_ARTIFACT_REGION-docker.pkg.dev/$PROJECT_ID/my-repo/dockercloudbuildyaml:latest"
```

Note how we're also adding an ID to every step of the build to make it easier to
digest for the reader.

Let's issue another build:

```bash
LOCATION="$(gcloud config get-value artifacts/location)"
REGION="$(gcloud config get-value run/region)"
gcloud builds submit \
  --substitutions _ARTIFACT_REGION=${LOCATION},_RUN_REGION=${REGION}
```

Good, let's move on.

## SLSA: Security Levels for Software Artifacts

_It is pronounced 'salsa' (/ˈsɑːl.sə/) 🌶️🥑 🍅_

Think of [SLSA](https://slsa.dev/) as a security checklist or framework designed
to protect the integrity of software throughout its entire lifecycle. From the
moment code is written to when it's running in production, SLSA provides
guidelines and best practices to prevent tampering, boost integrity, and secure
both the software packages themselves and the infrastructure they run on.

SLSA is an open-source project backed by Google and other organizations,
fostering industry collaboration and promoting best practices. It is by all
means am evolving framework with ongoing efforts to refine the requirements and
make it easier to adopt. A growing ecosystem of tools is emerging to help you
implement SLSA practices, from build systems to vulnerability scanners.

Cloud Build is one of those already integrated build systems that automatically
attempts to give a SLSA rating and proactively surfaces relevant information.

SLSA defines four levels of increasing security maturity:

1. **SLSA 1**: Requires basic build practices like using a version control
   system and generating provenance (a record of how the software was built).
2. **SLSA 2**: Adds requirements for a more controlled build environment, like
   using a build service and tracking dependencies.
3. **SLSA 3**: Demands a higher level of security for the build process,
   including using a hardened build service and auditing build configurations.
4. **SLSA 4**: Represents the most stringent level, requiring practices like
   two-person review for code changes and hermetic builds (where the build
   environment is isolated from external influences).

Navigate to the
[Cloud Build section in the Google Cloud Console](https://console.cloud.google.com/cloud-build/builds)
to and explore the most recent build at the top of the list. Click into the
build and familiarize yourself with the build details view. The right section of
the screen shows the _Build Log_ tab.

Switch to the _Build Artifacts_ tab and make sure that the step "Build Summary"
is selected.

Cloud Build lists the artifacts created by this build. _View_ the _Security
Insights_ for one of the image tags. Cloud Build automatically rates our
pipeline with a SLSA level and shows relevant details on vulnerabilities,
software dependencies and all sorts of other relevant metadata about the build.

Right now there are not many issues to report. Good.

## Binary Authorization

At its core, binary authorization is a security mechanism that ensures only
trusted software is deployed and executed in your environment. It's like a
bouncer at a club, but instead of checking IDs, it verifies the authenticity and
integrity of software artifacts before they're allowed to run.

Here's how it typically works:

1. Signing: Before deployment, trusted parties (developers, build systems, etc.)
   digitally sign the software artifact (container image, executable, etc.)
   using cryptographic keys. This signature is like a seal of approval.
2. Policy Enforcement: A policy engine is configured to define what signatures
   are considered valid. This could be based on specific signers, signing keys,
   or even attributes of the software itself.
3. Verification: When the software is about to be deployed or executed, the
   system checks if it has a valid signature that matches the policy. If not,
   it's blocked.

For our pipeline, we'd have Cloud Build as the trusted party that will sign the
container artifacts it creates. We can then setup a policy in Binary
Authorization, which says something like: "only trust and deploy images built by
Cloud Build". Finally, the policy is enforced by Cloud Run, so every time a new
image gets deployed as a part of a new Cloud Run revision, the presence and
validity of a signature created by Cloud Build will be checked.

Now, Cloud Build already automatically signs every container image it produces,
so the signatures are already present.

Next, Binary Authorization comes with a _default_ policy. The default policy
doesn't enforce anything as this policy states: _"allow all images"_. Let's have
a look at what the policy looks like:

```bash
gcloud container binauthz policy export
```

We need to update this policy to that it requires an attestation, namely the
attestation called
`projects/<walkthrough-project-id />/attestors/built-by-cloud-build`. This will
enforce checking for a valid Cloud Build signature on any image being deployed.

Create a new file called <walkthrough-editor-open-file
  filePath="cloudshell_open/serverless/binauthz_policy.yaml">
`binauthz_policy.yaml`</walkthrough-editor-open-file> file. Let's put in the
following contents:

```yaml
defaultAdmissionRule:
  evaluationMode: REQUIRE_ATTESTATION
  enforcementMode: ENFORCED_BLOCK_AND_AUDIT_LOG
  requireAttestationsBy:
    - projects/<walkthrough-project-id />/attestors/built-by-cloud-build
name: projects/<walkthrough-project-id />/policy
```

Let's replace the remote default policy definition with the local
`binauthz_policy.yaml` by running the following command:

```bash
gcloud container binauthz policy import binauthz_policy.yaml
```

We know need to specify using the _default_ policy by explicitly adding it as an
option to all Cloud Run deployments. Let's deploy the previously built image:

```bash
LOCATION="$(gcloud config get-value artifacts/location)"
PROJECT="$(gcloud config get-value project)"

gcloud run deploy jokes \
  --binary-authorization default \
  --image "${LOCATION}-docker.pkg.dev/${PROJECT}/my-repo/dockercloudbuildyaml"
```

Navigate to the
[Cloud Run section in the Google Cloud Console](https://console.cloud.google.com/run)
select the `jokes` services and have a look at the _Security_ tab. It should
tell you the Binary Authorization is enabled and the _default_ policy is being
enforced.

<walkthrough-info-message>Did you know that you can also enforce the use of
certain Binary Authorization policies? Using Organization Policies you can
enforce the usage of policies at the level of an organization, a folder or
individual projects in Resource Manager. You can allow-list which policies
should be used by setting up the
[`run.allowedBinaryAuthorizationPolicies`](https://cloud.google.com/binary-authorization/docs/run/requiring-binauthz-cloud-run#gcloud)
policy.</walkthrough-info-message>

Binary Authorization is now enabled and protects our services _jokes_. It
doesn't allow any image that has not been built by Cloud Build in our project.

Let's verify this by attempting to deploy a different image, the failing one
from Module 4:

```bash
gcloud run deploy jokes \
  --binary-authorization default \
  --image gcr.io/serverless-journey/failing
```

This step _fails_, because the provided image is not signed by our project local
Cloud Build. Perfect. We are now blocking images that have not gone through our
vetted CI pipeline.

## Inspecting problematic container images

With Binary Authorization enabled we are now certain that only images built by
the build system we trust, like Cloud Build, end up in production on our compute
infrastructure, such as Cloud Run.

We have also made sure that our build processes are very precise in explicitly
specifying which dependencies to load and to reject.

However, we haven't yet really discussed dealing with vulnerabilities in our
stack. So let's address this quickly.

In the beginning of this Module, we've enabled the
[Container Scanning and Analysis APIs](https://cloud.google.com/artifact-registry/docs/analysis)
which will conveniently scan container images and other artifacts for known
vulnerabilities and surface the results in the Google Cloud Console.

Let's demo this by building a vulnerable image on purpose.

Create a new file called <walkthrough-editor-open-file
  filePath="cloudshell_open/serverless/cloudbuild-insecure.yaml">
`cloudbuild-insecure.yaml`</walkthrough-editor-open-file> file. Let's put in the
following contents:

```yaml
steps:
  - id: "Pulling image"
    name: "gcr.io/cloud-builders/docker"
    args:
      - "pull"
      - "gcr.io/serverless-journey/insecure"
  - id: "Tagging image"
    name: "gcr.io/cloud-builders/docker"
    args:
      - "tag"
      - "gcr.io/serverless-journey/insecure"
      - "$_ARTIFACT_REGION-docker.pkg.dev/$PROJECT_ID/my-repo/dockerinsecure:latest"
  - id: "Pushing image"
    name: "gcr.io/cloud-builders/docker"
    args:
      - "push"
      - "$_ARTIFACT_REGION-docker.pkg.dev/$PROJECT_ID/my-repo/dockerinsecure:latest"
images:
  - "$_ARTIFACT_REGION-docker.pkg.dev/$PROJECT_ID/my-repo/dockerinsecure:latest"
```

Take a look at the file and see what it does...

Apparently, it just pulls in a foreign image, applies a local name and stores it
in the project-local artifact registry. Let's execute the build by running:

```bash
LOCATION="$(gcloud config get-value artifacts/location)"
gcloud builds submit \
  --substitutions _ARTIFACT_REGION=${LOCATION} \
  --config cloudbuild-insecure.yaml \
  --no-source
```

Once the build completes take a look at the
[top build in Cloud Build](https://console.cloud.google.com/cloud-build/builds;region=global).
Click into the build, select the build summary and view the _Security Insights_
of the build artifact
`europe-north1-docker.pkg.dev/<walkthrough-project-id />/my-repo/dockerinsecure:latest`.

You'll immediately see many vulnerabilities detected with actionable fixes
listed as well.

Go to [Artifact Registry](https://console.cloud.google.com/artifacts) and find
the image
`europe-north1-docker.pkg.dev/<walkthrough-project-id />/my-repo/dockerinsecure@sha256:8b3934e38ea64021c4323060150cbc877511fd8d52c81165d9f8ea1e189b0a06`.

First, let's take a closer look at all the dependencies used by our image. The
complete list of dependencies is the most crucial component of our artifact's
software bill-of-materials (SBOM). You can explore it by clicking on the
_Dependencies_ tab.

It seems the container images we built is a Java application using Maven
dependencies running on an Ubuntu base.

Next, let's take a look at the detected list of vulnerabilities by clicking on
the _Vulnerabilities_ tab.

Oh oh! There seem to be quite some serious vulnerabilities, including two with
critical severity. As it turns out, the application is using a quite outdated
version of the Spring framework.

Explore the top vulnerabilities for a bit and have a look at the fixes. Note how
the Container Scanning APIs also provides you with insights around
vulnerabilities detected at the OS level as well as many other programming
language tool chains, such as Go, Python, JavaScript NPM, C# NuGet, Ruby Gems,
PHP Composer and Rust Crates.

In general, most fixes will require you to patch dependencies in your
application by upgrading the to newer versions.

## Summary

Amazing, you've learned how you can increase the security posture of your
software delivery process by understanding and applying the SLSA framework,
following best practices around explicitly specifying which software
dependencies your application requires and continuously scanning and analyzing
produced software artifact to proactively patch vulnerabilities.

Congratulations! You are ready to develop, build, deploy, operate and secure
serverless applications in production!

<walkthrough-conclusion-trophy></walkthrough-conclusion-trophy>

You have completed the tutorial, well done! Please take a moment and let us know
what you think.

<walkthrough-inline-feedback></walkthrough-inline-feedback>

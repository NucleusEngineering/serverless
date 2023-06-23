<walkthrough-metadata>
  <meta name="title" content="Serverless Journey>: Run a freshly built container on Cloud Run" />
  <meta name="description" content="Learn how to build, containerize, store and deploy a container image to Google Cloud Run." />
  <meta name="keywords" content="deploy, containers, console, run" />
</walkthrough-metadata>

# Serverless Journey: Run a freshly built container on Cloud Run

## ![Cloud Run Logo][intro image]

In this tutorial we'll learn how to remotely build a container image from source code, store the image in Artifact Registry and deploy it to Cloud Run. After that, we'll familiarize ourselves with the Cloud Run UI, it's API and explore some of the avaliable options to tweak our service.

<walkthrough-tutorial-difficulty difficulty="2"></walkthrough-tutorial-difficulty>

Estimated time:
<walkthrough-tutorial-duration duration="15"></walkthrough-tutorial-duration>

To get started, click **Start**.

[intro image]: https://walkthroughs.googleusercontent.com/content/images/run.png

## Project setup

First, let's make sure we got the correct project selected. Go ahead and select the provided project ID.

<walkthrough-project-setup billing="true"></walkthrough-project-setup>

<walkthrough-enable-apis apis="cloudbuild.googleapis.com,
run.googleapis.com,
artifactregistry.googleapis.com">
</walkthrough-enable-apis>

## Exploring the code

There is a sample application in `container/`. Let's change directory into that:

```bash
cd container/
```

Take some time and 
<walkthrough-editor-open-file filePath="cloudshell_open/serverless/container/main.go">
familiarize yourself with the code.
</walkthrough-editor-open-file>

Also have a look at the dependencies reference in the GO module. You can navigate to it's [upstream repository](https://github.com/helloworlddan/tortune) to figure out what it does.

Once you've understood what's going on, you can try to run that app directly in Cloud Shell, running the following in the terminal:

```bash
go run main.go
```

This compiles and starts the web server. You can now use Cloud Shell Web Preview on port 8080 to check out your application. You can find the Web Preview at the top right in Cloud Shell.

<walkthrough-editor-spotlight spotlightId="cloud-shell-web-preview-button" target="cloudshell">Web Preview</walkthrough-editor-spotlight>

Finally, focus the terminal again and terminate the web server with `Ctrl-C`.

## Building and deploying the code

We established that our code runs locally, great! Let's put in on the cloud.

The `gcloud` CLI has a convenient short cut for Cloud Run which quickly allows us to do so. We can use a single command to easily:
- tarball a directory (build context)
- upload it to Google Cloud Storage
- use Cloud Build to untar the context, build the app, containerize and store it on Artifact Registry
- create a new Cloud Run service with a fresh revision
- bind an IAM policy to the the service
- route traffic to the new endpoint

```bash
gcloud run deploy --source .
```

The command requires additional information and hence switches to an interactive mode. You can have a look at the `gcloud` CLI [reference](https://cloud.google.com/sdk/gcloud/reference/run/deploy) in case some of the option may be unclear to you.

We want our application to be publicly available on the internet.

Cloud Run is available in all Google Cloud Platform regions globally. If you are unsure where to deploy to, you can use the [Region Picker](https://cloud.withgoogle.com/region-picker/) tool to find the most suitable one.



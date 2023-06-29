<walkthrough-metadata>
  <meta name="title" content="Serverless Journey>: Run a freshly built container on Cloud Run" />
  <meta name="description" content="Learn how to build, containerize, store and deploy a container image to Google Cloud Run." />
  <meta name="keywords" content="deploy, containers, console, run" />
</walkthrough-metadata>

# Serverless Journey: Run a freshly built container on Cloud Run

## ![Cloud Run Logo][https://walkthroughs.googleusercontent.com/content/images/run.png]

In this tutorial we'll learn how to remotely build a container image from source code, store the image in Artifact Registry and deploy it to Cloud Run. After that, we'll familiarize ourselves with the Cloud Run UI, it's API and explore some of the avaliable options to tweak our service.

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

## Exploring the code

Take some time and 
<walkthrough-editor-open-file filePath="cloudshell_open/serverless/main.go">
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

<!-- TODO, verify OCI -->
In the following steps, we'll deploy our app to Cloud Run, which is a compute platform for running modern, cloud-first, containerized applications. Cloud Run schedules and runs container images that expose services over HTTP. You can use any programming language or framework, as long as you can bind an HTTP server on a port provided by the `$PORT` environment variable.  Please note, that container images must be in Docker or OCI format and will be run on Linux x86_64.

Cloud Run is available in all Google Cloud Platform regions globally. If you are unsure where to deploy to, you can use the [Region Picker](https://cloud.withgoogle.com/region-picker/) tool to find the most suitable one.

You can configure preferred regions and zones in `gcloud` so its invokaction get more convenient.

```bash
# Set default locations for Cloud Run and Artifact Registry to europe-north1, Finland.

gcloud config set run/region europe-north1 
gcloud config set artifacts/location europe-north1 
```

Note that our code is not yet build or containerized, yet Cloud Run requires that.
The `gcloud` CLI has a convenient short cut for deploying Cloud Run which quickly allows us to do so.

We can use a single command to easily:
- tarball a directory (build context)
- upload it to Google Cloud Storage
- use Cloud Build to untar the context, build the app, containerize and store it on Artifact Registry
- create a new Cloud Run service with a fresh revision
- bind an IAM policy to the the service
- route traffic to the new endpoint

```bash
gcloud run deploy server --source .
```

The command requires additional information and hence switches to an interactive mode. You can have a look at the `gcloud` CLI [reference](https://cloud.google.com/sdk/gcloud/reference/run/deploy) in case some of the option may be unclear to you.

We want our application to be publicly available on the internet.

Wait for the deployment to finish and then navigate to the `*.a.run.app` endpoint it created. You should be able to call the endpoint from your browser or by using any other HTTP client like cURL,

```bash
# Retrieve auto-generated URL and cURL it.

curl $(gcloud run services describe server --format 'status.url')
```

The first deployment to Cloud Run is complete!

## Using Cloud Code 

Alternatively, you can deploy and otherwise manage the lifecycle of your Cloud Run services and other resources using Cloud Code.

Cloud Code is an IDE plugin for both VS Code and Intellj-based IDEs and is designed to make you interaction with Google Cloud more convenient.

Cloud Code is a pre-installed in Cloud Shell.

Click the
<walkthrough-editor-spotlight spotlightId="activity-bar-cloud-code">
Cloud Code icon in the activity bar
</walkthrough-editor-spotlight>
to expand it.

Within the activity bar
<walkthrough-editor-spotlight spotlightId="cloud-code-cloud-run-deploy">
you can directly deploy to Cloud Run
</walkthrough-editor-spotlight>
too, with out the need to use the CLI.


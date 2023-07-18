<walkthrough-metadata>
  <meta name="title" content="Serverless Journey>: Extend your code to call Cloud APIs" />
  <meta name="description" content="Learn how to call Cloud APIs directly from your code and use service accounts to secure your app." />
  <meta name="keywords" content="deploy, containers, console, run" />
</walkthrough-metadata>

# Serverless Journey: Extend your code to call Cloud APIs

![Tutorial header image](https://storage.googleapis.com/gweb-cloudblog-publish/images/databases_2022_HTRs5Tr.max-700x700.jpg)

In this tutorial we'll learn how to extend the existing code to call Cloud APIs directly. Currently, the deployed application uses a library which contains a static set of jokes. Whenever the library is used it randomly selects a joke and returns it. After a while we will surely start to see the same jokes again and the only way to see new jokes is when a human would actually implement them in the library.

Luckliy, there is a thing called _generative AI_ now. Google Cloud Vertex AI contains a Google-built, pre-trained, PaLM model which is a general-purpose, generative LLM that can be query with free-texts prompts to generate all sorts of text-based outputs. In this tutorial we'll implement the `model:predict` endpoint of Vertex AI to execute this model in order to new dad jokes in a generative matter.

Additionally, we'll learn a little bit about custom service accounts, IAM permissions and how to use the principle of least privilege to secure our services on Cloud Run.

<walkthrough-tutorial-difficulty difficulty="3"></walkthrough-tutorial-difficulty>

Estimated time:
<walkthrough-tutorial-duration duration="15"></walkthrough-tutorial-duration>

To get started, click **Start**.

## Project setup

First, let's make sure we got the correct project selected. Go ahead and select the provided project ID.

<walkthrough-project-setup billing="true"></walkthrough-project-setup>

Run the following to make sure all required APIs are enabled. Note that `aiplatform.googleapis.com` is added.

<walkthrough-enable-apis apis="cloudbuild.googleapis.com,
run.googleapis.com,aiplatform.googleapis.com,
artifactregistry.googleapis.com">
</walkthrough-enable-apis>

## Extending the code

<!-- TODO change deps to tortuneai -->
<!-- TODO extending code base (hitme wants params now) -->
<!-- TODO new SA with perms for model:predict, tightening IAM -->


```bash
cloudshell launch-tutorial .journey/03-operate.neos.md
```


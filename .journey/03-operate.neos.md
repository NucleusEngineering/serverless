<walkthrough-metadata>
  <meta name="title" content="Serverless Journey>: Operate your services in production" />
  <meta name="description" content="Learn how to safely deploy changes and operate your services in production." />
  <meta name="keywords" content="deploy, containers, console, run" />
</walkthrough-metadata>

# Serverless Journey: Operate your services in production

![Tutorial header image](https://storage.googleapis.com/gweb-cloudblog-publish/images/web3_2022_7d2BUsw.max-700x700.jpg)

In the final part of the Serverless Journey, we are going to look at some basic principles of how to operate your services in production. We'll learn about SRE and define some custom SLOs to keep track of the health of our services.

After that we'll have a look at Cloud Run's traffic splitting capabilities. This allows us to deploy new revisions of a Cloud Run service and the gradually move portion of production traffic over. Cloud Run's programmable network control plane allows you to split traffic between revisions; you can use this built-in feature to implement strategies like blue/green deployments, canaries releases, or rollback to previous revisions in seconds should you ever push a bad release. 

To bring everything to life, we'll be deploying a faulty version as a canary release with some amount of live traffic, get an alert, observe how this will burn through the [error budget](https://cloud.google.com/blog/products/management-tools/sre-error-budgets-and-maintenance-windows) and quickly rollback to a safe state.

<walkthrough-tutorial-difficulty difficulty="3"></walkthrough-tutorial-difficulty>

Estimated time:
<walkthrough-tutorial-duration duration="35"></walkthrough-tutorial-duration>

To get started, click **Start**.

## Project setup

First, let's make sure we got the correct project selected. Go ahead and select the provided project ID.

<walkthrough-project-setup billing="true"></walkthrough-project-setup>

Run the following to make sure all required APIs are enabled.

<walkthrough-enable-apis apis="cloudbuild.googleapis.com,
run.googleapis.com,aiplatform.googleapis.com,
artifactregistry.googleapis.com">
</walkthrough-enable-apis>

## Hope is not a strategy

Google's Site Reliability Engineering discipline encompasses... yadda yadda

<!-- TODO basics of SRE -->

## Defining SLOs and alerts

<!-- TODO defining SLO and alerts -->

## Deploying a canary with traffic splitting

<!-- TODO traffic split -->

<!-- TODO load test -->

## Responding to the incident

<!-- TODO email alert -->

<!-- TODO roll back -->

## Summary

<!-- TODO Learned about traffic splitting, and SLOs -->

<walkthrough-conclusion-trophy></walkthrough-conclusion-trophy>

You have completed the Serverless Journey, congratulations! Please take a moment and let us know what you think.

<walkthrough-inline-feedback></walkthrough-inline-feedback>


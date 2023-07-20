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

## SRE: Hope is not a strategy

Google's Site Reliability Engineering discipline encompasses... yadda yadda

<!-- TODO basics of SRE -->

<!-- TODO SLI, SLO -->

<!-- TODO four golden signals -->
[Read more about the four golden signals](https://sre.google/sre-book/monitoring-distributed-systems/#xref_monitoring_golden-signals)

## Defining SLOs and alerts

Let's begin by defining an SLO for our service that keeps track of it's availability, specifically of it's error rate. A service that is unable to serve requests because requests are errors, are unavailable. We are about to define an objective that says: We aim at serving 95% of all requests without and error. Cloud Run already keeps track of application errors. You can see the log-based metric _Request count_ in the metrics section of your Cloud Run service. You can find it by clicking [into your service in the Cloud Run section](https://console.cloud.google.com/run) of the Google Cloud Console. The metrics counts HTTP status codes and `5xx` status codes are consider to be server-side errors.

Right next to the _Metrics_ tab, you'll find the tab that says _SLOs_. Get in there and start the new SLO creation wizard. Define an SLO, that:
* uses a **windows-based, availability metric** as SLI
* this should track the metric `run.googleapis.com/request_count`
* define **goodness as 95% of requests, trailing over a 5 minute windows**
* the SLO definition should look at a **rolling 1 day period** with and **objective of 95%**

Finally, also create an alert that will notify you in case a fast-burn of your error budget gets detected. Create an new alert, that:
* uses a **lookback duration of 5 minutes**
* breaches with a **fast-burn threshold 10(x baseline)**

<walkthrough-info-message>Create a new alert notification channel in Cloud Monitoring using your email.</walkthrough-info-message>

Done. You should now be immediately notified as soon as the as your service is responding with too many errors in a given moment and your about to burn through your error budget quickly.

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


<walkthrough-metadata>
  <meta name="title" content="Serverless Journey>: Operate your services in production" />
  <meta name="description" content="Learn how to safely deploy changes and operate your services in production." />
  <meta name="keywords" content="deploy, containers, console, run" />
</walkthrough-metadata>

# Serverless Journey Module 4: Operate your services in production

![Tutorial header image](https://storage.googleapis.com/gweb-cloudblog-publish/images/web3_2022_7d2BUsw.max-700x700.jpg)

<!-- TODO include on-demand clip for module 3 -->

In the final part of the Serverless Journey, we are going to look at some basic principles of how to operate your services in production. We'll learn about SRE and define some custom SLOs to keep track of the health of our services.

After that we'll have a look at Cloud Run's traffic splitting capabilities. This allows us to deploy new revisions of a Cloud Run service and the gradually move portion of production traffic over. Cloud Run's programmable network control plane allows you to split traffic between revisions; you can use this built-in feature to implement strategies like blue/green deployments, canaries releases, or rollback to previous revisions in seconds should you ever push a bad release. 

To bring everything to life, we'll be deploying a faulty version as a canary release with some amount of live traffic, get an alert, observe how this will burn through the [error budget](https://cloud.google.com/blog/products/management-tools/sre-error-budgets-and-maintenance-windows) and quickly rollback to a safe state.

<walkthrough-tutorial-difficulty difficulty="3"></walkthrough-tutorial-difficulty>

Estimated time:
<walkthrough-tutorial-duration duration="45"></walkthrough-tutorial-duration>

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

Google's Site Reliability Engineering is what you get when you treat operations as if it’s a software problem. The mission is to protect, provide for, and progress the software and systems behind all of Google’s public services — Google Search, Ads, Gmail, Android, YouTube, and App Engine, to name just a few — with an ever-watchful eye on their availability, latency, performance, and capacity.

Many of the practices and tools of SRE have been integrated into Google Cloud Platform so that everyone running services in production can benefit from what Google learned over decades.

SRE begins with the idea that a prerequisite to success is availability. A system that is unavailable cannot perform its function and will fail by default. Availability, in SRE terms, defines whether a system is able to fulfill its intended function at a point in time. In addition to being used as a reporting tool, the historical availability measurement can also describe the probability that your system will perform as expected in the future.

When we set out to define the terms of SRE, we wanted to set a precise numerical target for system availability. We term this target the availability Service-Level Objective (SLO) of our system. Any discussion we have in the future about whether the system is running sufficiently reliably and what design or architectural changes we should make to it must be framed in terms of our system continuing to meet this SLO.

We also have a direct measurement of a service’s behavior: the frequency of successful probes of our system. This is a Service-Level Indicator (SLI). When we evaluate whether our system has been running within SLO for the past week, we look at the SLI to get the service availability percentage. If it goes below the specified SLO, we have a problem and may need to make the system more available in some way, such as running a second instance of the service in a different city and load-balancing between the two. If you want to know how reliable your service is, you must be able to measure the rates of successful and unsuccessful queries as your SLIs.

### The four golden signals

SRE defines _four golden signals_ to always keep track of in most systems.

**Latency**: The time it takes to service a request. It’s important to distinguish between the latency of successful requests and the latency of failed requests. For example, an HTTP 500 error triggered due to loss of connection to a database or other critical backend might be served very quickly; however, as an HTTP 500 error indicates a failed request, factoring 500s into your overall latency might result in misleading calculations. On the other hand, a slow error is even worse than a fast error! Therefore, it’s important to track error latency, as opposed to just filtering out errors.

**Traffic**: A measure of how much demand is being placed on your system, measured in a high-level system-specific metric. For a web service, this measurement is usually HTTP requests per second, perhaps broken out by the nature of the requests (e.g., static versus dynamic content). For an audio streaming system, this measurement might focus on network I/O rate or concurrent sessions. For a key-value storage system, this measurement might be transactions and retrievals per second.

**Errors**: The rate of requests that fail, either explicitly (e.g., HTTP 500s), implicitly (for example, an HTTP 200 success response, but coupled with the wrong content), or by policy (for example, "If you committed to one-second response times, any request over one second is an error"). Where protocol response codes are insufficient to express all failure conditions, secondary (internal) protocols may be necessary to track partial failure modes. Monitoring these cases can be drastically different: catching HTTP 500s at your load balancer can do a decent job of catching all completely failed requests, while only end-to-end system tests can detect that you’re serving the wrong content.

**Saturation**: How "full" your service is. A measure of your system fraction, emphasizing the resources that are most constrained (e.g., in a memory-constrained system, show memory; in an I/O-constrained system, show I/O). Note that many systems degrade in performance before they achieve 100% utilization, so having a utilization target is essential.

In complex systems, saturation can be supplemented with higher-level load measurement: can your service properly handle double the traffic, handle only 10% more traffic, or handle even less traffic than it currently receives? For very simple services that have no parameters that alter the complexity of the request (e.g., "Give me a nonce" or "I need a globally unique monotonic integer") that rarely change configuration, a static value from a load test might be adequate. As discussed in the previous paragraph, however, most services need to use indirect signals like CPU utilization or network bandwidth that have a known upper bound. Latency increases are often a leading indicator of saturation. Measuring your 99th percentile response time over some small window (e.g., one minute) can give a very early signal of saturation.

Finally, saturation is also concerned with predictions of impending saturation, such as "It looks like your database will fill its hard drive in 4 hours."

You can [read more about the four golden signals](https://sre.google/sre-book/monitoring-distributed-systems/#xref_monitoring_golden-signals).

## Defining SLOs and alerts

Let's begin by defining an SLO for our service that keeps track of it's availability, specifically of it's error rate. A service that is unable to serve requests because requests are errors, are unavailable. We are about to define an objective that says: We aim at serving 95% of all requests without and error. Cloud Run already keeps track of application errors. You can see the log-based metric _Request count_ in the metrics section of your Cloud Run service. You can find it by clicking [into your service in the Cloud Run section](https://console.cloud.google.com/run) of the Google Cloud Console. The metrics counts HTTP status codes and `5xx` status codes are consider to be server-side errors.

Right next to the _Metrics_ tab, you'll find the tab that says _SLOs_. Get in there and start the new SLO creation wizard. Define an SLO, that:
* uses a **windows-based, availability metric** as SLI
* this should track the metric `run.googleapis.com/request_count`
* define **goodness as 95% of requests, trailing over a 5 minute windows**
* the SLO definition should look at a **rolling 1 day period** with and **objective of 95%**

Finally, also create an alert that will notify you in case a fast-burn of your error budget gets detected. Create an new alert, that:
* uses a **look back duration of 5 minutes**
* breaches with a **fast-burn threshold 10(x baseline)**
* creates a **new alert notification channel using your email** in Cloud Monitoring

Done. You should now be immediately notified as soon as the as your service is responding with too many errors in a given moment and your about to burn through your error budget quickly.

## Deploying a canary with traffic splitting

Before continuing, make sure `gcloud` configuration is set to our preferences:

```bash
gcloud config set project <walkthrough-project-id/>
gcloud config set run/region europe-north1 
gcloud config set artifacts/location europe-north1 
```

Okay! Good to go.

Cloud Run comes with a built-in traffic control plane, which lets operators programmatically assign traffic to individual revisions of the same service. This allows you to deploy new changes to your service with none or only single-digit traffic routed to them. Once you've slowly gained confidence in your new revision you can gradually increased traffic until you reach a 100% of traffic on the new revision and your rollout is complete. This strategy is commonly referred to as a [canary release](https://cloud.google.com/deploy/docs/deployment-strategies/canary).

<walkthrough-info-message>Cloud Run's default deployment strategy is to automatically route all traffic to the new revision if it should pass minimum health checks.</walkthrough-info-message>

In Cloud Run, you can tag your individual revisions. Each tagged revision can be easily refereed to by its tag and also get an individual endpoint, just for itself. That means we can preview the revision even without assigning production traffic to it.

Let's deploy a new revision of the jokes service using a different container image. We've been told that requirements have changed dramatically for this new version. We don't trust the image and decide to implement and canary release strategy. First of all, let's tag the current, good revision, like this:

```bash
gcloud run services update jokes \
    --tag good
```

Good! The revision is tagged and gets a separate URL that looks something like `https://good---jokes-xxxxxxxxxx-xx.a.run.app`. As long as we keep the tag around we can always directly reach this revision with this endpoint. You can try it by directly cURLing it:

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

Remember: we don't trust this new container image. So, let's begin by routing only 10% of traffic to the new revision tagged as `next`, by running the following:

```bash
gcloud run services update-traffic jokes \
    --to-tags next=10
```

Cloud Run will quickly reconfigure the traffic pattern and print out the current pattern. It's time to simulate some traffic to see how the canary behaves under traffic. Use `hey` again to do this:

```bash
hey -n 50000 $(gcloud run services describe jokes --format 'value(status.url)')
```

Eventually, the trajectory with which we will be burning though the error budget will cause an alert to fire.

<walkthrough-info-message>Now is probably a good time to catch a little break, get a snack, get hydrated, maybe take a stroll through the park. Don't forget: you are still on-call!</walkthrough-info-message>

Proceed to the next section, once you've received the alert email.

## Responding to the incident

Oh, no! We have an active incident!

Follow the link in the alert mail or navigate to the [incidents section of Cloud Monitoring](https://console.cloud.google.com/monitoring/alerting/incidents). You should you the active incident marked in red, explore it. Notice how the rate at which we are burning error budget is way past threshold.

Also look at the SLOs section of your Cloud Run service in the [Cloud Run section of the Google Cloud Console](https://console.cloud.google.com/run). The page tells you that the percentage of healthy, non-error responses from your service has dropped to 90%, which is below your SLO target of 95%!

Let's investigate a bit further. Head to the _Logs_ tab of the service. Next to your application logs, Cloud Run itself logs system events here. You can see how every invocation of your service produces and event. Most lines are marked with a blue `i` (Loglevel `info`). You see that these are the requests to which the service responded with `HTTP 200 OK`. About 10% of the log lines are marked with a red `!!` (Loglevel `error`), though!

Have a closer look at one of each log line and expand into the fields `resource` and further down into `labels`. The structured log message tells you exactly which revision of the service belongs to the message. Turns out that all the errors are caused by the revision we previously marked as `next`.

Now, before we tell the engineering team about the problem, we first keep calm and a remember Cloud Run traffic management. To bring the system back to a safe state, let's rollback to the previous revision, the one we tagged with `good`. To do so, simply reconfigured the traffic pattern on the service like this:

```bash
gcloud run services update-traffic jokes \
    --to-tags good=100
```

Good! Run a bit more simulated traffic over the service to check if this remediated the issue.

```bash
hey -n 50000 $(gcloud run services describe jokes --format 'value(status.url)')
```

Go back to the _SLOs_ section of your service and observe how the SLI for your error rate climbs back over 95%. Crisis averted!

## Summary

Great! You've learned the basic of SRE, defined meaningful SLOs, mastered Cloud Run traffic management and responded to a failing canary release by rolling back to a known good state of the system. 

You are ready to develop, build, deploy and run serverless applications in production!

<walkthrough-conclusion-trophy></walkthrough-conclusion-trophy>

You have completed the Serverless Journey, congratulations! Please take a moment and let us know what you think.

<walkthrough-inline-feedback></walkthrough-inline-feedback>


## sia_exporter, a lightweight and flexible exporter to send Sia metrics to
Prometheus
Introducing sia_exporter a lightweight and flexible metrics exporter to send Sia
metrics to Prometheus for analysis and Grafana for graphing and alerting.
Sia_exporter makes it easy to monitor one or many Sia instances using
Prometheus, Grafana, and other tools.

tl;dr: `sia_exporter` lets anyone create custom Sia graphs, alerts, and
dashboards.

![example dashboard](https://i.imgur.com/xHZvj5Z.png)

## What on Earth is Prometheus?
Prometheus is an extremely popular metric gathering and monitoring solution.
Prometheus server scrapes HTTP endpoints at regular intervals and stores that
scraped data in a super-efficient database. The result is a highly performant
repository of time-series metrics representing whatever data is served by the
exporters. Prometheus then presents this data for analysis and display to a
variety of other tools such as Grafana.

## What on Earth is Grafana?
Grafana is an extremely popular vizualization tool that interfaces with various
types of databases, including Prometheus. Here's how they describe themselves.
Grafana allows you to query, visualize, alert on and understand your metrics no
matter where they are stored. Create, explore, and share dashboards with your
team and foster a data driven culture.

## What exactly does sia_exporter do?
Sia_exporter queries Sia's API and gathers data (metrics), then serves it over
HTTP for Prometheus to consume and Grafana to graph. Prometheus consumes the
data, organizes it into time-series values, which Grafana can then query and
graph. The whole workflow can be broken down into 3step.
*  Sia_exporter queries the Sia API and generates a set of metrics for
   Prometheus to scrape. By default, this happens every 5 minutes.
*  Prometheus server scrapes the metrics from sia_exporter and saves the data
   into its internal database. The scrape interval can be customized, but I
recommend the same 5 minutes that the sia_exporter is being updated.
*  Grafana queries Prometheus for the data, then graphs it in customizable,
   beautiful graphs. Grafana can also generate alerts, such as emails, when
certain conditions are met.

![Allowance Graphs](https://i.imgur.com/0x7oHQx.png)

## How do I get sia_exporter?
`Sia_exporter` can be found on GitHub at
https://github.com/tbenz9/sia_exporter/releases. Simply download the binary for
your Operating System and run it on the same system as your Sia instance.
`Sia_exporter`, Prometheus, and Grafana all work on virtually all operating
systems and architectures. Whether you're a Windows, Mac, or Linux user running
on just about any hardware platform there should be a binary that works for you.

## How do I use sia_exporter?
`Sia_exporter` is a command-line tool and should be ready to use straight out of
the box for most users.
```
$> ./sia_exporter
INFO[0000] Beginning to metrics at http://<your ip address>:9983/metrics
```
Once you've got sia_exporter running you can start adding metrics to Grafana.
Metrics are organized by Sia modules. Adding charts, graphs, and alerts are as
simple as point and click!

![Add Metrics](https://i.imgur.com/FTN3U15.png)

For more advanced users sia_exporter does have a number of command-line flags to
turn functionality on/off, adjust options, and access a remote Sia instance.
```
$> ./sia_exporter -h
Usage of ./sia_exporter:
  -address string
        Sia's API address (default "127.0.0.1:9980")
  -agent string
        Sia agent (default "Sia-Agent")
  -debug
        Enable debug mode. Warning: generates a lot of output.
  -modules string
        Sia Modules to monitor (default "cghmrtw")
  -port int
        Port to serve Prometheus Metrics on (default 9983)
  -refresh int
        Frequency to get Metrics from Sia (minutes) (default 5)
```
        
## Troubleshooting and installation details
Verify that `sia_exporter` is gathering metrics and serving them over HTTP. This
step verifies that `sia_exporter` is working as expected. Make sure you enter
your private IP address of the node running sia_exporter wherever `<your ip
address>` is shown.
```
$> curl -s http://<your ip address>:9983/metrics
# HELP consensus_difficulty Consensus difficulty
# TYPE consensus_difficulty gauge
consensus_difficulty 1.8213302339204035e+18
# HELP consensus_height Consensus block height
# TYPE consensus_height gauge
consensus_height 229577
# HELP consensus_module_loaded Is the consensus module loaded. 0=not loaded.
1=loaded
# TYPE consensus_module_loaded gauge
consensus_module_loaded 1
# HELP consensus_synced Consensus sync status, 0=not synced.  1=synced
# TYPE consensus_synced gauge
consensus_synced 1

... truncated
```

This guide does not cover installing Prometheus or Grafana but both projects
have excellent installation guides on their websites. The rest of the
instructions assume you have Prometheus and Grafana installed and working.

Configure your prometheus.yaml file to scrape the new `sia_exporter` metrics.
Below is a sample prometheus.yaml file scraping a single `sia_exporter`
endpoint. Don't forget to change the IP address to the IP address of the node
running sia_exporter.
```
$> cat /etc/prometheus/prometheus.yaml
global:
        scrape_interval: 300s
scrape_configs:
        - job_name: 'sia_exporter'
          metrics_path: /metrics
          static_configs:
                  - targets: ['<your ip address>:9983']
```
Now log in to the Prometheus and verify its successfully scraping the
sia_exporter metrics.

![Prometheus is scraping sia_exporter](https://i.imgur.com/R8MP4HF.png)

Now log in to Grafana and start creating panels, dashboards, alerts, and acting
on your new Sia metrics!

![Sample Wallet](https://i.imgur.com/hWK2UBG.png)

## Example use cases
Prometheus is designed to scrape thousands of endpoints quickly if you have
multiple Sia instances running in your network you can install `sia_exporter` on
each of them and scrape them from a single Prometheus server. For example,
suppose you manage 10 Sia wallets, you could easily graph the balance of all 10
wallets on a single chart, and send an email alert when the balance gets too low
on any single wallet.

![Sample Alert Configuration](https://i.imgur.com/xMYw1R9.png)

Suppose you have 100 Sia instances and want to ensure they all have the correct
block height, simply set up `sia_exporter` for all 100 Sia instances and graph
the block height. Turn the chart red or send an alert if one of the block
heights is different than the others.

Grafana also integrates seamlessly with alerting tools such as PagerDuty, and
AlertManager for more advanced alert use cases. It also offers cool features
such as "kiosk" mode to keep your favorite dashboards on display in a public
way. Got an extra tablet lying around, set it up to display your Sia status,
wallet balance, free space, number of uploaded files, etc.

## Like what you see? Want to see more?
Contribute by opening a ticket or pull request, letting me know in the comments
section below, on reddit at /u/tbenz9 or mentioning me in the Sia Discord at
@tbenz9#2796.
Beer money - Siacoin:
`f63f6c5663efd3dcee50eb28ba520661b1cd68c3fe3e09bb16355d0c11523eebef454689d8cf`

---

<a href="https://sia.tech"><img
src="https://files.helpdocs.io/YzA4Zq3JuM/other/1571158167508/built-with-sia-color.png"
width="300"></a>

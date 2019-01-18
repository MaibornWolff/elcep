# Running example

This example `docker-compose.yml` runs a independent instance of elcep.
It contains
- elastic search
- a log generator (filling the elastic search instance with 'random' logs)
- kibana (showing the elastic search contents directly)
- elcep (providing metrics over the logs to prometheus)
- prometheus (showing the metrics exported by elcep)

This example is meant to provide a starting point for your own setup.
Feel free to copy it and adjust it to your needs.

## Adjustments

1. **Remove the log-generator**. In your environment you probably don't need this. It only pushes logs to elasticsearch in order to get some example data.
2. **Use your own elastic search instance**. Make sure to adjust hostname and port for elcep and kibana (if needed).
3. **Use your own kibana** - if you want/need it at all.
4. **Configure your prometheus**. If you use a different prometheus instance, you want to add the elcep metrics. See [`prometheus.yml`](prometheus.yml) for the configuration.

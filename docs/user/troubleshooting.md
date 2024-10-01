# Troubleshooting

## Viewing a Device's Effective Target Configuration

The device manifest returned by the `flightctl get device` command still only contains references to external configuration and secret objects. Only when the device agent queries the service, the service will replace the references with the actual configuration and secret data. While this better protects potentially sensitive data, it also makes troubleshooting faulty configurations hard. This is why a user can be authorized to query the effective configuration as rendered by the service to the agent.

To query that configuration, use the following command.

```console
flightctl get device/${device_name} --rendered | jq
```
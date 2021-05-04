# 389ds / FreeIPA Exporter

Started out as just a replication status exporter, and evolved to export more FreeIPA related objects.

To run:
```bash
go build
./389ds_exporter [flags]
```
## Exported Metrics

| Metric | Meaning | Labels |
| ------ | ------- | ------ |
| ldap_389ds_users | Number of FreeIPA users | active, staged, preserved |
| ldap_389ds_groups | Number of FreeIPA groups | |
| ldap_389ds_hosts | Number of FreeIPA hosts | |
| ldap_389ds_hostgroups | Number of FreeIPA hostgroups | |
| ldap_389ds_hbac_rules | Number of FreeIPA HBAC rules | |
| ldap_389ds_sudo_rules | Number of FreeIPA SUDO rules | |
| ldap_389ds_dns_zones | Number of FreeIPA DNS zones (including forward zones) | |
| ldap_389ds_replication_conflicts | Number of LDAP replication conflicts | |
| ldap_389ds_replication_status | Replication status of peered 389ds nodes (1 good, 0 bad) | server |
| ldap_389ds_scrape_count | Number of successful or unsuccessful scrapes | result |
| ldap_389ds_scrape_duration_seconds | How long the last scrape took |

### Flags

```bash
./389ds_exporter --help
Usage of ./389ds_exporter:
```

| CLI flag                  | Environment variable     | Description                                                           |
|---------------------------|--------------------------|-----------------------------------------------------------------------|
| `--config`                | `389DS_CONFIG`            | YAML format config file with the extension (i.e. /path/to/config.yaml) |
| `--interval`              | `389DS_INTERVAL`                    | Scrape interval (default 60s)                                         |
| `--ipa.dns`               | `389DS_IPA_DNS`                     | Should we scrape DNS stats? (default true)              |
| `--ipa.domain`            | `389DS_IPA_DOMAIN`                  | FreeIPA domain e.g. example.org                                       |
| `--ldap.addr`             | `389DS_LDAP_ADDR`                   | URI of 389ds server (default "ldap://localhost:389")                  |
| `--ldap.cert`             | `389DS_LDAP_CERT`                   | Certificate for LDAP with startTLS or TLS                             |
| `--ldap.cert-server-name` | `389DS_CERT_SERVER_NAME`            | ServerName for LDAP with startTLS or TLS                              |
| `--ldap.enablestarttls`   | `389DS_ENABLESTARTTLS`              | Use StartTLS for ldap:// connections                                  |
| `--ldap.pass`             | `389DS_LDAP_PASS`                   | 389ds Directory Manager password                                      |
| `--ldap.user`             | `389DS_LDAP_USER`                   | 389ds Directory Manager user (default "cn=Directory Manager")         |
| `--log.format`              | `389DS_LOG_FORMAT`                    | Log format (default or json) messages                                           |
| `--log.level`                 | `389DS_LOG_LEVEL`                       | Log level logging                                                         |
| `--web.listen-address`    | `389DS_WEB_LISTEN_ADDRESS`          | Bind address for prometheus HTTP metrics server (default ":9496")     |
| `--web.telemetry-path`    | `389DS_WEB_TELEMETRY_PATH`          | Path to expose metrics on (default "/metrics")                        |


### Credits

This repo essentially started off as a clone of the openldap_exporter modified to query
some FreeIPA DNs. The openldap_exporter was a great help in getting this started, as was
the consul_exporter which served as a great reference on how to package a prometheus
exporter.

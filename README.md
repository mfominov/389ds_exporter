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

* __`--config string`:__  YAML format config file with the extension (i.e. /path/to/config.yaml)
* __`--debug`:__  Debug logging
* __`--interval duration`:__  Scrape interval (default 60s)
* __`--ipa-domain string`:__  FreeIPA domain e.g. example.org
* __`--ldap.addr string`:__  URI of 389ds server (default "ldap://localhost:389")
* __`--ldap.cert string`:__  Certificate for LDAP with startTLS or TLS
* __`--ldap.cert-server-name string`:__  ServerName for LDAP with startTLS or TLS
* __`--ldap.enablestarttls`:__  Use StartTLS for ldap:// connections
* __`--ldap.pass string`:__  389ds Directory Manager password
* __`--ldap.user string`:__  389ds Directory Manager user (default "cn=Directory Manager")
* __`--log-json`:__  JSON formatted log messages
* __`--web.listen-address string`:__  Bind address for prometheus HTTP metrics server (default ":9496")
* __`--web.telemetry-path string`:__  Path to expose metrics on (default "/metrics")

### Credits

This repo essentially started off as a clone of the openldap_exporter modified to query
some FreeIPA DNs. The openldap_exporter was a great help in getting this started, as was
the consul_exporter which served as a great reference on how to package a prometheus
exporter.

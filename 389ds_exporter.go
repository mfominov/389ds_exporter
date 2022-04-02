package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/jessebl/389ds_exporter/exporter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	flag.String("web.listen-address", ":9496", "Bind address for prometheus HTTP metrics server")
	flag.String("web.telemetry-path", "/metrics", "Path to expose metrics on")
	flag.Bool("ipa.dns", true, "Should we scrape DNS stats?")
	flag.String("ldap.addr", "ldap://localhost:389", "URI of 389ds server")
	flag.String("ldap.user", "cn=Directory Manager", "389ds Directory Manager user")
	flag.String("ldap.pass", "", "389ds Directory Manager password")
	flag.String("ldap.cert", "", "Certificate for  LDAP with startTLS")
	flag.String("ldap.cert-server-name", "", "ServerName for LDAP with startTLS")
	flag.String("ipa.domain", "", "FreeIPA domain e.g. example.org")
	flag.Duration("interval", 60*time.Second, "Scrape interval")
	flag.Bool("ldap.enablestarttls", false, "Use StartTLS")
	flag.Bool("debug", false, "Debug logging")
	flag.String("log.level", "info", "Log level")
	flag.String("log.format", "default", "Log format (default, json)")
	flag.String("config", "", "YAML format config file with the extension (i.e. /path/to/config.yaml)")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.SetEnvPrefix("ds")
	viper.AutomaticEnv()
	configFile := viper.GetString("config")
	if configFile != "" {
		viper.SetConfigType("yaml")
		viper.SetConfigFile(configFile)
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatal("error reading config. the error is", err)
		}
	}

	listenPort := viper.GetString("web.listen-address")
	metricsPath := viper.GetString("web.telemetry-path")
	ldapAddr := viper.GetString("ldap.addr")
	ldapUser := viper.GetString("ldap.user")
	ldapPass := viper.GetString("ldap.pass")
	ldapCert := viper.GetString("ldap.cert")
	ldapCertServerName := viper.GetString("ldap.cert-server-name")
	ipaDomain := viper.GetString("ipa.domain")
	interval := viper.GetDuration("interval")
	enableStartTLS := viper.GetBool("ldap.enablestarttls")
	ipaDns := viper.GetBool("ipa.dns")
	logLevel := viper.GetString("log.level")
	logFormat := viper.GetString("log.format")

	if level, err := log.ParseLevel(logLevel); err != nil {
		log.Fatalf("log.level must be one of %v", log.AllLevels)
	} else {
		log.SetLevel(level)
	}

	switch logFormat {
	case "default":
		log.SetFormatter(&log.TextFormatter{})
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.Fatal("log.level must be one of [default json]")
	}

	if ldapPass == "" {
		log.Fatal("ldapPass cannot be empty")
	}
	if ipaDomain == "" {
		log.Fatal("ipaDomain cannot be empty")
	}

	if (ldapCert == "") != (ldapCertServerName == "") {
		log.Fatal("ldapCert & ldapCertServerName must come together")
	}
	tlsConf := &tls.Config{}
	if ldapCert != "" {
		roots := x509.NewCertPool()
		b, err := ioutil.ReadFile(ldapCert)
		if err != nil {
			log.Panic(err)
		}
		ok := roots.AppendCertsFromPEM(b)
		if !ok {
			log.Panic("failed to parse root cert")
		}
		tlsConf = &tls.Config{ServerName: ldapCertServerName, RootCAs: roots}
	} else {
		tlsConf.InsecureSkipVerify = true
	}

	log.Info("Starting prometheus HTTP metrics server on ", listenPort)
	go StartMetricsServer(listenPort, metricsPath)

	log.Info("Starting 389ds scraper for ", ldapAddr)
	log.Debug("Starting metrics scrape")
	exporter.ScrapeMetrics(ldapAddr, ldapUser, ldapPass, ipaDomain, tlsConf, enableStartTLS, ipaDns)
	for range time.Tick(interval) {
		log.Debug("Starting metrics scrape")
		exporter.ScrapeMetrics(ldapAddr, ldapUser, ldapPass, ipaDomain, tlsConf, enableStartTLS, ipaDns)
	}
}

func StartMetricsServer(bindAddr, metricsPath string) {
	d := http.NewServeMux()
	d.Handle(metricsPath, promhttp.Handler())
	d.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
<head><title>389ds Exporter</title></head>
						 <body>
						 <h1>389ds Exporter</h1>
						 <p><a href='` + metricsPath + `'>Metrics</a></p>
						 </dl>
						 <h2>Build</h2>
						 <pre>` + version.Info() + ` ` + version.BuildContext() + `</pre>
						 </body>
						 </html>`))
	})

	err := http.ListenAndServe(bindAddr, d)
	if err != nil {
		log.Fatal("Failed to start metrics server, error is:", err)
	}
}

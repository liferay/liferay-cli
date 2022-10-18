package mkcert

import (
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/spf13/viper"
	"liferay.com/lcectl/constants"
	lio "liferay.com/lcectl/io"
)

func init() {
	viper.SetDefault(constants.Const.TlsLfrdevDomain, "lfr.dev")
}

func VerifyRootCALoaded() {
	log.SetFlags(0)
	var m = &mkcert{}

	m.CAROOT = getCAROOT()
	if m.CAROOT == "" {
		log.Fatalln("ERROR: failed to find the default CA location, set one as the CAROOT env var")
	}
	fatalIfErr(os.MkdirAll(m.CAROOT, 0755), "failed to create the CAROOT")
	m.loadCA()
	log.Printf("Successfully loaded certicate authority. Issuer=[%s]", m.caCert.Issuer.String())
}

func InstallRootCA() {
	var (
		installFlag   = flag.Bool("install", true, "")
		uninstallFlag = flag.Bool("uninstall", false, "")
		pkcs12Flag    = flag.Bool("pkcs12", false, "")
		ecdsaFlag     = flag.Bool("ecdsa", false, "")
		clientFlag    = flag.Bool("client", false, "")
		csrFlag       = flag.String("csr", "", "")
		certFileFlag  = flag.String("cert-file", "", "")
		keyFileFlag   = flag.String("key-file", "", "")
		p12FileFlag   = flag.String("p12-file", "", "")
	)

	var m = &mkcert{
		installMode: *installFlag, uninstallMode: *uninstallFlag, csrPath: *csrFlag,
		pkcs12: *pkcs12Flag, ecdsa: *ecdsaFlag, client: *clientFlag,
		certFile: *certFileFlag, keyFile: *keyFileFlag, p12File: *p12FileFlag,
	}

	args := []string{}
	log.Println("In order to install your local CA into your truststores one-time admin privileges are needed.")
	m.Run(args)
}

func UninstallRootCA() {
	var (
		installFlag   = flag.Bool("install", false, "")
		uninstallFlag = flag.Bool("uninstall", true, "")
		pkcs12Flag    = flag.Bool("pkcs12", false, "")
		ecdsaFlag     = flag.Bool("ecdsa", false, "")
		clientFlag    = flag.Bool("client", false, "")
		csrFlag       = flag.String("csr", "", "")
		certFileFlag  = flag.String("cert-file", "", "")
		keyFileFlag   = flag.String("key-file", "", "")
		p12FileFlag   = flag.String("p12-file", "", "")
	)

	var m = &mkcert{
		installMode: *installFlag, uninstallMode: *uninstallFlag, csrPath: *csrFlag,
		pkcs12: *pkcs12Flag, ecdsa: *ecdsaFlag, client: *clientFlag,
		certFile: *certFileFlag, keyFile: *keyFileFlag, p12File: *p12FileFlag,
	}

	args := []string{}
	log.Println("In order to uninstall your local CA into your truststores one-time admin privileges are needed.")
	m.Run(args)
}

func MakeCert() {
	caroot := getCAROOT()
	lfrdevDomain := viper.GetString(constants.Const.TlsLfrdevDomain)

	var (
		installFlag   = flag.Bool("install", false, "")
		uninstallFlag = flag.Bool("uninstall", false, "")
		pkcs12Flag    = flag.Bool("pkcs12", false, "")
		ecdsaFlag     = flag.Bool("ecdsa", false, "")
		clientFlag    = flag.Bool("client", false, "")
		csrFlag       = flag.String("csr", "", "")
		certFileFlag  = flag.String("cert-file", path.Join(caroot, lfrdevDomain+".crt"), "")
		keyFileFlag   = flag.String("key-file", path.Join(caroot, lfrdevDomain+".key"), "")
		p12FileFlag   = flag.String("p12-file", "", "")
	)

	args := []string{fmt.Sprintf("*.%s", lfrdevDomain)}

	var m = &mkcert{
		installMode: *installFlag, uninstallMode: *uninstallFlag, csrPath: *csrFlag,
		pkcs12: *pkcs12Flag, ecdsa: *ecdsaFlag, client: *clientFlag,
		certFile: *certFileFlag, keyFile: *keyFileFlag, p12File: *p12FileFlag,
	}

	m.Run(args)
}

func CopyCerts(verbose bool) {
	caroot := getCAROOT()
	lfrdevDomain := viper.GetString(constants.Const.TlsLfrdevDomain)
	repoDir := viper.GetString(constants.Const.RepoDir)

	lfrdevCrtFile := path.Join(caroot, lfrdevDomain+".crt")
	lfrdevKeyFile := path.Join(caroot, lfrdevDomain+".key")
	lfrdevRootCA := path.Join(caroot, rootName)

	if !lio.Exists(lfrdevCrtFile) || !lio.Exists(lfrdevKeyFile) || !lio.Exists(lfrdevRootCA) {
		log.Fatalf("Missing one or more local certificates.  Execute 'runtime mkcert' command to generate one.")
	}

	crt, key, err := loadX509KeyPair(lfrdevCrtFile, lfrdevKeyFile)
	if crt == nil || key == nil || err != nil {
		log.Fatalf("Could not load x509 key pair: %s", err)
	}
	lrdevDomain := viper.GetString(constants.Const.TlsLfrdevDomain)
	if crt.DNSNames[0] != "*."+lrdevDomain {
		log.Fatalf("Generated certificate DNSName does not match configured domain: %s != %s\nPlease run 'runtime mkcert' command again.", crt.DNSNames[0], lrdevDomain)
	}

	lio.Copy(lfrdevCrtFile, path.Join(repoDir, fmt.Sprintf("/k8s/tls/%s.crt", lfrdevDomain)), 1024, verbose)
	lio.Copy(lfrdevKeyFile, path.Join(repoDir, fmt.Sprintf("/k8s/tls/%s.key", lfrdevDomain)), 1024, verbose)
	lio.Copy(lfrdevRootCA, path.Join(repoDir, "/k8s/tls/", rootName), 1024, verbose)
	lio.Copy(lfrdevRootCA, path.Join(repoDir, "/docker/images/dxp-server/", rootName), 1024, verbose)
	lio.Copy(lfrdevRootCA, path.Join(repoDir, "/docker/images/localdev-server/", rootName), 1024, verbose)
}

func loadX509KeyPair(certFile, keyFile string) (*x509.Certificate, any, error) {
	cf, e := ioutil.ReadFile(certFile)
	if e != nil {
		fmt.Println("cfload:", e.Error())
		return nil, nil, e
	}
	kf, e := ioutil.ReadFile(keyFile)
	if e != nil {
		fmt.Println("kfload:", e.Error())
		return nil, nil, e
	}
	cpb, _ := pem.Decode(cf)
	kpb, _ := pem.Decode(kf)
	crt, e := x509.ParseCertificate(cpb.Bytes)
	if e != nil {
		fmt.Println("parsex509:", e.Error())
		return nil, nil, e
	}
	key, e := x509.ParsePKCS8PrivateKey(kpb.Bytes)
	if e != nil {
		fmt.Println("parsekey:", e.Error())
		return nil, nil, e
	}
	return crt, key, nil
}

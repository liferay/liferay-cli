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
	var m = &mkcert{}
	m.loadCA()
	log.Println("In order to install your local CA into your truststores one-time admin privileges are needed.")
	m.install()
}

func UninstallRootCA() {
	var m = &mkcert{}
	m.loadCA()
	log.Println("In order to uninstall your local CA from your truststores one-time admin privileges are needed.")
	m.uninstall()
}

func MakeCert() {
	repoDir := viper.GetString(constants.Const.RepoDir)
	lfrdevDomain := viper.GetString(constants.Const.TlsLfrdevDomain)

	var (
		installFlag   = flag.Bool("install", false, "")
		uninstallFlag = flag.Bool("uninstall", false, "")
		pkcs12Flag    = flag.Bool("pkcs12", false, "")
		ecdsaFlag     = flag.Bool("ecdsa", false, "")
		clientFlag    = flag.Bool("client", false, "")
		csrFlag       = flag.String("csr", "", "")
		certFileFlag  = flag.String("cert-file", path.Join(repoDir, "/k8s/tls/lfrdev.crt"), "")
		keyFileFlag   = flag.String("key-file", path.Join(repoDir, "/k8s/tls/lfrdev.key"), "")
		p12FileFlag   = flag.String("p12-file", "", "")
	)

	args := []string{fmt.Sprintf("*.%s", lfrdevDomain)}

	var m = &mkcert{
		installMode: *installFlag, uninstallMode: *uninstallFlag, csrPath: *csrFlag,
		pkcs12: *pkcs12Flag, ecdsa: *ecdsaFlag, client: *clientFlag,
		certFile: *certFileFlag, keyFile: *keyFileFlag, p12File: *p12FileFlag,
	}

	m.Run(args)
	copyRootCA()
}

func copyRootCA() {
	repoDir := viper.GetString(constants.Const.RepoDir)
	caroot := getCAROOT()

	lio.Copy(path.Join(caroot, rootName), path.Join(repoDir, "/k8s/tls/", rootName), 1024)
	lio.Copy(path.Join(caroot, rootName), path.Join(repoDir, "/docker/images/dxp-server/", rootName), 1024)
	lio.Copy(path.Join(caroot, rootName), path.Join(repoDir, "/docker/images/localdev-server/", rootName), 1024)
}

func GetRootName() string {
	return rootName
}

func LoadX509KeyPair(certFile, keyFile string) (*x509.Certificate, any, error) {
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

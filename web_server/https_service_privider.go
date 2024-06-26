package web_server

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/big"
	"os"
	"owl"
	"time"
)

type HttpsOptions struct {
	*WebServerOptions
	Port     int    `json:"port"`
	KeyFile  string `json:"key-file"`
	CertFile string `json:"cert-file"`
}

func NewHttpsOptionFromConfigFile(cfgManager *owl.ConfManager, file string) (opt *HttpsOptions) {
	err := cfgManager.GetConfig(file+".https", &opt)
	if err != nil {
		return nil
	}
	return opt
}

func NewDefaultHttpsOption(stage *owl.Stage) (opt *HttpsOptions) {
	opt = &HttpsOptions{
		WebServerOptions: &WebServerOptions{
			Domain:       "",
			MaxCons:      1024,
			ReadTimeout:  100,
			WriteTimeout: 100,
			IdleTimeout:  100,
			Mode:         "release",
		},
		Port:     443,
		KeyFile:  stage.StoragePath() + "/certs/self-signed-key.pem",
		CertFile: stage.StoragePath() + "/certs/self-signed-cert.pem",
	}
	return opt
}

type HttpsService struct {
	*WebServer
	opt *HttpsOptions
}

func NewHttpsService(stage *owl.Stage, e *gin.Engine, opt *HttpsOptions) *HttpsService {

	if opt == nil {
		opt = NewDefaultHttpsOption(stage)
	}

	server := &HttpsService{
		opt:       opt,
		WebServer: NewWebServer(stage, e, opt.WebServerOptions),
	}

	go server.BlockRun()

	return server
}

func (i *HttpsService) BlockRun() {
	i.generateCerts()
	httpsPort := i.opt.Port
	key := i.opt.KeyFile
	cert := i.opt.CertFile

	// 设置 TLS 证书和密钥文件路径
	keyFile := key
	certFile := cert

	go func() {

		server, listener := i.getServerAndListener(httpsPort)
		err := server.ServeTLS(listener, certFile, keyFile) // 启动 HTTPS 服务器
		if err != nil {
			log.Fatal("Failed to start server: ", err)
		}
	}()
}

func (i *HttpsService) GetOptions() *HttpsOptions {
	return i.opt
}

func (i *HttpsService) generateCerts() {
	certsPath := i.stage.RuntimePath(owl.StoragePath + "/certs")
	os.Mkdir(certsPath, 0666)
	// 生成私钥
	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	// 生成自签名证书
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:       []string{"Lesent"},
			CommonName:         "乐诚",
			Province:           []string{"贵州"},
			OrganizationalUnit: []string{"乐诚技术"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour), // 有效期一年
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatal(err)
	}

	// 将私钥写入文件
	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		log.Fatal(err)
	}

	privFile, err := os.Create(certsPath + "/self-signed-key.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer privFile.Close()

	privBlock := &pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes}
	if err := pem.Encode(privFile, privBlock); err != nil {
		log.Fatal(err)
	}

	// 将证书写入文件
	certFile, err := os.Create(certsPath + "/self-signed-cert.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer certFile.Close()

	certBlock := &pem.Block{Type: "CERTIFICATE", Bytes: certDER}
	if err := pem.Encode(certFile, certBlock); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Self-signed certificate and private key generated successfully.")
}

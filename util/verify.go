package util

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"galaxy-guardian/logger"
	"io/ioutil"
	"os"
	"path"
	"time"
)

func RsaSignVerify(url, file, asc, agent string) error {
	agentPublic := path.Join(cachePath, agent, "public.pem")
	if _, err := os.Stat(agentPublic); os.IsNotExist(err) {
		logger.Logger.Info("need to download public.pem")
		err := DownloadPublicKey(url, agentPublic)
		if err != nil {
			return err
		}
	}
	for i := 0; i < 10; i++ {
		if i > 0 {
			_ = DownloadPublicKey(url, agentPublic)
			time.Sleep(time.Minute)
		}
		publicKey, _ := ioutil.ReadFile(agentPublic)
		signature, _ := ioutil.ReadFile(asc)
		tarFile, _ := ioutil.ReadFile(file)
		hashed := sha256.Sum256(tarFile)
		block, _ := pem.Decode(publicKey)
		if block == nil {
			logger.Logger.Error("public key error")
			continue
		}
		pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			logger.Logger.Error("parseKey Error")
			continue
		}
		pub := pubInterface.(*rsa.PublicKey)
		err = rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed[:], signature)
		if err != nil {
			logger.Logger.Error("VerifyPKCS1v15 Error")
			continue
		}
		return nil
	}
	return errors.New("failed to verify signature")
}

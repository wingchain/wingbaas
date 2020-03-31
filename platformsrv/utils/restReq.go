
package utils

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "crypto/tls"
    "sync"
    "github.com/wingbaas/platformsrv/logger"
)

const (
    REQ_GET    string = "GET" 
    REQ_POST   string = "POST"
) 

var keyPairLocker sync.Mutex //cert and key file locker

func RequestWithCert(reqUrl string,method string,certPath string,keyPath string) ([]byte,error) {
    keyPairLocker.Lock()
    cert, _ := tls.LoadX509KeyPair(certPath,keyPath)
    keyPairLocker.Unlock()  
    ssl := &tls.Config {
        Certificates: []tls.Certificate{cert},
        InsecureSkipVerify: true,
    }
    tr := &http.Transport{
        TLSClientConfig: ssl,
    }
    client := &http.Client{Transport: tr}
    req,err := http.NewRequest(method,reqUrl,nil)
    req.Close = true
    resp,err := client.Do(req)
    if err != nil {
        logger.Debug("RequestWithCert: request response error,url  " + reqUrl)
        return nil,fmt.Errorf("%v", err) 
    }
    defer resp.Body.Close()
    bodyText,err := ioutil.ReadAll(resp.Body)
    return bodyText,nil
}


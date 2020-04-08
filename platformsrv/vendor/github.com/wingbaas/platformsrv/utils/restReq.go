
package utils

import (
    "fmt"
    "strings"
    "io/ioutil"
    "net/http"
    "crypto/tls"
    "sync"
    "github.com/wingbaas/platformsrv/logger"
)

const (
    REQ_GET    string = "GET" 
    REQ_POST   string = "POST"
    REQ_DELETE string = "DELETE"
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
    req,err := http.NewRequest(method,reqUrl,nil)
    req.Close = true
    if err != nil {
        logger.Errorf("RequestWithCert: http NewRequest error,%v",err)
        return nil,fmt.Errorf("RequestWithCert: http NewRequest error,%v",err)
    }
    client := &http.Client{Transport: tr}
    resp,err := client.Do(req)
    defer resp.Body.Close()
    if err != nil {
        logger.Errorf("RequestWithCert: request response error,url=%s,err,%v",reqUrl,err)
        return nil,fmt.Errorf("RequestWithCert: request response error,url=%s,err,%v",reqUrl,err) 
    }
    bodyText,err := ioutil.ReadAll(resp.Body)
    if err != nil {
        logger.Errorf("RequestWithCert: read response body error,%v",err)
        return nil,fmt.Errorf("RequestWithCert: read response body error,%v",err)
    }
    return bodyText,nil
}

func RequestWithCertAndBody(reqUrl string,method string,certPath string,keyPath string,bodyStr string) ([]byte,error) {
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
    body := strings.NewReader(bodyStr)
    req,err := http.NewRequest(method,reqUrl,body)
    req.Close = true
    req.Header.Set("Content-Type", "application/yaml") 
    if err != nil {
        logger.Errorf("RequestWithCertAndBody: http NewRequest error,%v",err)
        return nil,fmt.Errorf("RequestWithCertAndBody: http NewRequest error,%v",err) 
    }
    client := &http.Client{Transport: tr}
    resp,err := client.Do(req)
    defer resp.Body.Close()
    if err != nil {
        logger.Errorf("RequestWithCertAndBody: request response error,url=%s,err,%v",reqUrl,err)
        return nil,fmt.Errorf("RequestWithCertAndBody: request response error,url=%s,err,%v",reqUrl,err) 
    }
    bodyText,err := ioutil.ReadAll(resp.Body)
    if err != nil {
        logger.Errorf("RequestWithCertAndBody: read response body error,%v",err)
        return nil,fmt.Errorf("RequestWithCertAndBody: read response body error,%v",err)
    }
    return bodyText,nil
}


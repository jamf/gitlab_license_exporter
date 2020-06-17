package main

import (
    "net/http"
	"log"
	"io/ioutil"
    "encoding/json"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "time"
    "errors"
    "os"
    "strconv"
    "strings"
)

var (
    expiration = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "gitlab_license_expires_at",
        Help: "Gitlab expiration day",
    })

    activeUsers = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "gitlab_active_users",
        Help: "Gitlab active users",
    })

    limitUsers = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "gitlab_user_limit",
        Help: "Users allowed by license",
    })

    scrapeSuccess = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "gitlab_scrape_success",
        Help: "Gitlab go exporter scrape status when try to read the API",
    })
)

type gitlab struct {
    ActiveUsers float64 `json:"active_users"`
    MaxUsers float64 `json:"user_limit"`
    ExpiresDate string `json:"expires_at"`
    ExpireSec float64
}

func main() {
    url, err := seturl()
    if err != nil{
        log.Fatal(err)
    }
    recordMetrics(url)
    http.Handle("/", promhttp.Handler())
    log.Fatal(http.ListenAndServe(":2222", nil))
}

func init() {
    prometheus.MustRegister(expiration)
    prometheus.MustRegister(activeUsers)
    prometheus.MustRegister(limitUsers)
    prometheus.MustRegister(scrapeSuccess)
}

func seturl() (string, error) {
    token := os.Getenv("TOKEN")
    if token == "" {
        return "", errors.New("Please set a token in env variable TOKEN")
    }

    var url string
    varURL := os.Getenv("URL")
    if varURL == "" {
        url = "http://gitlab-web:80"
    } else {
        url = strings.TrimRight(varURL, "/")
    }

    return string(url + "/api/v4/license?private_token=" + token), nil
}

func recordMetrics(url string) {
    go func() {
        for {
            body, err := getBody(url)
            if err != nil {
                log.Print(err)
                scrapeSuccess.Set(0)
            } else {
                glab, err := parseGitlab(body)
                if err != nil {
                    log.Print(err)
                    scrapeSuccess.Set(0)
                } else{
                    scrapeSuccess.Set(1)
                    expiration.Set(glab.ExpireSec)
                    activeUsers.Set(glab.ActiveUsers)
                    limitUsers.Set(glab.MaxUsers)
                    log.Print("Scrapped successfully")
                }
            }
            time.Sleep(60 * time.Minute)
        }
    }()
}

func getBody(url string) ([]byte, error){

    client := http.Client{
        Timeout: 30 * time.Second,
    }

    resp, err := client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        if resp.StatusCode == 401 {
            return nil, errors.New("401 unauthorised")
        } else if resp.StatusCode == 403 {
            return nil, errors.New("403 The request requires higher privileges than provided by the access token.")
        } else {
            return nil, errors.New("Unexpected error. HTTP status code: " + strconv.Itoa(resp.StatusCode))
        }
    }

    responseData, err := ioutil.ReadAll(resp.Body)
     if err != nil {
        return nil, err
    }

    return []byte(responseData), nil
}

func parseGitlab(textBytes []byte) (gitlab, error){
    const layoutISO = "2006-01-02"
    g := gitlab{}

    json.Unmarshal(textBytes, &g)

    t, err := time.Parse(layoutISO, g.ExpiresDate)
    if err != nil {
        return g, err
    }
    g.ExpireSec = float64(t.Unix())

    return g, nil
}
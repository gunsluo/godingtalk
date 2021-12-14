# Dingtalk
GO SDK for DingTalk, clone repo: https://github.com/icepy/go-dingtalk

## Example
```
package main

import (
  "os"
  "github.com/gunsluo/godingtalk"
)

func main() {
  c := getCompanyDingTalkClient()
  c.RefreshCompanyAccessToken()
}

func getCompanyDingTalkClient() *dingtalk.DingTalkClient {
  CorpID := os.Getenv("CorpId")
  CorpSecret := os.Getenv("CorpSecret")
  config := &dingtalk.DTConfig{
    CorpID:     CorpID,
    CorpSecret: CorpSecret,
  }
  c := dingtalk.NewDingTalkCompanyClient(config)
  return c
}
```

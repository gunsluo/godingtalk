# Dingtalk
GO SDK for DingTalk

## Example
```
package main

import (
  "github.com/gunsluo/godingtalk"
)

func main() {
	config := dingtalk.ISVConfig{
		CorpId: corpId,
		//CorpSecret: suiteSecret,
		SuiteKey:    suiteKey,
		SuiteSecret: suiteSecret,
	}
	client = dingtalk.NewISVClient(config)
	client.RefreshAndGetSuiteAccessToken()
}
```

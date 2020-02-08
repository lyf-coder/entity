# entity
Go json access like [viper](https://github.com/spf13/viper)

[![Actions](https://github.com/lyf-coder/entity/workflows/CI/badge.svg)](https://github.com/lyf-coder/entity)
[![GoDoc](https://godoc.org/github.com/lyf-coder/entity?status.svg)](https://godoc.org/github.com/lyf-coder/entity)
[![Go Report Card](https://goreportcard.com/badge/github.com/lyf-coder/entity)](https://goreportcard.com/report/github.com/lyf-coder/entity)

## Install

```console
go get github.com/lyf-coder/entity
```
## Usage
    import (
        "github.com/lyf-coder/entity"
    )
    // json string
    jsonStr := `{"IP": "127.0.0.1", "admin": {"name":"jack"}}`
    
    // json string to []byte
    b := []byte(jsonStr)
    
    // new entity
    entity := entity.NewByJSON(b)
    
    // Usage
    entity.GetString("IP")  // "127.0.0.1"
    entity.GetString("admin:name")  // "jack"

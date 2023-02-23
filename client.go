package gameoflife

import (
  "bytes"
  "encoding/json"
  "fmt"
  "net/http"
)

type Client struct {
  UpdateUrl string
  SyncUrl   string
  InitUrl   string
}

func (c *Client) SendChanges(changes []Change) {
  buf, err := json.Marshal(changes)
  if err != nil {
    fmt.Println("Error while create buffer from json:", err)
    return
  }

  req, err := http.NewRequest("POST", c.UpdateUrl, bytes.NewBuffer(buf))
  if err != nil {
    fmt.Println("Error creating request:", err)
    return
  }

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    fmt.Println("Error sending request:", err)
    return
  }
  defer resp.Body.Close()

  fmt.Println("Response status code:", resp.StatusCode)
}

func (c *Client) SendInit(init Init) {
  buf, err := json.Marshal(init)
  req, err := http.NewRequest("POST", c.InitUrl, bytes.NewBuffer(buf))
  if err != nil {
    fmt.Println("Error creating request:", err)
    return
  }

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    fmt.Println("Error sending request:", err)
    return
  }

  defer resp.Body.Close()

  fmt.Println("Response status code: ", resp.StatusCode)
}

func (c *Client) SendSync(sync Sync) {
  buf, err := json.Marshal(sync)
  if err != nil {
    fmt.Println("Error while create buffer from json:", err)
    return
  }

  req, err := http.NewRequest("POST", c.SyncUrl, bytes.NewBuffer(buf))
  if err != nil {
    fmt.Println("Error creating request:", err)
    return
  }

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    fmt.Println("Error sending request:", err)
    return
  }
  defer resp.Body.Close()

  fmt.Println("Response status code:", resp.StatusCode)
}

package model

import (
    "io/ioutil"
    "json"
    "os"
    "appengine"
    "appengine/memcache"
    "appengine/urlfetch"
    "bytes"
    "gob"
    "hvadsigeralex/config"
)

type Status struct {
    Id string
    Message string
}

type StatusList struct {
  Data []Status
}

func GetStatuses(c appengine.Context) ([]Status, os.Error) {
  var statusList []Status
  var err os.Error
  statusList, err = fetchMemcache(c)  
  if len(statusList) == 0 {
    statusList, err = fetchFacebookStatuses(c)
    updateMemcache(c, statusList)
    if err != nil {
      return nil, err
    }
  }
  return statusList, nil
}

func ForceUpdateStatuses(c appengine.Context) (os.Error) {
  statusList, err := fetchFacebookStatuses(c)
  updateMemcache(c, statusList)  
  return err
}

func fetchMemcache(c appengine.Context) ([]Status, os.Error) {
  // Get the item from the memcache
  item, memErr := memcache.Get(c, "statuses")
  if memErr != nil {
    c.Errorf("Error fetching item from memcache, %s", memErr)
    return nil, memErr
  }
  
  var data []byte = item.Value
  buffer := bytes.NewBuffer(data)
  dec := gob.NewDecoder(buffer)
      
  var statusList []Status
  dec.Decode(&statusList)
  
  return statusList, nil
}

func updateMemcache(c appengine.Context, statusList []Status) {
  c.Debugf("Updating cache")
  
  var buffer bytes.Buffer
  enc := gob.NewEncoder(&buffer)
  enc.Encode(statusList)
  
  var data []byte = buffer.Bytes()
  item := &memcache.Item{
      Key:   "statuses",
      Value: data,
  }
  if err := memcache.Set(c, item); err != nil {
      c.Errorf("Could not set item in memcache.")
  }
}

func fetchFacebookStatuses(c appengine.Context) ([]Status, os.Error) {
  graph_url := "https://graph.facebook.com/banckp/statuses" + "?limit=1000&access_token=" +config.AccessToken
  client := urlfetch.Client(c)
  response, err := client.Get(graph_url)

  if err != nil {
    c.Errorf("Error fetching item from facebook, %s", err)
    return nil, err
  }

  data, readErr := ioutil.ReadAll(response.Body);
  response.Body.Close();
  if readErr != nil {
    c.Errorf("Error reading bytes from response, %s", readErr)
    return nil, readErr
  }
  
  var m StatusList
  jsonErr := json.Unmarshal(data, &m)
  if jsonErr != nil {
    c.Errorf("Error unmarshalling json from facebook, %s", jsonErr)
    return nil, jsonErr
  }
  return m.Data, nil;
}
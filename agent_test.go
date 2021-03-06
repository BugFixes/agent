package agent_test

import (
  "database/sql"
  "fmt"
  "os"
  "testing"

  "github.com/bugfixes/agent"
  "github.com/joho/godotenv"
  _ "github.com/lib/pq"
  "github.com/stretchr/testify/assert"
)

var connectDetails = ""

func injectAgent(data agent.Data) error {
  db, err := sql.Open("postgres", connectDetails)
  if err != nil {
    return fmt.Errorf("injectAgent db.open: %w", err)
  }
  defer func() {
    err := db.Close()
    if err != nil {
      fmt.Printf("injectAgent db.close: %v", err)
    }
  }()
  _, err = db.Exec(
    "INSERT INTO agent (id, key, secret, company_id, name) VALUES ($1, $2, $3, $4, $5)",
    data.ID,
    data.Key,
    data.Secret,
    data.CompanyID,
    data.Name)
  if err != nil {
    return fmt.Errorf("injectAgent db.exec: %w", err)
  }

  return nil
}

func deleteAgent(id string) error {
  db, err := sql.Open("postgres", connectDetails)
  if err != nil {
    return fmt.Errorf("deleteAgent db.open: %w", err)
  }
  defer func() {
    err := db.Close()
    if err != nil {
      fmt.Printf("deleteAgent db.close: %v", err)
    }
  }()
  _, err = db.Exec("DELETE FROM agent WHERE id = $1", id)
  if err != nil {
    return fmt.Errorf("deleteAgent db.exec: %w", err)
  }

  return nil
}

func TestConnectDetails_FindAgentFromHeaders(t *testing.T) {
  if os.Getenv("GITHUB_ACTOR") == "" {
    err := godotenv.Load()
    if err != nil {
      t.Errorf("godotenv err: %w", err)
    }
  }

  connectDetails = fmt.Sprintf(
    "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
    os.Getenv("DB_HOSTNAME"),
    os.Getenv("DB_PORT"),
    os.Getenv("DB_USERNAME"),
    os.Getenv("DB_PASSWORD"),
    os.Getenv("DB_DATABASE"))

  tests := []struct {
    name string
    inject agent.Data
    headers map[string]string
    expect string
    err error
  }{
    {
      name: "valid headers",
      inject: agent.Data{
        ID:        "ad4b99e1-dec8-4682-862a-6b017e7c7c75",
        Key:       "94365b00-c6df-483f-804e-363312750500",
        Secret:    "f7356946-5814-4b5e-ad45-0348a89576ef",
        CompanyID: "b9e9153a-028c-4173-a7a8-e5063334416a",
        Name:      "valid headers",
      },
      headers: map[string]string{
          "x-api-key": "94365b00-c6df-483f-804e-363312750500",
          "x-api-secret": "f7356946-5814-4b5e-ad45-0348a89576ef",
      },
      expect: "ad4b99e1-dec8-4682-862a-6b017e7c7c75",
    },
    {
      name: "invalid secret",
      inject: agent.Data{
        ID:        "ad4b99e1-dec8-4682-862a-6b017e7c7c75",
        Key:       "94365b00-c6df-483f-804e-363312750500",
        Secret:    "f7356946-5814-4b5e-ad45-0348a89576ef",
        CompanyID: "b9e9153a-028c-4173-a7a8-e5063334416a",
        Name:      "invalid secret",
      },
      headers: map[string]string{
          "x-api-key": "94365b00-c6df-483f-804e-363312750500",
          "x-api-secret": "",
      },
      expect: "",
      err: fmt.Errorf("agent.FindAgentFromHeaders: no secret"),
    },
    {
      name: "invalid key",
      inject: agent.Data{
        ID:        "ad4b99e1-dec8-4682-862a-6b017e7c7c75",
        Key:       "94365b00-c6df-483f-804e-363312750500",
        Secret:    "f7356946-5814-4b5e-ad45-0348a89576ef",
        CompanyID: "b9e9153a-028c-4173-a7a8-e5063334416a",
        Name:      "invalid secret",
      },
      headers: map[string]string{
        "x-api-key": "",
        "x-api-secret": "f7356946-5814-4b5e-ad45-0348a89576ef",
      },
      expect: "",
      err: fmt.Errorf("agent.FindAgentFromHeaders: no key"),
    },
    {
      name: "no headers",
      headers: map[string]string{},
      expect: "",
      err: fmt.Errorf("agent.FindAgentFromHeaders: no agentID, agentKey, or agentSecret"),
    },
  }

  for _, test := range tests {
    if test.inject.ID != "" {
      err := injectAgent(test.inject)
      if err != nil {
        t.Errorf("inject err: %w", err)
      }
    }

    t.Run(test.name, func(t *testing.T) {
      c := agent.ConnectDetails{
        Full: connectDetails,
      }

      resp, err := c.FindAgentFromHeaders(test.headers)
      passed := assert.IsType(t, test.err, err)
      if !passed {
        t.Errorf("err failed: %w", err)
      }

      passed = assert.Equal(t, test.err, err)
      if !passed {
        t.Errorf("err equal failed: %+v, %v", test, err)
      }

      passed = assert.Equal(t, test.expect, resp)
      if !passed {
        t.Errorf("equal failed: %+v, %v", test, resp)
      }
    })

    if test.inject.ID != "" {
      err := deleteAgent(test.inject.ID)
      if err != nil {
        t.Errorf("delete err: %w", err)
      }
    }
  }
}

func TestConnectDetails_ValidateAgentID(t *testing.T) {
  if os.Getenv("GITHUB_ACTOR") == "" {
    err := godotenv.Load()
    if err != nil {
      t.Errorf("godotenv err: %w", err)
    }
  }

  connectDetails = fmt.Sprintf(
    "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
    os.Getenv("DB_HOSTNAME"),
    os.Getenv("DB_PORT"),
    os.Getenv("DB_USERNAME"),
    os.Getenv("DB_PASSWORD"),
    os.Getenv("DB_DATABASE"))

  tests := []struct {
    name    string
    request agent.Data
    expect  bool
    err     error
  }{
    {
      name: "agentid valid",
      request: agent.Data{
        ID:        "ad4b99e1-dec8-4682-862a-6b017e7c7c74",
        Key:       "94365b00-c6df-483f-804e-363312750500",
        Secret:    "f7356946-5814-4b5e-ad45-0348a89576ef",
        CompanyID: "b9e9153a-028c-4173-a7a8-e5063334416a",
        Name:      "bugfixes test frontend",
      },
      expect: true,
    },
    {
      name: "agentid invalid",
      request: agent.Data{
        ID:        "ad4b99e1-dec8-4682-862a-6b017e7c7c75",
        Key:       "94365b00-c6df-483f-804e-363312750500",
        Secret:    "f7356946-5814-4b5e-ad45-0348a89576ef",
        CompanyID: "b9e9153a-028c-4173-a7a8-e5063334416a",
        Name:      "bugfixes test frontend",
      },
      err: fmt.Errorf("AgentId no rows"),
    },
  }

  injErr := injectAgent(tests[0].request)
  if injErr != nil {
    t.Errorf("injection err: %w", injErr)
  }

  c := agent.ConnectDetails{
    Full:  connectDetails,
  }

  for _, test := range tests {
    t.Run(test.name, func(t *testing.T) {
      resp, err := c.ValidateAgentID(test.request.ID)
      passed := assert.IsType(t, test.err, err)
      if !passed {
        t.Errorf("agent err: %w", err)
      }
      passed = assert.IsType(t, test.expect, resp)
      if !passed {
        t.Errorf("agent type test failed: %+v", test.expect)
      }
      passed = assert.Equal(t, test.expect, resp)
      if !passed {
        t.Errorf("agent equal test failed: %+v, resp: %+v", test.expect, resp)
      }
    })
  }

  delErr := deleteAgent(tests[0].request.ID)
  if delErr != nil {
    t.Errorf("delete err: %w", delErr)
  }
}

func TestConnectDetails_LookupAgentID(t *testing.T) {
  if os.Getenv("GITHUB_ACTOR") == "" {
    err := godotenv.Load()
    if err != nil {
      t.Errorf("godotenv err: %w", err)
    }
  }

  connectDetails = fmt.Sprintf(
    "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
    os.Getenv("DB_HOSTNAME"),
    os.Getenv("DB_PORT"),
    os.Getenv("DB_USERNAME"),
    os.Getenv("DB_PASSWORD"),
    os.Getenv("DB_DATABASE"))

  tests := []struct {
    name    string
    request agent.Data
    expect  string
    err     error
  }{
    {
      name: "agentid found",
      request: agent.Data{
        ID:        "ad4b99e1-dec8-4682-862a-6b017e7c7c72",
        Key:       "94365b00-c6df-483f-804e-363312750500",
        Secret:    "f7356946-5814-4b5e-ad45-0348a89576ef",
        CompanyID: "b9e9153a-028c-4173-a7a8-e5063334416a",
        Name:      "bugfixes test frontend",
      },
      expect: "ad4b99e1-dec8-4682-862a-6b017e7c7c72",
    },
  }

  injErr := injectAgent(tests[0].request)
  if injErr != nil {
    t.Errorf("injection err: %w", injErr)
  }

  c := agent.ConnectDetails{
    Full:  connectDetails,
  }

  for _, test := range tests {
    t.Run(test.name, func(t *testing.T) {
      resp, err := c.LookupAgentID(test.request.Key, test.request.Secret)
      passed := assert.IsType(t, test.err, err)
      if !passed {
        t.Errorf("agent err: %w", err)
      }
      passed = assert.Equal(t, test.expect, resp)
      if !passed {
        t.Errorf("agent equal: %v, resp: %v", test.expect, resp)
      }
    })
  }

  delErr := deleteAgent(tests[0].request.ID)
  if delErr != nil {
    t.Errorf("delete err: %w", delErr)
  }
}

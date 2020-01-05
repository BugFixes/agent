package agent_test

import (
  "fmt"
  "os"
  "testing"

  "github.com/bugfixes/agent"
  "github.com/joho/godotenv"
  "github.com/stretchr/testify/assert"
)

func BenchmarkConnectDetails_FindAgentFromHeaders(b *testing.B) {
  b.ReportAllocs()

  if os.Getenv("GITHUB_ACTOR") == "" {
    err := godotenv.Load()
    if err != nil {
      b.Errorf("godotenv err: %w", err)
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

  b.ResetTimer()

  for _, test := range tests {
    if test.inject.ID != "" {
      err := injectAgent(test.inject)
      if err != nil {
        b.Errorf("inject err: %w", err)
      }
    }

    b.Run(test.name, func(t *testing.B) {
      b.StartTimer()
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

      b.StopTimer()
    })

    if test.inject.ID != "" {
      err := deleteAgent(test.inject.ID)
      if err != nil {
        b.Errorf("delete err: %w", err)
      }
    }
  }
}

func BenchmarkConnectDetails_LookupAgentID(b *testing.B) {
  b.ReportAllocs()

  if os.Getenv("GITHUB_ACTOR") == "" {
    err := godotenv.Load()
    if err != nil {
      b.Errorf("godotenv err: %w", err)
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
    b.Errorf("injection err: %w", injErr)
  }

  b.ResetTimer()

  c := agent.ConnectDetails{
    Full:  connectDetails,
  }

  for _, test := range tests {
    b.Run(test.name, func(t *testing.B) {
      t.StartTimer()
      resp, err := c.LookupAgentID(test.request.Key, test.request.Secret)
      passed := assert.IsType(t, test.err, err)
      if !passed {
        t.Errorf("validator err: %w", err)
      }
      passed = assert.Equal(t, test.expect, resp)
      if !passed {
        t.Errorf("validator equal: %v, resp: %v", test.expect, resp)
      }
      t.StopTimer()
    })
  }

  delErr := deleteAgent(tests[0].request.ID)
  if delErr != nil {
    b.Errorf("delete err: %w", delErr)
  }
}

func BenchmarkConnectDetails_ValidateAgentID(b *testing.B) {
  b.ReportAllocs()

  if os.Getenv("GITHUB_ACTOR") == "" {
    err := godotenv.Load()
    if err != nil {
      b.Errorf("godotenv err: %w", err)
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
    b.Errorf("injection err: %w", injErr)
  }

  b.ResetTimer()

  c := agent.ConnectDetails{
    Full:  connectDetails,
  }

  for _, test := range tests {
    b.Run(test.name, func(t *testing.B) {
      t.StartTimer()

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

      t.StopTimer()
    })
  }

  delErr := deleteAgent(tests[0].request.ID)
  if delErr != nil {
    b.Errorf("delete err: %w", delErr)
  }
}

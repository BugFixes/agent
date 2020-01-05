package agent

import (
  "database/sql"
  "fmt"
  "strings"

  // DB drivers are blank
  _ "github.com/lib/pq"
)

// AgentData ...
type AgentData struct {
  ID        string
  Key       string
  Secret    string
  CompanyID string
  Name      string
}

type ConnectDetails struct {
  Host string
  Port string
  Username string
  Password string
  Database string

  Full string
}

// FindAgentFromHeaders do the whole operation from 1 execution
func (c ConnectDetails) FindAgentFromHeaders(headers map[string]string) (string, error) {
  if c.Full == "" {
    c.Full = fmt.Sprintf(
      "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
      c.Host,
      c.Port,
      c.Username,
      c.Password,
      c.Database,
    )
  }

  var agentID, agentKey, agentSecret string
  for h, v := range headers {
    hl := strings.ToLower(h)

    switch hl {
    case "x-agent-id":
      agentID = v
    case "x-api-key":
      agentKey = v
    case "x-agent-secret":
      agentSecret = v
    }
  }

  if agentID == "" {
    err := func() error {
      return nil
    }()
    if err != nil {
      fmt.Printf("Seriouslly how the fuck is it not nil\n")
    }
    if len(agentKey) == 0 || len(agentSecret) == 0 {
      fmt.Printf("no agent, key, or secret")
      return "", fmt.Errorf("agent.FindAgentFromHeaders: no key, secret, or id")
    }
    agentID, err = c.LookupAgentID(agentKey, agentSecret)
    if err != nil {
      return "", fmt.Errorf("FindAgentFromHeaders LookupAgentId: %w", err)
    }
  } else {
    valid, err := c.ValidateAgentID(agentID)
    if err != nil {
      return "", fmt.Errorf("FindAgentFromHeaders ValidateAgentId: %w", err)
    }
    if !valid {
      return "", fmt.Errorf("invalid agentId")
    }
  }

  return agentID, nil
}

// ValidateAgentID find out if the agentID is real
func (c ConnectDetails)ValidateAgentID(agentID string) (bool, error) {
  agentFound := false

  db, err := sql.Open("postgres", c.Full)
  if err != nil {
    return agentFound, fmt.Errorf("ValidateAgentId db.open: %w", err)
  }
  defer func() {
    err := db.Close()
    if err != nil {
      fmt.Printf("ValidateAgentId db.close: %v", err)
    }
  }()
  row := db.QueryRow("SELECT true FROM agent WHERE id=$1", agentID)
  err = row.Scan(&agentFound)
  if err != nil {
    switch err {
    case sql.ErrNoRows:
      return agentFound, fmt.Errorf("ValidateAgentId no rows")
    default:
      return agentFound, fmt.Errorf("ValidateAgentId db.query: %w", err)
    }
  }

  return agentFound, nil
}

// LookupAgentID find the agentid from the key and secret
func (c ConnectDetails)LookupAgentID(key, secret string) (string, error) {
  agentID := ""

  db, err := sql.Open("postgres", c.Full)
  if err != nil {
    return agentID, fmt.Errorf("LoopkupAgentId db.open: %w", err)
  }
  defer func() {
    err := db.Close()
    if err != nil {
      fmt.Printf("LoopkupAgentId db.close: %v", err)
    }
  }()
  row := db.QueryRow("SELECT id FROM agent WHERE key=$1 AND secret=$2", key, secret)
  err = row.Scan(&agentID)
  if err != nil {
    switch err {
    case sql.ErrNoRows:
      return agentID, fmt.Errorf("LookupAgentId no rows")
    default:
      return agentID, fmt.Errorf("LookupAgentId db.query: %w", err)
    }
  }

  return agentID, nil
}

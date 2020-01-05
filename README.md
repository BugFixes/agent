# BugFix.es Agent

This is the agent library, just for splitting purposes so it can be used in multiple projects

## Example

```go
package main

import (
    "fmt"

    "github.com/bugfixes/agent"
)

func main() {
    agentID, err := agent.LookUpAgentID("tester", "tester")
    if err != nil {
        fmt.Printf("AgentLookup Failed: %+v\n", err)
    }

    fmt.Printf("AgentID: %+v\n", agentID)
}
```

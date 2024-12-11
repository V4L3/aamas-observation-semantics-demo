# Repository Overview

This repository contains three main services designed to interact in a hypermedia-based environment:

1. **Hypermedia Environment**: Provides the hypermedia environment with the artifacts.
2. **Proxy for WebSub**: Acts as a proxy for WebSub communication.
3. **Agent Application**: Represents an autonomous agent operating in the environment.

---

## Setup

### Step 1: Clone the Repository
Clone the repository recursively with submodules:
```bash
git clone --recursive <repository-url>
```
Then navigate to the repository:
```bash
cd aamas-observation-semantics-demo
```
### Step 2: Start Services
The **Hypermedia Environment** and the **Proxy** can be started using Docker Compose:
```bash
docker compose up --build
```

- The **Hypermedia Environment** will be available at: `http://localhost:8080`
- The **Proxy** will be available at: `http://localhost:3000`

---

## Running the Agent Application

Once the environment is ready:
1. Navigate to the `hyperagents` directory:
   ```bash
   cd hyperagents
   ```
2. Start the agent application:
   ```bash
   ./gradlew run examples:runBA
   ```

---

## Simulations

### Simulate Room Unlocking
You can simulate a room being unlocked with the following command:
```bash
curl --request POST \
  --url http://localhost:8080/workspaces/103/artifacts/r3/unlockRoom \
  --header 'X-Agent-WebID: http://localhost:8080/agents/aamasdemo'
```

### Simulate Fall Detection
You can simulate a fall being detected with the following command:
```bash
curl --request POST \
  --url http://localhost:8080/workspaces/103/artifacts/f3/triggerFallDetected \
  --header 'X-Agent-WebID: http://localhost:8080/agents/aamasdemo'
```

---

## Notes
- Ensure all services are running before executing the simulation commands.
- For additional configuration or debugging, consult the individual service directories.


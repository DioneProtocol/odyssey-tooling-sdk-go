# Odyssey Chain SDK Adaptation Plan

This document outlines the comprehensive plan for adapting the Avalanche Tooling SDK to work with Odyssey Chain (an Avalanche fork with different chain naming and IDs).

## Overview

The goal is to transform the Avalanche Tooling SDK into an Odyssey Chain Tooling SDK by:
1. Renaming all necessary folders and files
2. Adding necessary configurations (API endpoints, feature flags)
3. Updating dependencies to use Odyssey repositories
4. Renaming variables and constants to use Odyssey naming conventions
5. Updating all references throughout the codebase

## Step 1: Feature Flags System

### 1.1 Enhanced Feature Flags Configuration

**File: `constants/constants.go`**

**Enhanced Feature Flags for Odyssey Chain:**
```go
// Odyssey Chain Feature Flags
const (
    // Core Features (Always Enabled)
    OdysseyCoreEnabled = true
    
    // Network Features
    OdysseyTestnetEnabled = true
    OdysseyMainnetEnabled = true
    OdysseyDevnetEnabled  = false  // Disable by default
    
    // Node Management Features
    NodeCreationEnabled     = true
    NodeValidationEnabled   = true
    NodeMonitoringEnabled   = true
    NodeLoadTestEnabled     = false  // Disable by default
    
    // Cloud Infrastructure Features
    AWSCloudEnabled         = true
    GCPCloudEnabled         = true
    DockerCloudEnabled      = true
    
    // Subnet Features
    SubnetCreationEnabled   = true
    SubnetEVMEnabled        = true
    SubnetMultisigEnabled   = true
    
    // Interchain Features (Experimental)
    RelayerEnabled          = false  // Disable by default
    TeleporterEnabled       = false  // Disable by default
    InterchainMessengerEnabled = false  // Disable by default
    
    // Advanced Features
    LedgerSupportEnabled    = false  // Disable by default
    BLSMultiSigEnabled      = true
    CustomVMEnabled         = false  // Disable by default
    
    // Monitoring Features
    GrafanaEnabled          = true
    PrometheusEnabled       = true
    LokiEnabled             = true
    PromtailEnabled         = true
    
    // Development Features
    E2ETestingEnabled       = true
    LoadTestingEnabled      = false  // Disable by default
    DebugModeEnabled        = false  // Disable by default
)
```

### 1.2 Environment-Based Feature Configuration

**File: `constants/feature_config.go` (New File)**

```go
package constants

import (
    "os"
    "strconv"
    "strings"
)

// FeatureConfig holds runtime feature configuration
type FeatureConfig struct {
    Environment string
    Features    map[string]bool
}

// LoadFeatureConfig loads feature configuration from environment variables
func LoadFeatureConfig() *FeatureConfig {
    config := &FeatureConfig{
        Environment: getEnvOrDefault("ODYSSEY_ENV", "production"),
        Features:    make(map[string]bool),
    }
    
    // Load feature flags from environment variables
    featureFlags := []string{
        "ODYSSEY_CORE_ENABLED",
        "ODYSSEY_TESTNET_ENABLED", 
        "ODYSSEY_MAINNET_ENABLED",
        "ODYSSEY_DEVNET_ENABLED",
        "NODE_CREATION_ENABLED",
        "NODE_VALIDATION_ENABLED",
        "NODE_MONITORING_ENABLED",
        "NODE_LOADTEST_ENABLED",
        "AWS_CLOUD_ENABLED",
        "GCP_CLOUD_ENABLED",
        "DOCKER_CLOUD_ENABLED",
        "SUBNET_CREATION_ENABLED",
        "SUBNET_EVM_ENABLED",
        "SUBNET_MULTISIG_ENABLED",
        "RELAYER_ENABLED",
        "TELEPORTER_ENABLED",
        "INTERCHAIN_MESSENGER_ENABLED",
        "LEDGER_SUPPORT_ENABLED",
        "BLS_MULTISIG_ENABLED",
        "CUSTOM_VM_ENABLED",
        "GRAFANA_ENABLED",
        "PROMETHEUS_ENABLED",
        "LOKI_ENABLED",
        "PROMTAIL_ENABLED",
        "E2E_TESTING_ENABLED",
        "LOAD_TESTING_ENABLED",
        "DEBUG_MODE_ENABLED",
    }
    
    for _, flag := range featureFlags {
        featureName := strings.ToLower(strings.ReplaceAll(flag, "_", "-"))
        config.Features[featureName] = parseBoolEnv(flag, getDefaultValue(featureName))
    }
    
    return config
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func parseBoolEnv(key string, defaultValue bool) bool {
    if value := os.Getenv(key); value != "" {
        if boolValue, err := strconv.ParseBool(value); err == nil {
            return boolValue
        }
    }
    return defaultValue
}

func getDefaultValue(feature string) bool {
    // Return default values based on environment
    switch feature {
    case "odyssey-core", "odyssey-testnet", "odyssey-mainnet", "node-creation", 
         "node-validation", "node-monitoring", "aws-cloud", "gcp-cloud", 
         "docker-cloud", "subnet-creation", "subnet-evm", "subnet-multisig",
         "bls-multisig", "grafana", "prometheus", "loki", "promtail", "e2e-testing":
        return true
    case "odyssey-devnet", "node-loadtest", "relayer", "teleporter", 
         "interchain-messenger", "ledger-support", "custom-vm", "load-testing", "debug-mode":
        return false
    default:
        return false
    }
}
```

## Step 2: Rename Folders and Files

### 2.1 Directory Renaming
- `avalanche/` → `odyssey/`
- `node/config/templates/avalanche-*` → `node/config/templates/odyssey-*`
- `node/templates/avalanchego.*` → `node/templates/odysseygo.*`
- `node/monitoring/dashboards/avalanche-*` → `node/monitoring/dashboards/odyssey-*`

### 2.2 File Renaming
- `node/config/avalanche.go` → `node/config/odyssey.go`
- `avalanche/network.go` → `odyssey/network.go`
- `avalanche/vm.go` → `odyssey/vm.go`
- `avalanche/log.go` → `odyssey/log.go`
- `node/templates/avalanchego.docker-compose.yml` → `node/templates/odysseygo.docker-compose.yml`
- `node/config/templates/avalanche-node.tmpl` → `node/config/templates/odyssey-node.tmpl`
- `node/config/templates/avalanche-cchain.tmpl` → `node/config/templates/odyssey-dchain.tmpl`


## Step 3: Update Network Configuration

### 3.1 Update Network Configuration
**File: `odyssey/network.go` (Updated from avalanche/network.go)**

**Variable Renaming:**
```go
// Change from:
const (
    TestnetAPIEndpoint    = "https://api.avax-test.network"
    MainnetAPIEndpoint = "https://api.avax.network"
    FujiID            = 1
    MainnetID         = 1
    FujiHRP           = "fuji"
    MainnetHRP        = "avax"
)

// To:
const (
    TestnetAPIEndpoint    = "https://api.odyssey-test.network"  // Renamed from TestnetAPIEndpoint
    MainnetAPIEndpoint    = "https://api.odyssey.network"       // Updated endpoint
    OdysseyTestnetID      = 12345  // Renamed from FujiID, replace with actual ID
    OdysseyMainnetID      = 67890  // Renamed from MainnetID, replace with actual ID
    OdysseyTestnetHRP     = "odyssey-test"  // Renamed from FujiHRP
    OdysseyMainnetHRP     = "odyssey"       // Renamed from MainnetHRP
)
```

### 3.2 Update Network Functions
```go
// Change from:
func FujiNetwork() Network {
    return NewNetwork(Fuji, FujiID, TestnetAPIEndpoint)
}

func MainnetNetwork() Network {
    return NewNetwork(Mainnet, MainnetID, MainnetAPIEndpoint)
}

// To:
func OdysseyTestnetNetwork() Network {  // Renamed from FujiNetwork
    return NewNetwork(Testnet, OdysseyTestnetID, TestnetAPIEndpoint)  // Updated parameters
}

func OdysseyMainnetNetwork() Network {  // Renamed from MainnetNetwork
    return NewNetwork(Mainnet, OdysseyMainnetID, MainnetAPIEndpoint)  // Updated parameters
}
```

### 3.3 Update Network Names and Chain IDs
```go
// Update network names
func (nk NetworkKind) String() string {
    switch nk {
    case Mainnet:
        return "OdysseyMainnet"
    case Fuji:
        return "OdysseyTestnet"
    case Devnet:
        return "OdysseyDevnet"
    }
    return "invalid network"
}

// Update chain ID references
func (n Network) HRP() string {
    switch n.ID {
    case OdysseyTestnetID:      // Your testnet ID
        return OdysseyTestnetHRP
    case OdysseyMainnetID:   // Your mainnet ID
        return OdysseyMainnetHRP
    default:
        return OdysseyFallbackHRP
    }
}

func NetworkFromNetworkID(networkID uint32) Network {
    switch networkID {
    case OdysseyMainnetID:
        return OdysseyMainnetNetwork()
    case OdysseyTestnetID:
        return OdysseyTestnetNetwork()
    }
    return UndefinedNetwork
}
```

## Step 4: Update Dependencies

### 4.1 Update Module Name
**File: `go.mod`**

```go
// Change from:
module github.com/ava-labs/avalanche-tooling-sdk-go

// To:
module github.com/DioneProtocol/odyssey-tooling-sdk-go
```

### 4.2 Update External Dependencies (Odyssey repos)
```go
// Change from:
require (
    github.com/ava-labs/avalanchego vX.Y.Z
    github.com/ava-labs/coreth vX.Y.Z
    github.com/ava-labs/subnet-evm vX.Y.Z
    github.com/ava-labs/ledger-avalanche/go vX.Y.Z // indirect
)

// To (Odyssey forks):
require (
    github.com/DioneProtocol/odysseygo vX.Y.Z               // OdysseyGo (avalanchego fork)
    github.com/DioneProtocol/coreth vX.Y.Z                  // Coreth fork
    github.com/DioneProtocol/subnet-evm vX.Y.Z              // Subnet-EVM fork
    github.com/DioneProtocol/ledger-odyssey/go vX.Y.Z       // Ledger support (replaces ledger-avalanche)
    github.com/DioneProtocol/chains vX.Y.Z                  // Chain/params/configs (if used)
)

// Note:
// - If AWM Relayer/Teleporter are needed later, add forks or pin to ava-labs versions and keep
//   related feature flags disabled until repos exist.
```

### 4.3 Update Import Paths
```go
// Change all imports from:
"github.com/ava-labs/avalanche-tooling-sdk-go/avalanche/..."
"github.com/ava-labs/avalanchego/..."
"github.com/ava-labs/coreth/..."
"github.com/ava-labs/subnet-evm/..."
"github.com/ava-labs/ledger-avalanche/go"

// To:
"github.com/DioneProtocol/odyssey-tooling-sdk-go/odyssey/..."
"github.com/DioneProtocol/odysseygo/..."
"github.com/DioneProtocol/coreth/..."
"github.com/DioneProtocol/subnet-evm/..."
"github.com/DioneProtocol/ledger-odyssey/go"

// Optional/when available:
// "github.com/DioneProtocol/odyssey-awm-relayer/..."
// "github.com/DioneProtocol/odyssey-teleporter/..."
```

## Step 5: Update Constants and Configuration

### 5.1 Update Service Names and Paths
**File: `constants/constants.go`**

```go
// Change from:
const (
    ServiceAvalanchego = "avalanchego"
    CloudNodeCLIConfigBasePath = "/home/ubuntu/.avalanche-cli/"
    CloudNodeStakingPath       = "/home/ubuntu/.avalanchego/staking/"
    CloudNodeConfigPath        = "/home/ubuntu/.avalanchego/configs/"
    AvalancheGoDockerImage = "dionetech/avalanchego"
    AvalancheGoGitRepo     = "https://github.com/ava-labs/avalanchego"
    AvaLabsOrg      = "ava-labs"
)

// To:
const (
    ServiceOdysseygo = "odysseygo"  // Renamed from ServiceAvalanchego
    CloudNodeCLIConfigBasePath = "/home/ubuntu/.odyssey-cli/"      // Renamed from .avalanche-cli
    CloudNodeStakingPath       = "/home/ubuntu/.odysseygo/staking/" // Renamed from .avalanchego
    CloudNodeConfigPath        = "/home/ubuntu/.odysseygo/configs/" // Renamed from .avalanchego
    OdysseyGoDockerImage = "odysseylabs/odysseygo"  // Renamed from AvalancheGoDockerImage
    OdysseyGoGitRepo = "https://github.com/DioneProtocol/odysseygo"  // Renamed from AvalancheGoGitRepo
    OdysseyOrg = "DioneProtocol"  // Renamed from AvaLabsOrg
)
```

## Step 6: Update Template Files

### 6.1 Update Configuration Templates
**Files in `node/config/templates/` (Updated)**

**Template File Changes:**
- `avalanche-node.tmpl` → `odyssey-node.tmpl`
- `avalanche-cchain.tmpl` → `odyssey-dchain.tmpl`

**Content Updates:**
```json
// Change from:
{
    "log-dir": "/var/log/avalanchego",
    "db-dir": "/var/lib/avalanchego",
    "config-file": "/etc/avalanchego/config.json"
}

// To:
{
    "log-dir": "/var/log/odysseygo",
    "db-dir": "/var/lib/odysseygo", 
    "config-file": "/etc/odysseygo/config.json"
}
```

### 6.2 Update Docker Compose Templates
**File: `node/templates/odysseygo.docker-compose.yml` (Renamed from avalanchego.docker-compose.yml)**

```yaml
# Change from:
services:
  avalanchego:
    image: dionetech/avalanchego:{{ .AvalancheGoVersion }}
    container_name: avalanchego
    command: >
        ./avalanchego
        --config-file=/.avalanchego/configs/node.json

# To:
services:
  odysseygo:
    image: odysseylabs/odysseygo:{{ .OdysseyGoVersion }}
    container_name: odysseygo
    command: >
        ./odysseygo
        --config-file=/.odysseygo/configs/node.json
```

## Step 7: Update Package Names and Imports

### 7.1 Update Package Names
**Files requiring package name changes:**
- `avalanche/vm.go` → `odyssey/vm.go` (package avalanche → package odyssey)
- `avalanche/log.go` → `odyssey/log.go` (package avalanche → package odyssey)
- `avalanche/network.go` → `odyssey/network.go` (package avalanche → package odyssey)

### 7.2 Update Configuration Files
**File: `node/config/avalanche.go` → `node/config/odyssey.go`**

**Function and struct renames:**
```go
// Change from:
type AvalancheConfigInputs struct
func PrepareAvalancheConfig()
func RenderAvalancheTemplate()
func RenderAvalancheNodeConfig()
func RenderAvalancheCChainConfig()
func GetRemoteAvalancheNodeConfig()
func GetRemoteAvalancheCChainConfig()
func GetRemoteAvalancheGenesis()
func AvalancheFolderToCreate()

// To:
type OdysseyConfigInputs struct
func PrepareOdysseyConfig()
func RenderOdysseyTemplate()
func RenderOdysseyNodeConfig()
func RenderOdysseyCChainConfig()
func GetRemoteOdysseyNodeConfig()
func GetRemoteOdysseyCChainConfig()
func GetRemoteOdysseyGenesis()
func OdysseyFolderToCreate()
```

## Step 8: Update Cloud and Infrastructure

### 8.1 Update Cloud Configuration
**File: `cloud/aws/aws.go`**

**AWS-specific updates:**
```go
// Change from:
Value: aws.String("avalanche-tooling-sdk-node"),
Value: aws.String("avalanche-cli"),
descriptionFilterValue := fmt.Sprintf("Avalanche-CLI Ubuntu %s Docker", ubuntuVerLTS)

// To:
Value: aws.String("odyssey-tooling-sdk-node"),
Value: aws.String("odyssey-cli"),
descriptionFilterValue := fmt.Sprintf("Odyssey-CLI Ubuntu %s Docker", ubuntuVerLTS)
```

### 8.2 Update Monitoring Dashboards
**Files in `node/monitoring/dashboards/` (Updated)**

**Dashboard content updates:**
```json
// Change from:
"title": "Avalanche Node Dashboard",
"tags": ["avalanche", "blockchain", "monitoring"],
"expr": "avalanche_node_status"

// To:
"title": "Odyssey Node Dashboard",
"tags": ["odyssey", "blockchain", "monitoring"],
"expr": "odyssey_node_status"
```

### 8.3 Update Shell Scripts (content only; keep filenames)
**File: `node/shell/setupNode.sh`**

Update displayed messages and any references from Avalanche to Odyssey, but keep the filename:
- "Setting up Avalanche Node..." → "Setting up Odyssey Node..."
- Any path references from `/.avalanchego` → `/.odysseygo` (if you decide to rename binary/config paths)

## Step 9: Update Application Logic

### 9.1 Update Keychain and Wallet References
**Files: `keychain/keychain.go`, `wallet/wallet.go`**

**Variable and comment updates:**
```go
// Change from:
// To view Ledger addresses and their balances, you can use Avalanche CLI
// RequiredFunds is the minimum total AVAX that the selected addresses from Ledger should contain
config.AVAXKeychain

// To:
// To view Ledger addresses and their balances, you can use Odyssey CLI
// RequiredFunds is the minimum total ODYSSEY that the selected addresses from Ledger should contain
config.ODYSSEYKeychain
```

### 9.2 Update Validator and Staking References
**Files: `validator/validator.go`, `node/staking.go`**

**Comment and variable updates:**
```go
// Change from:
// StakeAmount is the amount of Avalanche tokens (AVAX) to stake in this validator
// For more information on delegation fee, please head to https://docs.avax.network/...

// To:
// StakeAmount is the amount of Odyssey tokens (ODYSSEY) to stake in this validator
// For more information on delegation fee, please head to https://docs.odyssey.network/...
```

### 9.3 Update Interchain and Relayer References
**Files: `interchain/relayer/conf.go`, `interchain/interchainmessenger/`**

**Function parameter updates:**
```go
// Change from:
network avalanche.Network,

// To:
network odyssey.Network,
```

### 9.4 Update EVM and Contract References
**Files: `evm/evm.go`, `evm/contract.go`**

**Import and comment updates:**
```go
// Change from:
"github.com/ava-labs/subnet-evm/accounts/abi/bind"
"github.com/ava-labs/subnet-evm/core/types"

// To:
"github.com/DioneProtocol/subnet-evm/accounts/abi/bind"
"github.com/DioneProtocol/subnet-evm/core/types"
```

## Step 10: Update Documentation and Examples

### 10.1 Update README.md
**File: `README.md` (Updated)**

```markdown
# Change from:
# Avalanche Tooling Go SDK
# The official Avalanche Tooling Go SDK library.

# To:
# Odyssey Chain Tooling Go SDK
# The official Odyssey Chain Tooling Go SDK library.

# Update all references:
# - "Avalanche" → "Odyssey Chain"
# - "avalanchego" → "odysseygo"
# - "Fuji" → "Testnet"
# - "avax" → "odyssey"
# - Update URLs and endpoints
```

### 10.2 Update Example Files
**Files in `examples/` directory (Updated)**

**Function Call Updates:**
```go
// Change from:
network := avalanche.FujiNetwork()
avalanchegoVersion := "v1.11.8"

// To:
network := odyssey.OdysseyTestnetNetwork()
odysseygoVersion := "v1.11.8"
```

### 10.3 Update Comments and Documentation
**Throughout the codebase:**

- Replace "Avalanche" with "Odyssey Chain"
- Replace "avalanchego" with "odysseygo"
- Replace "Fuji" with "Testnet"
- Replace "avax" with "odyssey"
- Update all URLs and endpoints
- Update all configuration examples

## Step 11: Update License Headers and URLs

### 11.1 Update License Headers
**All Go files throughout the codebase:**

```go
// Change from:
// Copyright (C) 2024, Ava Labs, Inc. All rights reserved.

// To:
// Copyright (C) 2024, Dione Protocol, Inc. All rights reserved.
```

### 11.2 Update Documentation URLs
**Throughout the codebase:**

```go
// Change from:
// https://docs.avax.network/...
// https://github.com/ava-labs/avalanchego/...

// To:
// https://docs.odyssey.network/...
// https://github.com/DioneProtocol/odysseygo/...
```

## Step 12: Feature Flag Integration

### 12.1 Feature Flag Manager
**File: `utils/feature_flags.go` (New File)**

```go
package utils

import (
    "fmt"
    "github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
)

// FeatureFlagManager manages feature flags throughout the application
type FeatureFlagManager struct {
    config *constants.FeatureConfig
}

// NewFeatureFlagManager creates a new feature flag manager
func NewFeatureFlagManager() *FeatureFlagManager {
    return &FeatureFlagManager{
        config: constants.LoadFeatureConfig(),
    }
}

// IsEnabled checks if a specific feature is enabled
func (ffm *FeatureFlagManager) IsEnabled(feature string) bool {
    if enabled, exists := ffm.config.Features[feature]; exists {
        return enabled
    }
    return constants.IsFeatureEnabled(feature)
}

// ValidateFeatures validates a combination of features
func (ffm *FeatureFlagManager) ValidateFeatures(features []string) error {
    for _, feature := range features {
        if !ffm.IsEnabled(feature) {
            return fmt.Errorf("feature '%s' is disabled", feature)
        }
    }
    return constants.ValidateFeatureCombination(features)
}

// GetEnabledFeatures returns all enabled features
func (ffm *FeatureFlagManager) GetEnabledFeatures() []string {
    var enabled []string
    for feature, enabled := range ffm.config.Features {
        if enabled {
            enabled = append(enabled, feature)
        }
    }
    return enabled
}

// GetEnvironment returns the current environment
func (ffm *FeatureFlagManager) GetEnvironment() string {
    return ffm.config.Environment
}
```

### 12.2 Integration Examples
**Integration in Node Creation:**
```go
// Example: Feature flag integration in node creation
func CreateNodesWithFeatureFlags(ctx context.Context, params *NodeParams) ([]*Node, error) {
    ffm := utils.NewFeatureFlagManager()
    
    // Validate required features
    requiredFeatures := []string{"node-creation"}
    if params.UseAWS {
        requiredFeatures = append(requiredFeatures, "aws-cloud")
    }
    if params.UseGCP {
        requiredFeatures = append(requiredFeatures, "gcp-cloud")
    }
    if params.UseDocker {
        requiredFeatures = append(requiredFeatures, "docker-cloud")
    }
    
    if err := ffm.ValidateFeatures(requiredFeatures); err != nil {
        return nil, fmt.Errorf("feature validation failed: %w", err)
    }
    
    // Proceed with node creation
    return CreateNodes(ctx, params)
}
```

## Step 13: Implementation Strategy

### 13.1 Repository Setup
```bash
# 1. Fork the repository
git clone https://github.com/DioneProtocol/odyssey-tooling-sdk-go.git
cd odyssey-tooling-sdk-go
```

### 13.2 Core Changes (No scripts, actionable checklist)
- Update module name in `go.mod`:
  - From: `module github.com/ava-labs/avalanche-tooling-sdk-go`
  - To:   `module github.com/DioneProtocol/odyssey-tooling-sdk-go`
- Update all import paths across `.go` files:
  - `github.com/ava-labs/avalanche-tooling-sdk-go/...` → `github.com/DioneProtocol/odyssey-tooling-sdk-go/...`
  - `github.com/ava-labs/avalanchego/...` → `github.com/DioneProtocol/odysseygo/...`
  - `github.com/ava-labs/coreth/...` → `github.com/DioneProtocol/coreth/...`
  - `github.com/ava-labs/subnet-evm/...` → `github.com/DioneProtocol/subnet-evm/...`
- Replace remaining `Ava Labs` organization and repo references where applicable.

### 13.3 File and Folder Renaming (No scripts, clear actions)
- Rename package directory: `avalanche/` → `odyssey/`
- Rename templates:
  - `node/config/templates/avalanche-*` → `node/config/templates/odyssey-*`
- Rename docker-compose files:
  - `node/templates/avalanchego.*` → `node/templates/odysseygo.*`
- Update dashboards naming under `node/monitoring/dashboards/` from `avalanche-*` to `odyssey-*` where applicable.
- Keep example filenames and shell scripts; only update content and imports.
- Update constants and configuration (IDs, HRP, endpoints) as specified in Steps 3 and 5.

### 13.4 Post-Renaming Validation (Checklist)
- Verify `odyssey/` directory exists and builds
- Search for remaining `avalanche` references (case-sensitive) in code, docs, configs
- Validate import paths resolve
- Run `go mod tidy` to sync dependencies
- Run tests: `go test ./...`
- Build verification: `go build ./...`
- Manually verify updated directories: `odyssey/`, `node/config/templates/`, `node/templates/`

## Step 14: Environment Configuration Examples

### 14.1 Development Environment
```bash
# Enable all features for development
export ODYSSEY_ENV=development
export ODYSSEY_DEVNET_ENABLED=true
export NODE_LOADTEST_ENABLED=true
export RELAYER_ENABLED=true
export TELEPORTER_ENABLED=true
export INTERCHAIN_MESSENGER_ENABLED=true
export LEDGER_SUPPORT_ENABLED=true
export CUSTOM_VM_ENABLED=true
export LOAD_TESTING_ENABLED=true
export DEBUG_MODE_ENABLED=true
```

### 14.2 Production Environment
```bash
# Enable only stable features for production
export ODYSSEY_ENV=production
export ODYSSEY_DEVNET_ENABLED=false
export NODE_LOADTEST_ENABLED=false
export RELAYER_ENABLED=false
export TELEPORTER_ENABLED=false
export INTERCHAIN_MESSENGER_ENABLED=false
export LEDGER_SUPPORT_ENABLED=false
export CUSTOM_VM_ENABLED=false
export LOAD_TESTING_ENABLED=false
export DEBUG_MODE_ENABLED=false
```

## Step 15: Testing and Validation

### 15.1 Network Tests
- [ ] Network connection to Odyssey testnet
- [ ] Network connection to Odyssey mainnet
- [ ] Chain ID validation
- [ ] HRP (address format) validation

### 15.2 Functionality Tests
- [ ] Node creation and deployment
- [ ] Subnet creation and management
- [ ] Validator operations
- [ ] Wallet and keychain operations
- [ ] Transaction building and signing

### 15.3 Renaming Validation Tests
- [ ] All directory names updated to Odyssey naming
- [ ] All file names updated to Odyssey naming
- [ ] All import paths updated to Odyssey references
- [ ] All package declarations updated to Odyssey
- [ ] All function calls updated to Odyssey naming
- [ ] All configuration files updated to Odyssey references
- [ ] All template files updated to Odyssey naming
- [ ] All documentation updated to Odyssey references
- [ ] All shell scripts updated to Odyssey naming
- [ ] All Docker configurations updated to Odyssey naming
- [ ] No remaining Avalanche references in codebase
- [ ] No remaining Avalanche references in documentation
- [ ] No remaining Avalanche references in configuration files

### 15.4 Integration Tests
- [ ] Docker container deployment
- [ ] Configuration file generation
- [ ] API endpoint connectivity
- [ ] Genesis file loading

### 15.5 Feature Flag Tests
- [ ] Feature flag validation and combination checking
- [ ] Environment-based feature configuration
- [ ] Feature flag integration in node creation
- [ ] Feature flag integration in subnet creation
- [ ] Feature flag integration in interchain features
- [ ] Feature flag integration in monitoring features
- [ ] Feature flag integration in cloud providers
- [ ] Feature flag integration in development tools

## Critical Considerations

### 1. File Renaming Considerations
- Handle case-sensitive file systems (Linux/macOS) that may have issues with case-only renames
- Use git mv commands to preserve git history while renaming files and directories
- Create backup before renaming using cp or git commit/tag
- Handle files with special characters or spaces in names during renaming
- Update any configuration files that reference old paths
- Update any environment variable references from AVALANCHE_ to ODYSSEY_

### 2. Chain ID Conflicts
- Ensure Odyssey Chain IDs don't conflict with existing Avalanche IDs
- Verify network ID uniqueness across the ecosystem
- Test chain ID validation in your forked dependencies

### 3. Genesis Compatibility
- Your genesis parameters must be compatible with the forked avalanchego
- Ensure all genesis fields are properly configured for Odyssey Chain
- Test genesis loading and validation

### 4. API Compatibility
- Odyssey Chain API endpoints must match the expected format
- Verify all API calls work with your network endpoints
- Test network connectivity and API responses

### 5. Dependency Compatibility
- All forked dependencies must be compatible with each other
- Ensure version compatibility between odysseygo, coreth, and subnet-evm
- Test integration between all components

### 6. Binary Names
- Decide whether to rename binaries (avalanchego → odysseygo)
- Update all references consistently throughout the codebase
- Ensure Docker images and deployment scripts are updated

## Success Criteria

The adaptation is successful when:

1. ✅ All tests pass with Odyssey Chain configuration
2. ✅ Network operations work with Odyssey endpoints
3. ✅ Node creation and management functions correctly
4. ✅ Subnet operations work with Odyssey Chain
5. ✅ Documentation is updated and accurate
6. ✅ Dependencies are properly forked and compatible

## Timeline Estimate

- **Steps 1-3 (Feature Flags & Network Config)**: 1-2 days
- **Steps 4-6 (Dependencies & Templates)**: 2-3 days
- **Steps 7-9 (Package Names & Application Logic)**: 2-3 days
- **Steps 10-11 (Documentation & License)**: 1 day
- **Steps 12-13 (Feature Flags & Implementation)**: 2-3 days
- **Steps 14-15 (Environment & Testing)**: 2-3 days

**Total Estimated Time**: 10-15 days

## Next Steps

1. **Fork Dependencies**: Create Odyssey versions of avalanchego, coreth, and subnet-evm
2. **Set Chain IDs**: Define unique chain IDs for Odyssey networks
3. **Configure Endpoints**: Set up API endpoints for Odyssey networks
4. **Begin Implementation**: Start with Step 1 changes
5. **Test Incrementally**: Test each step before proceeding to the next


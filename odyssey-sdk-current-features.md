# Odyssey Tooling Go SDK - Supported Features

## Core Features (Currently Supported - v0.3.0)

### 1. Subnet Management
- Subnet Creation: Create subnets on Odyssey testnet and Mainnet
- Blockchain Creation: Deploy blockchains within subnets
- Subnet Genesis: Generate and configure subnet genesis files
- Subnet EVM Support: Full Subnet-EVM integration with customizable parameters
- Multisig Subnet Control: Multi-signature control keys for subnet management
- Subnet Validator Management: Add validators to existing subnets

### 2. Node Management
- Node Creation: Deploy Odyssey nodes on cloud platforms
- Node Types Supported:
  - Validator Nodes: For validating Primary Network and Subnets
  - API Nodes: For providing API access to the network
  - Monitoring Nodes: Centralized monitoring with Grafana dashboards
  - Load Test Nodes: For performance testing
### 3. Primary Network Validation
- Validator Staking: Enable nodes to validate the Primary Network
- Stake Management: Configure staking amounts and durations
- BLS Key Management: Support for BLS multi-signatures

### 4. Key Management & Security
- Soft Key Support: Stored private keys for transaction signing
- Ledger Integration: Hardware wallet support (coming soon)
- Keychain Management: Secure key storage and management
- Multi-signature Support: Threshold-based transaction signing

### 5. Wallet & Transaction Management
- Wallet Creation: Multi-chain wallet support (O-Chain, D-Chain, A-Chain)
- Transaction Building: Create and sign various transaction types
- Fee Management: Automatic fee calculation and payment
- Change Address Management: Secure change UTXO handling

### 6. EVM Integration
- Smart Contract Deployment: Deploy and interact with EVM contracts
- Contract Interaction: Call contract functions and read state
- Gas Management: Configure gas limits and fee structures
- Precompiles Support: Access to Odyssey-specific precompiles

### 7. Monitoring & Observability
- Grafana Dashboards: Pre-configured monitoring dashboards
- Prometheus Metrics: Comprehensive metrics collection
- Loki Logging: Centralized log aggregation
- Promtail Configuration: Log shipping and processing
- Custom Dashboards: Support for D-Chain, O-Chain, A-Chain, and Subnet metrics

### 8. Network Support
- Odyssey Testnet: Full testnet support
- Mainnet: Production network support
- Custom Networks: Support for custom network configurations
- Network Switching: Easy switching between networks

### 9. Development & Testing
- E2E Testing: End-to-end testing framework
- Load Testing: Performance testing capabilities
- Docker Compose: Local development environment setup
- Configuration Templates: Pre-built configuration templates

## Future Features (Planned)
- Devnet Support: Additional development network features
- Custom Subnets: Enhanced custom subnet capabilities

## Architecture Support
- x86_64: Full support for x86_64 architecture
- ARM64: Support for ARM64 architecture (Apple Silicon)

## Dependencies & Integrations
- OdysseyGo: Core Odyssey node software
- Subnet-EVM: EVM-compatible subnet implementation
- Coreth: D-Chain implementation

## Summary
This SDK provides a comprehensive toolkit for building, deploying, and managing Odyssey infrastructure, from simple node deployment to complex multi-chain applications.

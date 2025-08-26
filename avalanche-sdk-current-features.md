# Avalanche Tooling Go SDK - Supported Features

## Core Features (Currently Supported - v0.3.0)

### 1. Subnet Management
- Subnet Creation: Create subnets on Fuji testnet and Mainnet
- Blockchain Creation: Deploy blockchains within subnets
- Subnet Genesis: Generate and configure subnet genesis files
- Subnet EVM Support: Full Subnet-EVM integration with customizable parameters
- Multisig Subnet Control: Multi-signature control keys for subnet management
- Subnet Validator Management: Add validators to existing subnets

### 2. Node Management
- Node Creation: Deploy Avalanche nodes on cloud platforms
- Node Types Supported:
  - Validator Nodes: For validating Primary Network and Subnets
  - API Nodes: For providing API access to the network
  - Monitoring Nodes: Centralized monitoring with Grafana dashboards
  - Load Test Nodes: For performance testing
  - AWM Relayer Nodes: For cross-chain messaging (experimental)

### 3. Cloud Infrastructure
- AWS Integration: Full AWS EC2 support for node deployment
- GCP Integration: Google Cloud Platform support
- Docker Support: Local development and testing with Docker containers
- Instance Management: Automatic instance provisioning and configuration
- Security Groups: Automated security group creation and management
- SSH Key Management: Secure SSH access to deployed nodes

### 4. Primary Network Validation
- Validator Staking: Enable nodes to validate the Primary Network
- Stake Management: Configure staking amounts and durations
- BLS Key Management: Support for BLS multi-signatures

### 5. Key Management & Security
- Soft Key Support: Stored private keys for transaction signing
- Ledger Integration: Hardware wallet support (coming soon)
- Keychain Management: Secure key storage and management
- Multi-signature Support: Threshold-based transaction signing

### 6. Wallet & Transaction Management
- Wallet Creation: Multi-chain wallet support (P-Chain, C-Chain, X-Chain)
- Transaction Building: Create and sign various transaction types
- Fee Management: Automatic fee calculation and payment
- Change Address Management: Secure change UTXO handling

### 7. EVM Integration
- Smart Contract Deployment: Deploy and interact with EVM contracts
- Contract Interaction: Call contract functions and read state
- Gas Management: Configure gas limits and fee structures
- Precompiles Support: Access to Avalanche-specific precompiles

### 8. Monitoring & Observability
- Grafana Dashboards: Pre-configured monitoring dashboards
- Prometheus Metrics: Comprehensive metrics collection
- Loki Logging: Centralized log aggregation
- Promtail Configuration: Log shipping and processing
- Custom Dashboards: Support for C-Chain, P-Chain, X-Chain, and Subnet metrics

### 9. Interchain Features (Experimental)
- Teleporter Integration: Cross-chain messaging infrastructure
- Interchain Messenger (ICM): Deploy and configure ICM contracts
- AWM Relayer: Cross-chain message relaying
- Multi-chain Support: Connect multiple blockchains

### 10. Network Support
- Fuji Testnet: Full testnet support
- Mainnet: Production network support
- Custom Networks: Support for custom network configurations
- Network Switching: Easy switching between networks

### 11. Development & Testing
- E2E Testing: End-to-end testing framework
- Load Testing: Performance testing capabilities
- Docker Compose: Local development environment setup
- Configuration Templates: Pre-built configuration templates

## Future Features (Planned)
- Devnet Support: Additional development network features
- Custom Subnets: Enhanced custom subnet capabilities
- Enhanced Teleporter SDK: More comprehensive cross-chain features
- Additional Cloud Providers: Support for more cloud platforms

## Architecture Support
- x86_64: Full support for x86_64 architecture
- ARM64: Support for ARM64 architecture (Apple Silicon, AWS Graviton)

## Dependencies & Integrations
- AvalancheGo: Core Avalanche node software
- Subnet-EVM: EVM-compatible subnet implementation
- Coreth: C-Chain implementation
- AWM Relayer: Cross-chain messaging relay
- Teleporter: Cross-chain messaging protocol

## Summary
This SDK provides a comprehensive toolkit for building, deploying, and managing Avalanche infrastructure, from simple node deployment to complex multi-chain applications with cross-chain messaging capabilities.

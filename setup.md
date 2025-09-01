# Avalanche Tooling SDK Setup Guide

This guide will walk you through installing Go 1.21.11 on WSL and setting up the Avalanche Tooling SDK in a hybrid Windows 11 + WSL environment.

## Prerequisites

- Windows 11 with WSL2 installed
- Ubuntu or similar Linux distribution in WSL
- Internet connection for downloading Go and dependencies
- Project repository stored on Windows 11 filesystem (accessible via `/mnt/c/` or `/mnt/d/` in WSL)

## Step 1: Install Go 1.21.11 on WSL

### 1.1 Remove any existing Go installations

First, check if Go is already installed and remove any existing versions:

```bash
# Check current Go installation
go version

# Remove existing Go installations
sudo rm -rf /usr/local/go
sudo rm -rf /usr/lib/go-*
```

### 1.2 Download and install Go 1.21.11

```bash
# Navigate to temporary directory
cd /tmp

# Download Go 1.21.11 for Linux AMD64
wget https://go.dev/dl/go1.21.11.linux-amd64.tar.gz

# Extract to /usr/local
sudo tar -C /usr/local -xzf go1.21.11.linux-amd64.tar.gz

# Verify installation
/usr/local/go/bin/go version
```

Expected output:
```
go version go1.21.11 linux/amd64
```

### 1.3 Configure PATH

Add Go to your PATH by editing your shell configuration:

```bash
# Add Go to PATH in bashrc
echo "export PATH=/usr/local/go/bin:\$PATH" >> ~/.bashrc

# Reload shell configuration
source ~/.bashrc

# Verify Go is in PATH
go version
```

## Step 2: Access the Avalanche Tooling SDK

### Option A: If you already have the repository on Windows 11

```bash
# Navigate to your Windows 11 project directory
cd /mnt/d/Projects/Dione\ Protocol/odyssey-tooling-sdk-go/avalanche-tooling-sdk-go

# Verify you're in the correct directory
ls -la
```

### Option B: If you need to clone the repository

```bash
# Navigate to your desired Windows 11 directory
cd /mnt/d/Projects/Dione\ Protocol/

# Clone the repository
git clone https://github.com/ava-labs/avalanche-tooling-sdk-go.git
cd avalanche-tooling-sdk-go
```

### Important Notes for Windows 11 + WSL Setup

- The project is stored on Windows 11 filesystem (accessible via `/mnt/d/` in WSL)
- All Go operations (build, test, run) are executed in WSL
- File I/O performance may be slower when accessing Windows filesystem from WSL
- Consider using WSL filesystem for better performance if working extensively with the project

## Step 3: Set up the SDK

### 3.1 Clear any existing Go cache

```bash
# Clear all Go caches
go clean -cache -modcache -testcache
```

### 3.2 Download dependencies

```bash
# Download all module dependencies
go mod download
```

### 3.3 Verify module integrity

```bash
# Verify all modules are correct
go mod verify
```

Expected output:
```
all modules verified
```

## Step 4: Run Tests

### 4.1 Run all tests

```bash
# Run all tests in the project
go test ./...
```

### 4.2 Run specific test suites

```bash
# Test key management functionality
go test ./key -v

# Test EVM functionality
go test ./evm -v

# Test utilities
go test ./utils -v

# Test AWS cloud integration
go test ./cloud/aws -v
```

## Expected Test Results

When running `go test ./...`, you should see output similar to:

```
?       github.com/ava-labs/avalanche-tooling-sdk-go/avalanche  [no test files]
?       github.com/ava-labs/avalanche-tooling-sdk-go/cloud/gcp  [no test files]
?       github.com/ava-labs/avalanche-tooling-sdk-go/constants  [no test files]
ok      github.com/ava-labs/avalanche-tooling-sdk-go/cloud/aws  0.010s
ok      github.com/ava-labs/avalanche-tooling-sdk-go/evm        0.102s
?       github.com/ava-labs/avalanche-tooling-sdk-go/examples   [no test files]
ok      github.com/ava-labs/avalanche-tooling-sdk-go/key        0.112s
ok      github.com/ava-labs/avalanche-tooling-sdk-go/utils      6.515s
```

### Note about failing tests

Some tests may fail due to missing configuration or external dependencies:

- **Node tests**: May fail due to missing AWS credentials or configuration
- **Subnet tests**: May fail due to insufficient funds or missing genesis files
- **Ledger tests**: May fail if no Ledger device is connected

These failures are expected in a development environment without proper configuration.

## Step 5: Verify Installation

### 5.1 Check Go environment

```bash
# Check Go environment variables
go env

# Check Go version
go version
```

### 5.2 Check module status

```bash
# Check module status
go mod tidy

# List all dependencies
go list -m all
```

## Troubleshooting

### Go command not found

If you get "go: command not found", ensure the PATH is set correctly:

```bash
# Check if Go is in PATH
which go

# If not found, manually set PATH
export PATH=/usr/local/go/bin:$PATH

# Add to bashrc permanently
echo "export PATH=/usr/local/go/bin:\$PATH" >> ~/.bashrc
source ~/.bashrc
```

### Windows 11 + WSL specific issues

#### File permission issues

When working with files on Windows 11 filesystem from WSL, you might encounter permission issues:

```bash
# Check file permissions
ls -la

# If you see permission issues, you may need to adjust Windows file permissions
# or work from WSL filesystem instead
```

#### Performance issues with Windows filesystem

If you experience slow performance when accessing files on Windows 11 filesystem:

```bash
# Consider copying project to WSL filesystem for better performance
cp -r /mnt/d/Projects/Dione\ Protocol/odyssey-tooling-sdk-go ~/avalanche-tooling-sdk-go
cd ~/avalanche-tooling-sdk-go

# Or use WSL filesystem for development
mkdir -p ~/projects
cd ~/projects
git clone https://github.com/ava-labs/avalanche-tooling-sdk-go.git
cd avalanche-tooling-sdk-go
```

#### Path issues with spaces

If your Windows path contains spaces, use proper escaping:

```bash
# Correct way to navigate paths with spaces
cd "/mnt/d/Projects/Dione Protocol/odyssey-tooling-sdk-go/avalanche-tooling-sdk-go"

# Or escape spaces
cd /mnt/d/Projects/Dione\ Protocol/odyssey-tooling-sdk-go/avalanche-tooling-sdk-go
```

### Permission denied errors

If you encounter permission errors:

```bash
# Fix ownership of Go installation
sudo chown -R $USER:$USER /usr/local/go

# Or install Go in your home directory
mkdir ~/go
cd ~/go
wget https://go.dev/dl/go1.21.11.linux-amd64.tar.gz
tar -xzf go1.21.11.linux-amd64.tar.gz
echo "export PATH=~/go/go/bin:\$PATH" >> ~/.bashrc
```

### Module download issues

If you have issues downloading modules:

```bash
# Set GOPROXY if needed
export GOPROXY=https://proxy.golang.org,direct

# Clear module cache and retry
go clean -modcache
go mod download
```

## Next Steps

After successful setup, you can:

1. Explore the examples in the `examples/` directory
2. Read the documentation in the project
3. Configure AWS credentials for cloud functionality
4. Set up a local Avalanche network for testing
5. Consider moving the project to WSL filesystem for better performance if you'll be working extensively with it

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Avalanche Documentation](https://docs.dione.network/)
- [WSL Documentation](https://docs.microsoft.com/en-us/windows/wsl/)

## Support

If you encounter issues:

1. Check the troubleshooting section above
2. Verify your Go installation with `go version`
3. Ensure all dependencies are downloaded with `go mod verify`
4. Check the project's GitHub issues for known problems

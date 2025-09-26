# Feature Flags Testing Guide

This guide explains how to test the new feature flags that control various functionalities in the Odyssey Tooling SDK.

## Overview

The following feature flags have been implemented:

- **AWS Integration**: `AWSIntegrationEnabled`
- **GCP Integration**: `GCPIntegrationEnabled`
- **Docker Support**: `DockerSupportEnabled`
- **Instance Management**: `InstanceManagementEnabled`
- **Security Groups**: `SecurityGroupsEnabled`
- **SSH Key Management**: `SSHKeyManagementEnabled`

## Running Tests

### Quick Test (Windows)
```bash
test_feature_flags.bat
```

### Quick Test (Linux/macOS)
```bash
./test_feature_flags.sh
```

### Manual Testing

#### 1. Test Individual Feature Flags
```bash
# Test AWS integration flag
go test -v ./node -run TestFeatureFlags_AWSIntegration

# Test Docker support flag
go test -v ./node -run TestFeatureFlags_DockerSupport

# Test security groups flag
go test -v ./cloud/aws -run TestFeatureFlags_SecurityGroups
```

#### 2. Test All Feature Flags
```bash
# Test all feature flag tests
go test -v ./node -run TestFeatureFlags
go test -v ./cloud/aws -run TestFeatureFlags
go test -v ./examples -run TestFeatureFlags
```

#### 3. Test Integration Scenarios
```bash
# Test example scenarios
go test -v ./examples -run TestFeatureFlags_Example
```

## Test Structure

### 1. Unit Tests (`node/feature_flags_test.go`)

Tests individual feature flags in isolation:

- `TestFeatureFlags_AWSIntegration` - Tests AWS integration gating
- `TestFeatureFlags_GCPIntegration` - Tests GCP integration gating
- `TestFeatureFlags_DockerSupport` - Tests Docker support gating
- `TestFeatureFlags_InstanceManagement` - Tests instance management gating
- `TestFeatureFlags_SSHKeyManagement` - Tests SSH key management gating
- `TestFeatureFlags_DockerFunctions` - Tests individual Docker function gating
- `TestFeatureFlags_SSHConnect` - Tests SSH connection gating
- `TestFeatureFlags_MultipleFlags` - Tests multiple flags disabled
- `TestFeatureFlags_AllEnabled` - Tests all flags enabled

### 2. AWS-Specific Tests (`cloud/aws/feature_flags_test.go`)

Tests AWS-specific functionality:

- `TestFeatureFlags_SecurityGroups` - Tests security group creation gating
- `TestFeatureFlags_SecurityGroupsInstance` - Tests security group instance method gating

### 3. Integration Tests (`examples/feature_flags_example_test.go`)

Tests real-world scenarios:

- `TestFeatureFlags_Example` - Demonstrates feature flag testing patterns
- `TestFeatureFlags_IndividualFunctions` - Tests individual function gating

## Testing Patterns

### 1. Save and Restore Pattern

```go
func TestFeatureFlag(t *testing.T) {
    // Save original value
    originalValue := constants.FeatureFlagName
    defer func() {
        constants.FeatureFlagName = originalValue
    }()
    
    // Test with flag disabled
    constants.FeatureFlagName = false
    // ... test logic ...
    
    // Test with flag enabled
    constants.FeatureFlagName = true
    // ... test logic ...
}
```

### 2. Error Message Validation

```go
// Should fail with specific error message
_, err := SomeFunction()
assert.Error(t, err)
assert.Contains(t, err.Error(), "functionality is disabled")
assert.Contains(t, err.Error(), "Set constants.FeatureFlagName = true to enable")
```

### 3. Multiple Flag Testing

```go
func TestMultipleFlags(t *testing.T) {
    // Save all original values
    originalAWS := constants.AWSIntegrationEnabled
    originalDocker := constants.DockerSupportEnabled
    // ... save others ...
    
    defer func() {
        constants.AWSIntegrationEnabled = originalAWS
        constants.DockerSupportEnabled = originalDocker
        // ... restore others ...
    }()
    
    // Test combinations
    constants.AWSIntegrationEnabled = false
    constants.DockerSupportEnabled = true
    // ... test logic ...
}
```

## Expected Test Results

### When Flags are Disabled
- Functions should return errors with clear messages
- Error messages should include instructions on how to enable the functionality
- No actual functionality should be executed

### When Flags are Enabled
- Functions should pass the feature flag check
- Functions may still fail due to other missing dependencies (credentials, etc.)
- Error messages should NOT contain "functionality is disabled"

## Manual Testing

### 1. Disable a Feature Flag

Edit `constants/constants.go`:
```go
var (
    AWSIntegrationEnabled = false  // Disable AWS integration
    // ... other flags ...
)
```

### 2. Test Your Code

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
    "github.com/DioneProtocol/odyssey-tooling-sdk-go/node"
    "github.com/DioneProtocol/odyssey-tooling-sdk-go/odyssey"
)

func main() {
    ctx := context.Background()
    nodeParams := &node.NodeParams{
        CloudParams: &node.CloudParams{
            AWSConfig: &node.AWSConfig{
                // ... AWS config ...
            },
            // ... other config ...
        },
        // ... other params ...
    }
    
    nodes, err := node.CreateNodes(ctx, nodeParams)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        // Should contain "AWS integration functionality is disabled"
    }
}
```

### 3. Verify Error Messages

The error messages should be clear and actionable:
- "AWS integration functionality is disabled. Set constants.AWSIntegrationEnabled = true to enable"
- "Docker support functionality is disabled. Set constants.DockerSupportEnabled = true to enable"
- "SSH key management functionality is disabled. Set constants.SSHKeyManagementEnabled = true to enable"

## Continuous Integration

Add these tests to your CI pipeline:

```yaml
# .github/workflows/test-feature-flags.yml
name: Test Feature Flags
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - uses: actions/setup-go@v2
      with:
        go-version: 1.20
    - name: Test Feature Flags
      run: |
        go test -v ./node -run TestFeatureFlags
        go test -v ./cloud/aws -run TestFeatureFlags
        go test -v ./examples -run TestFeatureFlags
```

## Troubleshooting

### Common Issues

1. **"cannot assign to constants.FeatureFlagName"**
   - Solution: Feature flags are now `var` declarations, not `const`
   - Make sure you're using the latest version of the constants file

2. **Tests pass when they should fail**
   - Check that the feature flag is actually being checked in the function
   - Verify the flag is set to `false` in your test

3. **Tests fail when they should pass**
   - Check that all required flags are set to `true`
   - Verify the function is not failing for other reasons (missing credentials, etc.)

### Debug Tips

1. **Add logging to see flag values**:
   ```go
   fmt.Printf("AWSIntegrationEnabled: %v\n", constants.AWSIntegrationEnabled)
   ```

2. **Test individual functions**:
   ```go
   // Test just the flag check
   if !constants.AWSIntegrationEnabled {
       t.Log("AWS integration is disabled")
   }
   ```

3. **Use table-driven tests for multiple scenarios**:
   ```go
   tests := []struct {
       name     string
       flag     bool
       wantErr  bool
       errorMsg string
   }{
       {"enabled", true, false, ""},
       {"disabled", false, true, "functionality is disabled"},
   }
   ```

## Best Practices

1. **Always restore original values** in tests using `defer`
2. **Test both enabled and disabled states** for each flag
3. **Validate error messages** contain helpful instructions
4. **Test combinations** of multiple flags
5. **Use descriptive test names** that indicate what's being tested
6. **Keep tests focused** - one flag per test when possible
7. **Test edge cases** like all flags disabled/enabled

## Conclusion

The feature flags provide a robust way to control functionality in the Odyssey Tooling SDK. The comprehensive test suite ensures that:

- Flags work correctly when enabled/disabled
- Error messages are clear and actionable
- No regressions are introduced
- Integration scenarios work as expected

Use this guide to test the feature flags in your own code and ensure they work correctly for your use cases.

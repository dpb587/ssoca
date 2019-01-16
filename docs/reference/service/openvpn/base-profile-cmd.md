---
title: ssoca openvpn base-profile
aliases:
- /service/openvpn/base-profile-cmd/
---

# `ssoca openvpn base-profile ...`

Show the base connection profile of the OpenVPN server

    Usage:
      ssoca [OPTIONS] openvpn [openvpn-OPTIONS] base-profile [base-profile-OPTIONS]
    
    Application Options:
          --config=              Configuration file path (default: ~/.config/ssoca/config) [$SSOCA_CONFIG]
      -e, --environment=         Environment name [$SSOCA_ENVIRONMENT]
          --log-level=           Log level (default: WARN) [$SSOCA_LOG_LEVEL]
    
    Help Options:
      -h, --help                 Show this help message
    
    [openvpn command options]
    
        Establish OpenVPN connections to remote servers:
          -s, --service=         Service name (default: openvpn) [$SSOCA_SERVICE]
    
    [base-profile command options]
              --skip-auth-retry  Skip interactive authentication retries when logged out
    

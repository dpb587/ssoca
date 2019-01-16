---
title: ssoca openvpn create-profile
aliases:
- /service/openvpn/create-profile-cmd/
---

# `ssoca openvpn create-profile ...`

Create and sign an OpenVPN configuration profile

    Usage:
      ssoca [OPTIONS] openvpn [openvpn-OPTIONS] create-profile [create-profile-OPTIONS]
    
    Application Options:
          --config=              Configuration file path (default: ~/.config/ssoca/config) [$SSOCA_CONFIG]
      -e, --environment=         Environment name [$SSOCA_ENVIRONMENT]
          --log-level=           Log level (default: WARN) [$SSOCA_LOG_LEVEL]
    
    Help Options:
      -h, --help                 Show this help message
    
    [openvpn command options]
    
        Establish OpenVPN connections to remote servers:
          -s, --service=         Service name (default: openvpn) [$SSOCA_SERVICE]
    
    [create-profile command options]
              --skip-auth-retry  Skip interactive authentication retries when logged out
    

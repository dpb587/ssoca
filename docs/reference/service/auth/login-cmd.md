---
title: ssoca auth login
aliases:
- /service/auth/login-cmd/
---

# `ssoca auth login ...`

Authenticate for a new token

    Usage:
      ssoca [OPTIONS] auth [auth-OPTIONS] login [login-OPTIONS]
    
    Application Options:
          --config=           Configuration file path (default: ~/.config/ssoca/config) [$SSOCA_CONFIG]
      -e, --environment=      Environment name [$SSOCA_ENVIRONMENT]
          --log-level=        Log level (default: WARN) [$SSOCA_LOG_LEVEL]
    
    Help Options:
      -h, --help              Show this help message
    
    [auth command options]
    
        Manage authentication:
          -s, --service=      Service name [$SSOCA_SERVICE]
    
    [login command options]
              --skip-verify   Skip verification of authentication, once complete
              --wait-timeout= Timeout to wait for authentication before erroring (default: 15m)
    

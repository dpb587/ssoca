---
title: ssoca ssh sign-public-key
aliases:
- /service/ssh/sign-public-key-cmd/
---

# `ssoca ssh sign-public-key ...`

Create a certificate for a specific public key

    Usage:
      ssoca [OPTIONS] ssh [ssh-OPTIONS] sign-public-key [sign-public-key-OPTIONS] PATH
    
    Application Options:
          --config=              Configuration file path (default: ~/.config/ssoca/config) [$SSOCA_CONFIG]
      -e, --environment=         Environment name [$SSOCA_ENVIRONMENT]
          --log-level=           Log level (default: WARN) [$SSOCA_LOG_LEVEL]
    
    Help Options:
      -h, --help                 Show this help message
    
    [ssh command options]
    
        Establish SSH connections to remote servers:
          -s, --service=         Service name (default: ssh) [$SSOCA_SERVICE]
    
    [sign-public-key command options]
              --skip-auth-retry  Skip interactive authentication retries when logged out
    

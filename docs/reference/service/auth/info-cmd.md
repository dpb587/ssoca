---
title: ssoca auth info
aliases:
- /service/auth/info-cmd/
---

# `ssoca auth info ...`

Show current authentication information

    Usage:
      ssoca [OPTIONS] auth [auth-OPTIONS] info [info-OPTIONS]
    
    Application Options:
          --config=            Configuration file path (default: ~/.config/ssoca/config) [$SSOCA_CONFIG]
      -e, --environment=       Environment name [$SSOCA_ENVIRONMENT]
          --log-level=         Log level (default: WARN) [$SSOCA_LOG_LEVEL]
    
    Help Options:
      -h, --help               Show this help message
    
    [auth command options]
    
        Manage authentication:
          -s, --service=       Service name [$SSOCA_SERVICE]
    
    [info command options]
              --authenticated  Show only whether the user is authenticated
              --id             Show only the ID of the authenticated user
              --groups         Show only the groups of the authenticated user
    

---
title: ssoca env set
aliases:
- /service/env/set-cmd/
---

# `ssoca env set ...`

Configure the connection to an environment (aliases: add)

    Usage:
      ssoca [OPTIONS] env set [set-OPTIONS] [URL]
    
    Application Options:
          --config=          Configuration file path (default: ~/.config/ssoca/config) [$SSOCA_CONFIG]
      -e, --environment=     Environment name [$SSOCA_ENVIRONMENT]
          --log-level=       Log level (default: WARN) [$SSOCA_LOG_LEVEL]
    
    Help Options:
      -h, --help             Show this help message
    
    [set command options]
              --ca-cert=     Environment CA certificate path
              --skip-verify  Skip verification of environment availability
    
    [set command arguments]
      URL:                   Environment URL
    

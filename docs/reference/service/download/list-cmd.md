---
title: ssoca download list
aliases:
- /service/download/list-cmd/
---

# `ssoca download list ...`

List available artifacts

    Usage:
      ssoca [OPTIONS] download [download-OPTIONS] list [list-OPTIONS]
    
    Application Options:
          --config=              Configuration file path (default: ~/.config/ssoca/config) [$SSOCA_CONFIG]
      -e, --environment=         Environment name [$SSOCA_ENVIRONMENT]
          --log-level=           Log level (default: WARN) [$SSOCA_LOG_LEVEL]
    
    Help Options:
      -h, --help                 Show this help message
    
    [download command options]
    
        Download environment artifacts:
          -s, --service=         Service name (default: download) [$SSOCA_SERVICE]
    
    [list command options]
              --skip-auth-retry  Skip interactive authentication retries when logged out
    

---
title: ssoca download get
aliases:
- /service/download/get-cmd/
---

# `ssoca download get ...`

Get an artifact

    Usage:
      ssoca [OPTIONS] download [download-OPTIONS] get [get-OPTIONS] [FILE] [TARGET-FILE]
    
    Application Options:
          --config=              Configuration file path (default: ~/.config/ssoca/config) [$SSOCA_CONFIG]
      -e, --environment=         Environment name [$SSOCA_ENVIRONMENT]
          --log-level=           Log level (default: WARN) [$SSOCA_LOG_LEVEL]
    
    Help Options:
      -h, --help                 Show this help message
    
    [download command options]
    
        Download environment artifacts:
          -s, --service=         Service name (default: download) [$SSOCA_SERVICE]
    
    [get command options]
              --skip-auth-retry  Skip interactive authentication retries when logged out
    
    [get command arguments]
      FILE:                      File name
      TARGET-FILE:               Target path to write download (use '-' for STDOUT)
    

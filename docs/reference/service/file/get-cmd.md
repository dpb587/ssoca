---
title: ssoca file get
aliases:
- /service/file/get-cmd/
- /reference/service/download/get-cmd/
---

# `ssoca file get ...`

Download a file and verify its checksum (aliases: download)

    Usage:
      ssoca [OPTIONS] file [file-OPTIONS] get [get-OPTIONS] [FILE] [TARGET-FILE]
    
    Application Options:
          --config=              Configuration file path (default: ~/.config/ssoca/config) [$SSOCA_CONFIG]
      -e, --environment=         Environment name [$SSOCA_ENVIRONMENT]
          --log-level=           Log level (default: WARN) [$SSOCA_LOG_LEVEL]
    
    Help Options:
      -h, --help                 Show this help message
    
    [file command options]
    
        Access files from the environment:
          -s, --service=         Service name (default: file) [$SSOCA_SERVICE]
    
    [get command options]
              --skip-auth-retry  Skip interactive authentication retries when logged out
    
    [get command arguments]
      FILE:                      File name
      TARGET-FILE:               Target path to write download (use '-' for STDOUT)
    

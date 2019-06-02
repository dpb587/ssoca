---
title: ssoca file exec
aliases:
- /service/file/exec-cmd/
---

# `ssoca file exec ...`

Temporarily get and then execute a file

    Usage:
      ssoca [OPTIONS] file [file-OPTIONS] exec [exec-OPTIONS] [FILE] [EXTRA...]
    
    Application Options:
          --config=              Configuration file path (default: ~/.config/ssoca/config) [$SSOCA_CONFIG]
      -e, --environment=         Environment name [$SSOCA_ENVIRONMENT]
          --log-level=           Log level (default: WARN) [$SSOCA_LOG_LEVEL]
    
    Help Options:
      -h, --help                 Show this help message
    
    [file command options]
    
        Access files from the environment:
          -s, --service=         Service name (default: file) [$SSOCA_SERVICE]
    
    [exec command options]
              --skip-auth-retry  Skip interactive authentication retries when logged out
    
    [exec command arguments]
      FILE:                      File name
      EXTRA:                     Additional arguments to pass
    

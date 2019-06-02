---
title: ssoca file list
aliases:
- /service/file/list-cmd/
- /reference/service/download/list-cmd/
---

# `ssoca file list ...`

List available files

    Usage:
      ssoca [OPTIONS] file [file-OPTIONS] list [list-OPTIONS]
    
    Application Options:
          --config=              Configuration file path (default: ~/.config/ssoca/config) [$SSOCA_CONFIG]
      -e, --environment=         Environment name [$SSOCA_ENVIRONMENT]
          --log-level=           Log level (default: WARN) [$SSOCA_LOG_LEVEL]
    
    Help Options:
      -h, --help                 Show this help message
    
    [file command options]
    
        Access files from the environment:
          -s, --service=         Service name (default: file) [$SSOCA_SERVICE]
    
    [list command options]
              --skip-auth-retry  Skip interactive authentication retries when logged out
    

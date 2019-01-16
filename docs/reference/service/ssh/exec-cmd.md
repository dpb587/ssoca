---
title: ssoca ssh exec
aliases:
- /service/ssh/exec-cmd/
---

# `ssoca ssh exec ...`

Connect to a remote SSH server

    Usage:
      ssoca [OPTIONS] ssh [ssh-OPTIONS] exec [exec-OPTIONS] [HOST]
    
    Application Options:
          --config=              Configuration file path (default: ~/.config/ssoca/config) [$SSOCA_CONFIG]
      -e, --environment=         Environment name [$SSOCA_ENVIRONMENT]
          --log-level=           Log level (default: WARN) [$SSOCA_LOG_LEVEL]
    
    Help Options:
      -h, --help                 Show this help message
    
    [ssh command options]
    
        Establish SSH connections to remote servers:
          -s, --service=         Service name (default: ssh) [$SSOCA_SERVICE]
    
    [exec command options]
              --skip-auth-retry  Skip interactive authentication retries when logged out
              --exec=            Path to the ssh binary
              --opt=             Additional option to pass to ssh
    

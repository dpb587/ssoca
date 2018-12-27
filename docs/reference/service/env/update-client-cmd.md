---
title: ssoca env update-client
aliases:
- /service/env/update-client-cmd/
---

# `ssoca env update-client ...`

Download the latest client from the environment

    Usage:
      ssoca [OPTIONS] env update-client [update-client-OPTIONS]
    
    Application Options:
          --config=              Configuration file path (default: ~/.config/ssoca/config) [$SSOCA_CONFIG]
      -e, --environment=         Environment name [$SSOCA_ENVIRONMENT]
          --log-level=           Log level (default: WARN) [$SSOCA_LOG_LEVEL]
    
    Help Options:
      -h, --help                 Show this help message
    
    [update-client command options]
          -s, --service=         Service name (default: env) [$SSOCA_SERVICE]
              --skip-auth-retry  Skip interactive authentication retries when logged out
    

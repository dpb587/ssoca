---
title: ssoca env rename
aliases:
- /service/env/rename-cmd/
---

# `ssoca env rename ...`

Set a new name for the environment

    Usage:
      ssoca [OPTIONS] env rename [rename-OPTIONS] [NEW-NAME]
    
    Application Options:
          --config=       Configuration file path (default: ~/.config/ssoca/config) [$SSOCA_CONFIG]
      -e, --environment=  Environment name [$SSOCA_ENVIRONMENT]
          --log-level=    Log level (default: WARN) [$SSOCA_LOG_LEVEL]
    
    Help Options:
      -h, --help          Show this help message
    
    [rename command options]
          -s, --service=  Service name (default: env) [$SSOCA_SERVICE]
    
    [rename command arguments]
      NEW-NAME:           New environment name
    

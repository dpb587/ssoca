# `ssoca download get ...`

Get an artifact

    Usage:
      ssoca [OPTIONS] download get [get-OPTIONS] [FILE] [TARGET-FILE]
    
    Application Options:
          --config=              Configuration file path (default: ~/.config/ssoca/config) [$SSOCA_CONFIG]
      -e, --environment=         Environment name [$SSOCA_ENVIRONMENT]
    
    Help Options:
      -h, --help                 Show this help message
    
    [get command options]
          -s, --service=         Service name (default: download) [$SSOCA_SERVICE]
              --skip-auth-retry  Skip interactive authentication retries when logged out
    
    [get command arguments]
      FILE:                      File name
      TARGET-FILE:               Target path to write download
    

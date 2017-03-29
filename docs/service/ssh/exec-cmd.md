# `ssoca ssh exec ...`

Connect to a remote SSH server

    Usage:
      ssoca [OPTIONS] ssh exec [exec-OPTIONS] [HOST]
    
    Application Options:
          --config=      Configuration file path (default: ~/.ssoca/config) [$SSOCA_CONFIG]
      -e, --environment= Environment name [$SSOCA_ENVIRONMENT]
    
    Help Options:
      -h, --help         Show this help message
    
    [exec command options]
          -s, --service= Service name (default: ssh) [$SSOCA_SERVICE]
              --opts=    Options to pass through to SSH
    

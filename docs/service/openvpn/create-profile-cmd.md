# `ssoca openvpn create-profile ...`

Create and sign an OpenVPN configuration profile

    Usage:
      ssoca [OPTIONS] openvpn create-profile [create-profile-OPTIONS]
    
    Application Options:
          --config=              Configuration file path (default: ~/.config/ssoca/config) [$SSOCA_CONFIG]
      -e, --environment=         Environment name [$SSOCA_ENVIRONMENT]
    
    Help Options:
      -h, --help                 Show this help message
    
    [create-profile command options]
          -s, --service=         Service name (default: openvpn) [$SSOCA_SERVICE]
              --skip-auth-retry  Skip interactive authentication retries when logged out
    

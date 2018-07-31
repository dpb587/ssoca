# `ssoca openvpn create-launchd-service ...`

Create a launchd service

    Usage:
      ssoca [OPTIONS] openvpn create-launchd-service [create-launchd-service-OPTIONS] [DESTINATION-DIR]
    
    Application Options:
          --config=              Configuration file path (default: ~/.config/ssoca/config) [$SSOCA_CONFIG]
      -e, --environment=         Environment name [$SSOCA_ENVIRONMENT]
          --log-level=           Log level (default: WARN) [$SSOCA_LOG_LEVEL]
    
    Help Options:
      -h, --help                 Show this help message
    
    [create-launchd-service command options]
          -s, --service=         Service name (default: openvpn) [$SSOCA_SERVICE]
              --skip-auth-retry  Skip interactive authentication retries when logged out
              --exec-ssoca=      Path to the ssoca binary (default: ssoca)
              --name=            Specific file name to use for *.tblk
              --exec-openvpn=    Path to the openvpn binary
              --run-at-load      Run the service at load
              --log-dir=         Log directory for the service (default: ~/Library/Logs)
              --start            Load and start the service after installation
    
    [create-launchd-service command arguments]
      DESTINATION-DIR:           Directory where the *.plist service will be created (default: ~/Library/LaunchAgents)
    

## Usage Details

To create and automatically start a VPN profile via [launchd](https://developer.apple.com/library/archive/documentation/MacOSX/Conceptual/BPSystemStartup/Chapters/CreatingLaunchdJobs.html#//apple_ref/doc/uid/10000172i-SW7-BCIEDDBJ) you may use the `--start` option...

    $ ssoca openvpn create-launchd-service --start
    The service 'acme-prod-aws-use1.openvpn.ssoca.dpb587.github.io' has successfully been started.

By default, service names are suffixed with a global ssoca-based domain. Use the `--name` flag to choose your service's own FQDN.

To remove a service, be sure to stop, unload, and remove it...

    $ launchctl stop acme-prod-aws-use1
    $ launchctl unload ~/Library/LaunchAgents/acme-prod-aws-use1.plist
    $ rm ~/Library/LaunchAgents/acme-prod-aws-use1.plist

If you are experiencing issues, you can find logs in `~/Library/Logs/$name.*.log`.


## Notes

Reminder: `ssoca` may require interactive authentication or `sudo` privileges. Depending on authentication strategies, browser-based access may work successfully. For avoiding `sudo` restrictions, you may want to consider specifying an `openvpn` wrapper with sudo or SUID privileges (via the `--exec-openvpn` flag).

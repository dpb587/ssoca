---
title: ssoca openvpn create-tunnelblick-profile
aliases:
- /service/openvpn/create-tunnelblick-profile-cmd/
---

# `ssoca openvpn create-tunnelblick-profile ...`

Create a Tunnelblick profile

    Usage:
      ssoca [OPTIONS] openvpn [openvpn-OPTIONS] create-tunnelblick-profile [create-tunnelblick-profile-OPTIONS] [DESTINATION-DIR]
    
    Application Options:
          --config=              Configuration file path (default: ~/.config/ssoca/config) [$SSOCA_CONFIG]
      -e, --environment=         Environment name [$SSOCA_ENVIRONMENT]
          --log-level=           Log level (default: WARN) [$SSOCA_LOG_LEVEL]
    
    Help Options:
      -h, --help                 Show this help message
    
    [openvpn command options]
    
        Establish OpenVPN connections to remote servers:
          -s, --service=         Service name (default: openvpn) [$SSOCA_SERVICE]
    
    [create-tunnelblick-profile command options]
              --skip-auth-retry  Skip interactive authentication retries when logged out
              --exec-ssoca=      Path to the ssoca binary (default: ssoca)
              --name=            Specific file name to use for *.tblk
              --install          Install the profile (sudo may prompt for privileges)
    
    [create-tunnelblick-profile command arguments]
      DESTINATION-DIR:           Directory where the *.tblk profile will be created (default: $PWD)
    

## Usage Details

To create and automatically install a Tunnelblick profile you may use the `--install` option...

    $ ssoca openvpn create-tunnelblick-profile --install
    The profile 'acme-prod-aws-use1' has successfully been installed.

To create only the profile files, leave off `--install` (it will create a `{service-name}.tblk` in your current directory)...

    $ ssoca openvpn create-tunnelblick-profile

Once you are ready to install the profile, open the file to register it with Tunnelblick...

    $ open openvpn.tblk

Tunnelblick will give two prompts regarding the security concern of profiles with custom scripts, followed by a prompt about which context to install the profile, and finally an authentication prompt to install the profile. Once installed, you can connect/disconnect like other OpenVPN profiles.


## Notes

The generated profile includes scripts to automatically regenerate connection certificates upon establishing a VPN connection. There are a few caveats...

 * The scripts execute as `root`, therefore Tunnelblick gives very clear warnings and confirmations when first installing the profile. It is referring to the `pre-connect.sh` file inside the generated `*.tblk` directory.
 * For security, Tunnelblick may keep a shadow copy of profile configuration which requires a user authentication prompt whenever the configuration is changed. The ssoca profiles have a short lifetime and are regenerated every connection. To avoid this user intervention, the background scripts automatically updates the shadow copy.
 * If your authentication session expires, Tunnelblick may error during the next connection attempt with an [ambiguous] error dialog.
 * When the connection is lost, the `openvpn` process will exit and, by default, not attempt to reconnect. To change this behavior, navigate to the Advanced settings of your profile and enable the "Keep connected" checkbox of the "Connecting & Disconnecting" tab.
